package flip

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/logger"
)

// PageHandler 阅读界面（翻页模式）
func PageHandler(c echo.Context) error {
	bookID := c.Param("id")
	logger.Info("Flip Mode Book ID:" + bookID)
	// 图片排序方式
	sortBy := "default"
	// c.Cookie("key") 没找到，那么就会取到空值（nil），没判断nil就直接访问 .Value 属性，会导致空指针引用错误。
	sortBookBy, err := c.Cookie("FlipSortBy")
	if err == nil {
		sortBy = sortBookBy.Value
	}
	// 读取url参数，获取书籍ID
	book, err := model.IStore.GetBook(bookID)
	if err != nil {
		logger.Infof("GetBook: %v", err)
		// 没有找到书，显示 HTTP 404 错误
		indexHtml := common.Html(
			c,
			error_page.NotFound404(c),
			[]string{},
		)
		// 渲染 404 页面
		if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
			// 渲染失败，返回 HTTP 500 错误。
			return c.NoContent(http.StatusInternalServerError)
		}
		return nil
	}
	book.SortPages(sortBy)

	// 翻页模式页面
	indexHtml := common.Html(
		c,
		FlipPage(c, book),
		[]string{"script/flip.js", "script/flip_sketch.js"})
	// 渲染翻页模式阅读页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 如果渲染失败，返回 HTTP 500 错误
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
