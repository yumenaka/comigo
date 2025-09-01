package settings

import (
	"reflect"
	"strconv"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/config_api"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
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
		// 此处需要不能导入routers，因为会循环引用  routers.RestartWebServer()
		// 使用广播的方式来通知 :
		*RestartWebServerBroadcast <- "restart_web_server"
	}
	if reScanStores {
		config_api.StartReScan()
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
	if !reflect.DeepEqual(oldConfig.StoreUrls, newConfig.StoreUrls) {
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
	return
}
