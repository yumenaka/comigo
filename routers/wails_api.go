//go:build wails

package routers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var wailsContext context.Context

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
}
