package routers

import (
	"embed"
	"github.com/gin-gonic/gin"
)

// TemplateString 模板文件
//
//go:embed static/index.html
var TemplateString string

//go:embed  static
var staticFS embed.FS

//go:embed  static/assets
var staticAssetFS embed.FS

//go:embed  static/images
var staticImageFS embed.FS

// gin-jwt相关 https://github.com/appleboy/gin-jwt

// 声明swagHandler，该参数不为空时才加入路由，以减少包体积
// 通过go build -tags "doc"来打包带文档的包，直接go build打包不带文档的包
var swagHandler gin.HandlerFunc

// 简单的路由组: api,方便管理部分相同的URL
var api *gin.RouterGroup

// StartWebServer 启动web服务
func StartWebServer() {
	//设置 gin
	gin.SetMode(gin.ReleaseMode)
	//// 创建带有默认中间件的路由: 日志与恢复中间件
	engine := gin.Default()
	////Recovery 中间件会恢复(recovers) 任何恐慌(panics) 如果存在恐慌，中间件将会写入500 gin.Default()已经默认启用了这个中间件
	//engine.Use(gin.Recovery())
	////Logger() 以默认配置创建日志中间件，将所有请求信息按指定格式打印到标准输出。 gin.Default()已经默认启用了这个中间件
	//engine.Use(gin.Logger())

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
	//6、printQRCodeInCMD
	printQRCodeInCMD()
	//7、StartGinEngine 监听并启动web服务
	StartGinEngine(engine)
}

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
