package router

import (
	"errors"
	"github.com/yumenaka/comigo/routers"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/htmx/comigo"
	"github.com/yumenaka/comigo/htmx/embed_files"
	"github.com/yumenaka/comigo/util/logger"
)

// noCache 中间件设置 HTTP 响应头，禁用缓存。
func noCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		c.Next()
	}
}

// RunServer 运行一个新的 HTTP 服务器。
func RunServer() (err error) {
	gin.SetMode(gin.ReleaseMode)
	// 创建一个新的Gin服务器。
	router := gin.Default()

	// 使用 noCache 中间件，会导致浏览器每次都重新加载页面，不使用缓存。与翻页模式的预加载功能冲突。
	//router.Use(noCache())

	//Recovery 中间件。返回 500 错误页面，避免程序直接崩溃，同时记录错误日志。
	router.Use(gin.Recovery())

	SetGinLogger(router)
	routers.SetWebServerPort()
	// 扫描漫画
	comigo.StartComigoServer(router)
	// 为模板引擎定义 HTML 渲染器。
	router.HTMLRender = &TemplRender{}
	// 设置上传文件的最大内存限制
	router.MaxMultipartMemory = 5000 << 20 // 5000 MB

	// 设置嵌入静态文件的文件系统
	embed_files.StaticFS, err = fs.Sub(embed_files.Static, "static")
	if err != nil {
		logger.Infof("%s", err)
	}
	router.StaticFS("/static/", http.FS(embed_files.StaticFS))
	//favicon.ico
	router.GET("/favicon.ico", func(c *gin.Context) {
		file, err := embed_files.Static.ReadFile("/images/favicon.ico")
		if err != nil {
			logger.Infof("%s", err)
		}
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})
	// 设置路由
	setURLs(router)

	// 发消息
	logger.Infof("Starting server... port %v", config.GetPort())

	// 是否对外服务
	webHost := ":"
	if config.GetDisableLAN() {
		webHost = "localhost:"
	}
	// 是否启用TLS
	enableTLS := config.GetCertFile() != "" && config.GetKeyFile() != ""
	server := &http.Server{
		Addr:         webHost + strconv.Itoa(config.GetPort()),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router, // gin.Engine本身可以作为一个Handler传递到http包,用于启动服务器
	}
	// 监听并启动服务(TLS)
	if enableTLS {
		if err = server.ListenAndServeTLS(config.GetCertFile(), config.GetKeyFile()); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
