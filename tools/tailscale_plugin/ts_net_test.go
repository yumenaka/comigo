//go:build !wails

package tailscale_plugin

import (
	"errors"
	"net"
	"testing"

	"github.com/labstack/echo/v4"
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
	// 再次关闭返回 net.ErrClosed，证明重载已释放旧监听而非只覆盖其引用。
	if err := listener.Close(); !errors.Is(err, net.ErrClosed) {
		t.Fatalf("旧监听未关闭: %v", err)
	}
}
