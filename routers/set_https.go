package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// SetAutoTLS 校验自动证书的启动前提，失败不降级成明文服务。
func SetAutoTLS(_ *echo.Echo) error {
	cfg := config.GetCfg()
	if !cfg.AutoTLSCertificate {
		return nil
	}
	if err := cfg.ValidateTLS(); err != nil {
		return err
	}
	// 端口绑定错误由 StartEcho 原样返回，包括无权限与端口被占用。
	cfg.Port = 443
	return nil
}
