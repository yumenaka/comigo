package routers

import (
	"errors"
	"io/fs"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/logger"
)

var engine *echo.Echo

func InitEcho() {
	// ***共通的 404 页面，需要在创建路由之前就替换***
	echo.NotFoundHandler = error_page.NotFoundCommon
	// 创建新的 Echo 实例
	engine = echo.New()
	// 禁用 Echo 的 banner
	engine.HideBanner = true
	SetHTTPErrorHandler(engine)
	// 设置中间件
	SetMiddleware()
	// 绑定路由：内嵌资源
	EmbedStaticFiles()
	// 绑定路由：页面与API
	BindURLs()
}

func SetHTTPErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		status := http.StatusInternalServerError
		code := "internal_error"
		message := locale.GetString("err_internal_server")
		var details interface{}

		if he, ok := errors.AsType[*echo.HTTPError](err); ok {
			status = he.Code
			switch v := he.Message.(type) {
			case string:
				message = v
			default:
				details = v
			}
			if status >= 400 && status < 500 {
				code = "bad_request"
			}
			if status == http.StatusUnauthorized {
				code = "unauthorized"
			}
			if status == http.StatusForbidden {
				code = "forbidden"
			}
			if status == http.StatusNotFound {
				code = "not_found"
			}
		}
		_ = apiresp.Error(c, status, code, message, details)
	}
}

// SetMiddleware 设置 Echo 的中间件等
func SetMiddleware() {
	// Recovery 中间件。返回 500 错误，避免程序直接崩溃，同时记录错误日志。
	engine.Use(middleware.Recover())
	engine.Use(middleware.RequestID())

	// 设置 Echo 的日志输出
	SetEchoLogger(engine)

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
		Level:   5, // 取值范围 -2～9；-2=DefaultCompression，0=NoCompression
		Skipper: gzipSkipper,
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

	engine.Use(middleware.CORSWithConfig(corsConfig()))
}

func gzipSkipper(c echo.Context) bool {
	urlPath := config.StripBasePath(c.Request().URL.Path)
	// SSE 与 WebSocket 是长连接，跳过 gzip 避免缓冲或握手异常；其他响应继续由 gzip 中间件判断是否压缩。
	return urlPath == "/api/sse" || urlPath == "/api/ws"
}

func corsConfig() middleware.CORSConfig {
	cors := middleware.CORSConfig{
		AllowMethods:  []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderContentLength, echo.HeaderAcceptEncoding, "X-CSRF-Token", echo.HeaderAuthorization, echo.HeaderXRequestID},
		ExposeHeaders: []string{echo.HeaderContentLength, echo.HeaderXRequestID},
	}
	if !config.GetCfg().RequiresAuth() {
		cors.AllowOrigins = []string{"*"}
		return cors
	}

	// 认证模式会携带 cookie/JWT，不能使用 "*" + AllowCredentials 的组合。
	cors.AllowOriginFunc = credentialedCORSOriginAllowed
	cors.AllowCredentials = true
	return cors
}

func credentialedCORSOriginAllowed(origin string) (bool, error) {
	host := originHost(origin)
	if host == "" {
		return false, nil
	}
	for _, allowedHost := range []string{"localhost", "127.0.0.1", "::1", originHost(config.GetCfg().Host)} {
		if allowedHost != "" && strings.EqualFold(host, allowedHost) {
			return true, nil
		}
	}
	return false, nil
}

func originHost(rawOrigin string) string {
	origin := strings.TrimSpace(rawOrigin)
	if origin == "" {
		return ""
	}
	if parsedURL, err := url.Parse(origin); err == nil && parsedURL.Hostname() != "" {
		return strings.ToLower(parsedURL.Hostname())
	}
	if parsedURL, err := url.Parse("//" + origin); err == nil && parsedURL.Hostname() != "" {
		return strings.ToLower(parsedURL.Hostname())
	}
	if strings.HasPrefix(origin, "[") && strings.HasSuffix(origin, "]") {
		return strings.ToLower(strings.Trim(origin, "[]"))
	}
	return strings.ToLower(origin)
}

// EmbedStaticFiles 绑定静态资源
func EmbedStaticFiles() {
	// 嵌入前端资源：/assets/dist 是编译产物，/assets/static 是页面级静态脚本。
	var err error = nil
	assets.FrontendFS, err = fs.Sub(assets.Frontend, ".")
	if err != nil {
		logger.Infof("%s", err)
	}
	engine.StaticFS(config.PrefixPath("/assets"), assets.FrontendFS)
	// 嵌入图片资源
	assets.ImagesFS, err = fs.Sub(assets.Images, "images")
	if err != nil {
		logger.Infof("%s", err)
	}
	engine.StaticFS(config.PrefixPath("/images"), assets.ImagesFS)
	// PWA manifest 使用标准 MIME，避免部分浏览器拒绝识别安装信息。
	engine.GET(config.PrefixPath("/images/manifest.webmanifest"), func(c echo.Context) error {
		data, err := fs.ReadFile(assets.Images, "images/manifest.webmanifest")
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "application/manifest+json", renderPwaManifest(data))
	})
	// PWA Service Worker 必须挂在根路径，才能覆盖 /reader 页面。
	engine.FileFS(config.PrefixPath("/reader-sw.js"), "pwa/reader-sw.js", assets.Pwa)
	// 暴露 robots.txt，供搜索引擎按标准路径读取
	engine.FileFS(config.PrefixPath("/robots.txt"), "robots.txt", assets.Robots)
}

func renderPwaManifest(data []byte) []byte {
	return []byte(rewritePwaManifestPaths(string(data)))
}

func rewritePwaManifestPaths(manifest string) string {
	// manifest 中的路径是按根路径编写的；这里统一加上 BasePath，保证反向代理子路径部署时 PWA 仍能定位资源。
	manifest = strings.ReplaceAll(manifest, `"start_url": "/reader"`, `"start_url": "`+config.PrefixPath("/reader")+`"`)
	manifest = strings.ReplaceAll(manifest, `"scope": "/"`, `"scope": "`+config.PrefixPath("/")+`"`)
	manifest = strings.ReplaceAll(manifest, `"src": "/images/`, `"src": "`+config.PrefixPath("/images/"))
	return manifest
}

// StartWebServer 启动web服务
func StartWebServer() error {
	// 初始化web服务器
	InitEcho()
	// 设置网页端口
	SetHttpPort()
	// 监听并启动web服务
	return StartEcho(engine)
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
