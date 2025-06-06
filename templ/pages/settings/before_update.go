package settings

import (
	"reflect"
	"strconv"

	"github.com/spf13/viper"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
	"github.com/yumenaka/comigo/util/scan"
)

var RestartWebServerBroadcast *chan string

// -------------------------
// 各种辅助函数
// -------------------------

// beforeConfigUpdate 根据配置的变化，判断是否需要打开浏览器重新扫描等
func beforeConfigUpdate(oldConfig *config.Config, newConfig *config.Config) {
	openBrowserIfNeeded(oldConfig, newConfig)
	reScanStores, reStartWebServer := checkServerActions(oldConfig, newConfig)
	logger.Infof("reScanDir: %v, reStartWebServer: %v ", reScanStores, reStartWebServer)
	if reStartWebServer {
		// 此处需要不能导入routers，因为会循环引用  routers.RestartWebServer() X
		// 使用广播的方式来通知 :
		*RestartWebServerBroadcast <- "restart_web_server"
	}
	if reScanStores {
		startReScan()
	} else {
		if newConfig.Debug {
			logger.Info("No changes in cfg, skipped rescan dir\n")
		}
	}
}

func openBrowserIfNeeded(oldConfig *config.Config, newConfig *config.Config) {
	if !oldConfig.OpenBrowser && newConfig.OpenBrowser {
		protocol := "http://"
		if newConfig.EnableTLS {
			protocol = "https://"
		}
		go util.OpenBrowser(protocol + "localhost:" + strconv.Itoa(newConfig.Port))
	}
}

// checkServerActions 检查旧的和新的配置是否需要更新，并返回需要重启网页服务器、重新扫描整个书库、重新扫描所有文件的布尔值
func checkServerActions(oldConfig *config.Config, newConfig *config.Config) (reScanStores bool, reStartWebServer bool) {
	// 下面这些值修改的时候，需要扫描整个书库、或重新扫描所有文件
	if !reflect.DeepEqual(oldConfig.LocalStores, newConfig.LocalStores) {
		reScanStores = true
	}
	if oldConfig.MaxScanDepth != newConfig.MaxScanDepth {
		reScanStores = true
	}
	if !reflect.DeepEqual(oldConfig.SupportMediaType, newConfig.SupportMediaType) {
		reScanStores = true
	}
	if !reflect.DeepEqual(oldConfig.SupportFileType, newConfig.SupportFileType) {
		reScanStores = true
	}
	if oldConfig.MinImageNum != newConfig.MinImageNum {
		reScanStores = true
	}
	if !reflect.DeepEqual(oldConfig.ExcludePath, newConfig.ExcludePath) {
		reScanStores = true
	}
	if oldConfig.EnableDatabase != newConfig.EnableDatabase {
		reScanStores = true
	}
	// 下面这些值修改的时候，需要重启网页服务器
	if oldConfig.Port != newConfig.Port {
		reStartWebServer = true
	}
	if oldConfig.Username != newConfig.Username {
		reStartWebServer = true
	}
	if oldConfig.Password != newConfig.Password {
		reStartWebServer = true
	}
	if oldConfig.Host != newConfig.Host {
		reStartWebServer = true
	}
	if oldConfig.Timeout != newConfig.Timeout {
		reStartWebServer = true
	}
	if oldConfig.DisableLAN != newConfig.DisableLAN {
		reStartWebServer = true
	}
	// if oldConfig.StaticFileMode != newConfig.StaticFileMode {
	// 	reStartWebServer = true
	// }
	return
}

// startReScan 扫描并相应地更新数据库
func startReScan() {
	config.InitCfgStores()
	if err := scan.InitAllStore(scan.NewOption(config.GetCfg())); err != nil {
		logger.Infof("Failed to scan store path: %v", err)
	}
	if config.GetEnableDatabase() {
		saveResultsToDatabase(viper.ConfigFileUsed(), config.GetClearDatabaseWhenExit())
	}
}

func saveResultsToDatabase(configPath string, clearDatabaseWhenExit bool) {
	if err := scan.SaveResultsToDatabase(configPath, clearDatabaseWhenExit); err != nil {
		logger.Infof("Failed to save results to database: %v", err)
	}
}
