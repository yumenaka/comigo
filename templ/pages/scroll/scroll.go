package scroll

import (
	"net/http"
	"strconv"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/templ/state"
	"github.com/yumenaka/comigo/util/logger"
)

// PageHandler 阅读界面（卷轴模式）
func PageHandler(c echo.Context) error {
	model.CheckAllBookFileExist()
	// 图片排序方式
	sortBy := "default"
	sortBookBy, err := c.Cookie("ScrollSortBy")
	if err == nil {
		// c.Cookie("key") 没找到，那么就会取到空值（nil），
		// 没处判断就直接访问 .Value 属性，会导致空指针引用错误。
		sortBy = sortBookBy.Value
		// fmt.Println("Scroll Mode Sort By:" + sortBy)
	}
	// 读取url参数，获取书籍ID
	bookID := c.Param("id")
	book, err := model.GetBookByID(bookID, sortBy)
	if err != nil {
		logger.Infof("GetBookByID: %v", err)
		// 没有找到书，显示 HTTP 404 错误
		indexHtml := common.Html(
			c,
			&state.Global,
			error_page.NotFound404(&state.Global),
			[]string{},
		)
		// 渲染 404 页面
		if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
			// 渲染失败，返回 HTTP 500 错误。
			return c.NoContent(http.StatusInternalServerError)
		}
		return nil
	}
	// 读取分页索引
	paginationIndex := -1
	page := c.QueryParam("page")
	if page != "" {
		index, err := strconv.Atoi(page)
		if err == nil {
			paginationIndex = index
		}
	}
	// 定义模板主体内容。
	scrollPage := ScrollPage(&state.Global, book, paginationIndex)
	// 拼接页面
	indexHtml := common.Html(
		c,
		&state.Global,
		scrollPage, // define body content
		[]string{"script/scroll.js"},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

// 跳转用分页链接 /scroll/4cTOjFm?page=1
func getScrollPaginationURL(book *model.Book, page int) string {
	readURL := `/scroll/` + book.BookID + `?page=` + strconv.Itoa(page)
	// href="javascript:void(0)" 是“点了什么也不发生”的老式写法
	if page < 1 {
		return `javascript:showToast(i18next.t('hint_first_page'), 'warning');`
	}
	if page > (book.PageCount/32 + 1) {
		return `javascript:showToast(i18next.t('hint_last_page'), 'warning')`
	}
	return readURL
}
