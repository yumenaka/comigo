package upload_page

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/state"
)

// PageHandler 上传文件页面
func PageHandler(c echo.Context) error {
	indexHtml := common.Html(
		c,
		&state.Global,
		UploadPage(c, &state.Global),
		[]string{},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
