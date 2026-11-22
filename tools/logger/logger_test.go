package logger

import (
	"io"
	"os"
	"path/filepath"
	"sync"

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
	// Echo 默认忽略转发头；仅在测试中注入原始字段，继续覆盖日志层的转义。
	e.IPExtractor = func(r *http.Request) string { return r.Header.Get(echo.HeaderXForwardedFor) }
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

// TestFileLoggingSingleOutput 验证首次创建、重配置、轮转、关闭和失败回退均不重复或丢失镜像日志。
func TestFileLoggingSingleOutput(t *testing.T) {
	oldOutput, oldMirror, oldSuppressed := baseOutput, mirrorOutput, IsStdoutSuppressed()
	t.Cleanup(func() {
		EchoLogHandler(false, "", "", false)
		SetOutput(oldOutput)
		SetMirrorOutput(oldMirror)
		SetSuppressStdout(oldSuppressed)
	})
	var output, mirror bytes.Buffer
	SetOutput(&output)
	SetMirrorOutput(&mirror)
	SetSuppressStdout(false)
	dir := t.TempDir()
	EchoLogHandler(true, dir, "app.log", true)
	Info("first-message")
	if !strings.Contains(output.String(), "logger_test.go:") {
		t.Fatal("Debug 日志缺少调用位置")
	}
	EchoLogHandler(true, dir, "app.log", false)
	Info("second-message")
	data, err := os.ReadFile(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"first-message", "second-message"} {
		if bytes.Count(data, []byte(text)) != 1 || strings.Count(mirror.String(), text) != 1 {
			t.Fatalf("日志重复或丢失: %s / %s", data, mirror.String())
		}
	}
	// 终端管道关闭后，后续镜像和文件仍各收到一次日志。
	reader, writer := io.Pipe()
	_ = reader.Close()
	SetOutput(writer)
	Info("closed-console")
	_ = writer.Close()
	data, err = os.ReadFile(filepath.Join(dir, "app.log"))
	if err != nil || bytes.Count(data, []byte("closed-console")) != 1 || strings.Count(mirror.String(), "closed-console") != 1 {
		t.Fatal("终端失败时丢失了文件或镜像日志")
	}
	SetOutput(&output)
	SetSuppressStdout(true)
	output.Reset()
	mirror.Reset()
	if err := fileOutput.Rotate(); err != nil {
		t.Fatal(err)
	}
	Info("rotated-message")
	data, err = os.ReadFile(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatal(err)
	}
	if output.Len() != 0 || !strings.Contains(mirror.String(), "rotated-message") || bytes.Count(data, []byte("rotated-message")) != 1 {
		t.Fatalf("抑制/轮转失效: stdout=%s mirror=%s file=%s", output.String(), mirror.String(), data)
	}
	EchoLogHandler(false, "", "", false)
	Info("file-disabled")
	data, err = os.ReadFile(filepath.Join(dir, "app.log"))
	if err != nil || bytes.Contains(data, []byte("file-disabled")) {
		t.Fatalf("禁用后仍写文件: %s %v", data, err)
	}
	// 用普通文件作为父目录，稳定触发创建失败。
	badPath := filepath.Join(dir, "not-a-directory")
	if err := os.WriteFile(badPath, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	EchoLogHandler(true, badPath, "app.log", false)
	Info("fallback-message")
	if fileOutput != nil || !strings.Contains(mirror.String(), "fallback-message") {
		t.Fatal("文件失败后镜像输出不可用")
	}
}

// TestFileLoggingConcurrentReconfigure 重配置期间的每条日志仍只写入一次，旧文件关闭不能打断写入。
func TestFileLoggingConcurrentReconfigure(t *testing.T) {
	oldOutput, oldMirror := baseOutput, mirrorOutput
	t.Cleanup(func() { EchoLogHandler(false, "", "", false); SetOutput(oldOutput); SetMirrorOutput(oldMirror) })
	SetOutput(io.Discard)
	SetMirrorOutput(nil)
	dir := t.TempDir()
	EchoLogHandler(true, dir, "app.log", false)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			Infof("record-%03d-end", i)
		}
	}()
	go func() {
		defer wg.Done()
		for range 10 {
			EchoLogHandler(true, dir, "app.log", false)
		}
	}()
	wg.Wait()
	data, err := os.ReadFile(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		if bytes.Count(data, []byte(fmt.Sprintf("record-%03d-end", i))) != 1 {
			t.Fatalf("重配置丢失或重复第 %d 条日志", i)
		}
	}
}
