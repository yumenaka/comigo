package settings

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

func getTranslations(value string) string {
	return "i18next.t(\"" + value + "\")"
}

// getTranslationsGO 返回 Go 语言的翻译字符串  not working
func getTranslationsGO(value string, lang string) string {
	return locale.GetStringByLocal(value, lang)
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
		//fmt.Printf("Error rendering settings page: %v\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
