package common

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// RenderHTML 将 templ 组件直接写入 Echo 响应。
// 页面级 handler 统一通过这里处理渲染错误，避免重复 500 响应逻辑。
func RenderHTML(c echo.Context, component templ.Component) error {
	if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
