package get_data_api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

func GetServerInfoHandler(c echo.Context) error {
	serverStatus := tools.GetServerInfo(
		tools.ServerInfoParams{
			Cfg:            config.GetCfg(),
			Version:        config.GetVersion(),
			AllBooksNumber: model.IStore.GetAllBooksNumber(),
			ClientIP:       c.RealIP(),
		})
	tailscaleStatus, err := tailscale_plugin.GetTailscaleStatus(c.Request().Context())
	if err == nil {
		// 设置 Tailscale 认证 URL
		if tailscaleStatus.AuthURL != "" {
			serverStatus.TailscaleAuthURL = tailscaleStatus.AuthURL
		}
		// 设置 Tailscale 访问 URL
		if tailscaleStatus.FQDN != "" {
			proto := "http://"
			if config.GetCfg().TailscalePort == 443 {
				proto = "https://"
			}
			if config.GetCfg().FunnelTunnel && (config.GetCfg().TailscalePort == 8443 || config.GetCfg().TailscalePort == 10000) {
				proto = "https://"
			}
			href := proto + tailscaleStatus.FQDN
			if config.GetCfg().TailscalePort != 443 && config.GetCfg().TailscalePort != 80 {
				href += ":" + strconv.Itoa(config.GetCfg().TailscalePort)
			}
			serverStatus.TailscaleUrl = href
		}
	}
	return c.JSON(http.StatusOK, serverStatus)
}
