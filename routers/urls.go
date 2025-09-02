package routers

import (
	"net/http"

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
	"github.com/yumenaka/comigo/templ/pages/login_page"
	"github.com/yumenaka/comigo/templ/pages/scroll"
	"github.com/yumenaka/comigo/templ/pages/settings"
	"github.com/yumenaka/comigo/templ/pages/shelf"
	"github.com/yumenaka/comigo/templ/pages/upload_page"
)

// BindURLs 为前端绑定 API 路由
func BindURLs() {
	// 绑定公开页面与api
	publicViewGroup := engine.Group("")
	// API 的URL统一以 /api 开头
	publicAPI := engine.Group("/api")
	bindPublicView(publicViewGroup)
	bindPublicAPI(publicAPI)

	// 可能需要登录的页面
	privateViewGroup := publicViewGroup.Group("")
	privateAPI := publicAPI.Group("")

	// echo jwt简明教程，还有google登录示例：https://echo.labstack.com/docs/cookbook/jwt
	if config.GetUsername() != "" && config.GetPassword() != "" {
		// jwtConfig格式参考：https://echo.labstack.com/docs/middleware/jwt#configuration
		jwtConfig := echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(login.JwtCustomClaims)
			},
			SigningKey: []byte(config.GetJwtSigningKey()),
			// 从Cookie中获取token
			TokenLookup: "cookie:" + login.CookieName + ",header:Authorization:Bearer ",
			// 处理验证错误
			ErrorHandler: func(c echo.Context, err error) error {
				// 更安全的方式判断API请求
				path := c.Request().URL.Path
				if len(path) >= 4 && path[:4] == "/api" {
					return echo.NewHTTPError(http.StatusUnauthorized, "请先登录")
				}
				// 页面请求重定向到登录页
				return c.Redirect(http.StatusFound, "/login")
			},
		}
		privateAPI.Use(echojwt.WithConfig(jwtConfig))
		privateViewGroup.Use(echojwt.WithConfig(jwtConfig))
	}
	bindProtectedView(privateViewGroup)
	bindProtectedAPI(privateAPI)
}

// bindPublicView 注册公共页面
func bindPublicView(group *echo.Group) {
	group.GET("/login", login_page.Handler)
}

// bindPublicAPI 注册公共路由
func bindPublicAPI(group *echo.Group) {
	// 生成QRCode
	group.GET("/qrcode.png", get_data_api.GetQrcode)
	// 查看服务器状态
	group.GET("/server_info", get_data_api.GetServerInfoHandler)
	group.POST("/login", login.Login)
	group.POST("/logout", login.Logout)
}

// bindProtectedView 注册需要登录的页面
func bindProtectedView(group *echo.Group) {
	// 主页
	group.GET("/", shelf.PageHandler)
	group.GET("/index.html", shelf.PageHandler)
	// 书架
	group.GET("/shelf/:id", shelf.PageHandler)
	// 卷轴模式
	group.GET("/scroll/:id", scroll.PageHandler)
	// 翻页模式
	group.GET("/flip/:id", flip.PageHandler)
	// 上传页面
	group.GET("/upload", upload_page.PageHandler)
	// 设置页面
	group.GET("/settings", settings.PageHandler)
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
	group.GET("/top_shelf", get_data_api.GetTopOfShelfInfo)
	// 查询书籍信息
	group.GET("/get_book", get_data_api.GetBook)
	// 查询父书籍信息
	group.GET("/parent_book_info", get_data_api.GetParentBook)
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
	group.GET("/htmx/settings/tab-book", settings.TabBook)
	group.GET("/htmx/settings/tab-net", settings.TabNetwork)
	group.GET("/htmx/settings/tab-labs", settings.TabLabs)
	group.POST("/update-string-config", settings.UpdateStringConfigHandler)
	group.POST("/update-bool-config", settings.UpdateBoolConfigHandler)
	group.POST("/update-number-config", settings.UpdateNumberConfigHandler)
	group.POST("/update-user-info", settings.UpdateUserInfoConfigHandler)
	group.POST("/delete-array-config", settings.DeleteArrayConfigHandler)
	group.POST("/add-array-config", settings.AddArrayConfigHandler)
	group.POST("/config-save", settings.HandleConfigSave)
	group.POST("/config-delete", settings.HandleConfigDelete)
	group.GET("/tailscale_status", get_data_api.GetTailscaleStatus)
}
