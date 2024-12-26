package routers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
)

// startEngine 7、启动网页服务
func startEngine(engine *gin.Engine) {
	//是否对外服务
	webHost := ":"
	if config.Cfg.DisableLAN {
		webHost = "localhost:"
	}
	//是否启用TLS
	enableTls := config.Cfg.CertFile != "" && config.Cfg.KeyFile != ""
	config.Srv = &http.Server{
		Addr:    webHost + strconv.Itoa(config.GetPort()),
		Handler: engine, //gin.Engine本身可以作为一个Handler传递到http包,用于启动服务器
	}
	//在 goroutine 中初始化服务器，这样它就不会阻塞关闭处理
	//从端口启动开始,后续的所有工作都是http包来完成的 https://go.sai.show/PART02.%20Server/0.1-server-xiang-jie-yu-mian-shi-yao-dian
	go func() {
		// 监听并启动服务(TLS)
		if enableTls {
			if err := config.Srv.ListenAndServeTLS(config.Cfg.CertFile, config.Cfg.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
