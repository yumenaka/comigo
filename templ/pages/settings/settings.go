package settings

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/templ/common"
)

func getTranslations(value string) string {
	return "i18next.t(\"" + value + "\")"
}

// PageHandler 设定页面
func PageHandler(c echo.Context) error {
	indexHtml := common.Html(
		c,
		SettingsPage(c),
		[]string{},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
