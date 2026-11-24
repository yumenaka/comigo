package settings

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/scan"
)

// CleanupDataHandler 预览不写盘；执行必须携带当前预览指纹，并受只读设置限制。
func CleanupDataHandler(c echo.Context) error {
	token := ""
	if c.Request().Method == http.MethodPost {
		if err := ensureWritableConfig(); err != nil {
			return err
		}
		var request struct {
			Token string `json:"token"`
		}
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		if request.Token == "" {
			return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("cleanup_changed"))
		}
		token = request.Token
	}
	preview, err := scan.CleanupInvalidData(token)
	if errors.Is(err, scan.ErrCleanupChanged) {
		return echo.NewHTTPError(http.StatusConflict, locale.GetString("cleanup_changed"))
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("cleanup_failed")).SetInternal(err)
	}
	c.Response().Header().Set("Cache-Control", "no-store")
	return c.JSON(http.StatusOK, preview)
}
