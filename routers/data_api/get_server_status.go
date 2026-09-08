package data_api

import (
	"net"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/releases"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
	"github.com/yumenaka/comigo/tools/traffic"
)

// GetServerInfoHandler 获取服务器信息的API处理函数
func GetServerInfoHandler(c echo.Context) error {
	serverStatus := tools.GetServerInfo(
		tools.ServerInfoParams{
			Cfg:                   config.GetCfg(),
			Version:               config.GetVersion(),
			AllBooksNumber:        model.GetAllBooksNumber(),
			ClientIP:              c.RealIP(),
			ReScanServiceEnable:   config.GlobalLibraryScanner.IsRunning(),
			ReScanServiceInterval: config.GlobalLibraryScanner.GetInterval(),
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
	connections, users, devices := connectionSnapshot()
	return c.JSON(http.StatusOK, struct {
		*tools.ServerStatus
		ReadingURL      string           `json:"readingURL"`
		LocalBrowserURL string           `json:"localBrowserURL"`
		LocalIPs        []string         `json:"localIPs"`
		Traffic         traffic.Snapshot `json:"traffic"`
		Update          releases.Status  `json:"update"`
		OnlineUsers     int              `json:"onlineUsers"`
		OnlineDevices   int              `json:"onlineDevices"`
		Connections     int              `json:"connections"`
		ExternalAccess  bool             `json:"externalAccess"`
		ListenAddress   string           `json:"listenAddress"`
	}{serverStatus, config.GetQrcodeURL(), config.GetLocalBrowserURL(), localIPs(), traffic.Default.Snapshot(), releases.Default.Snapshot(), users, devices, len(connections), !config.GetCfg().DisableLAN, config.GetListenHost()})
}

// localIPs 只列出启用网卡的可用单播地址，避免把客户端地址当成服务地址。
func localIPs() []string {
	result := []string{}
	seen := map[string]bool{}
	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addresses, _ := iface.Addrs()
		for _, address := range addresses {
			ip, _, err := net.ParseCIDR(address.String())
			if err != nil || !ip.IsGlobalUnicast() || seen[ip.String()] {
				continue
			}
			seen[ip.String()] = true
			result = append(result, ip.String())
		}
	}
	sort.Strings(result)
	return result
}

// GetServerUpdateHandler 只检查版本，不下载或执行更新。
func GetServerUpdateHandler(c echo.Context) error {
	status, err := releases.Default.Check(c.Request().Context(), config.GetVersion())
	if err != nil {
		return c.JSON(http.StatusBadGateway, status)
	}
	return c.JSON(http.StatusOK, status)
}

// GetServerTrafficHandler 为桌面面板提供轻量快照，避免高频查询系统信息和 Tailscale。
func GetServerTrafficHandler(c echo.Context) error {
	c.Response().Header().Set("Cache-Control", "no-store")
	return c.JSON(http.StatusOK, traffic.Default.Snapshot())
}
