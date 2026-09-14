package routers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/apiresp"
)

// controlAccess 集中保护管理路由；公开信息和阅读路由不经过此中间件。
func controlAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method != http.MethodGet && c.Request().Method != http.MethodHead && config.GetCfg().ReadOnlyMode {
			return apiresp.Forbidden(c, "config_locked", "Config is locked, cannot be modified", nil)
		}
		return next(c)
	}
}
