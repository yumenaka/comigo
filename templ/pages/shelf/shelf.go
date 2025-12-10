package shelf

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/logger"
)

func GetBookmarks(bookID string) model.BookMarks {
	bookMarks, _ := model.IStore.GetBookMarks(bookID)
	return *bookMarks
}

// ShelfHandler 书架页面的处理程序。
func ShelfHandler(c echo.Context) error {
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
	var storeBookInfos []model.StoreBookInfo
	if bookID == "" {
		storeBookInfos, err = store.TopOfShelfInfo(sortBy)
		if err != nil {
			logger.Errorf(locale.GetString("err_getbookshelf_error"), err)
		}
	}
	var childBookInfos model.BookInfos
	// 如果指定了书籍ID，获取子书架信息。
	if bookID != "" {
		model.ClearBookNotExist()
		logger.Infof(locale.GetString("log_get_child_books_for_bookid"), bookID)
		childBooks, err := store.GetChildBooksInfo(bookID)
		if err == nil {
			logger.Infof(locale.GetString("log_get_child_books_count"), bookID, len(*childBooks))
			childBookInfos = *childBooks
		}
		childBookInfos.SortBooks(sortBy)
		// 无图书的提示（返回主页\上传压缩包\远程下载示例漫画）
		if err != nil {
			logger.Infof(locale.GetString("log_get_bookshelf_error"), err)
			// 渲染 404 页面
			indexHtml := common.Html(
				c,
				error_page.NotFound404(c),
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
		ShelfPage(c, model.GetAllBooksNumber(), storeBookInfos, childBookInfos), // define body content
		[]string{"script/shelf.js"},
	)
	// 用模板渲染书架页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 如果出错，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

func generateReadURL(book model.BookInfo, lastReadPage int) string {
	// 如果是书籍组，就跳转到子书架
	if book.Type == model.TypeBooksGroup {
		return fmt.Sprintf("\"/shelf/%s\"", book.BookID)
	}
	// 如果是视频、音频、未知文件，就在新窗口打开
	if book.Type == model.TypeVideo || book.Type == model.TypeAudio || book.Type == model.TypeHTML || book.Type == model.TypeUnknownFile {
		return fmt.Sprintf("\"/api/raw/%s/%s\"", book.BookID, url.QueryEscape(book.Title))
	}
	// 其他情况，跳转到阅读页面，类似 /scroll/4cTOjFm?page=1
	return fmt.Sprintf("$store.global.getReadURL(\"%s\",%d)", book.BookID, lastReadPage)
}

// getBookCardTarget 根据书籍类型，返回书籍卡片的 target 属性值。似乎有点问题，暂时未使用。
func getBookCardTarget(book model.BookInfo) string {
	// 新页面打开
	if book.Type == model.TypeVideo || book.Type == model.TypeAudio || book.Type == model.TypeUnknownFile {
		return `_blank`
	}
	// 按照用户设置决定
	if book.Type == model.TypeZip || book.Type == model.TypeCbz || book.Type == model.TypeEpub || book.Type == model.TypeRar || book.Type == model.TypeCbr || book.Type == model.TypeTar {
		return `$store.shelf.openInNewTab ? '_blank' : '_self'`
	}
	// 当前页面打开
	return `_self`
}
