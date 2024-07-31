package templates

import (
	"errors"
	"log"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/entity"
	"github.com/yumenaka/comi/htmx/state"
	"github.com/yumenaka/comi/htmx/templates/pages"
	"github.com/yumenaka/comi/util/logger"
)

// MainHandler handles a view for the index page.
func MainHandler(c *gin.Context) {
	// 书籍排列的方式，默认name
	//sortBy := c.DefaultQuery("sort_by", "default")
	// 如果传了maxDepth这个参数
	var err error
	state.Global.BooksList, err = entity.TopOfShelfInfo("name")
	if err != nil {
		logger.Infof("TopOfShelfInfo: %v", err)
	}

	// 定义模板元标签。TODO：用书籍的元标签替换。
	metaTags := pages.MetaTags(
		"Comigo  Comic Manga Reader 在线漫画 阅读器",         // define meta keywords
		"Simple Manga Reader in Linux，Windows，Mac OS", // define meta description
	)

	// 定义模板主体内容。
	scrollPage := pages.ScrollPage(&state.Global)

	// 为首页定义模板布局。
	indexTemplate := MainLayout(
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

// ShowContentAPIHandler 处理一个用于显示内容的 API
func ShowContentAPIHandler(c *gin.Context) {
	// 检查当前请求是否有 'HX-Request' 头部。
	// 更多信息请见 https://htmx.org/docs/#request-headers”
	if !htmx.IsHTMX(c.Request) {
		// If not, return HTTP 400 error.
		err := c.AbortWithError(http.StatusBadRequest, errors.New("non-htmx request"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 编写 HTML内容。
	_, err := c.Writer.Write([]byte("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>"))
	if err != nil {
		log.Println(err)
	}

	// 发送 htmx 响应。
	err = htmx.NewResponse().Write(c.Writer)
	if err != nil {
		log.Println(err)
	}
}
