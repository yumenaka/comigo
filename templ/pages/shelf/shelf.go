package shelf

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

// SearchHandler 书架搜索页面处理程序。
// 支持两种范围：
// 1. /search?keyword=xxx               -> 全部顶层书架搜索
// 2. /search?keyword=xxx&parent=bookID -> 指定子书架搜索
func SearchHandler(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	keyword := strings.TrimSpace(c.QueryParam("keyword"))
	parentID := strings.TrimSpace(c.QueryParam("parent"))
	if keyword == "" {
		// 空关键词时回到对应书架，避免渲染无意义的搜索页。
		if parentID != "" {
			return c.Redirect(http.StatusFound, "/shelf/"+parentID)
		}
		return c.Redirect(http.StatusFound, "/")
	}

	sortBy := getShelfSortBy(c)

	var (
		storeBookInfos []model.StoreBookInfo
		childBookInfos model.BookInfos
		nowBookNum     int
		err            error
	)

	// parent 为空：在顶层书架中搜索
	if parentID == "" {
		storeBookInfos, err = store.TopOfShelfInfo(sortBy)
		if err != nil {
			logger.Errorf(locale.GetString("err_getbookshelf_error"), err)
		}
		storeBookInfos = filterStoreBookInfosByKeyword(storeBookInfos, keyword)
		nowBookNum = countStoreBookInfos(storeBookInfos)
	}

	// parent 不为空：仅在子书架中搜索
	if parentID != "" {
		logger.Infof(locale.GetString("log_get_child_books_for_bookid"), parentID)
		childBooks, childErr := store.GetChildBooksInfo(parentID)
		if childErr != nil {
			logger.Infof(locale.GetString("log_get_bookshelf_error"), childErr)
			return renderShelfNotFound(c)
		}

		childBookInfos = filterBookInfosByKeyword(*childBooks, keyword)
		childBookInfos.SortBooks(sortBy)
		nowBookNum = len(childBookInfos)
	}

	indexHtml := common.Html(
		c,
		ShelfPage(c, nowBookNum, storeBookInfos, childBookInfos),
		[]string{"script/shelf.js"},
	)
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

// getShelfSortBy 获取书架排序方式（与原书架页面保持一致）。
func getShelfSortBy(c echo.Context) string {
	sortBy := "default"
	sortBookBy, err := c.Cookie("ShelfSortBy")
	if err == nil {
		sortBy = sortBookBy.Value
	}
	return sortBy
}

// renderShelfNotFound 统一渲染 404 页面，复用 ShelfHandler 的行为。
func renderShelfNotFound(c echo.Context) error {
	indexHtml := common.Html(
		c,
		error_page.NotFound404(c),
		[]string{},
	)
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

// filterStoreBookInfosByKeyword 在顶层书架中按标题过滤。
func filterStoreBookInfosByKeyword(storeBookInfos []model.StoreBookInfo, keyword string) []model.StoreBookInfo {
	normalized := normalizeSearchKeyword(keyword)
	if normalized == "" {
		return storeBookInfos
	}

	filteredStoreBookInfos := make([]model.StoreBookInfo, 0, len(storeBookInfos))
	for _, shelf := range storeBookInfos {
		filteredBooks := filterBookInfosByKeyword(shelf.BookInfos, normalized)
		if len(filteredBooks) == 0 {
			continue
		}
		shelf.BookInfos = filteredBooks
		// 搜索结果中显示命中条目数量，而不是原始总数。
		shelf.ChildBookNum = len(filteredBooks)
		filteredStoreBookInfos = append(filteredStoreBookInfos, shelf)
	}
	return filteredStoreBookInfos
}

// filterBookInfosByKeyword 在书籍列表中按标题过滤。
func filterBookInfosByKeyword(bookInfos model.BookInfos, keyword string) model.BookInfos {
	normalized := normalizeSearchKeyword(keyword)
	if normalized == "" {
		return bookInfos
	}

	filtered := make(model.BookInfos, 0, len(bookInfos))
	for _, book := range bookInfos {
		if strings.Contains(strings.ToLower(book.Title), normalized) {
			filtered = append(filtered, book)
		}
	}
	return filtered
}

func normalizeSearchKeyword(keyword string) string {
	return strings.ToLower(strings.TrimSpace(keyword))
}

func countStoreBookInfos(storeBookInfos []model.StoreBookInfo) int {
	total := 0
	for _, shelf := range storeBookInfos {
		total += len(shelf.BookInfos)
	}
	return total
}

// getShelfParentID 获取当前书架上下文中的父级书籍 ID。
// 优先使用路径参数 /shelf/:id，其次读取 /search?parent=xxx。
func getShelfParentID(c echo.Context) string {
	pathParentID := strings.TrimSpace(c.Param("id"))
	if pathParentID != "" {
		return pathParentID
	}
	return strings.TrimSpace(c.QueryParam("parent"))
}

func isShelfSearchPage(c echo.Context) bool {
	return c.Path() == "/search"
}

// getShelfHeaderTitle 根据当前页面上下文生成标题文本。
func getShelfHeaderTitle(c echo.Context, nowBookNum int, storeBookInfos []model.StoreBookInfo, childBookInfos []model.BookInfo) string {
	if !isShelfSearchPage(c) {
		return common.GetPageTitle(c.Param("id"), nowBookNum, storeBookInfos, childBookInfos)
	}
	keyword := strings.TrimSpace(c.QueryParam("keyword"))
	if keyword == "" {
		return fmt.Sprintf(locale.GetString("search_result_title"), nowBookNum)
	}
	return fmt.Sprintf(locale.GetString("search_result_title_with_keyword"), keyword, nowBookNum)
}

func shouldShowShelfReturnIcon(c echo.Context) bool {
	return getShelfParentID(c) != "" || isShelfSearchPage(c)
}

func getShelfReturnURL(c echo.Context) string {
	if isShelfSearchPage(c) {
		parentID := getShelfParentID(c)
		if parentID != "" {
			return "/shelf/" + parentID
		}
		return "/"
	}
	return common.GetReturnUrl(c.Param("id"))
}

func generateReadURL(book model.BookInfo, lastReadPage int) string {
	// 如果是书籍组，就跳转到子书架
	if book.Type == model.TypeBooksGroup {
		return fmt.Sprintf("\"/shelf/%s\"", book.BookID)
	}
	// 如果是视频、音频，跳转到播放器页面
	if book.Type == model.TypeVideo || book.Type == model.TypeAudio {
		return fmt.Sprintf("\"/player/%s\"", book.BookID)
	}
	// 如果是 HTML、未知文件，就在新窗口打开原始文件
	if book.Type == model.TypeHTML || book.Type == model.TypeUnknownFile {
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
