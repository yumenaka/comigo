package routers

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

// 验证运行中启用密码立即保护所有控制资源，密码变更立即使旧 Bearer 令牌失效。
func TestControlAuthChangesWithoutRebinding(t *testing.T) {
	restore := withRouterAuthTestConfig(t)
	defer restore()
	config.GetCfg().Password = ""
	oldEngine := engine
	t.Cleanup(func() { engine = oldEngine })
	engine = echo.New()
	BindURLs()
	request := func(method, path, token string) int {
		req := httptest.NewRequest(method, path, nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)
		return rec.Code
	}
	if got := request(http.MethodGet, "/api/connections", ""); got != http.StatusOK {
		t.Fatalf("公开模式: %d", got)
	}
	config.GetCfg().Password = "secret"
	for _, route := range []struct{ method, path string }{
		{http.MethodGet, "/api/server"}, {http.MethodGet, "/api/connections"},
		{http.MethodGet, "/api/stores"}, {http.MethodGet, "/api/stores/test"},
		{http.MethodGet, "/api/books/test"}, {http.MethodPost, "/api/restart"},
		{http.MethodPost, "/api/stores/test/refresh"}, {http.MethodDelete, "/api/stores/test"},
	} {
		if got := request(route.method, route.path, ""); got != http.StatusUnauthorized {
			t.Fatalf("%s %s: %d", route.method, route.path, got)
		}
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": "comigo", "exp": time.Now().Add(time.Hour).Unix()}).SignedString([]byte(config.GetJwtSigningKey()))
	if err != nil {
		t.Fatal(err)
	}
	if got := request(http.MethodGet, "/api/connections", token); got != http.StatusOK {
		t.Fatalf("Bearer: %d", got)
	}
	config.GetCfg().Password = "changed"
	if got := request(http.MethodGet, "/api/connections", token); got != http.StatusUnauthorized {
		t.Fatalf("旧令牌: %d", got)
	}
}
