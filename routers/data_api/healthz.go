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
