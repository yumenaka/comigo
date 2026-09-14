package settings

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/tools/autostart"
)

// 共享入口便于权限测试确认拒绝请求没有调用系统服务管理。
var startupStatus = autostart.Status
var setStartup = autostart.Set

func GetAutostartHandler(c echo.Context) error {
	state, err := startupStatus()
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Unable to read startup state")
	}
	return c.JSON(http.StatusOK, state)
}

// SetAutostartHandler 自启动状态不属于配置文件；只接受显式布尔值。
func SetAutostartHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
	}
	var input struct {
		Enabled *bool `json:"enabled"`
	}
	if err := c.Bind(&input); err != nil || input.Enabled == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "enabled must be a boolean")
	}
	if err := setStartup(*input.Enabled); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return GetAutostartHandler(c)
}
