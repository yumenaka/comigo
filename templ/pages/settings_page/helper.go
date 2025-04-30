package settings_page

import (
	"reflect"
	"strconv"

	"github.com/spf13/viper"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
	"github.com/yumenaka/comigo/util/scan"
)

// -------------------------
// 各种辅助函数
// -------------------------

// beforeConfigUpdate 根据配置的变化，判断是否需要打开浏览器重新扫描等
func beforeConfigUpdate(oldConfig *config.Config, newConfig *config.Config) {
	openBrowserIfNeeded(oldConfig, newConfig)
	reScanDir, reScanFile := checkReScanStatus(oldConfig, newConfig)
	logger.Infof("reScanDir: %v, reScanFile: %v\n", reScanDir, reScanFile)
	if reScanDir {
		startReScan(reScanFile)
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

// checkReScanStatus 检查旧的和新的配置是否需要更新，并返回需要重新扫描和重新扫描文件的布尔值
func checkReScanStatus(oldConfig *config.Config, newConfig *config.Config) (reScanDir bool, reScanFile bool) {
	if !reflect.DeepEqual(oldConfig.LocalStores, newConfig.LocalStores) {
		reScanDir = true
	}
	if oldConfig.MaxScanDepth != newConfig.MaxScanDepth {
		reScanDir = true
	}
	if !reflect.DeepEqual(oldConfig.SupportMediaType, newConfig.SupportMediaType) {
		reScanDir = true
		reScanFile = true
	}
	if !reflect.DeepEqual(oldConfig.SupportFileType, newConfig.SupportFileType) {
		reScanDir = true
	}
	if oldConfig.MinImageNum != newConfig.MinImageNum {
		reScanDir = true
		reScanFile = true
	}
	if !reflect.DeepEqual(oldConfig.ExcludePath, newConfig.ExcludePath) {
		reScanDir = true
	}
	if oldConfig.EnableDatabase != newConfig.EnableDatabase {
		reScanDir = true
	}
	return
}

// startReScan 扫描并相应地更新数据库
func startReScan(reScanFile bool) {
	config.InitCfgStores()
	option := scan.NewOption(
		reScanFile,
		config.GetCfg(),
	)
	if err := scan.InitAllStore(option); err != nil {
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
