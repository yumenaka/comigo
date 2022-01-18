package routers

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanity-io/litter"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/routers/reverse_proxy"
	"github.com/yumenaka/comi/tools"
	"html/template"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

//go:embed static/index.html
var TemplateString string

//go:embed  static
var staticFS embed.FS

//go:embed  static/assets
var staticAssetFS embed.FS

//go:embed  static/images
var staticImageFS embed.FS

//退出时清理
func init() {
	common.SetupCloseHander()
}

// ParseCommands 解析命令
func ParseCommands(args []string) {
	//通过“可执行文件名”设置默认阅读模板
	common.Config.SetByExecutableFilename()
	//决定如何扫描，扫描哪个路径
	if len(args) == 0 { //没有指定路径或文件的情况下
		cmdPath := path.Dir(os.Args[0]) //当前执行路径
		err := common.ScanBookPath(cmdPath)
		if err != nil {
			fmt.Println(locale.GetString("scan_error"), cmdPath)
		}
		if len(common.BookList) > 0 {
			common.ReadingBook = common.BookList[0]
		}
	} else {
		//指定了多个参数的话，都扫描
		for _, p := range args {
			err := common.ScanBookPath(p)
			if err != nil {
				fmt.Println(locale.GetString("scan_error"), p)
			}
		}
	}
	//扫描完路径之后，选择第一本书开始解压
	switch len(common.BookList) {
	case 0:
		fmt.Println(locale.GetString("book_not_found"))
		os.Exit(0)
	default:
		common.ReadingBook = common.BookList[0]
	}
	//解压图片，分析分辨率（并发）
	if common.Config.CheckImage {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			err := common.InitReadingBook()
			if err != nil {
				return
			}
			defer wg.Done()
		}()
		wg.Wait()
	} else {
		err := common.InitReadingBook()
		if err != nil {
			fmt.Println(locale.GetString("can_not_init_book"), err, common.ReadingBook)
		}
	}
	StartWebServer()
}

//单独设定某个文件
func setStaticFiles(engine *gin.Engine, fileUrl string, filePath string, contentType string) {
	engine.GET(fileUrl, func(c *gin.Context) {
		file, _ := staticFS.ReadFile(filePath)
		c.Data(
			http.StatusOK,
			contentType,
			file,
		)
	})
}

// StartWebServer 启动web服务
func StartWebServer() {

	//获取模板，命名为"template-data"，同时把左右分隔符改为 [[ ]]
	tmpl := template.Must(template.New("template-data").Delims("[[", "]]").Parse(TemplateString))
	//设置 gin
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	//使用模板
	engine.SetHTMLTemplate(tmpl)
	if common.Config.LogToFile {
		// 关闭 log 打印的字体颜色。输出到文件不需要颜色
		//gin.DisableConsoleColor()
		// 输出 log 到文件
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
	//单独一张静态图片，没必要再定义一个FS
	setStaticFiles(engine, "/favicon.ico", "static/images/favicon.ico", "image/x-icon")

	//Download archive file
	if !common.ReadingBook.IsDir {
		engine.StaticFile("/raw/"+common.ReadingBook.Name, common.ReadingBook.FilePath)
	}
	//解析模板到HTML
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "template-data", gin.H{
			"title": common.ReadingBook.Name, //页面标题
		})
	})
	//简单http认证测试
	authorized := engine.Group("/", gin.BasicAuth(gin.Accounts{
		"comi": "go",
	}))
	authorized.GET("/secrets", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"secret": "这个路径需要认证。",
		})
	})
	//认证相关，还没写
	if common.Config.Auth != "" {

	}
	//解析json
	engine.GET("api/book.json", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, common.ReadingBook)
	})
	//解析书架json
	engine.GET("api/bookshelf.json", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, common.BookList)
	})
	//服务器设定
	engine.GET("api/setting.json", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, common.Config)
	})
	//服务器设定
	engine.GET("/config.yaml", func(c *gin.Context) {
		c.YAML(http.StatusOK, common.Config)
	})
	//初始化websocket
	engine.GET("/ws", wsHandler)
	//是否同时对外服务
	webHost := ":"
	if common.Config.DisableLAN {
		webHost = "localhost:"
	}
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
	//webp反向代理
	if common.Config.EnableWebpServer {
		webpError := common.StartWebPServer(common.CacheFilePath+"/webp_config.json", common.ReadingBook.ExtractPath, common.CacheFilePath+"/webp", common.Config.Port+1)
		if webpError != nil {
			fmt.Println(locale.GetString("webp_server_error"), webpError.Error())
			engine.Static("/cache", common.CacheFilePath)
		} else {
			fmt.Println(locale.GetString("webp_server_start"))
			engine.Use(reverse_proxy.ReverseProxyHandle("/cache", reverse_proxy.ReverseProxyOptions{
				TargetHost:  "http://localhost",
				TargetPort:  strconv.Itoa(common.Config.Port + 1),
				RewritePath: "/cache",
			}))
		}
	} else {
		if common.ReadingBook.IsDir {
			common.ReadingBook.SetBookID()
			engine.Static("/cache/"+common.ReadingBook.BookID, common.ReadingBook.FilePath)
		} else {
			engine.Static("/cache", common.CacheFilePath)
		}

		//具体的图片文件
		//直接建立一个zipfs，但非UTF文件有编码问题，待改进
		//ext := path.Ext(common.ReadingBook.FilePath)
		//if ext == ".zip" {
		//	fsys, zip_err := zip.OpenReader(common.ReadingBook.FilePath)
		//	if zip_err != nil {
		//		fmt.Println(zip_err)
		//	}
		//	engine.StaticFS("/cache", http.FS(fsys))
		//} else {
		//	//图片目录
		//	engine.Static("/cache", common.ExtractPath)
		//}
	}
	//cmd打印链接二维码
	tools.PrintAllReaderURL(common.Config.Port, common.Config.OpenBrowser, common.Config.EnableFrpcServer, common.Config.PrintAllIP, common.Config.Host, common.Config.FrpConfig.ServerAddr, common.Config.FrpConfig.RemotePort, common.Config.DisableLAN)
	//开始服务
	if common.Config.EnableFrpcServer {
		if common.Config.FrpConfig.RandomRemotePort {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			common.Config.FrpConfig.RemotePort = 50000 + r.Intn(10000)
		} else {
			if common.Config.FrpConfig.RemotePort <= 0 || common.Config.FrpConfig.RemotePort > 65535 {
				common.Config.FrpConfig.RemotePort = common.Config.Port
			}
		}
		frpcError := common.StartFrpC(common.CacheFilePath)
		if frpcError != nil {
			fmt.Println(locale.GetString("frpc_server_error"), frpcError.Error())
		} else {
			fmt.Println(locale.GetString("frpc_server_start"))
		}
	}
	//打印配置，调试用
	if common.Config.Debug {
		litter.Dump(common.Config)
	}
	fmt.Println(locale.GetString("ctrl_c_hint"))
	err = engine.Run(webHost + strconv.Itoa(common.Config.Port))
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, locale.GetString("web_server_error")+"%q\n", common.Config.Port)
		if err != nil {
			return
		}
	}
}
