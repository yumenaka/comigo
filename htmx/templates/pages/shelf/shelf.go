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
		//TODO: 没有图书的提示（返回主页\上传压缩包\远程下载示例漫画）
		if err != nil {
			logger.Infof("TopOfShelfInfo: %v", err)
			// 显示 HTTP 404 错误信息，文本“404 not found”
			c.String(http.StatusNotFound, "404 not found")
		}
	}
	// 如果指定了书籍ID，获取子书架信息。
	if bookID != "" {
		var err error
		state.Global.TopBooks, err = entity.GetBookInfoListByID(bookID, sortBy)
		//TODO: 没有图书的提示（返回主页\上传压缩包\远程下载示例漫画）
		if err != nil {
			logger.Infof("GetBookShelf: %v", err)
			// 显示 HTTP 404 错误信息，文本“404 not found”
			c.String(http.StatusNotFound, "404 not found")
			return
		}
	}

	// 为首页定义模板布局。
	indexTemplate := common.MainLayout(
		c,
		&state.Global,
		ShelfPage(c, &state.Global), // define body content
		"static/shelf.js")

	// 用模板渲染书架页面(htmx-go)
	if err := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, indexTemplate); err != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
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
