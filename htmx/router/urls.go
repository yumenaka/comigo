package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/htmx/templates/pages/flip"
	"github.com/yumenaka/comigo/htmx/templates/pages/scroll"
	"github.com/yumenaka/comigo/htmx/templates/pages/settings_page"
	"github.com/yumenaka/comigo/htmx/templates/pages/shelf"
	"github.com/yumenaka/comigo/htmx/templates/pages/upload_page"
)

func setURLs(e *echo.Echo) {
	bindView(e)
	bindAPI(e)
}

func bindView(e *echo.Echo) {
	// 主页
	e.GET("/", shelf.Handler)
	// 书架
	e.GET("/shelf/:id", shelf.Handler)
	// 卷轴模式
	e.GET("/scroll/:id", scroll.Handler)
	// 翻页模式
	e.GET("/flip/:id", flip.Handler)
	// 上传页面
	e.GET("/upload", upload_page.Handler)
	// 设置页面
	e.GET("/settings", settings_page.Handler)
}

// 注册路由
func bindAPI(e *echo.Echo) {
	e.GET("/api/shelf/:id", shelf.GetBookListHandler)

	e.GET("/htmx/settings/tab1", settings_page.Tab1)
	e.GET("/htmx/settings/tab2", settings_page.Tab2)
	e.GET("/htmx/settings/tab3", settings_page.Tab3)
	// Htmx api
	e.POST("/api/update-string_config", settings_page.UpdateStringConfigHandler)
	e.POST("/api/update-bool-config", settings_page.UpdateBoolConfigHandler)
	e.POST("/api/update-number-config", settings_page.UpdateNumberConfigHandler)
	e.POST("/api/delete-array-config", settings_page.DeleteArrayConfigHandler)
	e.POST("/api/add-array-config", settings_page.AddArrayConfigHandler)

	e.POST("/api/config-save", settings_page.HandleConfigSave)
	e.POST("/api/config-delete", settings_page.HandleConfigDelete)
}
