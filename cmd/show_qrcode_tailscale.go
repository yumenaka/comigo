package cmd

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

var tailscaleQRMu sync.Mutex
var tailscaleQRCancel context.CancelFunc

// stopTailscaleQRCode 取消上一轮轮询，避免重载后遗留定时任务。
func stopTailscaleQRCode() {
	tailscaleQRMu.Lock()
	defer tailscaleQRMu.Unlock()
	if tailscaleQRCancel != nil {
		tailscaleQRCancel()
		tailscaleQRCancel = nil
	}
}

// ShowQRCodeTailscale 每次调用替换上一轮轮询，计数和打印状态只属于本次调用。
func ShowQRCodeTailscale(ctx context.Context) {
	tailscaleQRMu.Lock()
	defer tailscaleQRMu.Unlock()
	if tailscaleQRCancel != nil {
		tailscaleQRCancel()
	}
	if !config.GetCfg().EnableTailscale {
		tailscaleQRCancel = nil
		return
	}
	ctx, tailscaleQRCancel = context.WithCancel(ctx)
	go pollTailscaleQRCode(ctx)
}

// pollTailscaleQRCode 使用单条可取消的轮询链，获取阅读地址后结束。
func pollTailscaleQRCode(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	authURLPrinted := false
	for counter := 0; counter < 101; counter++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		cfg := config.GetCfg()
		if !cfg.EnableTailscale {
			return
		}
		st, err := tailscale_plugin.GetTailscaleStatus(ctx)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			if cfg.Debug {
				logger.Infof(locale.GetString("log_tailscale_status_not_available"), err)
			}
			return
		}
		if st.Online {
			if st.FQDN == "" {
				logger.Info(locale.GetString("log_tailscale_not_yet_fqdn"))
				continue
			}
			proto := "http://"
			if cfg.TailscalePort == 443 || cfg.FunnelTunnel {
				proto = "https://"
			}
			readURL := proto + st.FQDN
			if cfg.TailscalePort != 80 && cfg.TailscalePort != 443 {
				readURL += ":" + strconv.Itoa(cfg.TailscalePort)
			}
			readURL += config.PrefixPath("/")
			logger.Info(locale.GetString("tailscale_reading_url") + readURL)
			tools.PrintQRCode(readURL)
			return
		}
		if st.AuthURL != "" && !authURLPrinted {
			logger.Info(locale.GetString("tailscale_auth_url_is") + st.AuthURL)
			tools.PrintQRCode(st.AuthURL)
			authURLPrinted = true
		}
	}
	logger.Info(locale.GetString("log_tailscale_status_check_exceeded"))
}
