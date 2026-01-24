package config_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

const (
	HomeDirectory    = "HomeDirectory"
	WorkingDirectory = "WorkingDirectory"
	ProgramDirectory = "ProgramDirectory"
)

// DeleteConfig 删除配置文件
func DeleteConfig(c echo.Context) error {
	// 如果配置被锁定，返回错误
	if config.GetCfg().ReadOnlyMode {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Config is locked, cannot be modified"})
	}
	in := c.Param("in")
	validDirs := []string{WorkingDirectory, HomeDirectory, ProgramDirectory}

	if !contains(validDirs, in) {
		logger.Infof(locale.GetString("log_error_failed_to_delete_config"), in)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to delete config in " + in + " directory",
		})
	}

	err := config.DeleteConfigIn(in)
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error": "Failed to delete config",
		})
	}

	return GetConfigStatus(c)
}

// contains 检查切片是否包含特定字符串
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
