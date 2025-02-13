package router

import (
	"errors"
	"github.com/yumenaka/comigo/htmx/comigo"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/htmx/embed"
	"github.com/yumenaka/comigo/util/logger"
)

// noCache 中间件设置 HTTP 响应头，禁用缓存。
func noCache() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Response().Header().Set("Pragma", "no-cache")
			c.Response().Header().Set("Expires", "0")
			c.Response().Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
			return next(c)
		}
	}
}

// RunServer 运行一个新的 HTTP 服务器。
func RunServer() (err error) {
	// 创建一个新的Echo服务器
	router := echo.New()
	router.HideBanner = true

	// 使用 noCache 中间件，会导致浏览器每次都重新加载页面，不使用缓存。与翻页模式的预加载功能冲突。
	// router.Use(noCache())

	// Recovery 中间件。返回 500 错误，避免程序直接崩溃，同时记录错误日志。
	router.Use(middleware.Recover())

	// 设置 Echo 的日志输出
	SetEchoLogger(router)

	// 扫描漫画
	comigo.SetComigoServer(router)

	// 设置嵌入静态文件的文件系统
	embed.StaticFS, err = fs.Sub(embed.Static, "static")
	if err != nil {
		logger.Infof("%s", err)
	}
	router.StaticFS("/static/", embed.StaticFS)

	// favicon.ico
	router.GET("/favicon.ico", func(c echo.Context) error {
		file, err := embed.Static.ReadFile("/images/favicon.ico")
		if err != nil {
			logger.Infof("%s", err)
			return err
		}
		return c.Blob(http.StatusOK, "image/x-icon", file)
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

	// 配置服务器
	router.Server = &http.Server{
		Addr:         webHost + strconv.Itoa(config.GetPort()),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 监听并启动服务
	if enableTLS {
		if err = router.StartTLS(router.Server.Addr, config.GetCertFile(), config.GetKeyFile()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			time.Sleep(3 * time.Second)
			logger.Fatalf("listen: %s\n", err)
		}
	} else {
		if err = router.Start(router.Server.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			time.Sleep(3 * time.Second)
			logger.Fatalf("listen: %s\n", err)
		}
	}

	return err
}
