package config_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// GetConfigStatus 获取json格式的当前配置
func GetConfigStatus(c echo.Context) error {
	err := config.CfgStatus.SetConfigStatus()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get config"})
	}
	return c.JSON(http.StatusOK, config.CfgStatus)
}
