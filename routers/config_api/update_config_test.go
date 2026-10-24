package config_api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// 验证只读模式下更新配置会返回统一错误结构。
func TestUpdateConfig_ReadOnly_ReturnsContract(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(`{"Debug":true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	oldReadOnly := config.GetCfg().ReadOnlyMode
	config.GetCfg().ReadOnlyMode = true
	t.Cleanup(func() {
		config.GetCfg().ReadOnlyMode = oldReadOnly
	})

	if err := UpdateConfig(c); err != nil {
		t.Fatalf("UpdateConfig returned error: %v", err)
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if body["code"] != "config_locked" {
		t.Fatalf("expected code config_locked, got %v", body["code"])
	}
}

// REST 写盘失败必须报错，不能修改正在使用的配置。
func TestUpdateConfigSaveFailure(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { config.ReplaceConfig(old) })
	next := config.CopyCfg()
	next.ReadOnlyMode = false
	next.TemporaryReaderMode = false
	next.ConfigFile = t.TempDir()
	next.Debug = false
	config.ReplaceConfig(next)
	e := echo.New()
	e.PATCH("/api/configs", UpdateConfig)
	req := httptest.NewRequest(http.MethodPatch, "/api/configs", strings.NewReader(`{"Debug":true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != 500 || config.GetCfg().Debug {
		t.Fatalf("status=%d debug=%v", rec.Code, config.GetCfg().Debug)
	}
}
