package config_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/entity"
	"github.com/yumenaka/comi/util"
	"github.com/yumenaka/comi/util/file/scan"
	"github.com/yumenaka/comi/util/logger"
	"io"
	"net/http"
	"reflect"
	"strconv"
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
	oldConfig, err := entity.UpdateConfig(&config.Config, jsonString)
	if err != nil {
		logger.Infof("%s", err.Error())
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Failed to parse JSON data"})
		return
	}
	err = config.UpdateLocalConfig()
	if err != nil {
		logger.Infof("Failed to update local config: %v", err)
	}
	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	BeforeConfigUpdate(oldConfig, &config.Config)
	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Server settings updated successfully"})
}

// BeforeConfigUpdate 根据配置的变化，判断是否需要打开浏览器重新扫描等
func BeforeConfigUpdate(oldConfig *entity.ComigoConfig, newConfig *entity.ComigoConfig) {
	openBrowserIfNeeded(oldConfig, newConfig)
	needScan, reScanFile := updateConfigIfNeeded(oldConfig, newConfig)
	if needScan {
		performScanAndUpdateDBIfNeeded(newConfig, reScanFile)
	} else {
		if newConfig.Debug {
			logger.Info("No changes in Config, skipped scan store path\n")
		}
	}
}

func openBrowserIfNeeded(oldConfig *entity.ComigoConfig, newConfig *entity.ComigoConfig) {
	if (oldConfig.OpenBrowser == false) && (newConfig.OpenBrowser == true) {
		protocol := "http://"
		if newConfig.EnableTLS {
			protocol = "https://"
		}
		util.OpenBrowser(protocol + "127.0.0.1:" + strconv.Itoa(newConfig.Port))
	}
}

// updateConfigIfNeeded 检查旧的和新的配置是否需要更新，并返回需要重新扫描和重新扫描文件的布尔值
func updateConfigIfNeeded(oldConfig *entity.ComigoConfig, newConfig *entity.ComigoConfig) (needScan bool, reScanFile bool) {
	if !reflect.DeepEqual(oldConfig.LocalStores, newConfig.LocalStores) {
		needScan = true
		oldConfig.LocalStores = newConfig.LocalStores
	}
	if oldConfig.MaxScanDepth != newConfig.MaxScanDepth {
		needScan = true
		oldConfig.MaxScanDepth = newConfig.MaxScanDepth
	}
	if !reflect.DeepEqual(oldConfig.SupportMediaType, newConfig.SupportMediaType) {
		needScan = true
		reScanFile = true
		oldConfig.SupportMediaType = newConfig.SupportMediaType
	}
	if !reflect.DeepEqual(oldConfig.SupportFileType, newConfig.SupportFileType) {
		needScan = true
		oldConfig.SupportFileType = newConfig.SupportFileType
	}
	if oldConfig.MinImageNum != newConfig.MinImageNum {
		needScan = true
		reScanFile = true
		oldConfig.MinImageNum = newConfig.MinImageNum
	}
	if !reflect.DeepEqual(oldConfig.ExcludePath, newConfig.ExcludePath) {
		needScan = true
		oldConfig.ExcludePath = newConfig.ExcludePath
	}
	if oldConfig.EnableDatabase != newConfig.EnableDatabase {
		needScan = true
		oldConfig.EnableDatabase = newConfig.EnableDatabase
	}
	return
}

// performScanAndUpdateDBIfNeeded 扫描并相应地更新数据库
func performScanAndUpdateDBIfNeeded(newConfig *entity.ComigoConfig, reScanFile bool) {
	option := scan.NewScanOption(
		reScanFile,
		newConfig.LocalStores,
		config.Config.BookStores,
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
	if err := scan.InitStore(option); err != nil {
		logger.Infof("Failed to scan store path: %v", err)
	}
	if newConfig.EnableDatabase {
		saveResultsToDatabase(viper.ConfigFileUsed(), config.Config.ClearDatabaseWhenExit)
	}
}

func saveResultsToDatabase(configPath string, clearDatabaseWhenExit bool) {
	if err := scan.SaveResultsToDatabase(configPath, clearDatabaseWhenExit); err != nil {
		logger.Infof("Failed to save results to database: %v", err)
	}
}
