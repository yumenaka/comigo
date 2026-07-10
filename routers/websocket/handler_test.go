package websocket

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gorilla "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func TestShouldLogWebSocketReadError(t *testing.T) {
	if shouldLogWebSocketReadError(net.ErrClosed) {
		t.Fatal("local shutdown should not be logged as websocket failure")
	}
	if !shouldLogWebSocketReadError(&gorilla.CloseError{Code: gorilla.CloseProtocolError}) {
		t.Fatal("protocol error should still be logged")
	}
}

// TestWebSocketOriginAllowed 验证浏览器连接必须同源，同时保留反代、原生与 Wails 客户端。
func TestWebSocketOriginAllowed(t *testing.T) {
	cases := []struct {
		origin        string
		forwardedHost string
		want          bool
	}{
		{origin: "", want: true},
		{origin: "http://localhost:1234", want: true},
		{origin: "https://reader.example.com", forwardedHost: "reader.example.com", want: true},
		{origin: "https://reader.example.com", forwardedHost: "reader.example.com, proxy.internal", want: true},
		{origin: "https://example.com", want: false},
		{origin: "wails://wails", want: true},
		{origin: "not a url", want: false},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:1234/api/ws", nil)
		req.Header.Set("Origin", tc.origin)
		req.Header.Set("X-Forwarded-Host", tc.forwardedHost)
		if got := websocketOriginAllowed(req); got != tc.want {
			t.Fatalf("origin %q allowed = %v, want %v", tc.origin, got, tc.want)
		}
	}
}

// TestCloseAllClosesActiveClients 验证退出流程能主动断开阅读页 WebSocket 长连接。
func TestCloseAllClosesActiveClients(t *testing.T) {
	CloseAll()
	t.Cleanup(CloseAll)

	e := echo.New()
	e.GET("/api/ws", WsHandler)
	server := httptest.NewServer(e)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/ws"
	conn, _, err := gorilla.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer conn.Close()

	waitForClientCount(t, 1)
	CloseAll()
	waitForClientCount(t, 0)

	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("set read deadline: %v", err)
	}
	if _, _, err := conn.ReadMessage(); err == nil {
		t.Fatal("expected websocket read to fail after CloseAll")
	}
}

// TestWebSocketRejectsOversizedMessage 验证异常客户端不能让服务端无界读取消息。
func TestWebSocketRejectsOversizedMessage(t *testing.T) {
	CloseAll()
	t.Cleanup(CloseAll)

	e := echo.New()
	e.GET("/api/ws", WsHandler)
	server := httptest.NewServer(e)
	defer server.Close()
	conn, _, err := gorilla.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/api/ws", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := conn.WriteMessage(gorilla.TextMessage, []byte(`{"detail":"`+strings.Repeat("x", maxWebSocketMessageSize)+`"}`)); err != nil {
		t.Fatal(err)
	}
	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if _, _, err := conn.ReadMessage(); err == nil {
		t.Fatal("oversized websocket message should close the connection")
	}
}

func waitForClientCount(t *testing.T, want int) {
	t.Helper()
	for range 50 {
		if got := clientCount(); got == want {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("client count = %d, want %d", clientCount(), want)
}

func clientCount() int {
	clientsMu.RLock()
	defer clientsMu.RUnlock()
	return len(clients)
}
