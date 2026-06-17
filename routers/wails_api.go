//go:build wails

package routers

import (
	"context"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var wailsContext context.Context

type wailsOpenURLRequest struct {
	URL string `json:"url"`
}

// SetWailsContext 保存桌面壳上下文，供普通 HTTP 页面触发窗口操作。
func SetWailsContext(ctx context.Context) {
	wailsContext = ctx
}

func bindWailsAPI(group *echo.Group) {
	group.POST("/wails/toggle-fullscreen", func(c echo.Context) error {
		if wailsContext == nil {
			return c.NoContent(http.StatusServiceUnavailable)
		}
		if wailsruntime.WindowIsFullscreen(wailsContext) {
			wailsruntime.WindowUnfullscreen(wailsContext)
		} else {
			wailsruntime.WindowFullscreen(wailsContext)
		}
		return c.NoContent(http.StatusNoContent)
	})
	group.POST("/wails/open-url", func(c echo.Context) error {
		if wailsContext == nil {
			return c.NoContent(http.StatusServiceUnavailable)
		}
		var req wailsOpenURLRequest
		if err := c.Bind(&req); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		parsed, err := url.Parse(req.URL)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return c.NoContent(http.StatusBadRequest)
		}
		// Wails WebView 里 target=_blank 不一定会交给系统浏览器，需由宿主显式打开外部 URL。
		wailsruntime.BrowserOpenURL(wailsContext, req.URL)
		return c.NoContent(http.StatusNoContent)
	})
}
