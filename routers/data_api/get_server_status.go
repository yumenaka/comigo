package data_api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

// GetBooksNumber 获取书籍总数，不包括 BookGroup 类型
func GetAllBooksNumber() int {
	// 用于计数的变量
	var count int
	// 遍历 map 并递增计数器
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
		if b.Type == model.TypeBooksGroup {
			continue // 跳过书组类型
		}
		count++
	}
	return count
}

// GetServerInfoHandler 获取服务器信息的API处理函数
func GetServerInfoHandler(c echo.Context) error {
	serverStatus := tools.GetServerInfo(
		tools.ServerInfoParams{
			Cfg:            config.GetCfg(),
			Version:        config.GetVersion(),
			AllBooksNumber: GetAllBooksNumber(),
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
