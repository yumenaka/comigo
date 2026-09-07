package routers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// 验证配置账号密码后实时接口也必须经过登录认证。
func TestRealtimeAPIRequiresAuthWhenPasswordConfigured(t *testing.T) {
	restore := withRouterAuthTestConfig(t)
	defer restore()

	oldEngine := engine
	t.Cleanup(func() {
		engine = oldEngine
	})
	engine = echo.New()
	BindURLs()

	for _, path := range []string{"/api/ws", "/api/sse"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			engine.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("%s status = %d, want %d", path, rec.Code, http.StatusUnauthorized)
			}
		})
	}
}

// 验证 Manual 即使启用登录也保持公开。
func TestManualIsPublicWhenPasswordConfigured(t *testing.T) {
	restore := withRouterAuthTestConfig(t)
	defer restore()
	config.GetCfg().BasePath = "/comics"

	oldEngine := engine
	t.Cleanup(func() { engine = oldEngine })
	engine = echo.New()
	BindURLs()
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/manual/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("/manual/ status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `/comics/assets/static/js/manual.js`) {
		t.Fatal("Manual 静态资源没有跟随 BasePath")
	}
}

// 验证阅读器离线缓存名会带上当前版本号，避免旧缓存长期命中。
func TestRenderReaderServiceWorkerUsesComigoVersion(t *testing.T) {
	got := string(renderReaderServiceWorker([]byte("const CACHE_NAME = __COMIGO_READER_PWA_CACHE_NAME__")))

	if !strings.Contains(got, `const CACHE_NAME = "comigo-reader-pwa-`+config.GetVersion()+`"`) {
		t.Fatalf("service worker cache name 未使用 Comigo 版本: %s", got)
	}
	if strings.Contains(got, "__COMIGO_READER_PWA_CACHE_NAME__") {
		t.Fatalf("service worker cache name 占位符未替换: %s", got)
	}
}

func withRouterAuthTestConfig(t *testing.T) func() {
	t.Helper()
	cfg := config.GetCfg()
	oldUsername := cfg.Username
	oldPassword := cfg.Password
	oldTimeout := cfg.Timeout
	oldBasePath := cfg.BasePath
	cfg.Username = "comigo"
	cfg.Password = "secret"
	cfg.Timeout = 60
	cfg.BasePath = ""
	return func() {
		cfg.Username = oldUsername
		cfg.Password = oldPassword
		cfg.Timeout = oldTimeout
		cfg.BasePath = oldBasePath
	}
}
