package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/htmx/templates/pages"
)

func bindURL(router *gin.Engine) {
	// 主页视图
	router.GET("/", pages.TopPageHandler)
	// 书架视图
	router.GET("/show/:id", pages.ShelfHandler)
	// 漫画视图
	router.GET("/read", pages.ReadHandler)
}
