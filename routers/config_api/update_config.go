package config_api

import (
	"io"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
)

// UpdateConfig 修改服务器配置(post json)
func UpdateConfig(c echo.Context) error {
	// 读取请求体中的 JSON 数据
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
	}
	// 将 JSON 数据转换为字符串并打印
	jsonString := string(body)
	logger.Infof(locale.GetString("log_received_json_data"), jsonString)
	// 如果配置被锁定，返回错误
	if config.GetCfg().ReadOnlyMode {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Config is locked, cannot be modified"})
	}
	// 复制当前配置以便后续比较
	oldConfig := config.CopyCfg()
	// 解析 JSON 数据并更新服务器配置
	err = config.UpdateConfigByJson(jsonString)
	if err != nil {
		logger.Infof("%s", err.Error())
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": locale.GetString("log_failed_to_parse_json")})
	}
	// 如果 Language 配置发生变化，重新初始化语言设置
	if oldConfig.Language != config.GetCfg().Language {
		locale.InitLanguageFromConfig(config.GetCfg().Language)
	}
	err = config.UpdateConfigFile()
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), err)
	}
	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描书库等
	BeforeConfigUpdate(&oldConfig, config.GetCfg())
	// 返回成功消息
	return c.JSON(http.StatusOK, map[string]string{"message": "Server settings updated successfully"})
}

// BeforeConfigUpdate 根据配置的变化，判断是否需要打开浏览器重新扫描等
func BeforeConfigUpdate(oldConfig *config.Config, newConfig *config.Config) {
	OpenBrowserIfNeeded(oldConfig, newConfig)
	needReScan := checkNeedReScan(oldConfig, newConfig)
	if needReScan {
		StartReScan()
	} else {
		if newConfig.Debug {
			logger.Info(locale.GetString("log_no_changes_skipped_scan"))
		}
	}
}

func OpenBrowserIfNeeded(oldConfig *config.Config, newConfig *config.Config) {
	if (oldConfig.OpenBrowser == false) && (newConfig.OpenBrowser == true) {
		go tools.OpenBrowser(newConfig.EnableTLS, "127.0.0.1", newConfig.Port)
	}
}

// checkNeedReScan 检查旧的和新的配置是否需要更新，并返回需要重新扫描和重新扫描文件的布尔值
func checkNeedReScan(oldConfig *config.Config, newConfig *config.Config) (reScanStores bool) {
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
	return
}

// StartReScan 扫描并相应地更新数据库
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
