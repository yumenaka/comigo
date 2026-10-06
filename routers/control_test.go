package routers

import (
	"fmt"
	"net"
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

// 锁定配置时 REST 不能改变对外服务状态。
func TestExternalAccessReadOnly(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { *config.GetCfg() = old; restarting.Store(false) })
	config.GetCfg().ReadOnlyMode = true
	config.GetCfg().DisableLAN = true
	e := echo.New()
	e.PATCH("/configs", updateConfigHandler)
	req := httptest.NewRequest(http.MethodPatch, "/configs", strings.NewReader(`{"DisableLAN":false}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden || !config.GetCfg().DisableLAN {
		t.Fatalf("锁定配置被修改: %d", rec.Code)
	}
}

// REST 无效 TLS 更新必须返回错误，并保留当前协议与其他配置。
func TestConfigUpdateRejectsInvalidTLS(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { *config.GetCfg() = old; restarting.Store(false) })
	config.GetCfg().TemporaryReaderMode = true
	config.GetCfg().ReadOnlyMode = false
	config.GetCfg().EnableTLS = false
	config.GetCfg().AutoTLSCertificate = false
	config.GetCfg().Debug = false
	config.GetCfg().CertFile, config.GetCfg().KeyFile = "", ""
	e := echo.New()
	e.PATCH("/configs", updateConfigHandler)
	req := httptest.NewRequest(http.MethodPatch, "/configs", strings.NewReader(`{"EnableTLS":true,"Debug":true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest || config.GetCfg().EnableTLS || config.GetCfg().Debug {
		t.Fatalf("status=%d cfg=%+v", rec.Code, config.GetCfg())
	}
}

// 配置导致端口绑定失败时必须恢复原监听，不能让网页控制永久断开。
func TestConfigRestartFailureRestoresPreviousListener(t *testing.T) {
	oldCfg, oldEngine, oldServer := config.CopyCfg(), engine, config.Server
	t.Cleanup(func() {
		_ = StopWebServer()
		*config.GetCfg() = oldCfg
		engine = oldEngine
		config.Server = oldServer
		restarting.Store(false)
	})
	cfg := config.GetCfg()
	cfg.TemporaryReaderMode, cfg.DisableLAN = true, true
	cfg.EnableTLS, cfg.AutoTLSCertificate, cfg.EnableTailscale, cfg.LogToFile = false, false, false, false
	cfg.Port = 0
	if err := StartEcho(echo.New()); err != nil {
		t.Fatal(err)
	}
	previous := config.CopyCfg()
	occupied, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer occupied.Close()
	cfg.Port = occupied.Addr().(*net.TCPAddr).Port
	restartService(&previous)
	if config.GetCfg().Port != previous.Port || config.Server == nil {
		t.Fatal("previous listener not restored")
	}
	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/healthz", previous.Port))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("health=%d", response.StatusCode)
	}
}
