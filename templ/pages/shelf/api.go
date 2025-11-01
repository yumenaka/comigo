package shelf

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/templ/state"
	"github.com/yumenaka/comigo/tools/logger"
)

// GetBookListHandler 返回排序完毕的book list
func GetBookListHandler(c echo.Context) error {
	// 书籍重排方式
	sortBy := "default"
	shelfBookByCookie, err := c.Cookie("ShelfSortBy")
	if err == nil {
		sortBy = shelfBookByCookie.Value
	} else {
		// TODO: 加密链接的时候，设置Secure为true
		// Secure 表示：Cookie 必须使用类似 HTTPS 的加密环境下才能读取
		// HttpOnly 表示：不能通过非HTTP方式来访问，拒绝 JavaScript 访问 Cookie！(例如引用 document.cookie）
		// SameSite 表示：所有和 Cookie 來源不同的請求都不會帶上 Cookie
		cookie := new(http.Cookie)
		cookie.Name = "SortBookBy"
		cookie.Value = "default"
		cookie.MaxAge = 3600000
		cookie.Path = "/"
		cookie.Domain = c.Request().Host
		cookie.Secure = false
		cookie.HttpOnly = false
		c.SetCookie(cookie)
	}

	// 读取url参数，获取书籍ID
	bookID := c.Param("id")
	// 如果没有指定书籍ID，获取顶层书架信息。
	if bookID == "" {
		var err error
		state.NowBookInfos, err = store.TopOfShelfInfo(sortBy)
		if err != nil {
			logger.Infof("TopOfShelfInfo: %v", err)
		}
	}
	// 如果指定了书籍ID，获取子书架信息。
	if bookID != "" {
		var err error
		state.NowBookInfos, err = store.GetChildBooksInfo(bookID)
		if err != nil {
			logger.Infof("GetBookShelf: %v", err)
		}
		state.NowBookInfos.SortBooks(sortBy)
	}

	if state.NowBookInfos == nil {
		state.NowBookInfos = &model.BookInfos{}
	}

	// https://github.com/angelofallars/htmx-go#templ-integration
	// 主体内容的模板(书籍列表)
	template := MainArea(c) // define body content

	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
