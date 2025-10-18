package routers

import (
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

// StartTailscale 根据配置启动或停止 Tailscale 服务
func StartTailscale() {
	// 启动 Tailscale 服务的前提条件
	if engine == nil || config.GetTailscaleEnable() == false {
		return
	}
	// 启动或重启 Tailscale 服务
	if err := tailscale_plugin.RunTailscale(
		engine,
		tailscale_plugin.TailscaleConfig{
			Hostname:   config.GetCfg().TailscaleHostname,
			Port:       uint16(config.GetCfg().TailscalePort),
			FunnelMode: config.GetCfg().FunnelTunnel, // Tailscale Funnel 模式，外网可访问，建议设置用户名与密码
			ConfigDir:  config.GetCfg().ConfigPath,
			AuthKey:    config.GetCfg().TailscaleAuthKey,
		},
	); err != nil {
		logger.Errorf("Failed to run Tailscale: %v", err)
	}
}

// StopTailscale 停止 Tailscale 服务
func StopTailscale() {
	// 停止 Tailscale 服务
	err := tailscale_plugin.StopTailscale()
	if err != nil {
		logger.Errorf("Error stopping Tailscale server: %v", err)
	}
}
