package get_data_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

// GetTailscaleStatus 处理Tailscale网络的身份信息查询请求
func GetTailscaleStatus(c echo.Context) error {
	tailscaleStatus, err := tailscale_plugin.GetTailscaleStatus(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tailscaleStatus)
}
