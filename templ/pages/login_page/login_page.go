package login_page

import (
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/templ/common"
)

// Handler 上传文件页面
func Handler(c echo.Context) error {
	indexHtml := common.Html(
		c,
		LoginPage(),
		[]string{},
	)
	// 渲染页面
	return common.RenderHTML(c, indexHtml)
}
