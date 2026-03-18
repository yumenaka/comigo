package service

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
)

type ConfigChangeAction struct {
	ReScanStores     bool `json:"reScanStores"`
	ReStartWebServer bool `json:"reStartWebServer"`
	StartTailscale   bool `json:"startTailscale"`
	StopTailscale    bool `json:"stopTailscale"`
	ReStartTailscale bool `json:"reStartTailscale"`
	UpdateAutoRescan bool `json:"updateAutoRescan"`
}

// BuildConfigChangeAction 比较新旧配置，计算后续需要执行的动作。
func BuildConfigChangeAction(oldConfig config.Config, newConfig *config.Config) (action ConfigChangeAction) {
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
	if oldConfig.AutoRescanIntervalMinutes != newConfig.AutoRescanIntervalMinutes {
		action.UpdateAutoRescan = true
	}

	if oldConfig.EnableTailscale != newConfig.EnableTailscale && newConfig.EnableTailscale {
		action.StartTailscale = true
	}
	if oldConfig.EnableTailscale != newConfig.EnableTailscale && !newConfig.EnableTailscale {
		action.StopTailscale = true
	}

	if newConfig.EnableTailscale {
		if oldConfig.TailscaleAuthKey != newConfig.TailscaleAuthKey ||
			oldConfig.FunnelTunnel != newConfig.FunnelTunnel ||
			oldConfig.TailscaleHostname != newConfig.TailscaleHostname ||
			oldConfig.TailscalePort != newConfig.TailscalePort {
			action.ReStartTailscale = true
			action.StartTailscale = false
			action.StopTailscale = false
			logger.Info(locale.GetString("log_tailscale_config_changed_restart"))
		}
	}
	return action
}

// ApplyConfigChange 执行配置变更副作用。
// restartSignal 可为空；为空时将跳过网页服务重启与 tailscale 信号广播。
func ApplyConfigChange(oldConfig config.Config, newConfig *config.Config, restartSignal chan<- string) {
	openBrowserIfNeeded(oldConfig, newConfig)
	action := BuildConfigChangeAction(oldConfig, newConfig)
	logAction(action)

	if restartSignal != nil && action.ReStartWebServer {
		restartSignal <- "restart_web_server"
		tools.WaitUntilServerReady("localhost", uint16(newConfig.Port), 15*time.Second)
	}
	if restartSignal != nil && action.StartTailscale {
		restartSignal <- "start_tailscale"
	}
	if restartSignal != nil && action.StopTailscale {
		restartSignal <- "stop_tailscale"
	}
	if restartSignal != nil && action.ReStartTailscale {
		restartSignal <- "restart_tailscale"
	}

	if action.ReScanStores {
		StartReScan()
	}
	if action.UpdateAutoRescan {
		config.StartOrStopAutoRescan()
	}
	if newConfig.Debug && !action.ReScanStores && !action.ReStartWebServer && !action.UpdateAutoRescan {
		logger.Info(locale.GetString("log_no_changes_skipped_rescan"))
	}
}

func StartReScan() {
	if err := scan.InitAllStore(config.GetCfg()); err != nil {
		logger.Infof(locale.GetString("log_failed_to_scan_store_path"), err)
	}
	if config.GetCfg().EnableDatabase {
		if err := scan.SaveBooksToDatabase(config.GetCfg()); err != nil {
			logger.Infof(locale.GetString("log_failed_to_save_results_to_database"), err)
		}
	}
}

func logAction(action ConfigChangeAction) {
	actionString, err := json.Marshal(action)
	if err != nil {
		logger.Infof(locale.GetString("log_server_action")+",", action)
		return
	}
	logger.Infof(locale.GetString("log_server_action_string"), string(actionString))
}

func openBrowserIfNeeded(oldConfig config.Config, newConfig *config.Config) {
	if !oldConfig.OpenBrowser && newConfig.OpenBrowser {
		go tools.OpenBrowser(newConfig.EnableTLS, "127.0.0.1", newConfig.Port)
	}
}
