package flip

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/util/logger"
)

// PageHandler 阅读界面（翻页模式）
func PageHandler(c echo.Context) error {
	model.CheckAllBookFileExist()
	bookID := c.Param("id")
	logger.Info("Flip Mode Book ID:" + bookID)
	// 图片排序方式
	sortBy := "default"
	// c.Cookie("key") 没找到，那么就会取到空值（nil），没判断nil就直接访问 .Value 属性，会导致空指针引用错误。
	sortBookBy, err := c.Cookie("FlipSortBy")
	if err == nil {
		sortBy = sortBookBy.Value
	}
	// // 给cookie设置默认值
	// if err != nil {
	// 	sortBookBy.Value = "default"
	// 	cookie := new(http.Cookie)
	// 	cookie.Name = "FlipSortBy"
	// 	cookie.Value = sortBookBy.Value
	// 	cookie.MaxAge = 3600000
	// 	cookie.Path = "/"
	// 	cookie.Domain = domain
	// 	cookie.Secure = false
	// 	cookie.HttpOnly = true
	// 	c.SetCookie(cookie)
	// }
	// 读取url参数，获取书籍ID
	book, err := model.GetBookByID(bookID, sortBy)
	if err != nil {
		logger.Infof("GetBookByID: %v", err)
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

	// // TODO：加密链接的时候，设置Secure为true
	// readingProgressStr.Value = `{"nowPageNum":0,"nowChapterNum":0,"readingTime":0}`
	// cookie := new(http.Cookie)
	// cookie.Name = "bookID:" + bookID
	// cookie.Value = readingProgressStr.Value
	// cookie.MaxAge = 60 * 60 * 24 * 356
	// cookie.Path = "/"
	// cookie.Domain = domain
	// cookie.Secure = false
	// cookie.HttpOnly = false
	// c.SetCookie(cookie)

	// // TODO：当前书籍的阅读进度，存储在cookie里面，与服务器共享与交互 readingProgress
	// readingProgressStr, err := c.Cookie("bookID:" + bookID)
	// // 获取纯域名部分，不带端口号 ////Cookie.Domain 的规范：根据 RFC 6265，Cookie.Domain 不应该包含端口号。它只能包含域名或 IP 地址
	// domain := c.Request().Host
	// if idx := strings.IndexByte(domain, ':'); idx != -1 {
	//	domain = domain[:idx] // 去掉端口号
	// }
	//
	// readingProgress, err := model.GetReadingProgress(readingProgressStr.Value)
	// if err != nil {
	//	logger.Infof("GetReadingProgress: %v readingProgressStr: %s", err, readingProgressStr.Value)
	// }
	//
	// state.NowBookList, err = model.TopOfShelfInfo("name")
	// if err != nil {
	//	logger.Infof("TopOfShelfInfo: %v", err)
	// }

	// 翻页模式页面
	indexHtml := common.Html(
		c,

		FlipPage(c, book),
		[]string{"script/flip.js", "script/flip_sketch.js"})
	// 渲染正常页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 如果渲染失败，返回 HTTP 500 错误
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
