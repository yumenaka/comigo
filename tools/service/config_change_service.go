package service

import (
	"encoding/json"
	"slices"
	"strconv"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
	"github.com/yumenaka/comigo/tools/sse_hub"
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
	action.ReScanStores = !slices.Equal(oldConfig.StoreUrls, newConfig.StoreUrls) ||
		oldConfig.MaxScanDepth != newConfig.MaxScanDepth ||
		!slices.Equal(oldConfig.SupportMediaType, newConfig.SupportMediaType) ||
		!slices.Equal(oldConfig.SupportFileType, newConfig.SupportFileType) ||
		oldConfig.MinImageNum != newConfig.MinImageNum ||
		!slices.Equal(oldConfig.ExcludePath, newConfig.ExcludePath)

	action.ReStartWebServer = oldConfig.Port != newConfig.Port ||
		oldConfig.Username != newConfig.Username ||
		oldConfig.Password != newConfig.Password ||
		oldConfig.Host != newConfig.Host ||
		config.NormalizeBasePath(oldConfig.BasePath) != config.NormalizeBasePath(newConfig.BasePath) ||
		oldConfig.Timeout != newConfig.Timeout ||
		oldConfig.DisableLAN != newConfig.DisableLAN

	action.UpdateAutoRescan = oldConfig.AutoRescanIntervalMinutes != newConfig.AutoRescanIntervalMinutes

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
		return
	}
	if config.GetCfg().EnableDatabase {
		if err := scan.SaveBooksToDatabase(config.GetCfg()); err != nil {
			logger.Infof(locale.GetString("log_failed_to_save_results_to_database"), err)
		}
	}
	sse_hub.BroadcastUISuggestReload(sse_hub.UISuggestReasonLibraryRescan)
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
		protocol := "http://"
		if newConfig.EnableTLS {
			protocol = "https://"
		}
		basePath := config.NormalizeBasePath(newConfig.BasePath)
		if basePath == "" {
			basePath = "/"
		} else {
			basePath += "/"
		}
		go tools.OpenBrowserByURL(protocol + "127.0.0.1:" + strconv.Itoa(newConfig.Port) + basePath)
	}
}
