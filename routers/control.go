package routers

import (
	"net/http"
	"sync/atomic"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/config_api"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/service"
)

var restarting atomic.Bool

// restartHandler 先返回接受状态，再重启 HTTP 服务，避免等待自身请求结束造成死锁。
// 嵌入式宿主进程继续运行；并发重启请求返回冲突。
func restartHandler(c echo.Context) error {
	if config.GetCfg().ReadOnlyMode {
		return echo.ErrForbidden
	}
	if !restarting.CompareAndSwap(false, true) {
		return echo.ErrConflict
	}
	if err := c.NoContent(http.StatusAccepted); err != nil {
		restarting.Store(false)
		return err
	}
	go restartService()
	return nil
}

// updateConfigHandler 在配置更新成功后重新绑定有变化的服务监听，覆盖 CLI 和嵌入式宿主。
func updateConfigHandler(c echo.Context) error {
	// 配置写入与重启共用占用标志，避免接受更新后遗漏所需的重启。
	if !restarting.CompareAndSwap(false, true) {
		return echo.ErrConflict
	}
	restart := false
	defer func() {
		if !restart {
			restarting.Store(false)
		}
	}()
	old := config.CopyCfg()
	if err := config_api.UpdateConfig(c); err != nil {
		return err
	}
	if c.Response().Status >= http.StatusBadRequest {
		return nil
	}
	action := service.BuildConfigChangeAction(old, config.GetCfg())
	if action.ReStartWebServer {
		restart = true
		go restartService()
	} else if action.StartTailscale || action.ReStartTailscale {
		go StartTailscale()
	} else if action.StopTailscale {
		go StopTailscale()
	}
	return nil
}

// restartService 统一执行控制接口触发的重启并释放并发标志。
func restartService() {
	defer restarting.Store(false)
	if err := RestartWebServer(); err != nil {
		logger.Errorf(locale.GetString("err_restart_web_server_failed"), err)
		return
	}
	StartTailscale()
}
