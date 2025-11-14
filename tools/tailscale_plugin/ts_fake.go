package tailscale_plugin

import (
	"context"

	"github.com/labstack/echo/v4"
)

type TailscaleConfig struct {
	Hostname   string
	Port       uint16
	FunnelMode bool
	ConfigDir  string
	AuthKey    string
}

// RunTailscale 启动Tailscale网络服务器，统一使用echo处理请求
func RunTailscale(e *echo.Echo, c TailscaleConfig) error {
	return nil
}

func (t *TailscaleStatus) GetTailscaleIP() string {
	return ""
}

// StopTailscale 停止Tailscale服务器
func StopTailscale() error {
	return nil
}

// GetTailscaleStatus 获取Tailscale服务状态 https://github.com/tailscale/golink/blob/b54cbbbb609ce8425193e7171a35af023cb5066d/golink.go#L787
func GetTailscaleStatus(ctx context.Context) (*TailscaleStatus, error) {
	return nil, nil
}
