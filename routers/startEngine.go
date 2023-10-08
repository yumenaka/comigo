package routers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/common"
	"log"
	"net/http"
	"strconv"
	"time"
)

// startEngine 7、启动网页服务
func startEngine(engine *gin.Engine) {
	//是否对外服务
	webHost := ":"
	if common.Config.DisableLAN {
		webHost = "localhost:"
	}
	//是否启用TLS
	enableTls := common.Config.CertFile != "" && common.Config.KeyFile != ""
	common.Srv = &http.Server{
		Addr:    webHost + strconv.Itoa(common.Config.Port),
		Handler: engine,
	}

	//在 goroutine 中初始化服务器，这样它就不会阻塞下面的正常关闭处理
	go func() {
		// 监听并启动服务(TLS)
		if enableTls {
			if err := common.Srv.ListenAndServeTLS(common.Config.CertFile, common.Config.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				time.Sleep(3 * time.Second)
				log.Fatalf("listen: %s\n", err)
			}
		}
		if !enableTls {
			// 监听并启动服务(HTTP)
			if err := common.Srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				time.Sleep(3 * time.Second)
				log.Fatalf("listen: %s\n", err)
			}
		}
	}()
}
