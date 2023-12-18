package routers

import (
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/routers/handlers"
	"github.com/yumenaka/comi/routers/token"
	"github.com/yumenaka/comi/routers/websocket"
	"github.com/yumenaka/comi/types"
)

var protectedAPI *gin.RouterGroup

// 前端需要的 API
func setWebAPI(engine *gin.Engine) {
	// 路由组,方便管理部分相同的URL
	api := engine.Group("/api")

	// 无需认证，不受保护的路由
	publicRoutes := func(rg *gin.RouterGroup) {
		rg.GET("/qrcode.png", handlers.GetQrcode)
		rg.GET("/server_info", handlers.GetServerInfoPublic)
		websocket.WsDebug = &config.Config.Debug
		rg.GET("/ws", websocket.WsHandler)
	}
	publicRoutes(api)

	// 可能需要认证的路由
	protectedAPI = api.Group("/")
	// 初始化 jwtMiddleware 一次，无论是否设置了密码。
	var jwtMiddleware *jwt.GinJWTMiddleware
	if config.Config.Password != "" {
		var err error
		jwtMiddleware, err = token.NewJwtMiddleware()
		if err != nil {
			log.Fatalf("JWT Error: %s", err.Error()) // 终止程序或其他错误处理
		}
	}
	if jwtMiddleware != nil {
		// 登录、注销和 token 刷新路由只有在设置了密码时才添加
		api.POST("/login", jwtMiddleware.LoginHandler)
		api.POST("/logout", jwtMiddleware.LogoutHandler)
		api.GET("/refresh_token", jwtMiddleware.RefreshHandler)
		// 如果设置了密码，则应用 JWT 中间件到一个新的路由组
		protectedAPI = api.Group("/", jwtMiddleware.MiddlewareFunc())
	}
	//RESTful API
	//Create	POST/PUT
	//Read	    GET
	//Update	PUT
	//Delete	DELETE
	//谷歌API设计规范：
	//https://cloud.google.com/apis/design/standard_methods?hl=zh-cn

	//文件上传
	protectedAPI.POST("/upload", handlers.UploadFile)
	//通过URL字符串参数获取特定文件
	protectedAPI.GET("/get_file", handlers.GetFile)
	//登录后才能查看的服务器状态，包括标题、机器状态等
	protectedAPI.GET("/server_info_all", handlers.GetServerInfo)
	//获取书架信息，不包含每页信息
	protectedAPI.GET("/book_infos", handlers.GetBookInfos)
	//通过URL字符串参数查询书籍信息
	protectedAPI.GET("/get_book", handlers.GetBook)
	//查询父书籍信息
	protectedAPI.GET("/parent_book_info", handlers.GetParentBookInfo)
	//返回同一文件夹的书籍ID列表
	protectedAPI.GET("/group_info", handlers.GroupInfo)
	//返回同一文件夹的书籍ID列表
	protectedAPI.GET("/group_info_filter", handlers.GroupInfoFilter)
	//通过链接下载reg配置
	protectedAPI.GET("/comigo.reg", handlers.GetRegFile)
	//通过链接下载toml格式的示例配置
	protectedAPI.GET("/config.toml", handlers.GetConfigToml)

	//config操作
	protectedAPI.GET("/config", handlers.GetConfig)
	protectedAPI.GET("/config/status", handlers.GetConfigStatus)
	protectedAPI.PUT("/config", handlers.UpdateConfig)
	protectedAPI.POST("/config/:to", handlers.SaveConfig)
	protectedAPI.DELETE("/config/:in", handlers.DeleteConfig)

	//压缩包直接下载链接
	SetDownloadLink()
}

// SetDownloadLink 压缩包直接下载链接
func SetDownloadLink() {
	if types.GetBooksNumber() >= 1 {
		allBook, err := types.GetAllBookInfoList("name")
		if err != nil {
			logger.Info("设置文件下载失败")
		} else {
			for _, info := range allBook.BookInfos {
				//下载文件
				if info.Type != types.TypeBooksGroup && info.Type != types.TypeDir {
					//staticUrl := "/raw/" + info.BookID + "/" + url.QueryEscape(info.title)
					staticUrl := "/raw/" + info.BookID + "/" + info.Title
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
