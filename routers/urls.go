package routers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/config_api"
	"github.com/yumenaka/comigo/routers/data_api"
	"github.com/yumenaka/comigo/routers/login"
	"github.com/yumenaka/comigo/routers/opds"
	"github.com/yumenaka/comigo/routers/reverse_proxy"
	"github.com/yumenaka/comigo/routers/upload_api"
	"github.com/yumenaka/comigo/routers/websocket"
	"github.com/yumenaka/comigo/templ/pages/flip"
	"github.com/yumenaka/comigo/templ/pages/login_page"
	"github.com/yumenaka/comigo/templ/pages/manual"
	"github.com/yumenaka/comigo/templ/pages/player"
	"github.com/yumenaka/comigo/templ/pages/reader"
	"github.com/yumenaka/comigo/templ/pages/scroll"
	"github.com/yumenaka/comigo/templ/pages/settings"
	"github.com/yumenaka/comigo/templ/pages/shelf"
	"github.com/yumenaka/comigo/templ/pages/upload_page"
	"github.com/yumenaka/comigo/tools/sse_hub"
)

// BindURLs 为前端绑定 API 路由
func BindURLs() {
	// Manual 固定在根路径且始终公开，不受 BasePath 与登录保护影响。
	engine.GET("/manual", manual.Handler)
	engine.GET("/manual/", manual.Handler)
	engine.GET("/manual/*", manual.Handler)

	// 绑定公开页面与api
	basePath := config.GetBasePath()
	if basePath != "" {
		engine.GET(basePath, func(c echo.Context) error {
			return c.Redirect(http.StatusMovedPermanently, basePath+"/")
		})
	}
	publicViewGroup := engine.Group(basePath)
	// API 的URL统一以 /api 开头；配置 BasePath 后实际挂载在 /base/api 下。
	publicAPI := engine.Group(config.PrefixPath("/api"))
	bindPublicView(publicViewGroup)
	bindPublicAPI(publicAPI)

	// 可以设置登录保护的页面与api
	privateViewGroup := publicViewGroup.Group("")
	privateAPI := publicAPI.Group("")

	// echo jwt简明教程，还有google登录示例：https://echo.labstack.com/docs/cookbook/jwt
	{
		// jwtConfig格式参考：https://echo.labstack.com/docs/middleware/jwt#configuration
		jwtConfig := echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(login.JwtCustomClaims)
			},
			// 每次请求读取认证配置，开启密码立即生效，修改密码使旧令牌失效。
			Skipper: func(c echo.Context) bool { return !config.GetCfg().RequiresAuth() },
			KeyFunc: func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, echo.ErrUnauthorized
				}
				return []byte(config.GetJwtSigningKey()), nil
			},
			// 从Cookie中获取token
			TokenLookup: "cookie:" + login.CookieName + ",header:Authorization:Bearer ",
			// 处理验证错误
			ErrorHandler: func(c echo.Context, err error) error {
				// 更安全的方式判断API请求
				path := config.StripBasePath(c.Request().URL.Path)
				if len(path) >= 4 && path[:4] == "/api" {
					return echo.NewHTTPError(http.StatusUnauthorized, locale.GetString("err_login_required"))
				}
				// 页面请求重定向到登录页
				return c.Redirect(http.StatusFound, config.PrefixPath("/login"))
			},
		}
		privateAPI.Use(echojwt.WithConfig(jwtConfig))
		privateViewGroup.Use(echojwt.WithConfig(jwtConfig))
	}
	privateViewGroup.Use(data_api.DeviceCookie)
	privateAPI.Use(data_api.DeviceCookie)
	bindProtectedView(privateViewGroup)
	bindProtectedAPI(privateAPI)
}

// bindPublicView 注册公共页面
func bindPublicView(group *echo.Group) {
	group.GET("/login", login_page.Handler)
	group.GET("/healthz", data_api.Healthz)
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
	group.GET("/info", data_api.GetPublicInfo)
}

// bindProtectedView 注册需要登录的页面
func bindProtectedView(group *echo.Group) {
	// 书架相关页面
	group.GET("/", shelf.ShelfHandler) //主书架
	group.GET("/index.html", shelf.ShelfHandler)
	group.GET("/shelf/:id", shelf.ShelfHandler) //子书架
	group.GET("/search", shelf.SearchHandler)   //书架搜索
	// 卷轴阅读
	group.GET("/scroll/:id", scroll.ScrollModeHandler)
	// 翻页阅读
	group.GET("/flip/:id", flip.FlipModeHandler)
	// 播放器模式
	group.GET("/player/:id", player.PlayerModeHandler)
	// 上传页面
	group.GET("/upload", upload_page.PageHandler)
	// 本地压缩包阅读页面（文件不上传服务器）
	group.GET("/reader", reader.PageHandler)
	// 设置页面
	group.GET("/settings", settings.PageHandler)
	// OPDS 1.2 目录
	group.GET("/opds", opds.RootHandler)
	group.GET("/opds/books", opds.BooksHandler)
	group.GET("/opds/books/:id", opds.BooksHandler)
}

// bindProtectedAPI 注册需要认证的路由
func bindProtectedAPI(group *echo.Group) {
	bindServerAPI(group)
	bindBookAPI(group)
	bindBookmarkAPI(group)
	bindConfigAPI(group.Group("", controlAccess))
	bindSettingsAPI(group.Group("", controlAccess))
	bindRealtimeAPI(group)
}

// bindServerAPI 注册服务状态、上传和书库管理等服务级 API。
func bindServerAPI(group *echo.Group) {
	control := group.Group("", controlAccess)
	bindWailsAPI(control)
	// 服务器状态
	control.GET("/server", data_api.GetServerInfoHandler)
	control.GET("/server/update", data_api.GetServerUpdateHandler)
	control.GET("/server/traffic", data_api.GetServerTrafficHandler)
	control.GET("/connections", data_api.GetConnections)
	control.POST("/restart", restartHandler)
	control.GET("/autostart", settings.GetAutostartHandler)
	control.PUT("/autostart", settings.SetAutostartHandler)
	// 获取书库列表
	group.GET("/stores", data_api.GetStores)
	group.GET("/stores/:id", data_api.GetStore)
	// 文件上传
	group.POST("/upload", upload_api.UploadFile)
	// 获取 tailscale 状态
	control.GET("/tailscale-status", data_api.GetTailscaleStatus)
	control.GET("/tailscale-status-sse", data_api.GetTailscaleStatusSSE)
}

// bindBookAPI 注册书籍读取、封面、下载和缓存相关 API。
func bindBookAPI(group *echo.Group) {
	// 获取特定文件
	group.GET("/get-file", data_api.GetFile)
	// 获取书籍封面
	group.GET("/get-cover", data_api.GetCover)
	// 直接下载原始文件
	group.GET("/raw/:book_id/:file_name", data_api.GetRawFile)
	// 获取书架信息
	group.GET("/top-shelf", data_api.GetTopOfShelfInfo)
	// 查询书籍信息
	group.GET("/books", data_api.GetBooks)
	group.GET("/books/:id", data_api.GetBook)
	// 查询父书籍信息
	group.GET("/books/:id/parent", data_api.GetParentBook)
	// 下载 reg 设置文件
	group.GET("/comigo.reg", data_api.GetRegFile)
	// 生成图片 http://localhost:1234/api/generate-image?height=220&width=160&text=12345&font_size=32
	group.GET("/generate-image", data_api.GetGeneratedImage)
	// 下载 TypeDir 书籍为 zip 文件
	group.GET("/download-zip", data_api.DownloadZip)
	// 下载书籍为 EPUB 文件
	group.GET("/download-epub", data_api.DownloadEpub)
	// 删除书籍的元数据和缓存文件
	group.DELETE("/books/:id/cache", data_api.DeleteBookCache)
}

// bindBookmarkAPI 注册阅读历史和书签相关 API。
func bindBookmarkAPI(group *echo.Group) {
	// 获取所有书签的API
	group.GET("/bookmarks", data_api.GetAllBookmarks)
	// 获取阅读历史（支持limit和分页参数）
	group.GET("/reading-history", data_api.GetReadingHistory)
	// 更新书签信息
	group.POST("/bookmarks", data_api.StoreBookmark)
	// 删除特定书签
	group.DELETE("/bookmarks", data_api.DeleteBookmark)
}

// bindConfigAPI 注册配置读写 API。
func bindConfigAPI(group *echo.Group) {
	// 获取配置状态
	group.GET("/configs/status", config_api.GetConfigStatus)
	group.GET("/configs", config_api.GetConfig)
	// 更新配置
	group.PATCH("/configs", updateConfigHandler)
}

// bindSettingsAPI 注册设置页使用的轻量操作 API。
func bindSettingsAPI(group *echo.Group) {
	// 字符串、布尔值、数字配置的更改
	group.PUT("/configs/:name", settings.UpdateValueConfigHandler)
	// 更改Comigo登录设置
	group.PATCH("/configs/login", settings.UpdateLoginSettingsHandler)
	// Tailscale配置更新JSON API
	group.PATCH("/configs/tailscale", settings.UpdateTailscaleConfigHandler)
	// 字符串数组配置的增删改
	group.DELETE("/configs/:name/items", settings.DeleteArrayConfigHandler)
	group.POST("/configs/:name/items", settings.AddArrayConfigHandler)
	// 书库管理
	group.POST("/stores/:id/refresh", settings.RescanStoreHandler)
	group.POST("/stores/refresh", settings.RescanAllStoresHandler)
	group.DELETE("/stores/:id", settings.DeleteStoreHandler)
	// 插件管理
	group.PUT("/plugins/:name", settings.EnablePluginHandler)
	group.DELETE("/plugins/:name", settings.DisablePluginHandler)
	// 保存和删除配置
	group.PUT("/configs/files/:location", settings.HandleConfigSave)
	group.DELETE("/configs/files/:location", settings.HandleConfigDelete)
}

// bindRealtimeAPI 注册 WebSocket/SSE 等实时通信 API。
func bindRealtimeAPI(group *echo.Group) {
	websocket.WsDebug = &config.GetCfg().Debug
	group.GET("/ws", websocket.WsHandler, data_api.TrackConnection)
	// SSE 服务器发送事件
	group.GET("/sse", sse_hub.SSEHandler, data_api.TrackConnection)
}
