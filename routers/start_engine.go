package routers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// startEngine 启动网页服务
func startEngine(e *echo.Echo) {
	// 是否对外服务
	webHost := ":"
	if config.GetDisableLAN() {
		webHost = "localhost:"
	}
	// 是否启用TLS
	enableTls := config.GetCertFile() != "" && config.GetKeyFile() != ""
	config.Srv = &http.Server{
		Addr:    webHost + strconv.Itoa(config.GetPort()),
		Handler: e, // echo.Echo 实现了 http.Handler 接口
	}
	// 在 goroutine 中初始化服务器，这样它就不会阻塞关闭处理
	go func() {
		// 监听并启动服务(TLS)
		if enableTls {
			if err := config.Srv.ListenAndServeTLS(config.GetCertFile(), config.GetKeyFile()); err != nil && !errors.Is(err, http.ErrServerClosed) {
				time.Sleep(3 * time.Second)
				log.Fatalf("listen: %s\n", err)
			}
		}
		if !enableTls {
			// 监听并启动服务(HTTP)
			if err := config.Srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				time.Sleep(3 * time.Second)
				log.Fatalf("listen: %s\n", err)
			}
		}
	}()
}
