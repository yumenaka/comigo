package routers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
)

// 测试真实登录路由与动态权限组合，公开元数据不能泄露管理信息。
func TestControlAccessMatrix(t *testing.T) {
	old, oldEngine := config.CopyCfg(), engine
	oldStore := model.IStore
	model.IStore = store.RamStore
	t.Cleanup(func() { model.IStore = oldStore })
	t.Cleanup(func() { *config.GetCfg() = old; engine = oldEngine })
	t.Setenv("HOME", t.TempDir())
	cfg := config.GetCfg()
	cfg.BasePath = "/books"
	cfg.Timeout = 60
	cfg.Username = "reader"
	cfg.Password = "test-password"
	cfg.TemporaryReaderMode = true
	cfg.ReadOnlyMode = false
	engine = echo.New()
	BindURLs()
	request := func(method, path, body, token string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, "/books"+path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)
		return rec
	}
	if rec := request("POST", "/api/login", `{"username":"reader","password":"wrong"}`, ""); rec.Code != http.StatusUnauthorized {
		t.Fatalf("wrong password: %d", rec.Code)
	}
	rec := request("POST", "/api/login", `{"username":"reader","password":"test-password"}`, "")
	var login struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &login); err != nil || login.Token == "" {
		t.Fatalf("login failed: %d", rec.Code)
	}
	for _, readOnly := range []bool{false, true} {
		cfg.ReadOnlyMode = readOnly
		for _, token := range []string{"", login.Token} {
			for _, path := range []string{"/healthz", "/api/info"} {
				rec := request("GET", path, "", token)
				if rec.Code != 200 {
					t.Fatalf("public %s: %d", path, rec.Code)
				}
				if strings.Contains(rec.Body.String(), "StoreUrls") || strings.Contains(rec.Body.String(), "test-password") || strings.Contains(rec.Body.String(), "localIPs") {
					t.Fatal("public metadata leaked private data")
				}
			}
			for _, route := range []struct{ method, path, body string }{
				{"GET", "/api/configs", ""}, {"GET", "/api/autostart", ""},
				{"PUT", "/api/autostart", `{"enabled":true}`},
				{"PUT", "/api/configs/Debug", `{"value":true}`},
				{"PATCH", "/api/configs", `{"DisableLAN":true}`},
				{"PATCH", "/api/configs/login", `{}`},
				{"PUT", "/api/configs/files/HomeDirectory", `{}`},
				{"DELETE", "/api/configs/files/HomeDirectory", `{}`},
				{"POST", "/api/stores/refresh", `{}`},
			} {
				expected := 0
				if token == "" {
					expected = 401
				} else if readOnly && route.method != "GET" {
					expected = 403
				}
				if expected == 0 {
					continue
				} // 本矩阵只执行应被拒绝的系统写入，正常操作另测。
				before := config.CopyCfg()
				rec := request(route.method, route.path, route.body, token)
				if rec.Code != expected {
					t.Fatalf("readonly=%v authenticated=%v %s %s: got %d want %d", readOnly, token != "", route.method, route.path, rec.Code, expected)
				}
				if !reflect.DeepEqual(before, config.CopyCfg()) {
					t.Fatal("rejected request changed configuration")
				}
			}
		}
	}
	cfg.ReadOnlyMode = false
	if rec := request("PUT", "/api/configs/Debug", `{"value":true}`, login.Token); rec.Code != 200 || !cfg.Debug {
		t.Fatalf("authorized write: %d", rec.Code)
	}
	if rec := request("GET", "/api/books", "", login.Token); rec.Code != 200 {
		t.Fatalf("reading blocked: %d", rec.Code)
	}
	cfg.Password = "changed"
	if rec := request("GET", "/api/books", "", login.Token); rec.Code != 401 {
		t.Fatalf("old token accepted: %d", rec.Code)
	}
}
