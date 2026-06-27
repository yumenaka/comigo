//go:build wails

package routers

import (
	"net/http"
	"net/url"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/yumenaka/comigo/assets/locale"
)

var wailsApp *application.App
var wailsWindow application.Window

type wailsOpenURLRequest struct {
	URL string `json:"url"`
}

type wailsDeleteBookFileRequest struct {
	BookID string `json:"bookId"`
}

// SetWailsRuntime 保存桌面壳对象，供普通 HTTP 页面触发窗口操作。
func SetWailsRuntime(app *application.App, window application.Window) {
	wailsApp = app
	wailsWindow = window
}

func bindWailsAPI(group *echo.Group) {
	group.POST("/wails/toggle-fullscreen", func(c echo.Context) error {
		if wailsWindow == nil {
			return c.NoContent(http.StatusServiceUnavailable)
		}
		wailsWindow.ToggleFullscreen()
		return c.NoContent(http.StatusNoContent)
	})
	group.POST("/wails/open-url", func(c echo.Context) error {
		if wailsApp == nil {
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
		if err := wailsApp.Browser.OpenURL(req.URL); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusNoContent)
	})
	group.POST("/wails/select-directory", func(c echo.Context) error {
		if wailsApp == nil || wailsWindow == nil {
			return c.NoContent(http.StatusServiceUnavailable)
		}
		if runtime.GOOS == "android" {
			return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("wails_android_directory_picker_unsupported"))
		}
		// 只在 Wails 桌面壳里调用系统目录选择器，普通 Web 环境不暴露本机路径能力。
		path, err := wailsApp.Dialog.OpenFile().
			SetTitle(locale.GetString("select_store_folder")).
			CanChooseFiles(false).
			CanChooseDirectories(true).
			CanCreateDirectories(true).
			ResolvesAliases(true).
			AttachToWindow(wailsWindow).
			PromptForSingleSelection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if path == "" {
			return c.NoContent(http.StatusNoContent)
		}
		return c.JSON(http.StatusOK, map[string]string{"path": path})
	})
	group.POST("/wails/delete-book-file", func(c echo.Context) error {
		var req wailsDeleteBookFileRequest
		if err := c.Bind(&req); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		deleted, err := DeleteBookFileForWails(req.BookID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"deleted": deleted,
			"message": DeleteBookFileSuccessMessageForWails(),
		})
	})
}
