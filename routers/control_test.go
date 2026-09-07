package routers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// 重启期间不能先写入配置再跳过重启，否则监听地址与保存的配置会不一致。
func TestConfigUpdateWhileRestarting(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { *config.GetCfg() = old; restarting.Store(false) })
	config.GetCfg().TemporaryReaderMode = true
	config.GetCfg().ReadOnlyMode = false
	restarting.Store(true)
	e := echo.New()
	e.PATCH("/configs", updateConfigHandler)
	req := httptest.NewRequest(http.MethodPatch, "/configs", strings.NewReader(`{"Port":18769}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict || config.GetCfg().Port != old.Port {
		t.Fatalf("重启期间更新: status=%d, port=%d，期望409且配置不变", rec.Code, config.GetCfg().Port)
	}
}
