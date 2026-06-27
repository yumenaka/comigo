//go:build wails

package routers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

var wailsApp *application.App
var wailsWindow application.Window

type wailsOpenURLRequest struct {
	URL string `json:"url"`
}

type wailsDeleteBookFileRequest struct {
	BookID string `json:"bookId"`
}

type wailsAndroidFetchRequest struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// SetWailsRuntime 保存桌面壳对象，供普通 HTTP 页面触发窗口操作。
func SetWailsRuntime(app *application.App, window application.Window) {
	wailsApp = app
	wailsWindow = window
}

func bindWailsAPI(group *echo.Group) {
	group.GET("/wails/android-fetch/:payload", handleWailsAndroidFetch)
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

// handleWailsAndroidFetch 在 Android WebView 内重放 fetch 请求，绕过 Wails 资源拦截丢 method/body 的限制。
func handleWailsAndroidFetch(c echo.Context) error {
	data, err := base64.RawURLEncoding.DecodeString(c.Param("payload"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	var payload wailsAndroidFetchRequest
	if err := json.Unmarshal(data, &payload); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	targetPath := strings.TrimSpace(payload.Path)
	if !strings.HasPrefix(targetPath, "/api/") || strings.HasPrefix(targetPath, "/api/wails/android-fetch/") {
		return c.NoContent(http.StatusBadRequest)
	}
	method := strings.ToUpper(strings.TrimSpace(payload.Method))
	if method == "" {
		method = http.MethodGet
	}
	if strings.ContainsAny(method, " \t\r\n") {
		return c.NoContent(http.StatusBadRequest)
	}

	req := httptest.NewRequest(method, config.PrefixPath(targetPath), strings.NewReader(payload.Body))
	for key, value := range payload.Headers {
		if strings.EqualFold(key, echo.HeaderContentLength) || strings.EqualFold(key, "Host") {
			continue
		}
		req.Header.Set(key, value)
	}
	for _, cookie := range c.Request().Cookies() {
		req.AddCookie(cookie)
	}

	recorder := httptest.NewRecorder()
	c.Echo().ServeHTTP(recorder, req)
	result := recorder.Result()
	defer result.Body.Close()

	for key, values := range result.Header {
		if strings.EqualFold(key, echo.HeaderContentLength) {
			continue
		}
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}
	return c.Blob(result.StatusCode, result.Header.Get(echo.HeaderContentType), body)
}
