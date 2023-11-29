package routers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/config"
)

// startEngine 7、启动网页服务
func startEngine(engine *gin.Engine) {
	//是否对外服务
	webHost := ":"
	if config.Config.DisableLAN {
		webHost = "localhost:"
	}
	//是否启用TLS
	enableTls := config.Config.CertFile != "" && config.Config.KeyFile != ""
	config.Srv = &http.Server{
		Addr:    webHost + strconv.Itoa(config.Config.Port),
		Handler: engine,
	}
	//在 goroutine 中初始化服务器，这样它就不会阻塞关闭处理
	go func() {
		// 监听并启动服务(TLS)
		if enableTls {
			if err := config.Srv.ListenAndServeTLS(config.Config.CertFile, config.Config.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
