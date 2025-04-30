package routers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// noCache 中间件设置 HTTP 响应头，禁用缓存。
// 这将确保每次请求都从服务器获取最新的响应，而不是使用缓存的版本。
// 使用 noCache 中间件，会导强制浏览器每次都重新加载页面。与翻页模式的预加载功能冲突。所以除了测试和调试外，一般不启用。
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
