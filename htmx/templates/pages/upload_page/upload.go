package upload_page

import (
	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/htmx/templates/common"
	"net/http"
)

// Handler 上传文件页面
func Handler(c *gin.Context) {
	indexTemplate := common.MainLayout(
		c,
		&state.Global,
		common.UploadArea(&state.Global),
		"",
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, indexTemplate); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
