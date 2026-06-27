package routers

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/websocket"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/sse_hub"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// autoTLSNextProtos 返回自动 TLS 模式需要声明的 ALPN 协议。
// 需要同时保留 ACME 的 tls-alpn-01，以及常规 HTTPS 使用的 h2 / http/1.1，
// 否则普通浏览器或 curl 在 TLS 握手阶段会收到 "no application protocol"。
func autoTLSNextProtos() []string {
	return []string{acme.ALPNProto, "h2", "http/1.1"}
}

type webServeMode int

const (
	webServeHTTP webServeMode = iota
	webServeCustomTLS
	webServeAutoTLS
)

// StartEcho 启动网页服务
func StartEcho(e *echo.Echo) error {
	logger.Infof(locale.GetString("log_starting_server_on_port"), config.GetCfg().Port)
	server, serveMode, err := buildHTTPServer(e)
	if err != nil {
		return err
	}
	config.Server = server
	listener, err := net.Listen("tcp", config.Server.Addr)
	if err != nil {
		return err
	}
	// 在 goroutine 中初始化 HTTP 服务器，这样它就不会阻塞关闭处理。
	go serveHTTPServer(listener, config.Server, serveMode)
	return nil
}

func buildHTTPServer(e *echo.Echo) (*http.Server, webServeMode, error) {
	SetAutoTLS(e)
	if config.GetCfg().AutoTLSCertificate {
		return buildAutoTLSServer(e)
	}
	serveMode := webServeHTTP
	if hasCustomTLSCertificate() {
		serveMode = webServeCustomTLS
	}
	return &http.Server{
		Addr:    webServerAddr(),
		Handler: e, // echo.Echo 实现了 http.Handler 接口
	}, serveMode, nil
}

func webServerAddr() string {
	// 绑定到 127.0.0.1，避免 localhost 在不同平台解析到 IPv6 时导致 WebView 无法通过 127.0.0.1 访问。
	if config.GetCfg().DisableLAN {
		return "127.0.0.1:" + strconv.Itoa(config.GetCfg().Port)
	}
	return ":" + strconv.Itoa(config.GetCfg().Port)
}

func hasCustomTLSCertificate() bool {
	return config.GetCfg().CertFile != "" && config.GetCfg().KeyFile != ""
}

func buildAutoTLSServer(e *echo.Echo) (*http.Server, webServeMode, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		logger.Errorf(locale.GetString("err_failed_to_get_config_dir"), err)
		return nil, webServeHTTP, err
	}
	autoTLSManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		// 缓存证书，避免触发 Let's Encrypt 频率限制。
		Cache:      autocert.DirCache(path.Join(configDir, "autotls_cache")),
		HostPolicy: autocert.HostWhitelist(config.GetCfg().Host),
	}
	logger.Infof(locale.GetString("log_auto_tls_enabled_for_domain"), config.GetCfg().Host)
	return &http.Server{
		Addr:    ":443",
		Handler: e,
		TLSConfig: &tls.Config{
			GetCertificate: autoTLSManager.GetCertificate,
			NextProtos:     autoTLSNextProtos(),
		},
	}, webServeAutoTLS, nil
}

func serveHTTPServer(listener net.Listener, server *http.Server, serveMode webServeMode) {
	switch serveMode {
	case webServeCustomTLS:
		logger.Infof(locale.GetString("log_custom_tls_cert"), config.GetCfg().CertFile, config.GetCfg().KeyFile)
		logServeError(server.ServeTLS(listener, config.GetCfg().CertFile, config.GetCfg().KeyFile))
	case webServeAutoTLS:
		tlsListener := tls.NewListener(listener, server.TLSConfig)
		logServeError(server.Serve(tlsListener))
	default:
		logServeError(server.Serve(listener))
	}
}

func logServeError(err error) {
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("listen: %s", err)
	}
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
	// 主动关闭所有 WebSocket 客户端，阅读页打开时也能快速退出
	websocket.CloseAll()
	// 停止 Tailscale HTTP 服务器（如启用）
	err := tailscale_plugin.StopTailscale()
	if err != nil {
		logger.Errorf(locale.GetString("err_error_stopping_tailscale_server"), err)
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

// RestartWebServer 停止当前服务器并重新启动，失败时返回错误，由调用方决定如何处理。
func RestartWebServer() error {
	if err := StopWebServer(); err != nil {
		return fmt.Errorf("%s: %w", locale.GetString("err_server_shutdown_failed"), err)
	}
	logger.Infof(locale.GetString("log_server_shutdown_successfully"), config.GetCfg().Port)
	// 重新初始化web服务器
	InitEcho()
	// 重新启动web服务器
	if err := StartEcho(engine); err != nil {
		return fmt.Errorf("%s: %w", locale.GetString("err_server_start_failed"), err)
	}
	return nil
}
