package pages

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/entity"
	"github.com/yumenaka/comi/htmx/state"
	"github.com/yumenaka/comi/htmx/templates/components"
	"github.com/yumenaka/comi/util/logger"
)

// ScrollHandler 阅读界面（先做卷轴模式）
func ScrollHandler(c *gin.Context) {
	// 书籍排列的方式，默认name
	//sortBy := c.DefaultQuery("sort_by", "default")
	// 如果传了maxDepth这个参数
	var err error
	state.Global.BooksList, err = entity.TopOfShelfInfo("name")
	if err != nil {
		logger.Infof("TopOfShelfInfo: %v", err)
	}

	// 网页meta标签。
	metaTags := components.MetaTags(
		"Comigo  Comic Manga Reader 在线漫画 阅读器",         // define meta keywords
		"Simple Manga Reader in Linux，Windows，Mac OS", // define meta description
	)

	// 定义模板主体内容。
	scrollPage := ScrollPage(&state.Global)

	// 为首页定义模板布局。
	indexTemplate := components.MainLayout(
		"Comigo "+state.Global.Version, // define title text
		metaTags,                       // define meta tags
		scrollPage,                     // define body content
	)

	// 渲染索引页模板。
	if err := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, indexTemplate); err != nil {
		// 如果不是，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
