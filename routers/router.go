package routers

import (
	"github.com/gin-gonic/gin"
)

// gin-jwt相关 https://github.com/appleboy/gin-jwt

// 声明swagHandler，该参数不为空时才加入路由，以减少包体积
// 通过go build -tags "doc"来打包带文档的包，直接go build打包不带文档的包
var swagHandler gin.HandlerFunc

// 简单的路由组: api,方便管理部分相同的URL
var api *gin.RouterGroup
var protectedAPI *gin.RouterGroup

// StartWebServer 启动web服务
func StartWebServer() {
	//设置 gin
	gin.SetMode(gin.ReleaseMode)
	//使用 gin.New() 而不是 gin.Default() 以避免使用 Gin 的默认日志中间件
	engine := gin.New()
	////Logger() 以默认配置创建日志中间件，将所有请求信息按指定格式打印到标准输出。 gin.Default()默认启用这个中间件
	//engine.Use(gin.Logger())
	//使用 Recovery 中间件，避免程序崩溃，返回 500 错误页面
	engine.Use(gin.Recovery())

	//1、embedFile
	embedFile(engine)
	//2、setWebAPI
	setWebAPI(engine)
	//TODO：Go中调用外部命令的几种姿势 https://darjun.github.io/2022/11/01/godailylib/osexec/
	//3、showQRCode
	showQRCode()
	//4、startEngine 监听并启动web服务
	startEngine(engine)
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
//			logger.Info(zipErr)
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
//			logger.Info(err)
//		}
//		if book.IsDir {
//			engine.Static("/cache/"+book.BookID, book.GetFilePath())
//		} else {
//			engine.StaticFS("/cache/"+book.BookID, httpFS)
//		}
//	}
//}
//}
