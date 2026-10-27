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
	// 先完成一次请求，确保 Serve 已接管监听，再模拟运行中的配置更新。
	ready, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", previous.Port))
	if err != nil {
		t.Fatal(err)
	}
	ready.Body.Close()
	occupied, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer occupied.Close()
	// StartEcho 已发布新的配置快照，必须通过 setter 更新当前配置。
	config.SetPort(occupied.Addr().(*net.TCPAddr).Port)
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

// 原生绑定与 HTTP 控制接口共用配置锁和重启锁，拒绝时不写入共享状态。
func TestNativeSharingGuards(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { config.ReplaceConfig(old); restarting.Store(false) })
	next := config.CopyCfg()
	next.TemporaryReaderMode, next.ReadOnlyMode, next.DisableLAN = true, true, true
	config.ReplaceConfig(next)
	if err := SetLANSharing(true); err == nil || !config.GetCfg().DisableLAN || restarting.Load() {
		t.Fatal("只读配置被修改或重启锁未释放")
	}
	next.ReadOnlyMode = false
	config.ReplaceConfig(next)
	restarting.Store(true)
	if err := SetLANSharing(true); err == nil || !config.GetCfg().DisableLAN {
		t.Fatal("已有重启时仍修改局域网共享")
	}
}

// 原生面板关闭局域网共享后必须完成监听重启，并继续在原端口提供服务。
func TestNativeSharingRebindsLoopback(t *testing.T) {
	oldCfg, oldEngine, oldServer := config.CopyCfg(), engine, config.Server
	t.Cleanup(func() {
		_ = StopWebServer()
		config.ReplaceConfig(oldCfg)
		engine, config.Server = oldEngine, oldServer
		restarting.Store(false)
	})
	next := config.CopyCfg()
	next.TemporaryReaderMode, next.DisableLAN = true, true
	next.EnableTLS, next.AutoTLSCertificate, next.EnableTailscale, next.LogToFile, next.ReadOnlyMode = false, false, false, false, false
	next.Port = 0
	config.ReplaceConfig(next)
	if err := StartEcho(echo.New()); err != nil {
		t.Fatal(err)
	}
	port := config.GetCfg().Port
	ready, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if err != nil {
		t.Fatal(err)
	}
	ready.Body.Close()
	// 仅改变配置快照，不实际向局域网开放测试服务器。
	next = config.CopyCfg()
	next.DisableLAN = false
	config.ReplaceConfig(next)
	if err := SetLANSharing(false); err != nil {
		t.Fatal(err)
	}
	if !config.GetCfg().DisableLAN || config.GetCfg().Port != port || restarting.Load() {
		t.Fatal("共享状态或监听端口未正确更新")
	}
	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/healthz", port))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("health=%d", response.StatusCode)
	}
}
