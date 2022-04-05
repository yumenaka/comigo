package common

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir" //不使用 cgo 获取用户主目录的第三方库，支持交叉编译
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/tools"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}
	Config.LogFilePath = home
	Config.LogFileName = "comigo.log"
	slcBooks = make([]*Book, 0, 10) //make:为slice, map, channel分配内存，并返回一个初始化的值,第二参数指定的是切片的长度，第三个参数是用来指定预留的空间长度——避免二次分配内存带来的开销，提高程序的性能.
	mapBooks = make(map[string]*Book)
	mapBookGroups = make(map[string]*Book)
}

var (
	ConfigFile    = ""
	Version       = "v0.7.2"
	ReadingBook   *Book
	slcBooks      []*Book
	Srv           *http.Server
	mapBooks      map[string]*Book //实际存在的书
	mapBookGroups map[string]*Book //通过分析，生成的书籍分组
	Config        = ServerSettings{
		Port:                 1234,
		Host:                 "",
		StoresPath:           []string{},
		CacheFileEnable:      true,
		CacheFilePath:        "",
		CacheFileClean:       true,
		OpenBrowser:          true,
		DisableLAN:           false,
		GenerateMetaData:     false,
		LogToFile:            false,
		MaxDepth:             3,
		MinImageNum:          3,
		ZipFileTextEncoding:  "",
		EnableFrpcServer:     false,
		SupportFileType:      []string{".zip", ".tar", ".rar", ".cbr", ".cbz", ".epub", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz", ".tar.lz4", ".tlz4", ".tar.sz", ".tsz", ".bz2", ".gz", ".lz4", ".sz", ".xz", ".pdf", ".mp4", ".webm"},
		SupportMediaType:     []string{".jpg", ".jpeg", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".bmp", ".webp", ".ico", ".heic", ".avif"},
		ExcludeFileOrFolders: []string{".idea", ".vscode", ".git", "node_modules", "flutter_ui", ".local/share/Trash", "$RECYCLE.BIN", "Config.Msi", "System Volume Information", ".sys", " .DS_Store", ".dll", ".log", ".cache", ".exe"},
		Stores: Bookstores{
			mapBookstores: make(map[string]*singleBookstore),
			SortBy:        "name",
		},
		FrpConfig: FrpClientConfig{
			FrpcCommand:      "frpc",
			ServerAddr:       "localhost", //server_addr
			ServerPort:       7000,        //server_port
			Token:            "&&%%!2356",
			FrpType:          "tcp",
			RemotePort:       50000, //remote_port
			RandomRemotePort: true,
		},
		//WebpConfig: WebPServerConfig{
		//	WebpCommand:  "webp-server",
		//	HOST:         "127.0.0.1",
		//	PORT:         "3333",
		//	ImgPath:      "",
		//	QUALITY:      70,
		//	AllowedTypes: []string{".jpg", ".jpeg", ".JPEG", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".bmp"},
		//	ExhaustPath:  "",
		//},
		//Template:             "scroll", //multi、single、random etc.
	}
)

// CheckPathSkip 检查路径是否应该跳过（排除文件，文件夹列表）。
func CheckPathSkip(path string) bool {
	for _, substr := range Config.ExcludeFileOrFolders {
		if strings.HasSuffix(path, substr) {
			return true
		}
	}
	return false
}

// ServerSettings 服务器设置(config.toml)，配置文件放在被扫描的根目录中，或$HOME/config/comigo.可以用“comi --generate-config”生成本示例文件
type ServerSettings struct {
	Port                 int             `json:"port" comment:"Comigo设置(config.toml)，放在执行目录中，或$HOME/.config/comigo/下。可用“comi --generate-config”生成本文件\n# 网页端口"`
	Host                 string          `json:"host" comment:"自定义二维码显示的主机名"`
	StoresPath           []string        `json:"-"    comment:"设置默认扫描的书库文件夹"`
	MaxDepth             int             `json:"-" comment:"最大扫描深度"`
	OpenBrowser          bool            `json:"-" comment:"是否同时打开浏览器，windows默认true，其他默认false"`
	DisableLAN           bool            `json:"-" comment:"只在本机localhost提供服务，不对外共享"`
	UserName             string          `json:"-" comment:"访问限制：用户名。需要设置密码"`
	Password             string          `json:"-" comment:"访问限制：密码。需要设置用户名。"`
	CertFile             string          `json:"-" comment:"Https证书，同时设置KeyFile则启用HTTPS协议"`
	KeyFile              string          `json:"-" comment:"Https证书，同时设置CertFile则启用HTTPS协议"`
	CacheFileEnable      bool            `json:"-" comment:"是否保存web图片缓存，可以加快二次读取，但会占用硬盘空间"`
	CacheFilePath        string          `json:"-" comment:"web图片缓存存储位置，默认系统临时文件夹"`
	CacheFileClean       bool            `json:"-" comment:"退出程序的时候，清理web图片缓存"`
	ExcludeFileOrFolders []string        `json:"-" comment:"需要排除的文件或文件夹"`
	SupportMediaType     []string        `json:"-" comment:"需要扫描的图片文件后缀"`
	SupportFileType      []string        `json:"-" comment:"需要扫描的图书文件后缀"`
	MinImageNum          int             `json:"-" comment:"压缩包或文件夹内，至少有几张图片，才算作书籍"`
	PrintAllIP           bool            `json:"-" comment:"打印所有可能阅读链接的二维码"`
	Debug                bool            `json:"-" comment:"开启Debug模式"`
	LogToFile            bool            `json:"-" comment:"记录Log到本地文件"`
	LogFilePath          string          `json:"-" comment:"Log保存的位置"`
	LogFileName          string          `json:"-" comment:"Log文件名"`
	ZipFileTextEncoding  string          `json:"-" comment:"非utf-8编码的ZIP文件，尝试用什么编码解析，默认GBK"`
	GenerateConfig       bool            `toml:"-" comment:"生成示例配置文件的标志"`
	EnableFrpcServer     bool            `json:"-" comment:"后台启动FrpClient"`
	FrpConfig            FrpClientConfig `json:"-" comment:"FrpClient设置"`
	GenerateMetaData     bool            `json:"-" toml:"-" comment:"生成书籍元数据（TODO）"`
	Stores               Bookstores      `json:"stores" toml:"-"`
	//EnableWebpServer       bool             `json:"enable_webp_server"`
	//SketchCountSeconds     int              `json:"sketch_count_seconds"`
	//WebpConfig             WebPServerConfig `json:"-"  comment:" WebPServer设置"`
	//Template               string           `json:"-"`
}

//WebPServerConfig  WebPServer服务端配置
type WebPServerConfig struct {
	WebpCommand  string
	HOST         string
	PORT         string
	ImgPath      string `json:"IMG_PATH"`
	QUALITY      int
	AllowedTypes []string `json:"ALLOWED_TYPES"`
	ExhaustPath  string   `json:"EXHAUST_PATH"`
}

//FrpClientConfig frp客户端配置
type FrpClientConfig struct {
	FrpcCommand      string `comment:"手动设定frpc可执行程序的路径,默认为frpc"`
	ServerAddr       string
	ServerPort       int
	Token            string
	FrpType          string //本地转发端口设置
	RemotePort       int
	RandomRemotePort bool
}

// SetupCloseHander 中断处理：程序被中断的时候，清理临时文件
// 在一个新的 goroutine 上创建一个监听器。 如果接收到了一个 interrupt 信号，就会立即通知程序，做一些清理工作并退出
func SetupCloseHander() {
	//中断处理：程序被中断的时候，清理临时文件
	//容量2(capacity)代表Channel容纳的最多的元素的数量，代表Channel的缓存的大小。如果设置了缓存，就有可能不发生阻塞， 只有buffer满了后 send才会阻塞， 而只有缓存空了后receive才会阻塞。
	c := make(chan os.Signal, 2)
	//SIGHUP（挂起）, SIGINT（中断）或 SIGTERM（终止）默认会使得程序退出。
	//1、SIGHUP 信号在用户终端连接(正常或非正常)结束时发出。
	//2、syscall.SIGINT 和 os.Interrupt 是同义词,按下 CTRL+C 时发出。
	//3、SIGTERM（终止）:kill终止进程,允许程序处理问题后退出。
	//4.syscall.SIGHUP,终端控制进程结束(终端连接断开)
	//5、syscall.SIGQUIT，CTRL+\ 退出

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	//signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		var temp = <-c
		fmt.Println(temp)
		if Config.CacheFileClean {
			fmt.Println("\r" + locale.GetString("start_clear_file"))
			ClearTempFilesALL()
		}
		os.Exit(0)
	}()
}

// SetTempDir 设置临时文件夹，退出时会被清理
func SetTempDir() {
	//手动设置的临时文件夹
	if Config.CacheFilePath != "" && tools.ChickExists(Config.CacheFilePath) && tools.ChickIsDir(Config.CacheFilePath) {
		Config.CacheFilePath = path.Join(Config.CacheFilePath)
	} else {
		Config.CacheFilePath = path.Join(os.TempDir(), "comigo_cache") //直接使用系统文件夹
	}
	err := os.MkdirAll(Config.CacheFilePath, os.ModePerm)
	if err != nil {
		println(locale.GetString("temp_folder_error"))
	} else {
		fmt.Println(locale.GetString("temp_folder_path") + Config.CacheFilePath)
	}
}

// ClearTempFilesALL web加载时保存的临时图片，在在退出后清理
func ClearTempFilesALL() {
	//fmt.Println(locale.GetString("clear_temp_file_start"))
	for _, tempBook := range mapBooks {
		clearTempFilesOne(tempBook)
	}
}

// 清空某一本压缩漫画的解压缓存
func clearTempFilesOne(book *Book) {
	//fmt.Println(locale.GetString("clear_temp_file_start"))
	haveThisBook := false
	for _, tempBook := range mapBooks {
		if tempBook.GetBookID() == book.GetBookID() {
			haveThisBook = true
		}
	}
	if haveThisBook {
		cachePath := path.Join(Config.CacheFilePath, book.GetBookID())
		err := os.RemoveAll(cachePath)
		if err != nil {
			fmt.Println(locale.GetString("clear_temp_file_error") + cachePath)
		} else {
			if Config.Debug {
				fmt.Println(locale.GetString("clear_temp_file_completed") + cachePath)
			}
		}
	}
}

// SetByExecutableFilename 通过执行文件名设置默认网页模板参数
func (config *ServerSettings) SetByExecutableFilename() {
	// 当前执行目录
	//targetPath, _ := os.Getwd()
	//fmt.Println(locale.GetString("target_path"), targetPath)
	// 带后缀的执行文件名 comi.exe  sketch.exe
	filenameWithSuffix := path.Base(os.Args[0])
	// 执行文件名后缀
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
	ExtFileName := strings.TrimPrefix(filenameWithOutSuffix, extPath)
	fmt.Println("ExtFileName =", ExtFileName)
	////如果执行文件名包含 scroll 等关键字，选择卷轴模板
	//if haveKeyWord(ExtFileName, []string{"scroll", "スクロール", "默认", "下拉", "卷轴"}) {
	//	config.Template = "scroll"
	//}
	////如果执行文件名包含 sketch 等关键字，选择速写模板
	//if haveKeyWord(ExtFileName, []string{"sketch", "croquis", "クロッキー", "素描", "速写"}) {
	//	config.Template = "sketch"
	//}
	////根据文件名设定倒计时秒数,不管默认是不是sketch模式
	//Seconds, err := getNumberFromString(ExtFileName)
	//if err != nil {
	//	if config.Template == "sketch" {
	//		//fmt.Println(Seconds)
	//	}
	//} else {
	//	config.SketchCountSeconds = Seconds
	//}
	////如果执行文件名包含 single 等关键字，选择 flip分页漫画模板
	//if haveKeyWord(ExtFileName, []string{"flip", "翻页", "めく"}) {
	//	config.Template = "flip"
	//}
	////选择模式以后，打印提示
	//switch config.Template {
	//case "scroll":
	//	fmt.Println(locale.GetString("scroll_template"))
	//case "flip":
	//	fmt.Println(locale.GetString("single_page_template"))
	//case "sketch":
	//	fmt.Println(locale.GetString("sketch_template"))
	//	//速写倒计时秒数
	//	fmt.Println(locale.GetString("SKETCH_COUNT_SECONDS"), config.SketchCountSeconds)
	//default:
	//}
}

//从字符串中提取数字,如果有几个数字，就简单地加起来
func getNumberFromString(s string) (int, error) {
	var err error
	num := 0
	//同时设定倒计时秒数
	valid := regexp.MustCompile("[0-9]+")
	numbers := valid.FindAllStringSubmatch(s, -1)
	if len(numbers) > 0 {
		//循环取出多维数组
		for _, value := range numbers {
			for _, v := range value {
				temp, errTemp := strconv.Atoi(v)
				if errTemp != nil {
					fmt.Println("error num value:" + v)
				} else {
					num = num + temp
				}
			}
		}
		//fmt.Println("get Number:",num," form string:",s,"numbers[]=",numbers)
	} else {
		err = errors.New("number not found")
		return 0, err
	}
	return num, err
}

//检测字符串中是否有关键字
func haveKeyWord(checkString string, list []string) bool {
	//转换为小写，使Sketch、DOUBLE也生效
	checkString = strings.ToLower(checkString)
	for _, key := range list {
		if strings.Contains(checkString, key) {
			return true
		}
	}
	return false
}
