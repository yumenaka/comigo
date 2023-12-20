package handlers

import (
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/arch/scan"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/types"
	"github.com/yumenaka/comi/util"
)

// UpdateConfig 修改服务器配置(post json)
func UpdateConfig(c *gin.Context) {
	// 读取请求体中的JSON数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	// 将JSON数据转换为字符串并打印
	jsonString := string(body)
	logger.Infof("Received JSON data: %s \n", jsonString)

	// 解析JSON数据并更新服务器配置
	newConfig, err := types.UpdateConfig(&config.Config, jsonString)
	if err != nil {
		logger.Info(err.Error())
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Failed to parse JSON data"})
		return
	}
	// 根据配置的变化，做相应的操作。比如打开浏览器重新扫描等
	BeforeConfigUpdate(&config.Config, newConfig)
	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Server settings updated successfully"})
}

// BeforeConfigUpdate 根据配置的变化，判断是否需要打开浏览器重新扫描等
func BeforeConfigUpdate(oldConfig *types.ComigoConfig, newConfig *types.ComigoConfig) {
	err := config.UpdateLocalConfig()
	if err != nil {
		logger.Infof("Failed to update local config: %v", err)
	}
	openBrowserIfNeeded(oldConfig, newConfig)
	needScan, reScanFile := updateConfigIfNeeded(oldConfig, newConfig)
	if needScan {
		performScanAndUpdateDBIfNeeded(oldConfig, newConfig, reScanFile)
	} else {
		if oldConfig.Debug {
			logger.Infof("No changes in Config, skipped scan store path")
		}
	}
}

func openBrowserIfNeeded(oldConfig *types.ComigoConfig, newConfig *types.ComigoConfig) {
	if newConfig.OpenBrowser && !oldConfig.OpenBrowser {
		protocol := "http://"
		if newConfig.EnableTLS {
			protocol = "https://"
		}
		util.OpenBrowser(protocol + "127.0.0.1:" + strconv.Itoa(newConfig.Port))
	}
}

// updateConfigIfNeeded 检查旧的和新的配置是否需要更新，并返回需要重新扫描和重新扫描文件的布尔值
func updateConfigIfNeeded(oldConfig *types.ComigoConfig, newConfig *types.ComigoConfig) (needScan bool, reScanFile bool) {
	if differentConfig(oldConfig.StoresPath, newConfig.StoresPath) {
		needScan = true
		oldConfig.StoresPath = newConfig.StoresPath
	}
	if oldConfig.MaxScanDepth != newConfig.MaxScanDepth {
		needScan = true
		oldConfig.MaxScanDepth = newConfig.MaxScanDepth
	}
	if differentConfig(oldConfig.SupportMediaType, newConfig.SupportMediaType) {
		needScan = true
		reScanFile = true
		oldConfig.SupportMediaType = newConfig.SupportMediaType
	}
	if differentConfig(oldConfig.SupportFileType, newConfig.SupportFileType) {
		needScan = true
		oldConfig.SupportFileType = newConfig.SupportFileType
	}
	if oldConfig.MinImageNum != newConfig.MinImageNum {
		needScan = true
		reScanFile = true
		oldConfig.MinImageNum = newConfig.MinImageNum
	}
	if differentConfig(oldConfig.ExcludePath, newConfig.ExcludePath) {
		needScan = true
		oldConfig.ExcludePath = newConfig.ExcludePath
	}
	if oldConfig.EnableDatabase != newConfig.EnableDatabase {
		needScan = true
		oldConfig.EnableDatabase = newConfig.EnableDatabase
	}
	return
}

func differentConfig(old, new interface{}) bool {
	return !reflect.DeepEqual(old, new)
}

// performScanAndUpdateDBIfNeeded 扫描并相应地更新数据库
func performScanAndUpdateDBIfNeeded(oldConfig *types.ComigoConfig, newConfig *types.ComigoConfig, reScanFile bool) {
	option := scan.NewScanOption(
		reScanFile,
		newConfig.StoresPath,
		newConfig.MaxScanDepth,
		newConfig.MinImageNum,
		newConfig.TimeoutLimitForScan,
		newConfig.ExcludePath,
		newConfig.SupportMediaType,
		newConfig.SupportFileType,
		newConfig.ZipFileTextEncoding,
		newConfig.EnableDatabase,
		newConfig.ClearDatabaseWhenExit,
		newConfig.Debug,
	)
	if err := scan.ScanStorePath(option); err != nil {
		logger.Infof("Failed to scan store path: %v", err)
	}
	if oldConfig.EnableDatabase {
		saveResultsToDatabase(config.Config.ConfigPath, config.Config.ClearDatabaseWhenExit)
	}
}

func saveResultsToDatabase(configPath string, clearDatabaseWhenExit bool) {
	if err := scan.SaveResultsToDatabase(configPath, clearDatabaseWhenExit); err != nil {
		logger.Infof("Failed to save results to database: %v", err)
	}
}
