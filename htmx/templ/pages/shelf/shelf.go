package shelf

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/htmx/templ/common"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// Handler 书架页面的处理程序。
func Handler(c echo.Context) error {
	// Set the response content type to HTML.
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	// 书籍重排方式
	sortBy := "default"
	sortBookBy, err := c.Cookie("ShelfSortBy")
	if err == nil {
		// 如果 c.Cookie("ShelfSortBy") 返回错误，那么 sortBookBy 可能是空值（nil），
		// 没有正确处理这种情况就直接访问了 .Value 属性，会导致了空指针引用错误。
		sortBy = sortBookBy.Value
		fmt.Println("sortBookBy:", sortBy)
	} else {
		fmt.Println("未设置书籍排序方式，使用默认排序:", sortBy)
	}
	// 读取url参数，获取书籍ID
	bookID := c.Param("id")
	// 如果没有指定书籍ID，获取顶层书架信息。
	if bookID == "" {
		var err error
		model.CheckAllBookFileExist()
		state.Global.ShelfBookList, err = model.TopOfShelfInfo(sortBy)
		// TODO: 没有图书的提示（上传压缩包\远程下载示例漫画）
		if err != nil {
			logger.Infof("TopOfShelfInfo: %v", err)
		}
	}
	// 如果指定了书籍ID，获取子书架信息。
	if bookID != "" {
		var err error
		model.CheckAllBookFileExist()
		state.Global.ShelfBookList, err = model.GetBookInfoListByID(bookID, sortBy)
		// TODO: 没有图书的提示（返回主页\上传压缩包\远程下载示例漫画）
		if err != nil {
			logger.Infof("GetBookShelf Error: %v", err)
			return nil
		}
	}

	// 为首页定义模板布局。
	indexTemplate := common.MainLayout(
		c,
		&state.Global,
		ShelfPage(c, &state.Global), // define body content
		"static/shelf.js")

	// 用模板渲染书架页面(htmx-go)
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexTemplate); err != nil {
		// 如果出错，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

func getHref(book model.BookInfo) string {
	// 如果是书籍组，就跳转到子书架
	if book.Type == model.TypeBooksGroup {
		return "\"/shelf/" + book.BookID + "\""
	}
	// 如果是视频、音频、未知文件，就在新窗口打开
	if book.Type == model.TypeVideo || book.Type == model.TypeAudio || book.Type == model.TypeUnknownFile {
		return "\"/api/raw/" + book.BookID + "/" + url.QueryEscape(book.Title) + "\""
	}
	// 其他情况，跳转到阅读页面,
	return "'/'+$store.global.readMode+ '/' + BookID"
}

func getTarget(book model.BookInfo) string {
	if book.Type == model.TypeVideo || book.Type == model.TypeAudio || book.Type == model.TypeUnknownFile {
		return "_blank"
	}
	return "_self"
}
