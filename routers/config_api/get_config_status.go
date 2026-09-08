package config_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
)

// GetConfigStatus 获取json格式的当前配置
func GetConfigStatus(c echo.Context) error {
	var status config.Status
	err := status.SetConfigStatus()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get config"})
	}
	c.Response().Header().Set("Cache-Control", "no-store")
	return c.JSON(http.StatusOK, status)
}

// GetConfig 返回当前配置的可读副本，凭据不回传到控制面板。
func GetConfig(c echo.Context) error {
	cfg := config.CopyCfg()
	cfg.Password = ""
	cfg.TailscaleAuthKey = ""
	cfg.DBDSN = ""
	cfg.StoreUrls = append([]string(nil), cfg.StoreUrls...)
	for i, url := range cfg.StoreUrls {
		cfg.StoreUrls[i] = tools.NormalizeStoreURLKey(url)
	}
	return c.JSON(http.StatusOK, cfg)
}
