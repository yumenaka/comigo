package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/htmx/templates/pages"
)

func bindURL(router *gin.Engine) {
	// 主页视图
	router.GET("/", pages.ShelfHandler)
	// 书架视图
	router.GET("/shelf/:id", pages.ShelfHandler)
	// 卷轴模式视图
	router.GET("/scroll/:id", pages.ScrollHandler)
	// 翻页模式视图
	router.GET("/flip/:id", pages.FlipHandler)
}
