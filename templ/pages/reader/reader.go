package reader

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/templ/common"
)

// PageHandler 本地压缩包阅读页面。
// 这个页面不依赖后端保存书籍数据：用户选择的压缩包会在浏览器内由 WASM 解包，
// 后端只负责输出页面、公共脚本以及 reader 专用脚本。static 模式下这些脚本会被内联，
// 用于生成可离线打开的便携 HTML。
func PageHandler(c echo.Context) error {
	indexHtml := common.Html(
		c,
		ReaderPage(c),
		[]string{
			"static/wasm/wasm_exec.js",
			"static/js/flip_modules/pagination_utils.js",
			"static/js/flip_modules/interaction_utils.js",
			"static/js/reader.js",
			"static/js/reader_pwa.js",
		},
	)
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
