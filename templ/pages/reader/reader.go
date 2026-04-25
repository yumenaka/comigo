package reader

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/templ/common"
)

// PageHandler 本地压缩包阅读页面。
func PageHandler(c echo.Context) error {
	indexHtml := common.Html(
		c,
		ReaderPage(c),
		[]string{"script/wasm/wasm_exec.js", "script/reader.js"},
	)
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
