package routers

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/util/logger"
)

var engine *echo.Echo

func InitEcho() {
	// ***共通的 404 页面，需要在创建路由之前就替换***
	echo.NotFoundHandler = error_page.NotFoundCommon
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

	// Auto TLS
	// https://echo.labstack.com/docs/cookbook/auto-tls
	// https://tatsuo.medium.com/lets-encrypt-aautotls-pitfalls-f0e278b265c4
	// https://qiita.com/smith-30/items/147ba45fa74b2fc265b6

	// 将 HTTP 流量重定向到 HTTPS，您可以使用重定向中间件
	// 支持重定向到www子域名或非www子域名
	// https://echo.labstack.com/docs/middleware/redirect#https-redirect
	// engine.Pre(middleware.HTTPSRedirect())

	// 流式处理 JSON 响应
	// https://echo.labstack.com/docs/cookbook/streaming-response
	// 可以用来试试重写时长比较长的API 0064155000

	// sse（服务器发送数据）
	// https://echo.labstack.com/docs/cookbook/sse

	// 子域名
	// https://echo.labstack.com/docs/cookbook/subdomain

	// 类似推特的简单用户系统
	// https://echo.labstack.com/docs/cookbook/twitter

	// 允许的源列表
	// https://echo.labstack.com/docs/cookbook/cors

	// 开启 Gzip
	// 等级越高，CPU 占用越大、体积越小；通常 3–5 是压缩率与性能的平衡点。
	// 如果后面还挂了 Nginx / Caddy，并由它们统一做压缩，可以在 Echo 内部关闭 gzip，以免双重压缩。
	engine.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5, // 取值范围 -2～9；-2=DefaultCompression，0=NoCompression
		Skipper: func(c echo.Context) bool {
			// 如果url里面包含了 .js 或者 .css 文件或 base64 这几个关键字
			// 那么就启用 gzip 压缩
			url := c.Request().URL.Path
			if strings.Contains(url, ".js") || strings.Contains(url, ".css") || strings.Contains(url, ".wasm") || strings.Contains(url, ".htm") || strings.Contains(url, "base64") {
				// 包含以上关键字，启用 gzip 压缩
				return false
			}
			// 否则就不启用 gzip 压缩
			return false
		},
	}))

	// 禁止缓存中间件。使用 noCache ，会导强制浏览器每次都重新加载页面。除了测试和调试，一般不启用。
	// router.Use(noCache())

	// 反向代理中间件。
	// 反向代理中间件会将请求转发到后端服务器，并将响应返回给客户端。
	// 以下示例将使用默认的内存存储将应用程序限制为 20 个请求/秒：
	// https://echo.labstack.com/docs/cookbook/reverse-proxy
	// engine.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{}))

	// 速率限制 中间件。
	// https://echo.labstack.com/docs/middleware/rate-limiter
	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

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
	// 初始化web服务器
	InitEcho()
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
