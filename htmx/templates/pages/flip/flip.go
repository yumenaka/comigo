package flip

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/htmx/templates/common"
	"github.com/yumenaka/comigo/util/logger"
)

// Handler 阅读界面（TODO：翻页模式）
func Handler(c *gin.Context) {
	bookID := c.Param("id")
	book, err := entity.GetBookByID(bookID, "default")
	//TODO: 没有图书的提示（返回主页\上传压缩包\远程下载示例漫画）
	if err != nil {
		logger.Infof("GetBookByID: %v", err)
		// 显示 HTTP 404 错误信息，文本“404 not found”
		c.String(http.StatusNotFound, "404 not found")
		return
	}
	// 当前书籍的阅读进度，存储在cookie里面，与服务器共享与交互 readingProgress
	readingProgressStr, err := c.Cookie("bookID:" + bookID)
	if err != nil {
		readingProgressStr = "{nowPageNum:0,nowChapterNum:0,readingTime:0}"
		// TODO：加密链接的时候，应该设置Secure为true
		//Secure 表示：Cookie 必须使用类似 HTTPS 的加密环境下才能读取
		//HttpOnly 表示：不能通过非HTTP方式来访问，拒绝 JavaScript 访问 Cookie！(例如引用 document.cookie）
		//SameSite 表示：所有和 Cookie 來源不同的請求都不會帶上 Cookie
		c.SetCookie("bookID:"+bookID, readingProgressStr, 60*60*24*356, "/", c.Request.Host, false, false)
	}
	readingProgress, err := entity.GetReadingProgress(readingProgressStr)
	if err != nil {
		logger.Infof("GetReadingProgress: %v", err)
	}

	state.Global.TopBooks, err = entity.TopOfShelfInfo("name")
	if err != nil {
		logger.Infof("TopOfShelfInfo: %v", err)
	}
	// 图片重排的方式，默认name
	sortPageBy, err := c.Cookie("SortPageBy")
	if err != nil {
		sortPageBy = "default"
		//Secure 表示：不讓 Cookie 在 HTTP 之外的環境下被存取
		//HttpOnly 表示：拒絕與 JavaScript 共享 Cookie！
		//SameSite 表示：所有和 Cookie 來源不同的請求都不會帶上 Cookie
		c.SetCookie("SortPageBy", sortPageBy, 3600000, "/", c.Request.Host, false, true)
	}

	// 定义模板主体内容。
	FlipPage := FlipPage(c, &state.Global, book, &readingProgress)
	// 为首页定义模板布局。
	indexTemplate := common.MainLayout(
		c,
		&state.Global,
		FlipPage, // define body content
		"static/flip.js")

	// 渲染索引页模板。
	if err := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, indexTemplate); err != nil {
		// 如果不是，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
