package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/htmx/templates/pages/flip"
	"github.com/yumenaka/comigo/htmx/templates/pages/scroll"
	"github.com/yumenaka/comigo/htmx/templates/pages/shelf"
)

func bind(router *gin.Engine) {
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
}

func bindAPI(router *gin.Engine) {
	router.GET("/api/shelf/:id", shelf.GetBookListHandler)
}
