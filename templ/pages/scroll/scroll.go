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

// Handler 阅读界面（卷轴模式）
func Handler(c echo.Context) error {
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
	paginationIndex := 0
	pagination, err := c.Cookie("PaginationIndex")
	if err == nil {
		index, err := strconv.Atoi(pagination.Value)
		if err == nil {
			logger.Infof("PaginationIndex: %v", err)
			paginationIndex = index
		}
	}

	// 读取url参数，获取书籍ID
	bookID := c.Param("id")
	// 没有找到书籍，显示 HTTP 404 错误
	indexTemplate := common.Html(
		c,
		&state.Global,
		error_page.NotFound404(&state.Global),
		[]string{},
	)
	book, err := model.GetBookByID(bookID, sortBy)
	if err != nil {
		logger.Infof("GetBookByID: %v", err)
		// 渲染页面
		if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexTemplate); err != nil {
			// 渲染失败，返回 HTTP 500 错误。
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	// 定义模板主体内容。
	scrollPage := ScrollPage(&state.Global, book, paginationIndex)
	// 拼接页面
	indexTemplate = common.Html(
		c,
		&state.Global,
		scrollPage, // define body content
		[]string{"script/scroll.js"},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexTemplate); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
