package router

import (
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/htmx/comigo"
	"github.com/yumenaka/comi/util/logger"
)

// RunServer 运行一个新的 HTTP 服务器。
func RunServer() (err error) {
	// 扫描漫画
	comigo.StartComigoServer()
	gin.SetMode(gin.ReleaseMode)
	// 创建一个新的Gin服务器。
	router := gin.Default()
	// 为模板引擎定义 HTML 渲染器。
	router.HTMLRender = &TemplRender{}
	// 静态文件。
	//router.Static("/static", "./router/static")
	// 嵌入静态文件。
	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		logger.Infof("%s", err)
	}
	router.StaticFS("/static/", http.FS(staticFS))
	// 设置路由
	setHandler(router)

	// 发消息
	slog.Info("Starting server...", "port", config.Config.Port)

	// 是否对外服务
	webHost := ":"
	if config.Config.DisableLAN {
		webHost = "localhost:"
	}
	// 是否启用TLS
	enableTLS := config.Config.CertFile != "" && config.Config.KeyFile != ""
	server := &http.Server{
		Addr:         webHost + strconv.Itoa(config.Config.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router, // gin.Engine本身可以作为一个Handler传递到http包,用于启动服务器
	}
	// 监听并启动服务(TLS)
	if enableTLS {
		if err = server.ListenAndServeTLS(config.Config.CertFile, config.Config.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			time.Sleep(3 * time.Second)
			logger.Fatalf("listen: %s\n", err)
		}
	}
	if !enableTLS {
		// 监听并启动服务(HTTP)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			time.Sleep(3 * time.Second)
			logger.Fatalf("listen: %s\n", err)
		}
	}
	return err
}
