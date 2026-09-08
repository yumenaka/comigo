package config_api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// 状态接口仅返回文件来源，不读取或返回配置内容。
func TestGetConfigFileStatus(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { *config.GetCfg() = old })
	file := filepath.Join(t.TempDir(), "reader.toml")
	if err := os.WriteFile(file, []byte(`Password = "private-value"`), 0600); err != nil {
		t.Fatal(err)
	}
	config.GetCfg().ConfigFile = file
	config.GetCfg().TemporaryReaderMode = false
	e := echo.New()
	rec := httptest.NewRecorder()
	if err := GetConfigStatus(e.NewContext(httptest.NewRequest(http.MethodGet, "/api/configs/status", nil), rec)); err != nil {
		t.Fatal(err)
	}
	var status config.Status
	if err := json.Unmarshal(rec.Body.Bytes(), &status); err != nil {
		t.Fatal(err)
	}
	if status.Current.Path != file || status.Current.Location != "Custom" || !status.Current.Exists {
		t.Fatalf("文件状态: %+v", status)
	}
	if rec.Header().Get("Cache-Control") != "no-store" || strings.Contains(rec.Body.String(), "private-value") {
		t.Fatal("缓存或内容泄漏")
	}
}
