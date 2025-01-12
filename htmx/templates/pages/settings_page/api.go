package settings_page

import (
	"errors"
	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/htmx/state"
	"net/http"
)

// 从模板中获取htmx请求，页面比较复杂的时候用
func Tab1(c *gin.Context) {
	//检查请求来源是不是htmx
	template := tab1(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func Tab2(c *gin.Context) {
	//检查请求来源是不是htmx
	template := tab2(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func Tab3(c *gin.Context) {
	//检查请求来源是不是htmx
	template := tab3(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// 比较简单的例子，直接返回一个字符串
func showContentAPIHandler(c *gin.Context) {
	// Check, if the current request has a 'HX-Request' header.
	// For more information, see https://htmx.org/docs/#request-headers
	if !htmx.IsHTMX(c.Request) {
		// If not, return HTTP 400 error.
		c.AbortWithError(http.StatusBadRequest, errors.New("non-htmx request"))
		return
	}

	// Write HTML content.
	c.Writer.Write([]byte("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>"))

	// Send htmx response.
	htmx.NewResponse().Write(c.Writer)
}
