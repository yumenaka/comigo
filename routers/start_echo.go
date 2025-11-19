package routers

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/sse_hub"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// StartEcho 启动网页服务
func StartEcho(e *echo.Echo) {
	// 是否仅监听本地回环地址
	webHost := ":"
	if config.GetCfg().DisableLAN {
		webHost = "localhost:"
	}
	// 记录日志并启动服务器
	logger.Infof("Starting Server...", "on port", config.GetCfg().Port, "...")
	// 初始化 HTTP 服务器
	config.Server = &http.Server{
		Addr:    webHost + strconv.Itoa(config.GetCfg().Port),
		Handler: e, // echo.Echo 实现了 http.Handler 接口
	}
	// 如果自定义 TLS 证书
	customCertTlS := config.GetCfg().CertFile != "" && config.GetCfg().KeyFile != ""
	// 如果自动申请 TLS 证书
	SetAutoTLS(e)
	if config.GetCfg().AutoTLSCertificate {
		configDir, err := config.GetConfigDir()
		if err != nil {
			configDir = ""
		}
		autoTLSManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
			Cache:      autocert.DirCache(path.Join(configDir, "autotls_cache")),
			HostPolicy: autocert.HostWhitelist(config.GetCfg().Host),
		}
		// 更新服务器配置以使用自动 TLS
		logger.Infof("Auto TLS enabled for domain:", config.GetCfg().Host)
		config.Server = &http.Server{
			Addr:    ":443",
			Handler: e, // set Echo as handler
			TLSConfig: &tls.Config{
				//Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
				GetCertificate: autoTLSManager.GetCertificate,
				NextProtos:     []string{acme.ALPNProto},
			},
			//ReadTimeout: 30 * time.Second, // use custom timeouts
		}
	}
	// 在 goroutine 中初始化 HTTP 服务器，这样它就不会阻塞关闭处理
	go func() {
		// 监听并启动服务(自定义TLS证书)
		if config.GetCfg().CertFile != "" && config.GetCfg().KeyFile != "" {
			logger.Infof("Custom TLS Cert", "CertFile:", config.GetCfg().CertFile, "KeyFile:", config.GetCfg().KeyFile)
			if err := config.Server.ListenAndServeTLS(config.GetCfg().CertFile, config.GetCfg().KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				time.Sleep(3 * time.Second)
				logger.Fatalf("listen: %s\n", err)
			}
		}
		// 监听并启动服务(自动TLS)
		if config.GetCfg().AutoTLSCertificate {
			if err := config.Server.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
				time.Sleep(3 * time.Second)
				logger.Fatalf("listen: %s\n", err)
			}
		}
		// 监听并启动服务(无TLS,普通HTTP)
		if !customCertTlS && !config.GetCfg().AutoTLSCertificate {
			if err := config.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				time.Sleep(3 * time.Second)
				logger.Fatalf("listen: %s\n", err)
			}
		}
	}()
}

// StopWebServer 停止当前的HTTP服务器
func StopWebServer() error {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	if config.Server == nil {
		return nil
	}
	// 主动关闭所有 SSE 客户端，避免优雅关闭时被长连接阻塞
	sse_hub.MessageHub.CloseAll()
	// 停止 Tailscale HTTP 服务器（如启用）
	err := tailscale_plugin.StopTailscale()
	if err != nil {
		logger.Errorf("Error stopping Tailscale server: %v", err)
	}
	// 关闭服务器（deadline 5秒）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Shutdown:一旦把服务器关机，就不能重复使用。将来的呼叫将返回”ErrServerClosed“。
	if err := config.Server.Shutdown(ctx); err != nil {
		// 超时或其他错误时回退到强制关闭，避免整个进程退出
		if closeErr := config.Server.Close(); closeErr != nil {
			return closeErr
		}
	}
	// 需要重新初始化服务器。也就是调用InitEcho()
	config.Server = nil
	return nil
}

// RestartWebServer 停止当前的服务器并重新启动(比如切换端口的时候)
func RestartWebServer() {
	if err := StopWebServer(); err != nil {
		logger.Fatalf("Server Shutdown Failed:%+v", err)
	}
	logger.Infof("Server Shutdown Successfully", "Starting Server...", "on port", config.GetCfg().Port, "...")
	// 重新初始化web服务器
	InitEcho()
	// 重新启动web服务器
	StartEcho(engine)
}
