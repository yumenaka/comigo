package routers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

// SetTailscale 根据配置启动或停止 Tailscale 服务
func SetTailscale(e *echo.Echo) {
	// 停止 Tailscale 服务
	if config.GetTailscaleEnable() == false {
		// 如果 Tailscale 服务正在运行,且Tailscale已经被禁用，则停止Tailscale服务
		tailscaleStatus, tsErr := tailscale_plugin.GetTailscaleStatus(context.Background())
		if tsErr != nil {
			logger.Errorf("Error getting Tailscale status: %v", tsErr)
		}
		if tsErr == nil && tailscaleStatus.Online && config.GetTailscaleEnable() == false {
			err := tailscale_plugin.StopTailscale()
			if err != nil {
				logger.Errorf("Error stopping Tailscale server: %v", err)
			}
		}
	}
	// 启动或重启 Tailscale 服务
	if config.GetTailscaleEnable() == true {
		go func() {
			// Funnel 模式，外网可访问，建议设置用户名与密码
			if config.GetTailscaleFunnelMode() {
				if err := tailscale_plugin.RunTailscale(
					e,
					tailscale_plugin.InitConfig{
						Hostname:   config.GetTailscaleHostname(),
						Port:       uint16(config.GetTailscalePort()),
						FunnelMode: config.GetTailscaleFunnelMode(),
						ConfigDir:  config.GetConfigPath(),
					},
				); err != nil {
					logger.Errorf("Failed to run Tailscale: %v", err)
				}
			}
			// 非Funnel模式，需要开启Tailscale服务才可以访问
			if !config.GetTailscaleFunnelMode() {
				if err := tailscale_plugin.RunTailscale(
					e,
					tailscale_plugin.InitConfig{
						Hostname:   config.GetTailscaleHostname(),
						Port:       uint16(config.GetTailscalePort()),
						FunnelMode: config.GetTailscaleFunnelMode(),
						ConfigDir:  config.GetConfigPath(),
					},
				); err != nil {
					logger.Errorf("Failed to run Tailscale: %v", err)
				}
			}
		}()
	}
}
