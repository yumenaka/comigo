package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/htmx/templates"
)

func setHandler(router *gin.Engine) {
	// 处理主页视图
	router.GET("/", templates.MainHandler)
	// 处理 API
	router.GET("/api/hello-world", templates.ShowContentAPIHandler)
}
