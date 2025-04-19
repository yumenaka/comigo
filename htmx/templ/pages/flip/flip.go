package flip

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/htmx/templ/common"
	"github.com/yumenaka/comigo/htmx/templ/pages/error_page"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// Handler 阅读界面（翻页模式）
func Handler(c echo.Context) error {
	model.CheckAllBookFileExist()
	bookID := c.Param("id")
	book, err := model.GetBookByID(bookID, "default")
	if err != nil {
		logger.Infof("GetBookByID: %v", err)
	}
	// 没有找到书，显示 HTTP 404 错误
	indexTemplate := common.Html(
		c,
		&state.Global,
		error_page.NotFound404(&state.Global),
		[]string{},
	)

	//// TODO：加密链接的时候，设置Secure为true
	//readingProgressStr.Value = `{"nowPageNum":0,"nowChapterNum":0,"readingTime":0}`
	//cookie := new(http.Cookie)
	//cookie.Name = "bookID:" + bookID
	//cookie.Value = readingProgressStr.Value
	//cookie.MaxAge = 60 * 60 * 24 * 356
	//cookie.Path = "/"
	//cookie.Domain = domain
	//cookie.Secure = false
	//cookie.HttpOnly = false
	//c.SetCookie(cookie)

	if err == nil {
		//// TODO：当前书籍的阅读进度，存储在cookie里面，与服务器共享与交互 readingProgress
		//readingProgressStr, err := c.Cookie("bookID:" + bookID)
		//// 获取纯域名部分，不带端口号 ////Cookie.Domain 的规范：根据 RFC 6265，Cookie.Domain 不应该包含端口号。它只能包含域名或 IP 地址
		//domain := c.Request().Host
		//if idx := strings.IndexByte(domain, ':'); idx != -1 {
		//	domain = domain[:idx] // 去掉端口号
		//}
		//
		//readingProgress, err := model.GetReadingProgress(readingProgressStr.Value)
		//if err != nil {
		//	logger.Infof("GetReadingProgress: %v readingProgressStr: %s", err, readingProgressStr.Value)
		//}
		//
		//state.Global.ShelfBookList, err = model.TopOfShelfInfo("name")
		//if err != nil {
		//	logger.Infof("TopOfShelfInfo: %v", err)
		//}
		//// 图片重排方式
		//sortPageBy, err := c.Cookie("SortPageBy")
		//if err != nil {
		//	sortPageBy.Value = "default"
		//	cookie := new(http.Cookie)
		//	cookie.Name = "SortPageBy"
		//	cookie.Value = sortPageBy.Value
		//	cookie.MaxAge = 3600000
		//	cookie.Path = "/"
		//	cookie.Domain = domain
		//	cookie.Secure = false
		//	cookie.HttpOnly = true
		//	c.SetCookie(cookie)
		//}

		// 翻页模式页面主体
		FlipPage := FlipPage(&state.Global, book)
		// 拼接页面
		indexTemplate = common.Html(
			c,
			&state.Global,
			FlipPage, // define body content
			[]string{"static/flip.js"})
	}

	// 渲染404或者正常页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexTemplate); err != nil {
		// 如果渲染失败，返回 HTTP 500 错误
		return c.NoContent(http.StatusInternalServerError)
	}

	return nil
}
