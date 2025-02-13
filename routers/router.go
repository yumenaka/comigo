package routers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/resource"
	"github.com/yumenaka/comigo/util/locale"
)

var engine *echo.Echo

func init() {
	// 创建新的 Echo 实例
	engine = echo.New()
	// 禁用 Echo 的 banner
	engine.HideBanner = true
}

func GetEngine() *echo.Echo {
	return engine
}

func SetEngine() {
	// 使用 Recovery 中间件
	engine.Use(middleware.Recover())

	// CORS 中间件
	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderContentLength, echo.HeaderAcceptEncoding, "X-CSRF-Token", echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true,
	}))

	// 嵌入静态文件到二进制文件
	resource.EmbedResoure(engine, locale.GetString("html_title")+config.GetVersion())

	// 设置各种API
	BindAPI(engine)
}

// StartWebServer 启动web服务
func StartWebServer() {
	SetEngine()
	// 显示QRCode
	showQRCode()
	// 监听并启动web服务
	startEngine(engine)
}
