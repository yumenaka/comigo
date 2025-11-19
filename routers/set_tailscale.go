package routers

import (
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

// StartTailscale 根据配置启动或停止 Tailscale 服务
func StartTailscale() {
	// 启动 Tailscale 服务的前提条件
	if engine == nil || config.GetCfg().EnableTailscale == false {
		return
	}
	// 如果启用了 Tailscale Funnel 模式且要求身份验证，但未设置用户名或密码，则记录错误并返回
	if config.GetCfg().EnableTailscale && config.GetCfg().FunnelTunnel && config.GetCfg().FunnelLoginCheck && !config.GetCfg().RequiresAuth() {
		logger.Errorf(locale.GetString("FunnelLoginCheckDescription"))
		return
	}
	// 启动或重启 Tailscale 服务
	configDir, err := config.GetConfigDir()
	if err != nil {
		configDir = ""
	}
	if tsError := tailscale_plugin.RunTailscale(
		engine,
		tailscale_plugin.TailscaleConfig{
			Hostname:   config.GetCfg().TailscaleHostname,
			Port:       uint16(config.GetCfg().TailscalePort),
			FunnelMode: config.GetCfg().FunnelTunnel, // Tailscale Funnel 模式，外网可访问，建议设置用户名与密码
			ConfigDir:  configDir,
			AuthKey:    config.GetCfg().TailscaleAuthKey,
		},
	); tsError != nil {
		logger.Errorf("Failed to run Tailscale: %v", tsError)
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
