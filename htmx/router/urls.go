package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/htmx/templates/pages/flip"
	"github.com/yumenaka/comigo/htmx/templates/pages/scroll"
	"github.com/yumenaka/comigo/htmx/templates/pages/settings_page"
	"github.com/yumenaka/comigo/htmx/templates/pages/shelf"
	"github.com/yumenaka/comigo/htmx/templates/pages/upload_page"
)

func setURLs(router *gin.Engine) {
	bindView(router)
	bindAPI(router)
}

func bindView(router *gin.Engine) {
	// 主页
	router.GET("/", shelf.Handler)
	// 书架
	router.GET("/shelf/:id", shelf.Handler)
	// 卷轴模式
	router.GET("/scroll/:id", scroll.Handler)
	// 翻页模式
	router.GET("/flip/:id", flip.Handler)
	// 上传页面
	router.GET("/upload", upload_page.Handler)
	// 设置页面
	router.GET("/settings", settings_page.Handler)
}

func bindAPI(router *gin.Engine) {
	router.GET("/api/shelf/:id", shelf.GetBookListHandler)

	router.GET("/htmx/settings/tab1", settings_page.Tab1)
	router.GET("/htmx/settings/tab2", settings_page.Tab2)
	router.GET("/htmx/settings/tab3", settings_page.Tab3)
	//Htmx api
	router.POST("/api/update-string_config", settings_page.UpdateStringConfigHandler)
	router.POST("/api/update-bool-config", settings_page.UpdateBoolConfigHandler)
	router.POST("/api/update-number-config", settings_page.UpdateNumberConfigHandler)
}
