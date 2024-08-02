package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/htmx/templates/pages"
)

func bindURL(router *gin.Engine) {
	// 主页视图
	router.GET("/", pages.ShelfHandler)
	// 书架视图
	router.GET("/shelf/:id", pages.ShelfHandler)
	// 漫画视图
	router.GET("/scroll/:id", pages.ScrollHandler)
}
