package logger

import (
	"bytes"
	"errors"
	"fmt"
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

// 轮询过滤同时覆盖终端、镜像与网页日志，Debug 模式同样静默，失败必须保留。
func TestEchoLogHandlerPolling(t *testing.T) {
	oldOutput, oldMirror := baseOutput, mirrorOutput
	oldFormatter, oldCaller, oldLevel := logger.Formatter, logger.ReportCaller, logger.GetLevel()
	t.Cleanup(func() {
		SetOutput(oldOutput)
		SetMirrorOutput(oldMirror)
		logger.SetFormatter(oldFormatter)
		logger.SetReportCaller(oldCaller)
		logger.SetLevel(oldLevel)
	})
	for _, tc := range []struct {
		name, route, method string
		debug               bool
		status              int
		err                 error
		logged              bool
	}{
		{"traffic", "/api/server/traffic", "GET", false, 200, nil, false},
		{"base_path", "/books/api/connections", "GET", false, 200, nil, false},
		{"debug", "/api/server/traffic", "GET", true, 200, nil, false},
		{"unauthorized", "/api/connections", "GET", false, 401, nil, true},
		{"http_error", "/api/server/traffic", "GET", false, 503, echo.ErrServiceUnavailable, true},
		{"error", "/api/connections", "GET", false, 500, errors.New("failed"), true},
		{"server", "/api/server", "GET", false, 200, nil, false},
		{"configs", "/api/configs", "GET", false, 200, nil, false},
		{"config_status", "/books/api/configs/status", "GET", false, 200, nil, false},
		{"config_write", "/api/configs", "PATCH", false, 200, nil, true},
		{"config_debug", "/api/configs/status", "GET", true, 200, nil, false},
		{"debug_error", "/api/server", "GET", true, 503, echo.ErrServiceUnavailable, true},
		{"config_forbidden", "/api/configs", "GET", false, 403, nil, true},
		{"other_route", "/api/server/update", "GET", false, 200, nil, true},
		{"other_method", "/api/connections", "POST", false, 200, nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output, mirror bytes.Buffer
			SetOutput(&output)
			SetMirrorOutput(&mirror)
			events := make(chan sse_hub.Event, 1)
			sse_hub.MessageHub.Add(t.Name(), events)
			defer sse_hub.MessageHub.Remove(t.Name())
			e := echo.New()
			c := e.NewContext(httptest.NewRequest(tc.method, tc.route+"?poll=1", nil), httptest.NewRecorder())
			c.SetPath(tc.route)
			handler := EchoLogHandler(false, "", "", tc.debug)(func(c echo.Context) error {
				if tc.err != nil {
					return tc.err
				}
				return c.NoContent(tc.status)
			})
			if err := handler(c); err != tc.err {
				t.Fatalf("中间件改变了返回错误: %v", err)
			}
			if (output.Len() > 0) != tc.logged || (mirror.Len() > 0) != tc.logged || (len(events) > 0) != tc.logged {
				t.Fatalf("日志输出不符: output=%q mirror=%q events=%d", output.String(), mirror.String(), len(events))
			}
			if tc.logged && !strings.Contains(output.String(), fmt.Sprintf("[%s:%d]", tc.method, tc.status)) {
				t.Fatalf("日志状态码不符: %s", output.String())
			}
		})
	}
}
