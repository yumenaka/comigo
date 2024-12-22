package scroll

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/htmx/templates/common"
	"github.com/yumenaka/comigo/htmx/templates/pages/error_page"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// Handler 阅读界面（卷轴模式）
func Handler(c *gin.Context) {
	model.CheckAllBookFileExist()
	//// 图片重排设定，存储在 cookie 里面，默认为“default”
	//sortPageBy, err := c.Cookie("SortPageBy")
	//if err != nil {
	//	sortPageBy = "default"
	//	//Secure 表示：不讓 Cookie 在 HTTP 之外的環境下被存取
	//	//HttpOnly 表示：拒絕與 JavaScript 共享 Cookie！
	//	//SameSite 表示：所有和 Cookie 來源不同的請求都不會帶上 Cookie
	//	c.SetCookie("SortPageBy", sortPageBy, 60*60*24*356, "/", c.Request.Host, false, true)
	//}
	// 获取查询参数并指定默认值 ?age=value
	sortBy := c.DefaultQuery("sortBy", "default")
	// 读取url参数，获取书籍ID
	bookID := c.Param("id")
	//没有找到书籍，显示 HTTP 404 错误
	indexTemplate := common.MainLayout(
		c,
		&state.Global,
		error_page.NotFound404(c, &state.Global),
		"",
	)
	book, err := model.GetBookByID(bookID, sortBy)
	if err != nil {
		c.Status(http.StatusNotFound)
		logger.Infof("GetBookByID: %v", err)
	}
	if err == nil {
		// 定义模板主体内容。
		scrollPage := ScrollPage(c, &state.Global, book)
		// 拼接页面
		indexTemplate = common.MainLayout(
			c,
			&state.Global,
			scrollPage, // define body content
			"static/scroll.js",
		)
	}
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, indexTemplate); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
