package routers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/config_api"
	"github.com/yumenaka/comigo/routers/data_api"
	"github.com/yumenaka/comigo/routers/login"
	"github.com/yumenaka/comigo/routers/reverse_proxy"
	"github.com/yumenaka/comigo/routers/upload_api"
	"github.com/yumenaka/comigo/routers/websocket"
	"github.com/yumenaka/comigo/templ/pages/flip"
	"github.com/yumenaka/comigo/templ/pages/login_page"
	"github.com/yumenaka/comigo/templ/pages/scroll"
	"github.com/yumenaka/comigo/templ/pages/settings"
	"github.com/yumenaka/comigo/templ/pages/shelf"
	"github.com/yumenaka/comigo/templ/pages/upload_page"
	"github.com/yumenaka/comigo/tools/sse_hub"
)

// BindURLs 为前端绑定 API 路由
func BindURLs() {
	// 绑定公开页面与api
	publicViewGroup := engine.Group("")
	// API 的URL统一以 /api 开头
	publicAPI := engine.Group("/api")
	bindPublicView(publicViewGroup)
	bindPublicAPI(publicAPI)

	// 可以设置登录保护的页面与api
	privateViewGroup := publicViewGroup.Group("")
	privateAPI := publicAPI.Group("")

	// echo jwt简明教程，还有google登录示例：https://echo.labstack.com/docs/cookbook/jwt
	if config.GetCfg().RequiresAuth() {
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
	// 简化路径：/get.sh -> https://raw.githubusercontent.com/yumenaka/comigo/master/get.sh
	group.GET("/get.sh", reverse_proxy.GetComigoScriptHandler)
	group.HEAD("/get.sh", reverse_proxy.GetComigoScriptHandler)
	// GitHub 反向代理：/yumenaka/* -> https://github.com/yumenaka/*
	group.GET("/yumenaka/*", reverse_proxy.ProxyHandler)
	group.HEAD("/yumenaka/*", reverse_proxy.ProxyHandler)
}

// bindPublicAPI 注册公共路由
func bindPublicAPI(group *echo.Group) {
	// 生成QRCode
	group.GET("/qrcode.png", data_api.GetQrcode)
	group.POST("/login", login.Login)
	group.POST("/logout", login.Logout)
}

// bindProtectedView 注册需要登录的页面
func bindProtectedView(group *echo.Group) {
	// 主页
	group.GET("/", shelf.ShelfHandler)
	group.GET("/index.html", shelf.ShelfHandler)
	group.GET("/shelf/:id", shelf.ShelfHandler)
	// 卷轴模式
	group.GET("/scroll/:id", scroll.ScrollModeHandler)
	// 翻页模式
	group.GET("/flip/:id", flip.FlipModeHandler)
	// 上传页面
	group.GET("/upload", upload_page.PageHandler)
	// 设置页面
	group.GET("/settings", settings.PageHandler)
}

// bindProtectedAPI 注册需要认证的路由
func bindProtectedAPI(group *echo.Group) {
	// 服务器状态
	group.GET("/server_info", data_api.GetServerInfoHandler)
	// 获取书库列表
	group.GET("/stores", data_api.GetStores)
	// 文件上传
	group.POST("/upload", upload_api.UploadFile)
	// 获取特定文件
	group.GET("/get_file", data_api.GetFile)
	// 获取书籍封面
	group.GET("/get_cover", data_api.GetCover)
	// 直接下载原始文件
	group.GET("/raw/:book_id/:file_name", data_api.GetRawFile)
	// 获取书架信息
	group.GET("/top_shelf", data_api.GetTopOfShelfInfo)
	// 查询书籍信息
	group.GET("/get_book", data_api.GetBook)
	// 获取所有书签的API
	group.GET("/all_bookmarks", data_api.GetAllBookmarks)
	// 更新书签信息
	group.POST("/store_bookmark", data_api.StoreBookmark)
	// 查询父书籍信息
	group.GET("/parent_book_info", data_api.GetParentBook)
	// 下载 reg 设置文件
	group.GET("/comigo.reg", data_api.GetRegFile)
	// 获取配置
	group.GET("/config", config_api.GetConfig)
	// 生成图片 http://localhost:1234/api/generate_image?height=220&width=160&text=12345&font_size=32
	group.GET("/generate_image", data_api.GetGeneratedImage)
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
	// 字符串、布尔值、数字配置的更改
	group.POST("/update-string-config", settings.UpdateStringConfigHandler)
	group.POST("/update-bool-config", settings.UpdateBoolConfigHandler)
	group.POST("/update-number-config", settings.UpdateNumberConfigHandler)
	// 更改Comigo登录设置
	group.POST("/update-login-settings", settings.UpdateLoginSettingsHandler)
	// Tailscale配置更新JSON API
	group.POST("/submit-tailscale-config", settings.UpdateTailscaleConfigHandler)
	// 字符串数组配置的增删改
	group.POST("/delete-array-config", settings.DeleteArrayConfigHandler)
	group.POST("/add-array-config", settings.AddArrayConfigHandler)
	// 书库管理
	group.POST("/rescan-store", settings.RescanStoreHandler)
	group.POST("/delete-store", settings.DeleteStoreHandler)
	// 插件管理
	group.POST("/enable-plugin", settings.EnablePluginHandler)
	group.POST("/disable-plugin", settings.DisablePluginHandler)
	// 保存和删除配置
	group.POST("/config-save", settings.HandleConfigSave)
	group.POST("/config-delete", settings.HandleConfigDelete)
	// 获取 tailscale 状态
	group.GET("/tailscale_status", data_api.GetTailscaleStatus)
	// SSE 服务器发送事件
	group.GET("/sse", sse_hub.SSEHandler)
	// SSE 广播接口
	group.POST("/push", sse_hub.PushHandler)
}
