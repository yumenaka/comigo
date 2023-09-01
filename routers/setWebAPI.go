package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/routers/handler"
	"github.com/yumenaka/comi/routers/token"
	"github.com/yumenaka/comi/routers/websocket"
)

// 2、设置获取书籍信息、图片文件的 API
func setWebAPI(engine *gin.Engine) {
	////TODO：第三方认证，可参考 https://darjun.github.io/2021/07/26/godailylib/goth/
	api = engine.Group("/api")
	//如果没有设置用户名和密码，则不需要验证
	if common.Config.Username == "" || common.Config.Password == "" {
		// 在需要验证的api中用jwt中间件
		//通过URL字符串参数获取特定文件
		api.GET("/getfile", handler.GetFileHandler)
		//文件上传
		api.POST("/upload", handler.UploadHandler)
		//登录后才能查看的服务器状态，包括标题、机器状态等
		api.GET("/get_status_all", handler.GetStatusAllHandler)
		//获取书架信息，不包含每页信息
		api.GET("/getlist", handler.GetBookListHandler)
		//通过URL字符串参数查询书籍信息
		api.GET("/getbook", handler.GetBookHandler)
		//通过链接下载reg配置
		api.GET("/comigo.reg", handler.GetRegFIleHandler)
		//通过链接下载toml格式的示例配置
		api.GET("/config.toml", handler.GetTomlConfigHandler)
		//获取json格式的当前配置
		api.GET("/config.json", handler.GetJsonConfigHandler)
		//修改服务器配置(post json)
		api.POST("/post_config", handler.PostConfigHandler)
		//修改服务器配置(post json)
		api.POST("/update_config", handler.UpdateConfigHandler)
	} else {
		// 创建 jwt 中间件
		jwtMiddleware, err := token.NewJwtMiddleware()
		if err != nil {
			fmt.Println("JWT Error:" + err.Error())
		}
		// 登录 api ，直接用 jwtMiddleware 中的 `LoginHandler`
		//这个函数中，会执行NewJwtMiddleware()中设置的Authenticator来验证用户权限，如果通过会返回token。
		api.POST("/login", jwtMiddleware.LoginHandler)
		//退出登录，会将用户的cookie中的token删除。
		api.POST("/logout", jwtMiddleware.LogoutHandler)
		// 刷新 token ，延长token的有效期
		api.GET("/refresh_token", jwtMiddleware.RefreshHandler)

		// 在需要验证的api中用jwt中间件
		//通过URL字符串参数获取特定文件
		api.GET("/getfile", jwtMiddleware.MiddlewareFunc(), handler.GetFileHandler)
		//文件上传
		api.POST("/upload", jwtMiddleware.MiddlewareFunc(), handler.UploadHandler)
		//登录后才能查看的服务器状态，包括标题、机器状态等
		api.GET("/get_status_all", jwtMiddleware.MiddlewareFunc(), handler.GetStatusAllHandler)
		//获取书架信息，不包含每页信息
		api.GET("/getlist", jwtMiddleware.MiddlewareFunc(), handler.GetBookListHandler)
		//通过URL字符串参数查询书籍信息
		api.GET("/getbook", jwtMiddleware.MiddlewareFunc(), handler.GetBookHandler)
		//通过链接下载reg配置
		api.GET("/comigo.reg", jwtMiddleware.MiddlewareFunc(), handler.GetRegFIleHandler)
		//通过链接下载示例配置
		api.GET("/config.toml", jwtMiddleware.MiddlewareFunc(), handler.GetTomlConfigHandler)
		//获取json格式的当前配置
		api.GET("/config.json", jwtMiddleware.MiddlewareFunc(), handler.GetJsonConfigHandler)

		//修改服务器配置(post json)
		api.POST("/post_config", jwtMiddleware.MiddlewareFunc(), handler.PostConfigHandler)
		//修改服务器配置(post json)
		api.POST("/update_config", jwtMiddleware.MiddlewareFunc(), handler.UpdateConfigHandler)
	}

	//web端公开的服务器状态，包括标题、端口等
	api.GET("/getstatus", handler.PublicServerInfoHandler)
	////通过URL字符串参数PDF文件里的图片，效率太低，注释掉
	//api.GET("/get_pdf_image", handler.GetPdfImageHandler)

	//通过链接下载qrcode
	api.GET("/qrcode.png", handler.GetQrcodeHandler)
	////301重定向跳转示例
	//api.GET("/redirect", func(c *gin.Context) {
	//	//支持内部和外部的重定向
	//	c.Redirect(http.StatusMovedPermanently, "https://www.google.com/")
	//})
	//初始化websocket
	websocket.WsDebug = &common.Config.Debug
	api.GET("/ws", websocket.WsHandler)
	SetDownloadLink()
	// swagger 自动生成文档用
	if swagHandler != nil {
		engine.GET("/swagger/*any", swagHandler)
	}
}
