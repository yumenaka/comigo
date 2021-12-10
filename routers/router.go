package routers

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanity-io/litter"
	"github.com/spf13/viper"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/routers/reverse_proxy"
	"github.com/yumenaka/comi/tools"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	yaml "github.com/goccy/go-yaml"
)

//go:embed index.html
var TemplateString string

//go:embed  favicon.ico js/* css/*
var EmbedFiles embed.FS

//退出时清理
func init() {
	common.SetupCloseHander()
}

// ParseCommands 解析命令
func ParseCommands(args []string) {

	//保存配置並退出
	if common.Config.NewConfig {
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println(err)
			return
		} // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
		err = viper.SafeWriteConfigAs("D:\\cvgo")
		if err != nil {
			return
		}
		err = viper.WriteConfigAs("D:\\cvgo")
		if err != nil {
			return
		}
		err = viper.SafeWriteConfigAs("D:\\cvgo")
		if err != nil {
			return
		} // will error since it has already been written
		err = viper.SafeWriteConfigAs("D:\\cvgo")
		if err != nil {
			return
		}

		bytes, err := yaml.Marshal(common.Config)
		if err != nil {
			fmt.Println("yaml.Marshal Error")
		}
		fmt.Println(string(bytes)) // "a: 1\nb: hello\n"
		err = ioutil.WriteFile("test.yaml", bytes, 0644)
		if err != nil {
			panic(err)
		}
		//		cfg := `#定义一个yaml配置文件
		//OpenBrowser:         false
		//DisableLAN:          false
		//Port:                1234
		//CheckImage:  true
		//LogToFile:           false
		//MaxDepth:            2
		//MinImageNum:         3
		//ZipFilenameEncoding: ""
		//EnableWebpServer: false
		//Host: "localhost"
		//`
		//		data := []byte(cfg)
		//		v := make(map[string]interface{})
		//		err := yaml.Unmarshal(data, v)
		//		if err != nil {
		//			fmt.Println("failed to unmarshal YAML")
		//		}
		os.Exit(0)
	}

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
	var wg sync.WaitGroup
	if common.Config.CheckImage {
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
		gin.DisableConsoleColor()
		// 输出 log 到文件
		engine.Use(tools.LoggerToFile(common.Config.LogFilePath, common.Config.LogFileName))
	}
	//自定义分隔符，避免与vue.js冲突
	engine.Delims("[[", "]]")
	//网站图标
	engine.GET("/resources/favicon.ico", func(c *gin.Context) {
		file, _ := EmbedFiles.ReadFile("favicon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})
	engine.StaticFS("/assets", http.FS(EmbedFiles))
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
		webpError := common.StartWebPServer(common.WebImagePath, common.WebImagePath, common.ComigoCacheFilePath+"/webp", common.Config.Port+1)
		if webpError != nil {
			fmt.Println(locale.GetString("webp_server_error"), webpError.Error())
			engine.Static("/cache", common.WebImagePath)
		} else {
			fmt.Println(locale.GetString("webp_server_start"))
			engine.Use(reverse_proxy.ReverseProxyHandle("/cache", reverse_proxy.ReverseProxyOptions{
				TargetHost:  "http://localhost",
				TargetPort:  strconv.Itoa(common.Config.Port + 1),
				RewritePath: "/cache",
			}))
		}
	} else {
		//具体的图片文件
		engine.Static("/cache", common.WebImagePath)
		//直接建立一个zipfs，但非UTF文件，会出现编码问题，待改进
		//ext := path.Ext(common.ReadingBook.FilePath)
		//if ext == ".zip" {
		//	fsys, zip_err := zip.OpenReader(common.ReadingBook.FilePath)
		//	if zip_err != nil {
		//		fmt.Println(zip_err)
		//	}
		//	engine.StaticFS("/cache", http.FS(fsys))
		//} else {
		//	//图片目录
		//	engine.Static("/cache", common.WebImagePath)
		//}
		//大概需要自己实现一个rar fs？  https://github.com/forensicanalysis/zipfs
		//// Error:*rardecode.ReadCloser does not implement fs.FS (missing Open method)
		//fsys2, rar_err := rar.OpenReader("test.rar","")
		//if rar_err != nil {
		//	fmt.Println(rar_err)
		//}
		//engine.StaticFS("/rar", http.FS(fsys2))
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
		frpcError := common.StartFrpC(common.ComigoCacheFilePath)
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
	err := engine.Run(webHost + strconv.Itoa(common.Config.Port))
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, locale.GetString("web_server_error")+"%q\n", common.Config.Port)
		if err != nil {
			return
		}
	}
}
