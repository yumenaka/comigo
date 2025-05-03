package routers

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/util/logger"
)

var engine *echo.Echo

func init() {
	InitEcho()
}

func InitEcho() {
	// 创建新的 Echo 实例
	engine = echo.New()
	// 禁用 Echo 的 banner
	engine.HideBanner = true
	// 设置中间件
	SetMiddleware()
	// 绑定路由：内嵌资源
	EmbedStaticFiles()
	// 绑定路由：页面与API
	BindURLs()
}

// SetMiddleware 设置 Echo 的中间件等
func SetMiddleware() {
	// Recovery 中间件。返回 500 错误，避免程序直接崩溃，同时记录错误日志。
	engine.Use(middleware.Recover())

	// 设置 Echo 的日志输出
	SetEchoLogger(engine)

	// 禁止缓存中间件。使用 noCache ，会导强制浏览器每次都重新加载页面。除了测试和调试，一般不启用。
	// router.Use(noCache())

	// CORS 中间件
	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderContentLength, echo.HeaderAcceptEncoding, "X-CSRF-Token", echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true,
	}))
}

// EmbedStaticFiles 绑定静态资源
func EmbedStaticFiles() {
	// 嵌入JavaScript与css脚本
	var err error = nil
	assets.ScriptFS, err = fs.Sub(assets.Script, "script")
	if err != nil {
		logger.Infof("%s", err)
	}
	engine.StaticFS("/script/", assets.ScriptFS)
	// 嵌入图片资源
	assets.ImagesFS, err = fs.Sub(assets.Images, "images")
	if err != nil {
		logger.Infof("%s", err)
	}
	engine.StaticFS("/images/", assets.ImagesFS)
}

// StartWebServer 启动web服务
func StartWebServer() {
	// 设置网页端口
	SetHttpPort()
	// 监听并启动web服务
	StartEcho(engine)
}

// GetWebServer 获取echo.Echo (实现了 http.Handler 接口)
func GetWebServer() *echo.Echo {
	// 设置网页端口
	SetHttpPort()
	EmbedStaticFiles()
	// 设置中间件，绑定资源
	BindURLs()
	SetMiddleware()
	return engine
}
