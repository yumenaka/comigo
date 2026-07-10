package logger

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/tools/sse_hub"
)

// TestEchoLogHandlerEscapesRequestFields 防止外部请求字段进入设置页日志 HTML。
func TestEchoLogHandlerEscapesRequestFields(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set(echo.HeaderXForwardedFor, `<img src=x onerror="alert(1)">`)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	events := make(chan sse_hub.Event, 1)
	sse_hub.MessageHub.Add(t.Name(), events)
	defer sse_hub.MessageHub.Remove(t.Name())
	handler := EchoLogHandler(false, "", "", false)(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	if err := handler(c); err != nil {
		t.Fatal(err)
	}

	event := <-events
	if strings.Contains(event.Data, "<img") || !strings.Contains(event.Data, "&lt;img") {
		t.Fatalf("request field was not HTML escaped: %s", event.Data)
	}
}
