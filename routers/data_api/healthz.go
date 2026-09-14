package data_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// Healthz 用于嵌入式宿主探测本地 HTTP 服务是否已可用。
func Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"port":   config.GetCfg().Port,
	})
}

// GetPublicInfo 只公开连接协商必需的信息，不包含书库、地址、设备或流量。
func GetPublicInfo(c echo.Context) error {
	c.Response().Header().Set("Cache-Control", "no-store")
	return c.JSON(http.StatusOK, map[string]any{
		"Version": config.GetVersion(), "desktopProtocol": 1,
		"requiresAuth": config.GetCfg().RequiresAuth(),
	})
}
