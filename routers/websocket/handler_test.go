package websocket

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gorilla "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

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
