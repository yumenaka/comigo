package routers

import (
	"github.com/gin-gonic/gin"
)

// StartWebServer 启动web服务
func StartWebServer() {
	//gin mode：ReleaseMode,DebugMode,TestMode
	gin.SetMode(gin.ReleaseMode)

	//不用 gin.Default()，避免使用 Gin 的默认日志中间件
	engine := gin.New()
	//Recovery 中间件。返回 500 错误页面，避免程序直接崩溃，同时记录错误日志。
	engine.Use(gin.Recovery())
	//日志中间件
	setLogger(engine)

	//嵌入静态文件到二进制文件
	embedFile(engine)
	//设置各种API
	setWebAPI(engine)
	//显示QRCode
	showQRCode()
	//监听并启动web服务
	startEngine(engine)
}
