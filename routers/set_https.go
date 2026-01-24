package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// Auto TLS 设置，参考资料：
// https://echo.labstack.com/docs/cookbook/auto-tls
func SetAutoTLS(e *echo.Echo) {
	// 未设定自动 TLS
	if !config.GetCfg().AutoTLSCertificate {
		return
	}
	// 禁用局域网访问时，禁用自动 TLS。
	if config.GetCfg().DisableLAN {
		logger.Infof(locale.GetString("auto_tls_disabled_lan_access_off"))
		config.GetCfg().AutoTLSCertificate = false
		return
	}
	// 检测域名是否有效
	if config.GetCfg().Host == "" || !tools.IsValidDomain(config.GetCfg().Host) {
		// 自动 TLS 需要有效的域名才能工作，已禁用自动 TLS。
		logger.Infof(locale.GetString("auto_tls_disabled_invalid_domain"))
		config.GetCfg().AutoTLSCertificate = false
		return
	}
	// 检测端口是否可用
	if !tools.CheckPort(443) {
		// 443 端口已被占用，已禁用自动 TLS。
		logger.Infof(locale.GetString("port443_busy_disable_auto_tls"))
		config.GetCfg().AutoTLSCertificate = false
		return
	}
	if config.GetCfg().CertFile != "" && config.GetCfg().KeyFile != "" {
		// 已设置自定义证书，已禁用自动 TLS。
		logger.Infof(locale.GetString("custom_cert_detected_disable_auto_tls"))
		config.GetCfg().AutoTLSCertificate = false
		return
	}
	// 设置端口为 443
	if config.GetCfg().AutoTLSCertificate {
		config.GetCfg().Port = 443
	}
}
