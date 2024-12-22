package shelf

import (
	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
	"net/http"
)

// GetBookListHandler 返回排序完毕的book list
func GetBookListHandler(c *gin.Context) {
	//检查请求来源是不是htmx   https://github.com/angelofallars/htmx-go#htmx-requests
	if htmx.IsHTMX(c.Request) {
		// logic for handling HTMX requests
	} else {
		// logic for handling non-HTMX requests (e.g. render a full page for first-time visitors)
	}
	// 书籍重排的方式，默认文件名
	sortBookBy, err := c.Cookie("SortBookBy")
	if err != nil {
		sortBookBy = "default"
		// TODO: 加密链接的时候，设置Secure为true
		//Secure 表示：Cookie 必须使用类似 HTTPS 的加密环境下才能读取
		//HttpOnly 表示：不能通过非HTTP方式来访问，拒绝 JavaScript 访问 Cookie！(例如引用 document.cookie）
		//SameSite 表示：所有和 Cookie 來源不同的請求都不會帶上 Cookie
		c.SetCookie("SortBookBy", sortBookBy, 3600000, "/", c.Request.Host, false, false)
	}

	// 读取url参数，获取书籍ID
	bookID := c.Param("id")
	// 如果没有指定书籍ID，获取顶层书架信息。
	if bookID == "" {
		var err error
		state.Global.ShelfBookList, err = model.TopOfShelfInfo(sortBookBy)
		if err != nil {
			logger.Infof("TopOfShelfInfo: %v", err)
			//TODO: 处理没有图书的情况（上传压缩包或远程下载示例漫画）
		}
	}
	// 如果指定了书籍ID，获取子书架信息。
	if bookID != "" {
		var err error
		state.Global.ShelfBookList, err = model.GetBookInfoListByID(bookID, sortBookBy)
		if err != nil {
			logger.Infof("GetBookShelf: %v", err)
		}
	}

	// https://github.com/angelofallars/htmx-go#templ-integration
	// 主体内容的模板(书籍列表)
	template := MainArea(c, &state.Global) // define body content

	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
