package settings_page

import (
	"errors"
	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/htmx/state"
	"net/http"
)

// 使用模板中响应htmx请求，页面比较复杂时用
func Tab1(c *gin.Context) {
	template := tab1(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func Tab2(c *gin.Context) {
	template := tab2(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func Tab3(c *gin.Context) {
	template := tab3(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// UpdateBoolConfigHandler 更新Config的htmx接口，返回变更后的html，布尔值专用
func UpdateBoolConfigHandler(c *gin.Context) {
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

// UpdateStringConfigHandler 处理 /api/update-string-config 请求
func UpdateStringConfigHandler(c *gin.Context) {
	// 仅接收 HTMX 请求
	if !htmx.IsHTMX(c.Request) {
		c.AbortWithError(http.StatusBadRequest, errors.New("non-htmx request"))
		return
	}

	// 解析表单
	if err := c.Request.ParseForm(); err != nil {
		c.String(http.StatusBadRequest, "ParseForm error: %v", err)
		return
	}

	// 假设只有一对数据 (key-value)
	formData := c.Request.PostForm
	if len(formData) == 0 {
		c.String(http.StatusBadRequest, "No form data")
		return
	}

	var (
		name     string
		newValue string
	)

	// 这里仅取第一对 key-value
	for key, values := range formData {
		name = key
		if len(values) > 0 {
			newValue = values[0] // values 是一个切片，通常只有一个值，但要注意可能有多个值
		}
		// 只需要取第一对就可以退出循环
		break
	}

	updatedHTML := StringConfig(name, newValue, name+"_Description")

	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
