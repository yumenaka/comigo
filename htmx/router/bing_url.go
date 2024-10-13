package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/htmx/templates/pages/flip"
	"github.com/yumenaka/comigo/htmx/templates/pages/scroll"
	"github.com/yumenaka/comigo/htmx/templates/pages/shelf"
)

func bindURL(router *gin.Engine) {
	// 主页视图
	router.GET("/", shelf.ShelfHandler)
	// 书架视图
	router.GET("/shelf/:id", shelf.ShelfHandler)
	// 卷轴模式视图
	router.GET("/scroll/:id", scroll.ScrollHandler)
	// 翻页模式视图
	router.GET("/flip/:id", flip.FlipHandler)
}
