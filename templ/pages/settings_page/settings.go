package settings_page

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/state"
)

func getTranslations(value string) string {
	return "i18next.t(\"" + value + "\")"
}

// Handler 设定页面
func Handler(c echo.Context) error {
	indexTemplate := common.Html(
		c,
		&state.Global,
		SettingsPage(c, &state.Global),
		[]string{},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexTemplate); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
