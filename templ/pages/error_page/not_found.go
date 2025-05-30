package error_page

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/templ/common"
)

// NotFoundCommon 共通的 404 页面
func NotFoundCommon(c echo.Context) error {
	// 没有找到书，显示 HTTP 404 错误
	indexHtml := common.Html(
		c,
		NotFound404(c),
		[]string{},
	)
	// 渲染 404 页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
