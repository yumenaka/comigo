package routers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/routers/handlers"
	"github.com/yumenaka/comi/routers/token"
	"github.com/yumenaka/comi/routers/websocket"
	"github.com/yumenaka/comi/types"
	"log"
)

// 2、设置获取书籍信息、图片文件的 API
func setWebAPI(engine *gin.Engine) {
	// swagger 文档路由（加编译参数才会启用）
	if swagHandler != nil {
		engine.GET("/swagger/*any", swagHandler)
	}
	api = engine.Group("/api")

	// 无需认证，不受保护的路由
	publicRoutes := func(rg *gin.RouterGroup) {
		rg.GET("/qrcode.png", handlers.GetQrcodeHandler)
		rg.GET("/server_info_public", handlers.HandlerServerInfoPublic)
		websocket.WsDebug = &config.Config.Debug
		rg.GET("/ws", websocket.WsHandler)
	}
	publicRoutes(api)

	// 初始化 jwtMiddleware 一次，无论是否设置了密码。
	var jwtMiddleware *jwt.GinJWTMiddleware
	if config.Config.Password != "" {
		var err error
		jwtMiddleware, err = token.NewJwtMiddleware()
		if err != nil {
			log.Fatalf("JWT Error: %s", err.Error()) // 终止程序或其他错误处理
		}
	}
	protectedAPI = api.Group("/")
	if jwtMiddleware != nil {
		// 登录、注销和 token 刷新路由只有在设置了密码时才添加
		api.POST("/login", jwtMiddleware.LoginHandler)
		api.POST("/logout", jwtMiddleware.LogoutHandler)
		api.GET("/refresh_token", jwtMiddleware.RefreshHandler)
		// 如果设置了密码，则应用 JWT 中间件到一个新的路由组
		protectedAPI = api.Group("/", jwtMiddleware.MiddlewareFunc())
	}
	// 以下路由都可能需要认证
	//文件上传
	protectedAPI.POST("/upload", handlers.HandlerUpload)
	//通过URL字符串参数获取特定文件
	protectedAPI.GET("/get_file", handlers.HandlerGetFile)
	//登录后才能查看的服务器状态，包括标题、机器状态等
	protectedAPI.GET("/server_info", handlers.HandlerServerInfo)
	//获取书架信息，不包含每页信息
	protectedAPI.GET("/get_shelf", handlers.HandlerGetShelf)
	//通过URL字符串参数查询书籍信息
	protectedAPI.GET("/get_book", handlers.HandlerGetBook)
	//返回同一文件夹的书籍ID列表
	protectedAPI.GET("/same_group_books", handlers.HandlerSameGroupBooks)
	//通过链接下载reg配置
	protectedAPI.GET("/comigo.reg", handlers.HandlerGetRegFile)
	//通过链接下载toml格式的示例配置
	protectedAPI.GET("/config.toml", handlers.HandlerGetConfigToml)
	//获取json格式的当前配置
	protectedAPI.GET("/config.json", handlers.HandlerGetConfigJson)
	//修改服务器配置
	protectedAPI.POST("/config_update", handlers.HandlerConfigUpdate)
	//保存服务器配置
	protectedAPI.POST("/config_save", handlers.HandlerConfigSave)
	// 其他初始化代码，如 SetDownloadLink()
	SetDownloadLink()
}

// SetDownloadLink 设定压缩包下载链接
func SetDownloadLink() {
	if types.GetBooksNumber() >= 1 {
		allBook, err := types.GetAllBookInfoList("name")
		if err != nil {
			logger.Info("设置文件下载失败")
		} else {
			for _, info := range allBook.BaseBooks {
				//下载文件
				if info.Type != types.TypeBooksGroup && info.Type != types.TypeDir {
					//staticUrl := "/raw/" + info.BookID + "/" + url.QueryEscape(info.Name)
					staticUrl := "/raw/" + info.BookID + "/" + info.Name
					if checkUrlRegistered(info.BookID) {
						if config.Config.Debug {
							logger.Info("路径已注册：", info)
						}
						continue
					} else {
						protectedAPI.StaticFile(staticUrl, info.FilePath)
						staticUrlMap[info.BookID] = staticUrl
					}
				}
			}
		}
	}
}
