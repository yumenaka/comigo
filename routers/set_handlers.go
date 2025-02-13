package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/config_handlers"
	"github.com/yumenaka/comigo/routers/handlers"
	"github.com/yumenaka/comigo/routers/websocket"
)

// BindAPI 为前端绑定 API 路由 (Echo 版本)
func BindAPI(e *echo.Echo) {
	// 路由组，方便管理部分相同的 URL
	api := e.Group("/api")

	// 注册公共路由
	publicRoutes(api)

	// 判断是否需要 JWT 认证
	if config.GetPassword() != "" {
		//jwtMiddleware, err := token.NewJwtMiddleware()
		//if err != nil {
		//	logger.Fatalf("JWT Error: %s", err.Error())
		//}

		//// 登录、注销和刷新 token 路由
		//// 假设 token.NewJwtMiddleware() 里已生成 Echo 版本的 handler
		//api.POST("/login", jwtMiddleware.LoginHandler)
		//api.POST("/logout", jwtMiddleware.LogoutHandler)
		//api.GET("/refresh_token", jwtMiddleware.RefreshHandler)

		// 受保护的路由，应用 JWT 中间件
		protectedAPI := api.Group("")
		//protectedAPI.Use(jwtMiddleware.MiddlewareFunc())
		protectedRoutes(protectedAPI)
	} else {
		// 如果不需要认证，直接注册受保护的路由
		protectedAPI := api.Group("")
		protectedRoutes(protectedAPI)
	}
}

// publicRoutes 注册公共路由 (Echo 版本)
func publicRoutes(rg *echo.Group) {
	rg.GET("/qrcode.png", handlers.GetQrcode)
	rg.GET("/server_info", handlers.GetServerInfoHandler)

	// 需要把 WsDebug 替换到正确位置
	websocket.WsDebug = &config.GetCfg().Debug
	rg.GET("/ws", websocket.WsHandler)
}

// protectedRoutes 注册需要认证的路由 (Echo 版本)
func protectedRoutes(rg *echo.Group) {
	// 文件上传
	rg.POST("/upload", handlers.UploadFile)
	// 获取特定文件
	rg.GET("/get_file", handlers.GetFile)
	// 直接下载原始文件
	rg.GET("/raw/:book_id/:file_name", handlers.GetRawFile)
	// 查看服务器状态
	rg.GET("/server_info_all", handlers.GetAllServerInfoHandler)
	// 获取书架信息
	rg.GET("/book_infos", handlers.GetBookInfos)
	// 获取书架信息 2.0
	rg.GET("/top_shelf", handlers.GetTopOfShelfInfo)
	// 查询书籍信息
	rg.GET("/get_book", handlers.GetBook)
	// 查询父书籍信息
	rg.GET("/parent_book_info", handlers.GetParentBookInfo)
	// 返回同一文件夹的书籍 ID 列表
	rg.GET("/group_info", handlers.GroupInfo)
	rg.GET("/group_info_filter", handlers.GroupInfoFilter)
	// 下载 reg 设置文件
	rg.GET("/comigo.reg", handlers.GetRegFile)
	// 获取配置
	rg.GET("/config", config_handlers.GetConfig)
	// 生成图片 http://localhost:1234/api/generate_image?height=220&width=160&text=12345&font_size=32
	rg.GET("/generate_image", handlers.GenerateImage)
	// 获取配置状态
	rg.GET("/config/status", config_handlers.GetConfigStatus)
	// 更新配置
	rg.PUT("/config", config_handlers.UpdateConfig)
	// 保存配置到文件
	rg.POST("/config/:to", config_handlers.SaveConfigHandler)
	// 删除特定路径下的配置
	rg.DELETE("/config/:in", config_handlers.DeleteConfig)
	// 下载 toml 格式的示例配置
	rg.GET("/config.toml", config_handlers.GetConfigToml)
}
