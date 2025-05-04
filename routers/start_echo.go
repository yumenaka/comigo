package routers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/logger"
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
		Addr:    webHost + strconv.Itoa(config.GetPort()),
		Handler: e, // echo.Echo 实现了 http.Handler 接口
	}
	// 在 goroutine 中初始化服务器，这样它就不会阻塞关闭处理
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
	// 关闭服务器（延时1秒）
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// Shutdown:一旦把服务器关机，就不能重复使用。将来的呼叫将返回”ErrServerClosed“。
	if err := config.Server.Shutdown(ctx); err != nil {
		return err
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
	logger.Infof("Server Shutdown Successfully", "Starting Server...", "on port", config.GetPort(), "...")
	// 重新初始化web服务器
	InitEcho()
	// 重新启动web服务器
	StartEcho(engine)
}
