package settings

import (
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/config_api"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

func getTranslations(value string) string {
	return "i18next.t(\"" + value + "\")"
}

// getTranslationsGO 返回 Go 语言的翻译字符串  not working
func getTranslationsGO(value string, lang string) string {
	return locale.GetStringByLocal(value, lang)
}

// PageHandler 设定页面
func PageHandler(c echo.Context) error {
	tsStatus, err := tailscale_plugin.GetTailscaleStatus(c.Request().Context())
	if err != nil {
		// 容错：当 Tailscale 未启用或尚未就绪时，不中断页面渲染，返回离线状态
		tsStatus = &tailscale_plugin.TailscaleStatus{}
	}
	indexHtml := common.Html(
		c,
		SettingsPage(c, tsStatus),
		[]string{},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

// -----------------------------------------
// BeforeConfigUpdate 在配置更新前执行的操作
// -------------------------

// RestartWebServerBroadcast 用于广播重启网页服务器的信号
var RestartWebServerBroadcast *chan string

// -------------------------
// 各种辅助函数
// -------------------------

// beforeConfigUpdate 根据配置的变化，判断是否需要打开浏览器重新扫描等
func beforeConfigUpdate(oldConfig *config.Config, newConfig *config.Config) {
	openBrowserIfNeeded(oldConfig, newConfig)
	action := checkServerActions(oldConfig, newConfig)
	// 记录需要执行的操作
	actionString, err := json.Marshal(action)
	if err != nil {
		logger.Infof(locale.GetString("log_server_action")+",", action)
	} else {
		logger.Infof(locale.GetString("log_server_action_string"), string(actionString))
	}
	// 重启网页服务器等，此处不能导入routers，因为会循环引用
	if action.ReStartWebServer {
		*RestartWebServerBroadcast <- "restart_web_server"
		// 等待服务器端口可用，确保重启完成后再继续
		tools.WaitUntilServerReady("localhost", uint16(newConfig.Port), 15*time.Second)
	}
	if action.StartTailscale {
		*RestartWebServerBroadcast <- "start_tailscale"
	}
	if action.StopTailscale {
		*RestartWebServerBroadcast <- "stop_tailscale"
	}
	if action.ReStartTailscale {
		*RestartWebServerBroadcast <- "restart_tailscale"
	}
	// 重新扫描书库
	if action.ReScanStores {
		config_api.StartReScan()
	}
	// 提示没有变化
	if newConfig.Debug && !action.ReScanStores && !action.ReStartWebServer {
		logger.Info(locale.GetString("log_no_changes_skipped_rescan"))
	}
}

func openBrowserIfNeeded(oldConfig *config.Config, newConfig *config.Config) {
	if !oldConfig.OpenBrowser && newConfig.OpenBrowser {
		go tools.OpenBrowser(newConfig.EnableTLS, "127.0.0.1", newConfig.Port)
	}
}

type ConfigChangeAction struct {
	ReScanStores     bool
	ReStartWebServer bool
	StartTailscale   bool
	StopTailscale    bool
	ReStartTailscale bool
}

// checkServerActions 检查旧的和新的配置是否需要更新，并返回需要重启网页服务器、重新扫描整个书库、重新扫描所有文件的布尔值
func checkServerActions(oldConfig *config.Config, newConfig *config.Config) (action ConfigChangeAction) {
	// 下面这些值修改的时候，需要扫描整个书库、或重新扫描所有文件
	if !reflect.DeepEqual(oldConfig.StoreUrls, newConfig.StoreUrls) {
		action.ReScanStores = true
	}
	if oldConfig.MaxScanDepth != newConfig.MaxScanDepth {
		action.ReScanStores = true
	}
	if !reflect.DeepEqual(oldConfig.SupportMediaType, newConfig.SupportMediaType) {
		action.ReScanStores = true
	}
	if !reflect.DeepEqual(oldConfig.SupportFileType, newConfig.SupportFileType) {
		action.ReScanStores = true
	}
	if oldConfig.MinImageNum != newConfig.MinImageNum {
		action.ReScanStores = true
	}
	if !reflect.DeepEqual(oldConfig.ExcludePath, newConfig.ExcludePath) {
		action.ReScanStores = true
	}
	if oldConfig.EnableDatabase != newConfig.EnableDatabase {
		action.ReScanStores = true
	}
	// 下面这些值修改的时候，需要重启网页服务器
	if oldConfig.Port != newConfig.Port {
		action.ReStartWebServer = true
	}
	if oldConfig.Username != newConfig.Username {
		action.ReStartWebServer = true
	}
	if oldConfig.Password != newConfig.Password {
		action.ReStartWebServer = true
	}
	if oldConfig.Host != newConfig.Host {
		action.ReStartWebServer = true
	}
	if oldConfig.Timeout != newConfig.Timeout {
		action.ReStartWebServer = true
	}
	if oldConfig.DisableLAN != newConfig.DisableLAN {
		action.ReStartWebServer = true
	}
	// 如果开启了Tailscale，且Tailscale设置有变化，则需要重启Tailscale服务器
	if newConfig.EnableTailscale == true {
		if oldConfig.TailscaleAuthKey != newConfig.TailscaleAuthKey {
			action.ReStartTailscale = true
		}
		if oldConfig.FunnelTunnel != newConfig.FunnelTunnel {
			action.ReStartTailscale = true
		}
		if oldConfig.TailscaleHostname != newConfig.TailscaleHostname {
			action.ReStartTailscale = true
		}
		if oldConfig.TailscalePort != newConfig.TailscalePort {
			action.ReStartTailscale = true
		}
		if action.ReStartTailscale == true {
			action.StartTailscale = false
			action.StopTailscale = false
			logger.Info(locale.GetString("log_tailscale_config_changed_restart"))
		}
	}
	// 什么情况下需要启动或停止Tailscale服务器
	if oldConfig.EnableTailscale != newConfig.EnableTailscale && newConfig.EnableTailscale == true {
		action.StartTailscale = true
		action.StopTailscale = false
		action.ReStartTailscale = false
		logger.Info(locale.GetString("log_tailscale_enabled_start"))
	}
	if oldConfig.EnableTailscale != newConfig.EnableTailscale && newConfig.EnableTailscale == false {
		action.StartTailscale = false
		action.StopTailscale = true
		action.ReStartTailscale = false
		logger.Info(locale.GetString("log_tailscale_disabled_stop"))
	}

	return action
}
