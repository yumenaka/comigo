package routers

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/routers/reverse_proxy"
	"github.com/yumenaka/comi/tools"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

//退出时清理
func init() {
	common.SetupCloseHander()
}

func StartServer(args []string) {

	switch {
	case common.Config.Template =="auto" :
		selectTemplate()
	case common.Config.Template =="multi":
		fmt.Println(locale.GetString("multi_page_template"))
	case common.Config.Template =="random":
		fmt.Println(locale.GetString("single_page_template"))
	case common.Config.Template =="single":
		fmt.Println(locale.GetString("single_page_template"))
	default:
		common.Config.Template ="multi"
	}

	if common.Config.Template =="auto" {

	}else if common.Config.Template =="auto"  {

	}
	cmdPath := path.Dir(os.Args[0]) //去除路径最后一个元素  /home/dir/comigo.exe -> /home/dir/
	if len(args) == 0 {
		err := common.ScanBookPath(cmdPath)
		if err != nil {
			fmt.Println(locale.GetString("scan_error"), cmdPath)
		}
	} else {
		for _, p := range args {
			if p == cmdPath {
				continue //指定参数的话，就不扫描当前目录
			}
			err := common.ScanBookPath(p)
			if err != nil {
				fmt.Println(locale.GetString("scan_error"), p)
			}
		}
	}
	switch len(common.BookList) {
	case 0:
		fmt.Println(locale.GetString("book_not_found"))
		os.Exit(0)
	default:
		setFirstBook(args)
	}
	var wg sync.WaitGroup
	//解压图片，分析分辨率
	if common.Config.CheckImageInServer {
		wg.Add(1)
		go func() {
			common.InitReadingBook()
			defer wg.Done()
		}()
		wg.Wait()
	} else {
		err := common.InitReadingBook()
		if err != nil {
			fmt.Println(locale.GetString("can_not_init_book"), err, common.ReadingBook)
		}
	}
	InitWebServer()
}

func selectTemplate() {
	// 当前执行目录
	targetPath,_ := os.Getwd()
	fmt.Println(locale.GetString("target_path"), targetPath)
	// 带后缀的执行文件名
	filenameWithSuffix := path.Base(os.Args[0])
	// 文件后缀
	fileSuffix := path.Ext(filenameWithSuffix)
	// 去掉后缀后的执行文件名
	filenameWithOutSuffix := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	//fmt.Println("filenameWithOutSuffix =", filenameWithOutSuffix)
	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	extPath := filepath.Dir(ex)
	//fmt.Println("extPath =",extPath)
	ExtFileName:=  strings.TrimPrefix(filenameWithOutSuffix, extPath)
	//fmt.Println("ExtFileName =", ExtFileName)
	common.Config.SetTemplateByName(ExtFileName)
}

func setFirstBook(args []string) {
	if len(common.BookList) == 0 {
		return
	}
	//多本书，读第一本
	if len(args) == 0 {
		if len(common.BookList) > 0 {
			common.ReadingBook = common.BookList[0]
		}
	}
	if len(args) > 0 {
		for _, b := range common.BookList {
			if b.FilePath == args[0] {
				common.ReadingBook = b
				break
			}
		}
	}
}

// 测试用的假数据
var secrets = gin.H{
	"comi":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"admin": gin.H{"email": "austin@example.com", "phone": "666"},
	"user1":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

//启动web服务
func InitWebServer() {
	//go:embed index.html
	var TemplateString string
	//获取模板，命名为"template-data"，同时把左右分隔符改为 [[ ]]
	tmpl := template.Must(template.New("template-data").Delims("[[", "]]").Parse(TemplateString))
	//设置 gin
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	//使用模板
	engine.SetHTMLTemplate(tmpl)
	if common.Config.LogToFile {
		// 关闭 log 打印的字体颜色。输出到文件不需要颜色
		gin.DisableConsoleColor()
		// 输出 log 到文件(logrus)
		engine.Use(tools.LoggerToFile(common.Config.LogFilePath, common.Config.LogFileName))
	}
	//自定义分隔符，避免与vue.js冲突
	engine.Delims("[[", "]]")
	//go:embed  favicon.ico js/* css/*
	var EmbedFiles embed.FS
	//网站图标
	engine.GET("/resources/favicon.ico", func(c *gin.Context) {
		file, _ := EmbedFiles.ReadFile("favicon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})
	engine.StaticFS("/assets",http.FS(EmbedFiles))
	if common.ReadingBook.IsFolder {
		engine.StaticFS("/raw/"+common.ReadingBook.Name, gin.Dir(common.ReadingBook.FilePath, true))
	} else {
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

	if common.Config.Auth != ""{

	}

	//解析json
	engine.GET("/book.json", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, common.ReadingBook)
	})
	//解析书架json
	engine.GET("/bookshelf.json", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, common.BookList)
	})
	//服务器设定
	engine.GET("/setting.json", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, common.Config)
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
			common.Config.Port = 50000 + r.Intn(10000)
		}
		fmt.Println(locale.GetString("port_busy") + strconv.Itoa(common.Config.Port))
	}
	//webp反向代理
	if common.Config.EnableWebpServer {
		webpError := common.StartWebPServer(common.PictureDir, common.PictureDir, common.TempDir+"/webp", common.Config.Port+1)
		if webpError != nil {
			fmt.Println(locale.GetString("webp_server_error"), webpError.Error())
			engine.Static("/cache", common.PictureDir)
		} else {
			fmt.Println(locale.GetString("werp_server_start"))
			engine.Use(reverse_proxy.ReverseProxyHandle("/cache", reverse_proxy.ReverseProxyOptions{
				TargetHost:  "http://localhost",
				TargetPort:  strconv.Itoa(common.Config.Port + 1),
				RewritePath: "/cache",
			}))
		}
	} else {
		//图片目录
		engine.Static("/cache", common.PictureDir)
	}
	if common.Config.EnableFrpcServer {
		if common.Config.FrpConfig.RandomRemotePort {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			common.Config.FrpConfig.RemotePort = 40000 + r.Intn(10000)
		} else {
			if common.Config.FrpConfig.RemotePort <= 0 || common.Config.FrpConfig.RemotePort > 65535 {
				common.Config.FrpConfig.RemotePort = common.Config.Port
			}
		}
		frpcError := common.StartFrpC(common.TempDir)
		if frpcError != nil {
			fmt.Println(locale.GetString("frpc_server_error"), frpcError.Error())
		} else {
			fmt.Println(locale.GetString("frpc_server_start"))
		}
	}
	//开始服务
	tools.PrintAllReaderURL(common.Config.Port,common.Config.OpenBrowser,common.Config.EnableFrpcServer,common.Config.PrintAllIP,common.Config.ServerHost,common.Config.FrpConfig.ServerAddr,common.Config.FrpConfig.RemotePort,common.Config.DisableLAN)
	//打印配置
	//fmt.Println(locale.GetString("print_config"))
	fmt.Println(common.Config)
	err := engine.Run(webHost + strconv.Itoa(common.Config.Port))
	if err != nil {
		fmt.Fprintf(os.Stderr, locale.GetString("web_server_error")+"%q\n", common.Config.Port)
	}
}
