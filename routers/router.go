package routers

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sanity-io/litter"

	"github.com/yumenaka/comi/book"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/plugin"
	"github.com/yumenaka/comi/routers/handler"
	"github.com/yumenaka/comi/tools"
)

// TemplateString 模板文件
//go:embed static/index.html
var TemplateString string

//go:embed  static
var staticFS embed.FS

//go:embed  static/assets
var staticAssetFS embed.FS

//go:embed  static/images
var staticImageFS embed.FS

//gin-jwt相关 https://github.com/appleboy/gin-jwt
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"text":     "Hello World.",
	})
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

// StartWebServer 启动web服务
func StartWebServer() {
	//设置 gin
	gin.SetMode(gin.ReleaseMode)
	//// 创建带有默认中间件的路由: 日志与恢复中间件
	engine := gin.Default()

	//TODO：登录、认证、鉴权etc //授权有些问题，会导致audio.png报401错误，暂时注释掉
	//// the jwt middleware
	//authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
	//	Realm:       "test zone",
	//	Key:         []byte("secret key"),
	//	Timeout:     time.Hour,
	//	MaxRefresh:  time.Hour,
	//	IdentityKey: identityKey,
	//	PayloadFunc: func(data interface{}) jwt.MapClaims {
	//		if v, ok := data.(*User); ok {
	//			return jwt.MapClaims{
	//				identityKey: v.UserName,
	//			}
	//		}
	//		return jwt.MapClaims{}
	//	},
	//	IdentityHandler: func(c *gin.Context) interface{} {
	//		claims := jwt.ExtractClaims(c)
	//		return &User{
	//			UserName: claims[identityKey].(string),
	//		}
	//	},
	//	Authenticator: func(c *gin.Context) (interface{}, error) {
	//		var loginVals login
	//		if err := c.ShouldBind(&loginVals); err != nil {
	//			return "", jwt.ErrMissingLoginValues
	//		}
	//		userID := loginVals.Username
	//		password := loginVals.Password
	//
	//		if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
	//			return &User{
	//				UserName:  userID,
	//				LastName:  "Bo-Yi",
	//				FirstName: "Wu",
	//			}, nil
	//		}
	//
	//		return nil, jwt.ErrFailedAuthentication
	//	},
	//	Authorizator: func(data interface{}, c *gin.Context) bool {
	//		if v, ok := data.(*User); ok && v.UserName == "admin" {
	//			return true
	//		}
	//
	//		return false
	//	},
	//	Unauthorized: func(c *gin.Context, code int, message string) {
	//		c.JSON(code, gin.H{
	//			"code":    code,
	//			"message": message,
	//		})
	//	},
	//	// TokenLookup is a string in the form of "<source>:<name>" that is used
	//	// to extract token from the request.
	//	// Optional. Default value "header:Authorization".
	//	// Possible values:
	//	// - "header:<name>"
	//	// - "query:<name>"
	//	// - "cookie:<name>"
	//	// - "param:<name>"
	//	TokenLookup: "header: Authorization, query: token, cookie: jwt",
	//	// TokenLookup: "query:token",
	//	// TokenLookup: "cookie:token",
	//
	//	// TokenHeadName is a string in the header. Default value is "Bearer"
	//	TokenHeadName: "Bearer",
	//
	//	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	//	TimeFunc: time.Now,
	//})
	//
	//if err != nil {
	//	time.Sleep(3 * time.Second)
	//	log.Fatal("JWT Error:" + err.Error())
	//}
	//
	//// When you use jwt.New(), the function is already automatically called for checking,
	//// which means you don't need to call it again.
	//errInit := authMiddleware.MiddlewareInit()
	//
	//if errInit != nil {
	//	time.Sleep(3 * time.Second)
	//	log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	//}
	//
	//engine.POST("/login", authMiddleware.LoginHandler)
	//
	//engine.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
	//	claims := jwt.ExtractClaims(c)
	//	log.Printf("NoRoute claims: %#v\n", claims)
	//	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	//})
	//
	//auth := engine.Group("/auth")
	//// Refresh time can be longer than token timeout
	//auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	//auth.Use(authMiddleware.MiddlewareFunc())
	//{
	//	auth.GET("/hello", helloHandler)
	//}

	//1、setStaticFiles
	setStaticFiles(engine)
	//2、setWebAPI
	setWebAPI(engine)
	//3、setPort
	setPort()
	//4、setWebpServer
	//setWebpServer(engine)
	//5、setFrpClient
	setFrpClient()
	//6、printCMDMessage
	printCMDMessage()
	//7、StartGinEngine 监听并启动web服务
	StartGinEngine(engine)
}

//1、设置静态文件
func setStaticFiles(engine *gin.Engine) {
	//使用自定义的模板引擎，命名为"template-data"，为了与VUE兼容，把左右分隔符自定义为 [[ ]]
	tmpl := template.Must(template.New("template-data").Delims("[[", "]]").Parse(TemplateString))
	//使用模板
	engine.SetHTMLTemplate(tmpl)
	if common.Config.LogToFile {
		// 关闭 log 打印的字体颜色。输出到文件不需要颜色
		//gin.DisableConsoleColor()
		// 中间件，输出 log 到文件
		engine.Use(tools.LoggerToFile(common.Config.LogFilePath, common.Config.LogFileName))
		//禁止控制台输出
		gin.DefaultWriter = ioutil.Discard
	}
	//自定义分隔符，避免与vue.js冲突
	engine.Delims("[[", "]]")
	//https://stackoverflow.com/questions/66248258/serve-embedded-filesystem-from-root-path-of-url
	assetsEmbedFS, err := fs.Sub(staticAssetFS, "static/assets")
	if err != nil {
		fmt.Println(err)
	}
	engine.StaticFS("/assets/", http.FS(assetsEmbedFS))
	imagesEmbedFS, errStaticImageFS := fs.Sub(staticImageFS, "static/images")
	if errStaticImageFS != nil {
		fmt.Println(errStaticImageFS)
	}
	engine.StaticFS("/images/", http.FS(imagesEmbedFS))
	//单独一张静态图片
	//singleStaticFiles(engine, "/favicon.ico", "static/images/favicon.ico", "image/x-icon")
	engine.GET("/favicon.ico", func(c *gin.Context) {
		file, _ := staticFS.ReadFile("static/images/favicon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})
	//解析模板到HTML
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "template-data", gin.H{
			"title": "Comigo 漫画阅读器 " + common.Version, //页面标题
		})
	})
}

type Login struct {
	User     string `form:"username" json:"user" uri:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}

//2、设置获取书籍信息、图片文件的 API
func setWebAPI(engine *gin.Engine) {
	////TODO：处理登陆 https://www.chaindesk.cn/witbook/19/329
	////TODO：实现第三方认证，可参考 https://darjun.github.io/2021/07/26/godailylib/goth/
	//engine.POST("/login", func(c *gin.Context) {
	//	RememberPassword := c.DefaultPostForm("RememberPassword", "true") //可设置默认值
	//	username := c.PostForm("username")
	//	password := c.PostForm("password")
	//	//bookList := c.PostFormMap("book_list")
	//	//bookList := c.QueryArray("book_list")
	//	bookList := c.PostFormArray("book_list")
	//	c.String(http.StatusOK, fmt.Sprintf("RememberPassword is %s, username is %s, password is %s,hobby is %v", RememberPassword, username, password, bookList))
	//})

	//1.binding JSON
	// Example for binding JSON ({"user": "admin", "password": "comigo"})
	engine.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		//其实就是将request中的Body中的数据按照JSON格式解析到json变量中
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if json.User != "admin" || json.Password != "comigo" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// 简单的路由组: api,方便管理部分相同的URL
	var api *gin.RouterGroup
	//简单http认证
	enableAuth := common.Config.UserName != "" && common.Config.Password != ""
	if enableAuth {
		// 路由组：https://learnku.com/docs/gin-gonic/1.7/examples-grouping-routes/11399
		//使用 BasicAuth 中间件  https://learnku.com/docs/gin-gonic/1.7/examples-using-basicauth-middleware/11377
		api = engine.Group("/api", gin.BasicAuth(gin.Accounts{
			common.Config.UserName: common.Config.Password,
		}))
	} else {
		api = engine.Group("/api")
	}

	//处理表单 https://www.chaindesk.cn/witbook/19/329
	api.POST("/form", func(c *gin.Context) {
		t := c.DefaultPostForm("template", "scroll") //可设置默认值
		username := c.PostForm("username")
		password := c.PostForm("password")

		//bookList := c.PostFormMap("book_list")
		//bookList := c.QueryArray("book_list")
		bookList := c.PostFormArray("book_list")
		c.String(http.StatusOK, fmt.Sprintf("template is %s, username is %s, password is %s,hobby is %v", t, username, password, bookList))
	})

	//文件上传
	// 除了设置头像以外，也可以做上传文件并阅读功能
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// https://github.com/gin-gonic/examples/blob/master/upload-file/single/main.go
	// 也能上传多个文件，示例：
	//https://github.com/gin-gonic/examples/blob/master/upload-file/multiple/main.go
	//engine.MaxMultipartMemory = 60 << 20  // 60 MiB  只限制程序在上传文件时可以使用多少内存，而不限制上传文件的大小。
	api.POST("/upload", func(c *gin.Context) {
		// single file
		file, err := c.FormFile("file")
		if err != nil { //没有传文件会报错，处理这个错误。
			fmt.Println(err)
		}
		log.Println(file.Filename)

		// Upload the file to specific dst.
		err = c.SaveUploadedFile(file, file.Filename)
		if err != nil {
			fmt.Println(err)
		}
		/*
		   也可以直接使用io操作，拷贝文件数据。
		   out, err := os.Create(filename)
		   defer out.Close()
		   _, err = io.Copy(out, file)
		*/
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	//web端需要的服务器状态，包括标题、机器状态等
	api.GET("/getstatus", handler.ServerStatusHandler)
	//获取书架信息，不包含每页信息
	api.GET("/getlist", handler.GetBookListHandler)
	//通过URL字符串参数查询书籍信息
	api.GET("/getbook", handler.GetBookHandler)
	//通过URL字符串参数获取特定文件
	api.GET("/getfile", handler.GetFileHandler)
	////通过URL字符串参数PDF文件里的图片，效率太低，注释掉
	//api.GET("/get_pdf_image", handler.GetPdfImageHandler)
	//通过链接下载示例配置
	api.GET("/config.toml", handler.GetConfigHandler)
	//通过链接下载示例配置
	api.GET("/qrcode.png", handler.GetQrcodeHandler)
	//301重定向跳转示例
	api.GET("/redirect", func(c *gin.Context) {
		//支持内部和外部的重定向
		c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
	})
	//初始化websocket
	api.GET("/ws", handler.WsHandler)

	//TODO：设定压缩包下载链接
	// panic: handlers are already registered for path
	if book.GetBooksNumber() >= 1 {
		allBook, err := book.GetAllBookInfoList("name")
		if err != nil {
			fmt.Println("设置文件下载失败")
		} else {
			for _, info := range allBook.BookInfos {
				//下载文件
				if info.Type != book.TypeBooksGroup && info.Type != book.TypeDir {
					api.StaticFile("/raw/"+info.BookID+"/"+info.Name, info.FilePath)
				}
			}
		}
	}

}

//3、选择服务端口
func setPort() {
	//检测端口
	if !tools.CheckPort(common.Config.Port) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if common.Config.Port+2000 > 65535 {
			common.Config.Port = common.Config.Port + r.Intn(1024)
		} else {
			common.Config.Port = 30000 + r.Intn(20000)
		}
		fmt.Println(locale.GetString("port_busy") + strconv.Itoa(common.Config.Port))
	}
}

//5、setFrpClient
func setFrpClient() {
	//frp服务
	if !common.Config.EnableFrpcServer {
		return
	}
	if common.Config.FrpConfig.RemotePort <= 0 || common.Config.FrpConfig.RemotePort > 65535 {
		common.Config.FrpConfig.RemotePort = common.Config.Port
	}
	if common.Config.FrpConfig.RandomRemotePort {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		common.Config.FrpConfig.RemotePort = 50000 + r.Intn(10000)
	}
	frpcError := plugin.StartFrpC(common.Config.CacheFilePath)
	if frpcError != nil {
		fmt.Println(locale.GetString("frpc_server_error"), frpcError.Error())
	} else {
		fmt.Println(locale.GetString("frpc_server_start"))
	}

}

//6、printCMDMessage
func printCMDMessage() {
	//cmd打印链接二维码.如果只有一本书，就直接打开那本书.
	etcStr := ""
	//只有一本书的时候，URL需要附加的参数
	if book.GetBooksNumber() == 1 {
		bookList, err := book.GetAllBookInfoList("name")
		if err != nil {
			fmt.Println(err)
		}
		if len(bookList.BookInfos) == 1 {
			etcStr = "/#/scroll/" + bookList.BookInfos[0].BookID
		}
		if common.Config.DefaultMode != "" {
			etcStr = "/#/" + strings.ToLower(common.Config.DefaultMode) + "/" + bookList.BookInfos[0].BookID
		}
	}
	enableTls := common.Config.CertFile != "" && common.Config.KeyFile != ""
	tools.PrintAllReaderURL(common.Config.Port, common.Config.OpenBrowser, common.Config.EnableFrpcServer, common.Config.PrintAllIP, common.Config.Host, common.Config.FrpConfig.ServerAddr, common.Config.FrpConfig.RemotePort, common.Config.DisableLAN, enableTls, etcStr)
	//打印配置，调试用
	if common.Config.Debug {
		litter.Dump(common.Config)
	}
	fmt.Println(locale.GetString("ctrl_c_hint"))
}

// StartGinEngine 7、启动网页服务
func StartGinEngine(engine *gin.Engine) {
	//是否对外服务
	webHost := ":"
	if common.Config.DisableLAN {
		webHost = "localhost:"
	}
	//是否启用TLS
	enableTls := common.Config.CertFile != "" && common.Config.KeyFile != ""
	common.Srv = &http.Server{
		Addr:    webHost + strconv.Itoa(common.Config.Port),
		Handler: engine,
	}

	//在 goroutine 中初始化服务器，这样它就不会阻塞下面的正常关闭处理
	go func() {
		// 监听并启动服务(TLS)
		if enableTls {
			if err := common.Srv.ListenAndServeTLS(common.Config.CertFile, common.Config.KeyFile); err != nil && err != http.ErrServerClosed {
				time.Sleep(3 * time.Second)
				log.Fatalf("listen: %s\n", err)
			}
		}
		if !enableTls {
			// 监听并启动服务(HTTP)
			if err := common.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				time.Sleep(3 * time.Second)
				log.Fatalf("listen: %s\n", err)
			}
		}
	}()
}

//SetShutdownHandler TODO:退出时清理临时文件的函数
func SetShutdownHandler() {
	//优雅地停止或重启： https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	// 创建侦听来自操作系统的中断信号的上下文。
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()
	// Listen for the interrupt signal.
	//监听中断信号。
	<-ctx.Done()
	//恢复中断信号的默认行为并通知用户关机。
	stop()
	log.Println("shutting down processing, press Ctrl+C again to force")
	//清理临时文件
	if common.Config.CacheFileClean {
		fmt.Println("\r" + locale.GetString("start_clear_file") + " CacheFilePath:" + common.Config.CacheFilePath)
		book.ClearTempFilesALL(common.Config.Debug, common.Config.CacheFilePath)
		fmt.Println(locale.GetString("clear_temp_file_completed"))
	}
	// 上下文用于通知服务器它有 5 秒的时间来完成它当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := common.Srv.Shutdown(ctx); err != nil {
		//fmt.Println("Comigo Server forced to shutdown: ", err)
		//time.Sleep(3 * time.Second)
		log.Fatal("Comigo Server forced to shutdown: ", err)
	}
	log.Println("Comigo Server exit.")
}

////4、setWebpServer TODO：新的webp模式：https://docs.webp.sh/usage/remote-backend/
//func setWebpServer(engine *gin.Engine) {
//	//webp反向代理
//	if common.Config.EnableWebpServer {
//		webpError := common.StartWebPServer(common.CacheFilePath+"/webp_config.json", common.ReadingBook.ExtractPath, common.CacheFilePath+"/webp", common.Config.Port+1)
//		if webpError != nil {
//			fmt.Println(locale.GetString("webp_server_error"), webpError.Error())
//			//engine.Static("/cache", common.CacheFilePath)
//
//		} else {
//			fmt.Println(locale.GetString("webp_server_start"))
//			engine.Use(reverse_proxy.ReverseProxyHandle("/cache", reverse_proxy.ReverseProxyOptions{
//				TargetHost:  "http://localhost",
//				TargetPort:  strconv.Itoa(common.Config.Port + 1),
//				RewritePath: "/cache",
//			}))
//		}
//	} else {
//		if common.ReadingBook.IsDir {
//			engine.Static("/cache/"+common.ReadingBook.BookID, common.ReadingBook.GetFilePath())
//		} else {
//			engine.Static("/cache", common.CacheFilePath)
//		}
//	}
//}

//// 静态文件服务 单独设定某个文件
//func singleStaticFiles(engine *gin.Engine, fileUrl string, filePath string, contentType string) {
//	engine.GET(fileUrl, func(c *gin.Context) {
//		file, _ := staticFS.ReadFile(filePath)
//		c.Data(
//			http.StatusOK,
//			contentType,
//			file,
//		)
//	})
//}

//// getFileApi正常运作，虚拟文件系统实现方式
//func set-archiverFileSystem(engine *gin.Engine) {
////使用虚拟文件系统，设置服务路径（每本书都设置一遍）
////参考了: https://bitfieldconsulting.com/golang/filesystems
//for _, book := range common.BookList {
//	if book.NonUTF8Zip {
//		continue
//	}
//	ext := path.Ext(book.GetFilePath())
//	if (ext == ".zip" || ext == ".epub" || ext == ".cbz") && !book.NonUTF8Zip {
//		fsys, zipErr := zip.OpenReader(book.GetFilePath())
//		if zipErr != nil {
//			fmt.Println(zipErr)
//		}
//		httpFS := http.FS(fsys)
//		if book.IsDir {
//			engine.Static("/cache/"+book.BookID, book.GetFilePath())
//		} else {
//			engine.StaticFS("/cache/"+book.BookID, httpFS)
//		}
//	} else {
//		// 通过archiver/v4，建立虚拟FS。非UTF zip文件有编码问题
//		fsys, err := archiver.FileSystem(book.GetFilePath())
//		httpFS := http.FS(fsys)
//		if err != nil {
//			fmt.Println(err)
//		}
//		if book.IsDir {
//			engine.Static("/cache/"+book.BookID, book.GetFilePath())
//		} else {
//			engine.StaticFS("/cache/"+book.BookID, httpFS)
//		}
//	}
//}
//}
