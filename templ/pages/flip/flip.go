package flip

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/logger"
)

// FlipModeHandler 阅读界面（翻页模式）
func FlipModeHandler(c echo.Context) error {
	bookID := c.Param("id")
	logger.Infof(locale.GetString("log_flip_mode_book_id"), bookID)
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
		logger.Infof(locale.GetString("log_get_book_error"), err)
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

	// 翻页模式
	indexHtml := common.Html(
		c,
		FlipPage(c, book),
		[]string{"script/flip.js", "script/flip_sketch.js"})
	// 静态模式
	staticMode := c.QueryParam("static") != ""
	if staticMode {
		// staticBook := *book 仅做浅拷贝，结构体中的 PageInfos 是切片，切片头被复制但底层数组仍与原对象共享，修改其元素会同步影响原书数据。
		staticBook := *book
		// 深拷贝 PageInfos 切片
		staticBook.PageInfos = make([]model.PageInfo, len(book.PageInfos))
		copy(staticBook.PageInfos, book.PageInfos)
		// 静态模式下，使用 base64 图片数据
		for i, image := range staticBook.PageInfos {
			staticBook.PageInfos[i].Url = common.GetFileBase64Text(book.BookInfo.BookID, image.Name)
		}
		indexHtml = common.Html(
			c,
			FlipPage(c, &staticBook),
			[]string{"script/flip.js", "script/flip_sketch.js"})
	}
	// 渲染翻页模式阅读页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 如果渲染失败，返回 HTTP 500 错误
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
