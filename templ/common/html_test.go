package common

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// TestHtmlLoadsPathHelpersBeforePageScripts 确认内联脚本执行前就能使用 ComiGoPath。
func TestHtmlLoadsPathHelpersBeforePageScripts(t *testing.T) {
	cfg := config.GetCfg()
	oldBasePath := cfg.BasePath
	t.Cleanup(func() { cfg.BasePath = oldBasePath })
	cfg.BasePath = "/proxy"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/proxy/settings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var html bytes.Buffer
	if err := Html(c, templ.NopComponent, nil).Render(context.Background(), &html); err != nil {
		t.Fatalf("render Html: %v", err)
	}

	output := html.String()
	pathScriptIndex := strings.Index(output, "window.ComiGoPath = function")
	mainScriptIndex := strings.Index(output, "/proxy/assets/dist/main.js")
	if pathScriptIndex < 0 {
		t.Fatalf("ComiGoPath script not found")
	}
	if mainScriptIndex < 0 {
		t.Fatalf("main.js script not found")
	}
	if pathScriptIndex > mainScriptIndex {
		t.Fatalf("ComiGoPath script should be rendered before main.js")
	}
	if count := strings.Count(output, "window.ComiGoPath = function"); count != 1 {
		t.Fatalf("ComiGoPath script count = %d, want 1", count)
	}
	if !strings.Contains(output, `window.ComiGoBasePath = "/proxy";`) {
		t.Fatalf("BasePath script did not include normalized base path")
	}
}
