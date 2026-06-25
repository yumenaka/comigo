package config_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

// SaveConfigHandler 保存服务器配置到文件
func SaveConfigHandler(c echo.Context) error {
	// 如果配置被锁定，返回错误
	if config.GetCfg().ReadOnlyMode {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Config is locked, cannot be modified"})
	}
	SaveTo := c.Param("to")
	// 如果不是三个目录之一，就不能保存
	if !(SaveTo == "WorkingDirectory" || SaveTo == "HomeDirectory" || SaveTo == "ProgramDirectory") {
		logger.Infof(locale.GetString("log_error_failed_save_to_directory"), SaveTo)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed save to " + SaveTo + " directory",
		})
	}
	// 保存配置
	err := config.SaveConfig(SaveTo)
	if err != nil {
		logger.Infof("%s", err.Error())
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error": "Failed to save config",
		})
	}
	// 返回成功消息
	return GetConfigStatus(c)
}
