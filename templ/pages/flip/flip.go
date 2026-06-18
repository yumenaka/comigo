package flip

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/comigo_remote"
	"github.com/yumenaka/comigo/tools/logger"
)

// FlipModeHandler 阅读界面（翻页阅读）
func FlipModeHandler(c echo.Context) error {
	bookID := c.Param("id")
	logger.Infof(locale.GetString("log_flip_mode_book_id"), bookID)
	// 图片排序方式
	sortBy := ""
	if assets.IsWailsWebViewRequest(c.Request()) {
		sortBy = c.QueryParam("sort_by")
	}
	if sortBy == "" {
		sortBy = "default"
		// c.Cookie("key") 没找到，那么就会取到空值（nil），没判断nil就直接访问 .Value 属性，会导致空指针引用错误。
		sortBookBy, err := c.Cookie("FlipSortBy")
		if err == nil {
			sortBy = sortBookBy.Value
		}
	}
	// 读取url参数，获取书籍ID
	var book *model.Book
	if remoteBook, ok, remoteErr := comigo_remote.ResolveBook(config.GetCfg().StoreUrls, c.QueryParam(comigo_remote.RemoteStoreQuery), bookID, sortBy, config.GetCfg().TimeoutLimitForScan); ok {
		if remoteErr != nil {
			logger.Infof(locale.GetString("log_get_book_error"), remoteErr)
			return renderFlipNotFound(c)
		}
		book = remoteBook
	} else {
		var err error
		book, err = model.IStore.GetBook(bookID)
		if err != nil {
			logger.Infof(locale.GetString("log_get_book_error"), err)
			return renderFlipNotFound(c)
		}
	}
	if book == nil {
		return renderFlipNotFound(c)
	}
	// HTML 单文件书籍直接返回源文件，避免翻页模板包裹后破坏原页面结构和脚本。
	if book.Type == model.TypeHTML {
		return c.Redirect(http.StatusTemporaryRedirect, common.RawBookURL(book))
	}
	book = book.CloneForView()
	book.SortPages(sortBy)

	// 翻页阅读（先加载共享 WebSocket 模块，再加载页面逻辑）
	indexHtml := common.Html(
		c,
		FlipPage(c, book),
		[]string{
			"static/js/ws_sync.js",
			"static/js/flip_modules/pagination_utils.js",
			"static/js/flip_modules/interaction_utils.js",
			"static/js/flip.js",
		})
	// 静态模式
	staticMode := c.QueryParam("static") != ""
	if staticMode {
		// 静态模式会把图片 URL 改成 base64，继续使用展示副本避免改到同一次渲染外的书籍数据。
		staticBook := book.CloneForView()
		// 静态模式下，使用 base64 图片数据
		for i, image := range staticBook.PageInfos {
			staticBook.PageInfos[i].Url = common.GetFileBase64Text(book.BookInfo.BookID, image.Name)
		}
		indexHtml = common.Html(
			c,
			FlipPage(c, staticBook),
			[]string{
				"static/js/ws_sync.js",
				"static/js/flip_modules/pagination_utils.js",
				"static/js/flip_modules/interaction_utils.js",
				"static/js/flip.js",
			})
	}
	// 渲染翻页阅读页面
	return common.RenderHTML(c, indexHtml)
}

func renderFlipNotFound(c echo.Context) error {
	indexHtml := common.Html(
		c,
		error_page.NotFound404(c),
		[]string{},
	)
	return common.RenderHTML(c, indexHtml)
}
