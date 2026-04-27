package settings

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/service"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

var localeKeyCollisionLabels = map[string]string{
	"EnableDatabase":  "enable_database_label",
	"FunnelTunnel":    "funnel_tunnel_label",
	"OpenBrowser":     "open_browser_label",
	"ReEnterPassword": "re_enter_password_label",
	"Timeout":         "timeout_label",
}

var (
	localeKeyUpperWordRegexp = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	localeKeyCamelRegexp     = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	localeKeySeparatorRegexp = regexp.MustCompile(`[^A-Za-z0-9]+`)
	localeKeyRepeatsRegexp   = regexp.MustCompile(`_+`)
)

func toLocaleKey(value string) string {
	if key, ok := localeKeyCollisionLabels[value]; ok {
		return key
	}
	key := localeKeyUpperWordRegexp.ReplaceAllString(value, `${1}_${2}`)
	key = localeKeyCamelRegexp.ReplaceAllString(key, `${1}_${2}`)
	key = localeKeySeparatorRegexp.ReplaceAllString(key, "_")
	key = localeKeyRepeatsRegexp.ReplaceAllString(key, "_")
	return strings.ToLower(strings.Trim(key, "_"))
}

func getTranslations(value string) string {
	return "i18next.t(\"" + toLocaleKey(value) + "\")"
}

// GetStoreBookCounts 获取每个书库的书籍数量
// 返回 map[storeUrl]bookCount
func GetStoreBookCounts() map[string]int {
	counts := make(map[string]int)

	// 获取所有书籍
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
		return counts
	}

	// 统计每个书库的书籍数量
	for _, book := range allBooks {
		// 只统计非书籍组的实际书籍
		if book.Type != model.TypeBooksGroup {
			// 将书库路径转换为绝对路径以便匹配
			storePathAbs, err := filepath.Abs(book.StoreUrl)
			if err != nil {
				storePathAbs = book.StoreUrl
			}
			counts[storePathAbs]++
		}
	}

	return counts
}

// PageHandler 设定页面
func PageHandler(c echo.Context) error {
	tsStatus, err := tailscale_plugin.GetTailscaleStatus(c.Request().Context())
	if err != nil {
		// 容错：当 Tailscale 未启用或尚未就绪时，不中断页面渲染，返回离线状态
		tsStatus = &tailscale_plugin.TailscaleStatus{}
	}
	indexHtml := common.Html(
		c,
		SettingsPage(c, tsStatus),
		[]string{},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

// -----------------------------------------
// BeforeConfigUpdate 在配置更新前执行的操作
// -------------------------

// RestartWebServerBroadcast 用于广播重启网页服务器的信号
var RestartWebServerBroadcast *chan string

// -------------------------
// 各种辅助函数
// -------------------------

// beforeConfigUpdate 根据配置的变化，判断是否需要打开浏览器重新扫描等
func beforeConfigUpdate(oldConfig config.Config, newConfig *config.Config) {
	if RestartWebServerBroadcast == nil {
		service.ApplyConfigChange(oldConfig, newConfig, nil)
		return
	}
	service.ApplyConfigChange(oldConfig, newConfig, *RestartWebServerBroadcast)
}
