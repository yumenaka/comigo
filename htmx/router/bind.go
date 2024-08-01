package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/htmx/templates/pages"
)

func bindURL(router *gin.Engine) {
	// 处理主页视图
	router.GET("/", pages.TopPageHandler)
	// 书架
	router.GET("/show/:id", pages.ShelfHandler)
	// 处理漫画视图
	router.GET("/read", pages.ReadHandler)
	// 处理 API
	router.GET("/api/hello-world", pages.ShowContentAPIHandler)
}
