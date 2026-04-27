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
