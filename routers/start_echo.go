package routers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/sse_hub"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

// StartEcho 启动网页服务
func StartEcho(e *echo.Echo) {
	// 是否对外服务
	webHost := ":"
	if config.GetDisableLAN() {
		webHost = "localhost:"
	}
	// 是否启用TLS
	enableTls := config.GetCertFile() != "" && config.GetKeyFile() != ""
	config.Server = &http.Server{
		Addr:    webHost + strconv.Itoa(config.GetCfg().Port),
		Handler: e, // echo.Echo 实现了 http.Handler 接口
	}
	// 记录日志并启动服务器
	logger.Infof("Starting Server...", "on port", config.GetCfg().Port, "...")
	if enableTls {
		logger.Infof("TLS enabled", "CertFile:", config.GetCertFile(), "KeyFile:", config.GetKeyFile())
	}
	// 在 goroutine 中初始化 HTTP 服务器，这样它就不会阻塞关闭处理
	go func() {
		// 监听并启动服务(TLS)
		if enableTls {
			if err := config.Server.ListenAndServeTLS(config.GetCertFile(), config.GetKeyFile()); err != nil && !errors.Is(err, http.ErrServerClosed) {
				time.Sleep(3 * time.Second)
				logger.Fatalf("listen: %s\n", err)
			}
		}
		if !enableTls {
			// 监听并启动服务(HTTP)
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
