package routers

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/config_api"
	"github.com/yumenaka/comigo/routers/get_data_api"
	"github.com/yumenaka/comigo/routers/login"
	"github.com/yumenaka/comigo/routers/upload_api"
	"github.com/yumenaka/comigo/routers/websocket"
	"github.com/yumenaka/comigo/templ/pages/flip"
	"github.com/yumenaka/comigo/templ/pages/scroll"
	"github.com/yumenaka/comigo/templ/pages/settings_page"
	"github.com/yumenaka/comigo/templ/pages/shelf"
	"github.com/yumenaka/comigo/templ/pages/upload_page"
)

// BindURLs 为前端绑定 API 路由
func BindURLs() {
	publicViewGroup := engine.Group("")
	bindView(publicViewGroup)

	// API 的URL统一以 /api 开头
	api := engine.Group("/api")
	// 无须登录的公开api
	bindPublicAPI(api)
	// 可能需要登录的api
	privateAPI := api.Group("")

	// echo 自带jwt简明教程，还有google登录示例：https://echo.labstack.com/docs/cookbook/jwt
	if config.GetPassword() != "" {
		// jwtConfig格式参考：https://echo.labstack.com/docs/middleware
		jwtConfig := echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(login.JwtCustomClaims)
			},
			SigningKey: []byte(config.GetJwtSigningKey()),
		}
		privateAPI.Use(echojwt.WithConfig(jwtConfig))
	}
	bindProtectedAPI(privateAPI)
}

// bindPublicAPI 注册公共路由
func bindPublicAPI(group *echo.Group) {
	// 生成QRCode
	group.GET("/qrcode.png", get_data_api.GetQrcode)
	// 查看服务器状态（TODO：限制信息范围）
	group.GET("/server_info", get_data_api.GetServerInfoHandler)
}

// bindProtectedAPI 注册需要认证的路由
func bindProtectedAPI(group *echo.Group) {
	// 文件上传
	group.POST("/upload", upload_api.UploadFile)
	// 获取特定文件
	group.GET("/get_file", get_data_api.GetFile)
	// 直接下载原始文件
	group.GET("/raw/:book_id/:file_name", get_data_api.GetRawFile)
	// 查看服务器状态
	group.GET("/server_info_all", get_data_api.GetAllServerInfoHandler)
	// 获取书架信息
	group.GET("/book_infos", get_data_api.GetBookInfos)
	// 获取书架信息 2.0
	group.GET("/top_shelf", get_data_api.GetTopOfShelfInfo)
	// 查询书籍信息
	group.GET("/get_book", get_data_api.GetBook)
	// 查询父书籍信息
	group.GET("/parent_book_info", get_data_api.GetParentBookInfo)
	// 返回同一文件夹的书籍 ID 列表
	group.GET("/group_info", get_data_api.GroupInfo)
	group.GET("/group_info_filter", get_data_api.GroupInfoFilter)
	// 下载 reg 设置文件
	group.GET("/comigo.reg", get_data_api.GetRegFile)
	// 获取配置
	group.GET("/config", config_api.GetConfig)
	// 生成图片 http://localhost:1234/api/generate_image?height=220&width=160&text=12345&font_size=32
	group.GET("/generate_image", get_data_api.GetGeneratedImage)
	// 获取配置状态
	group.GET("/config/status", config_api.GetConfigStatus)
	// 更新配置
	group.PUT("/config", config_api.UpdateConfig)
	// 保存配置到文件
	group.POST("/config/:to", config_api.SaveConfigHandler)
	// 删除特定路径下的配置
	group.DELETE("/config/:in", config_api.DeleteConfig)
	// 下载 toml 格式的示例配置
	group.GET("/config.toml", config_api.GetConfigToml)
	// TODO: 测试需要登录时的表现
	websocket.WsDebug = &config.GetCfg().Debug
	group.GET("/ws", websocket.WsHandler)
	// 新加的 HTMX 相关路由
	group.GET("/shelf/:id", shelf.GetBookListHandler)
	group.GET("/htmx/settings/tab1", settings_page.Tab1)
	group.GET("/htmx/settings/tab2", settings_page.Tab2)
	group.GET("/htmx/settings/tab3", settings_page.Tab3)
	group.POST("/update-string_config", settings_page.UpdateStringConfigHandler)
	group.POST("/update-bool-config", settings_page.UpdateBoolConfigHandler)
	group.POST("/update-number-config", settings_page.UpdateNumberConfigHandler)
	group.POST("/delete-array-config", settings_page.DeleteArrayConfigHandler)
	group.POST("/add-array-config", settings_page.AddArrayConfigHandler)
	group.POST("/config-save", settings_page.HandleConfigSave)
	group.POST("/config-delete", settings_page.HandleConfigDelete)
}

func bindView(group *echo.Group) {
	// 主页
	group.GET("/", shelf.Handler)
	group.GET("/index.html", shelf.Handler)
	// 书架
	group.GET("/shelf/:id", shelf.Handler)
	// 卷轴模式
	group.GET("/scroll/:id", scroll.Handler)
	// 翻页模式
	group.GET("/flip/:id", flip.Handler)
	// 上传页面
	group.GET("/upload", upload_page.Handler)
	// 设置页面
	group.GET("/settings", settings_page.Handler)
}
