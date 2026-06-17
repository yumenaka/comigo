package upload_page

import (
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/templ/common"
)

// PageHandler 上传文件页面
func PageHandler(c echo.Context) error {
	indexHtml := common.Html(
		c,
		UploadPage(c),
		[]string{},
	)
	// 渲染页面
	return common.RenderHTML(c, indexHtml)
}
