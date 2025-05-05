package shelf

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/templ/state"
	"github.com/yumenaka/comigo/util/logger"
)

// Handler 书架页面的处理程序。
func Handler(c echo.Context) error {
	// Set the response content type to HTML.
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	// 书籍排序方式
	sortBy := "default"
	sortBookBy, err := c.Cookie("ShelfSortBy")
	if err == nil {
		// c.Cookie("key") 没找到，那么就会取到空值（nil），
		// 没处判断就直接访问 .Value 属性，会导致空指针引用错误。
		sortBy = sortBookBy.Value
	}
	// 读取url参数，获取书籍ID
	bookID := c.Param("id")
	// 如果没有指定书籍ID，获取顶层书架信息。
	if bookID == "" {
		model.CheckAllBookFileExist()
		state.Global.ShelfBookList, _ = model.TopOfShelfInfo(sortBy)
	}

	// 如果指定了书籍ID，获取子书架信息。
	if bookID != "" {
		var err error
		model.CheckAllBookFileExist()
		state.Global.ShelfBookList, err = model.GetBookInfoListByID(bookID, sortBy)
		// TODO: 无图书的提示（返回主页\上传压缩包\远程下载示例漫画）
		if err != nil {
			logger.Infof("GetBookShelf Error: %v", err)
			// 渲染 404 页面
			indexHtml := common.Html(
				c,
				&state.Global,
				error_page.NotFound404(&state.Global),
				[]string{},
			)
			if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
				// 如果出错，返回 HTTP 500 错误。
				return c.NoContent(http.StatusInternalServerError)
			}
			return nil
		}
	}
	// 为首页定义模板布局。
	indexHtml := common.Html(
		c,
		&state.Global,
		ShelfPage(c, &state.Global), // define body content
		[]string{"script/shelf.js"},
	)
	// 用模板渲染书架页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 如果出错，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

func getReadURL(book model.BookInfo) string {
	// 如果是书籍组，就跳转到子书架
	if book.Type == model.TypeBooksGroup {
		return "\"/shelf/" + book.BookID + "\""
	}
	// 如果是视频、音频、未知文件，就在新窗口打开
	if book.Type == model.TypeVideo || book.Type == model.TypeAudio || book.Type == model.TypeUnknownFile {
		return fmt.Sprintf("/api/raw/%s/%s", book.BookID, url.QueryEscape(book.Title))
	}
	// 其他情况，跳转到阅读页面，类似 /scroll/4cTOjFm?page=1
	readURL := "'/'+$store.global.readMode+ '/' + BookID + ($store.global.readMode === 'scroll'?($store.scroll.fixedPagination?'?page=1':''):'')"
	return readURL
}

func getTarget(book model.BookInfo) string {
	if book.Type == model.TypeVideo || book.Type == model.TypeAudio || book.Type == model.TypeUnknownFile {
		return "_blank"
	}
	return "_self"
}
