package shelf

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/comigo_remote"
	"github.com/yumenaka/comigo/tools/logger"
)

// BookmarkResolver 在一次书架渲染中复用远端书签，避免每张卡片重复请求远端服务。
type BookmarkResolver struct {
	remoteMarks map[string]map[string]model.BookMarks
}

func NewBookmarkResolver(storeBookInfos []model.StoreBookInfo, childBookInfos []model.BookInfo) *BookmarkResolver {
	resolver := &BookmarkResolver{remoteMarks: map[string]map[string]model.BookMarks{}}
	remoteStoreKeys := map[string]bool{}
	for _, storeBooks := range storeBookInfos {
		for _, book := range storeBooks.BookInfos {
			if book.RemoteStoreKey != "" {
				remoteStoreKeys[book.RemoteStoreKey] = true
			}
		}
	}
	for _, book := range childBookInfos {
		if book.RemoteStoreKey != "" {
			remoteStoreKeys[book.RemoteStoreKey] = true
		}
	}
	for remoteStoreKey := range remoteStoreKeys {
		resolver.remoteMarks[remoteStoreKey] = loadRemoteComigoBookmarks(remoteStoreKey)
	}
	return resolver
}

// Get 返回书籍对应的本地或远端书签，远端书签只使用本次渲染已拉取的数据。
func (resolver *BookmarkResolver) Get(book model.BookInfo) model.BookMarks {
	if book.RemoteBookID != "" && book.RemoteStoreKey != "" {
		if resolver != nil && resolver.remoteMarks[book.RemoteStoreKey] != nil {
			return resolver.remoteMarks[book.RemoteStoreKey][book.RemoteBookID]
		}
		return nil
	}
	bookMarks, err := model.IStore.GetBookMarks(book.BookID)
	if err != nil || bookMarks == nil {
		return nil
	}
	return *bookMarks
}

// loadRemoteComigoBookmarks 实时拉取一个远端书库的全部书签，并按远端 BookID 归组。
func loadRemoteComigoBookmarks(remoteStoreKey string) map[string]model.BookMarks {
	bookMarksByRemoteID := map[string]model.BookMarks{}
	storeURL, err := comigo_remote.MatchStoreByKey(config.GetCfg().StoreUrls, remoteStoreKey)
	if err != nil {
		logger.Infof("%s", err)
		return bookMarksByRemoteID
	}
	client, err := comigo_remote.NewClient(storeURL, config.GetCfg().TimeoutLimitForScan)
	if err != nil {
		logger.Infof("%s", err)
		return bookMarksByRemoteID
	}
	remoteMarks, err := client.GetAllBookmarks()
	if err != nil {
		logger.Infof("%s", err)
		return bookMarksByRemoteID
	}
	for _, item := range remoteMarks {
		mark := item.BookMark
		bookMarksByRemoteID[item.BookInfo.BookID] = append(bookMarksByRemoteID[item.BookInfo.BookID], mark)
	}
	return bookMarksByRemoteID
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
			return common.RenderHTML(c, indexHtml)
		}
	}
	// 为首页定义模板布局。
	indexHtml := common.Html(
		c,
		ShelfPage(c, model.GetAllBooksNumber(), storeBookInfos, childBookInfos), // define body content
		[]string{"static/js/shelf.js"},
	)
	// 用模板渲染书架页面
	return common.RenderHTML(c, indexHtml)
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
			return c.Redirect(http.StatusFound, config.PrefixPath("/shelf/"+parentID))
		}
		return c.Redirect(http.StatusFound, config.PrefixPath("/"))
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
		[]string{"static/js/shelf.js"},
	)
	return common.RenderHTML(c, indexHtml)
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
	return common.RenderHTML(c, indexHtml)
}

// filterStoreBookInfosByKeyword 在顶层书架中按标题过滤。
func filterStoreBookInfosByKeyword(storeBookInfos []model.StoreBookInfo, keyword string) []model.StoreBookInfo {
	// 搜索关键词规范化：去除首尾空格并转换为小写，确保搜索不受大小写和多余空格影响。
	normalized := normalizeSearchKeyword(keyword)
	if normalized == "" {
		return storeBookInfos
	}
	// 过滤书架列表，保留至少有一个子书籍标题匹配关键词的书架。
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
			return config.PrefixPath("/shelf/" + parentID)
		}
		return config.PrefixPath("/")
	}
	return common.GetReturnUrl(c.Param("id"))
}

func generateReadURL(book model.BookInfo, lastReadPage int) string {
	remoteStoreArg := "undefined"
	if book.RemoteStoreKey != "" {
		remoteStoreArg = fmt.Sprintf("%q", book.RemoteStoreKey)
	}
	// 如果是书籍组，就跳转到子书架
	if book.Type == model.TypeBooksGroup {
		shelfURL := config.PrefixPath("/shelf/" + book.BookID)
		if book.RemoteStoreKey != "" {
			shelfURL += "?remote_store=" + url.QueryEscape(book.RemoteStoreKey)
		}
		return fmt.Sprintf("%q", shelfURL)
	}
	// 如果是视频、音频，跳转到播放器页面
	if book.Type == model.TypeVideo || book.Type == model.TypeAudio {
		playerURL := config.PrefixPath("/player/" + book.BookID)
		if book.RemoteStoreKey != "" {
			playerURL += "?remote_store=" + url.QueryEscape(book.RemoteStoreKey)
		}
		return fmt.Sprintf("%q", playerURL)
	}
	// 如果是 HTML、未知文件，就在新窗口打开原始文件
	if book.Type == model.TypeHTML || book.Type == model.TypeUnknownFile {
		rawURL := config.PrefixPath("/api/raw/" + book.BookID + "/" + url.QueryEscape(book.Title))
		if book.RemoteStoreKey != "" {
			rawURL += "?remote_store=" + url.QueryEscape(book.RemoteStoreKey)
		}
		return fmt.Sprintf("%q", rawURL)
	}
	// 其他情况，跳转到阅读页面，类似 /scroll/4cTOjFm?page=1
	return fmt.Sprintf("$store.global.getReadURL(\"%s\",%d,%s)", book.BookID, lastReadPage, remoteStoreArg)
}
