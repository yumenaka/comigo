package shelf

import (
	"net/http"
	"net/url"

	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/htmx/templates/common"
	"github.com/yumenaka/comigo/util/logger"
)

// Handler 书架页面的处理程序。
func Handler(c *gin.Context) {
	// 获取查询参数并指定默认值 ?age=value
	sortBy := c.DefaultQuery("sortBy", "default")

	// 读取url参数，获取书籍ID
	bookID := c.Param("id")
	// 如果没有指定书籍ID，获取顶层书架信息。
	if bookID == "" {
		var err error
		state.Global.TopBooks, err = entity.TopOfShelfInfo(sortBy)
		if err != nil {
			logger.Infof("TopOfShelfInfo: %v", err)
			//TODO: 处理没有图书的情况（上传压缩包或远程下载示例漫画）
		}
	}
	// 如果指定了书籍ID，获取子书架信息。
	if bookID != "" {
		var err error
		state.Global.TopBooks, err = entity.GetBookInfoListByID(bookID, sortBy)
		if err != nil {
			logger.Infof("GetBookShelf: %v", err)
		}
	}

	// 为首页定义模板布局。
	indexTemplate := common.MainLayout(
		c,
		ShelfPage(c, &state.Global), // define body content
	)

	// 用模板渲染书架页面(htmx-go)
	if err := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, indexTemplate); err != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//// 书籍重排方式，存储在 cookie 里面的方案，默认为“default”
	//sortBy, err := c.Cookie("SortBookBy")
	//if err != nil {
	//	sortBookBy = "default"
	//	// 加密链接的时候，应该设置Secure为true
	//	//Secure 表示：Cookie 必须使用类似 HTTPS 的加密环境下才能读取
	//	//HttpOnly 表示：不能通过非HTTP方式来访问，拒绝 JavaScript 访问 Cookie！(例如引用 document.cookie）
	//	//SameSite 表示：所有和 Cookie 來源不同的請求都不會帶上 Cookie
	//	c.SetCookie("SortBookBy", sortBy, 60*60*24*356, "/", c.Request.Host, false, false)
	//}
}

func getHref(book entity.BookInfo) string {
	// 如果是书籍组，就跳转到子书架
	if book.Type == entity.TypeBooksGroup {
		return "\"/shelf/" + book.BookID + "/\""
	}
	// 如果是视频、音频、未知文件，就在新窗口打开
	if book.Type == entity.TypeVideo || book.Type == entity.TypeAudio || book.Type == entity.TypeUnknownFile {
		return "\"/api/raw/" + book.BookID + "/" + url.QueryEscape(book.Title) + "\""
	}
	// 其他情况，跳转到阅读页面,
	return "'/'+$store.global.readMode+ '/' + BookID"
}

func getTarget(book entity.BookInfo) string {
	if book.Type == entity.TypeVideo || book.Type == entity.TypeAudio || book.Type == entity.TypeUnknownFile {
		return "_blank"
	}
	return "_self"
}
