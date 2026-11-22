//go:build !wails

package tailscale_plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"tailscale.com/client/local"
)

// 验证重载会释放旧监听；用无效 Funnel 端口终止新初始化，无需连接真实 Tailscale。
func TestRunTailscaleClosesPreviousListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	netListener = listener
	t.Cleanup(func() { _ = StopTailscale() })
	if err := RunTailscale(echo.New(), TailscaleConfig{FunnelMode: true, Port: 1}); err == nil {
		t.Fatal("无效 Funnel 端口应返回错误")
	}
	if tsServer != nil || tsHttpServer != nil || netListener != nil || localClient != nil {
		t.Fatal("失败初始化遗留实例")
	}
	// 再次关闭返回 net.ErrClosed，证明重载已释放旧监听而非只覆盖其引用。
	if err := listener.Close(); !errors.Is(err, net.ErrClosed) {
		t.Fatalf("旧监听未关闭: %v", err)
	}
}

// TestStatusSnapshotsAndLifecycle 在并发轮询、身份记录和停止过程中检查快照隔离。
func TestStatusSnapshotsAndLifecycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"BackendState":"Running","AuthURL":"","Self":{"DNSName":"reader.example.","Online":true},"TailscaleIPs":["100.64.0.1"]}`)
	}))
	defer server.Close()
	client := &local.Client{OmitAuth: true, Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "tcp", server.Listener.Addr().String())
	}}
	lifecycleMu.Lock()
	localClient = client
	lifecycleMu.Unlock()
	oldStatus := nowStatus
	nowStatus = &TailscaleStatus{}
	t.Cleanup(func() { _ = StopTailscale(); nowStatus = oldStatus })
	initial, err := GetTailscaleStatus(context.Background())
	if err != nil || initial.FQDN != "reader.example" || !initial.Online {
		t.Fatalf("在线状态错误: %+v %v", initial, err)
	}
	var wg sync.WaitGroup
	for range 6 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 30; i++ {
				nowStatus.AddClientInfo(&TailsClientInfo{LoginUserName: "user", NodeComputedName: "node"})
				status, err := GetTailscaleStatus(context.Background())
				if err != nil {
					t.Error(err)
					return
				}
				// 允许调用方编辑快照，不能影响其他读取者或共享客户端记录。
				status.FQDN = "changed"
				status.Clients[0].LoginUserName = "changed"
				if _, err := json.Marshal(status); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 30 {
			if err := StopTailscale(); err != nil {
				t.Error(err)
			}
		}
	}()
	wg.Wait()
	status, err := GetTailscaleStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if status.BackendState != "Stopped" || status.FQDN != "" || status.Clients[0].LoginUserName != "user" {
		t.Fatalf("共享状态被污染: %+v", status)
	}
	// 普通 HTTP 请求没有 tsnet 客户端，仍走原路由。
	called := false
	handler := TailscaleAuthMiddleware()(func(c echo.Context) error { called = true; return nil })
	req := httptest.NewRequest("GET", "/", nil)
	if err := handler(echo.New().NewContext(req, httptest.NewRecorder())); err != nil || !called {
		t.Fatalf("停止后的请求未继续处理: %v", err)
	}
}

// TestStopTailscaleClosesActiveHTTP 停止必须释放正在处理的连接，而不只是监听器。
func TestStopTailscaleClosesActiveHTTP(t *testing.T) {
	started, done := make(chan struct{}), make(chan struct{})
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(started); <-r.Context().Done(); close(done) })}
	lifecycleMu.Lock()
	netListener = listener
	tsHttpServer = server
	lifecycleMu.Unlock()
	t.Cleanup(func() { _ = StopTailscale() })
	go server.Serve(listener)
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	fmt.Fprint(conn, "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n")
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("请求未开始")
	}
	if err := StopTailscale(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("活动连接未关闭")
	}
}

// TestStopWhileFunnelAwaitsLogin 未登录时从本地状态接口查询并关闭正在等待的 Funnel。
func TestStopWhileFunnelAwaitsLogin(t *testing.T) {
	control := httptest.NewServer(http.NotFoundHandler())
	defer control.Close()
	t.Setenv("TS_CONTROL_URL", control.URL)
	t.Setenv("TS_NO_LOGS_NO_SUPPORT", "true")
	cfg := TailscaleConfig{Hostname: "test-funnel", ConfigDir: t.TempDir(), Port: 8443, FunnelMode: true}
	started := make(chan error, 1)
	go func() { started <- RunTailscale(echo.New(), cfg) }()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		status, err := GetTailscaleStatus(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if status.BackendState != "Stopped" {
			break
		}
		select {
		case err := <-started:
			t.Fatalf("登录前意外结束: %v", err)
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case <-time.After(10 * time.Millisecond):
		}
	}
	if err := StopTailscale(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-started:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("启动未取消: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("停止后 Funnel 启动未返回")
	}
}
