//go:build wails && !js && !bindings

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// TestServeWailsAssetBindsBody 验证 WebView 未提供长度时，Echo 仍能解析 JSON 和表单请求体。
func TestServeWailsAssetBindsBody(t *testing.T) {
	previous := config.Server
	t.Cleanup(func() { config.Server = previous })
	e := echo.New()
	e.POST("/bind", func(c echo.Context) error {
		var input struct {
			URL string `json:"url" form:"url"`
		}
		if err := c.Bind(&input); err != nil {
			return err
		}
		return c.String(http.StatusOK, input.URL)
	})
	config.Server = &http.Server{Handler: e}
	for _, tt := range []struct {
		name, contentType, body, want string
		knownLength                   bool
	}{
		{"JSON 未知长度", echo.MIMEApplicationJSON, `{"url":"https://example.com/"}`, "https://example.com/", false},
		{"表单未知长度", echo.MIMEApplicationForm, "url=https%3A%2F%2Fexample.com%2F", "https://example.com/", false},
		{"JSON 已知长度", echo.MIMEApplicationJSON, `{"url":"https://example.com/"}`, "https://example.com/", true},
		{"无请求体", echo.MIMEApplicationJSON, "", "", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Wails 使用 io.ReadCloser 创建请求；没有 Content-Length 时保留为 0。
			req, err := http.NewRequest(http.MethodPost, "/bind", io.NopCloser(strings.NewReader(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set(echo.HeaderContentType, tt.contentType)
			if tt.knownLength {
				req.ContentLength = int64(len(tt.body))
			}
			if tt.body == "" {
				req.Body = http.NoBody
			}
			rec := httptest.NewRecorder()
			serveWailsAsset(rec, req)
			if rec.Code != http.StatusOK || rec.Body.String() != tt.want {
				t.Fatalf("status=%d body=%q; want 200 %q", rec.Code, rec.Body.String(), tt.want)
			}
		})
	}
}
