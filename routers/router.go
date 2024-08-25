package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/resource"
	"github.com/yumenaka/comigo/util/locale"
	"net/http"
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

	// CORS 中间件
	engine.Use(func(c *gin.Context) {
		//在这个代码中，CORS 中间件设置了几个关键的 HTTP 头：
		//Access-Control-Allow-Origin: 指定允许跨域请求的域名。在示例中使用 * 表示允许所有域名。根据您的需求，您也可以指定具体的域名。
		//Access-Control-Allow-Methods: 指定允许的 HTTP 方法。
		//Access-Control-Allow-Headers: 指定允许的 HTTP 头部字段。
		//Access-Control-Expose-Headers: 指定可以暴露给浏览器的响应头部字段。
		//Access-Control-Allow-Credentials: 指定是否允许发送 cookies。
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有域名
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		//在中间件中，针对 OPTIONS 请求做了特殊处理。这是因为在发送实际请求之前，浏览器会先发送一个 OPTIONS 请求（预检请求），以确定服务器是否允许跨域请求。
		//对于这个预检请求，我们直接返回状态码 204 No Content 并结束请求处理。
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	//嵌入静态文件到二进制文件
	resource.EmbedResoure(engine, locale.GetString("HTML_TITLE")+config.Version)
	//设置各种API
	BindAPI(engine)
	//显示QRCode
	showQRCode()
	//监听并启动web服务
	startEngine(engine)
}
