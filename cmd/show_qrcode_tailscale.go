package cmd

import (
	"context"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

var authUrlPrinted bool = false
var readUrlPrinted bool = false
var counter int = 0

// ShowQRCodeTailscale 打印 Tailscale 访问地址的二维码
func ShowQRCodeTailscale(ctx context.Context) {
	if counter > 100 {
		logger.Info("Tailscale status check exceeded 100 times, stopping further checks.")
		return
	}
	// 1 秒后打印 Url（不会阻塞 main 线程）
	time.AfterFunc(1*time.Second, func() {
		// 前提： 启用 Tailscale 服务
		if config.GetCfg().GetEnableTailscale() == false {
			logger.Info(locale.GetString("tailscale_not_enabled"))
			return
		}
		// 获取 Tailscale 状态
		st, err := tailscale_plugin.GetTailscaleStatus(ctx)
		if err != nil {
			logger.Info("Failed to get Tailscale status:", err)
			return
		}
		// 计数器+1
		counter++
		// 打印 Tailscale 访问链接（只打印一次）
		if st.Online == true {
			if st.FQDN != "" {
				if readUrlPrinted == false {
					readURL := "https://" + st.FQDN
					if config.GetTailscalePort() != 443 && config.GetTailscalePort() != 8443 && config.GetTailscalePort() != 10000 {
						readURL = "http://" + st.FQDN
					}
					logger.Info(locale.GetString("tailscale_reading_url") + readURL)
					tools.PrintQRCode(readURL)
					readUrlPrinted = true
				}
			} else {
				logger.Info("tailscale_not_yet_fqdn")
				ShowQRCodeTailscale(ctx)
				return
			}
		}
		// 打印 Tailscale 授权验证链接（只打印一次）
		if st.Online == false {
			if st.AuthURL != "" {
				if authUrlPrinted == false {
					logger.Info(locale.GetString("tailscale_auth_url_is") + st.AuthURL)
					tools.PrintQRCode(st.AuthURL)
					authUrlPrinted = true
					ShowQRCodeTailscale(ctx)
				}
			}
			// 继续检查状态，直到上线
			ShowQRCodeTailscale(ctx)
			return
		}
	})
}
