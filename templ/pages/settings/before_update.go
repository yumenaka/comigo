package settings

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/config_api"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

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
		logger.Infof("Server action: %v,", action)
	} else {
		logger.Infof("Server action: %s", string(actionString))
	}
	// 重启网页服务器等，此处不能导入routers，因为会循环引用
	if action.ReStartWebServer {
		*RestartWebServerBroadcast <- "restart_web_server"
		// 等待服务器端口可用，确保重启完成后再继续
		tools.WaitUntilServerReady("localhost", newConfig.Port, 15*time.Second)
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
		logger.Info("No changes in cfg, skipped rescan dir\n")
	}
}

func openBrowserIfNeeded(oldConfig *config.Config, newConfig *config.Config) {
	if !oldConfig.OpenBrowser && newConfig.OpenBrowser {
		protocol := "http://"
		if newConfig.EnableTLS {
			protocol = "https://"
		}
		go tools.OpenBrowser(protocol + "localhost:" + strconv.Itoa(newConfig.Port))
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
			logger.Info("Tailscale config changed, will restart Tailscale server")
		}
	}
	// 什么情况下需要启动或停止Tailscale服务器
	if oldConfig.EnableTailscale != newConfig.EnableTailscale && newConfig.EnableTailscale == true {
		action.StartTailscale = true
		action.StopTailscale = false
		action.ReStartTailscale = false
		logger.Info("Tailscale enabled, will start Tailscale server")
	}
	if oldConfig.EnableTailscale != newConfig.EnableTailscale && newConfig.EnableTailscale == false {
		action.StartTailscale = false
		action.StopTailscale = true
		action.ReStartTailscale = false
		logger.Info("Tailscale disabled, will stop Tailscale server")
	}

	return action
}
