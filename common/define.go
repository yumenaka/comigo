package common

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir" //不使用 cgo 获取用户主目录的第三方库，支持交叉编译
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/tools"
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
	//退出时清理
	setupCloseHander()
	slcBooks = make([]*Book, 0, 10) //make:为slice, map, channel分配内存，并返回一个初始化的值,第二参数指定的是切片的长度，第三个参数是用来指定预留的空间长度——避免二次分配内存带来的开销，提高程序的性能.
	mapBooks = make(map[string]*Book)
}

var (
	ConfigFile  = ""
	Version     = "v0.6.0"
	ReadingBook *Book
	slcBooks    []*Book
	mapBooks    map[string]*Book
	Config      = ServerSettings{
		OpenBrowser: true,
		DisableLAN:  false,
		Stores: Bookstores{
			mapBookstores: make(map[string]*singleBookstore),
			SortBy:        "name",
		},
		Port:                 1234,
		GenerateMetaData:     false,
		LogToFile:            false,
		MaxDepth:             3,
		MinImageNum:          3,
		ZipFileTextEncoding:  "",
		CacheFilePath:        "",
		SupportFileType:      []string{".zip", ".tar", ".rar", ".cbr", ".cbz", ".epub", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz", ".tar.lz4", ".tlz4", ".tar.sz", ".tsz", ".bz2", ".gz", ".lz4", ".sz", ".xz", ".pdf", ".mp4", ".webm"},
		SupportMediaType:     []string{".jpg", ".jpeg", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".bmp", ".webp", ".ico", ".heic", ".avif"},
		ExcludeFileOrFolders: []string{".idea", ".vscode", ".git", "node_modules", "flutter_ui", ".local/share/Trash", "$RECYCLE.BIN", "Config.Msi", "System Volume Information", ".sys", " .DS_Store", ".dll", ".log", ".cache", ".exe"},
		WebpConfig: WebPServerConfig{
			WebpCommand:  "webp-server",
			HOST:         "127.0.0.1",
			PORT:         "3333",
			ImgPath:      "",
			QUALITY:      70,
			AllowedTypes: []string{".jpg", ".jpeg", ".JPEG", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".bmp"},
			ExhaustPath:  "",
		},
		EnableFrpcServer: false,
		FrpConfig: FrpClientConfig{
			FrpcCommand:      "frpc",
			ServerAddr:       "localhost", //server_addr
			ServerPort:       7000,        //server_port
			Token:            "&&%%!2356",
			FrpType:          "tcp",
			RemotePort:       50000, //remote_port
			RandomRemotePort: true,
		},
		Host:                   "",
		SketchCountSeconds:     90,
		TempPATH:               "",
		CleanAllTempFileOnExit: true,
		CleanAllTempFile:       true,
		//Template:             "scroll", //multi、single、random etc.
	}
)

// ServerStatus 服务器当前状况
type ServerStatus struct {
	//当前拥有的书籍总数
	NumberOfBooks int
	//在线用户数
	NumberOfOnLineUser int
	//在线设备数
	NumberOfOnLineDevices int
	OSInfo                tools.SystemStatus
}

// CheckPathSkip 检查路径是否应该跳过（排除文件，文件夹列表）。
func CheckPathSkip(path string) bool {
	for _, substr := range Config.ExcludeFileOrFolders {
		if strings.HasSuffix(path, substr) {
			return true
		}
	}
	return false
}

// GetServerStatus 获取服务器的状态
func GetServerStatus() *ServerStatus {
	return &ServerStatus{
		NumberOfBooks:         len(mapBooks),
		NumberOfOnLineUser:    1,
		NumberOfOnLineDevices: 1,
		OSInfo:                tools.GetSystemStatus(),
	}
}

type ServerSettings struct {
	Host                   string           `json:"host"`
	EnableWebpServer       bool             `json:"enable_webp_server"`
	EnableFrpcServer       bool             `json:"frpc_enable"`
	Port                   int              `json:"port"`
	SketchCountSeconds     int              `json:"sketch_count_seconds"`
	Stores                 Bookstores       `json:"stores"` //这个字段不解析
	CacheFilePath          string           `json:"-"`      //这个字段不解析
	ExcludeFileOrFolders   []string         `json:"-"`      //这个字段不解析
	SupportMediaType       []string         `json:"-"`      //这个字段不解析
	SupportFileType        []string         `json:"-"`      //这个字段不解析
	MinImageNum            int              `json:"-"`      //这个字段不解析
	GenerateMetaData       bool             `json:"-"`      //这个字段不解析
	UserName               string           `json:"-"`      //这个字段不解析
	Password               string           `json:"-"`      //这个字段不解析
	CertFile               string           `json:"-"`      //这个字段不解析
	KeyFile                string           `json:"-"`      //这个字段不解析
	OpenBrowser            bool             `json:"-"`      //这个字段不解析
	DisableLAN             bool             `json:"-"`      //这个字段不解析
	PrintAllIP             bool             `json:"-"`      //这个字段不解析
	Debug                  bool             `json:"-"`      //这个字段不解析
	LogToFile              bool             `json:"-"`      //这个字段不解析
	LogFilePath            string           `json:"-"`      //这个字段不解析
	LogFileName            string           `json:"-"`      //这个字段不解析
	MaxDepth               int              `json:"-"`      //这个字段不解析
	ZipFileTextEncoding    string           `json:"-"`      //这个字段不解析
	TempPATH               string           `json:"-"`      //这个字段不解析
	CleanAllTempFileOnExit bool             `json:"-"`      //这个字段不解析
	CleanAllTempFile       bool             `json:"-"`      //这个字段不解析
	GenerateConfig         bool             `json:"-"`      //这个字段不解析
	WebpConfig             WebPServerConfig `json:"-"`      //这个字段不解析
	FrpConfig              FrpClientConfig  `json:"-"`      //这个字段不解析
	//Template               string           `json:"-"` //这个字段不解析
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
	FrpcCommand      string
	ServerAddr       string
	ServerPort       int
	Token            string
	FrpType          string //本地转发端口设置
	RemotePort       int
	RandomRemotePort bool
}

// setupCloseHander 中断处理：程序被中断的时候，清理临时文件
func setupCloseHander() {
	c := make(chan os.Signal, 2)
	//SIGHUP（挂起）, SIGINT（中断）或 SIGTERM（终止）默认会使得程序退出。
	//1、SIGHUP 信号在用户终端连接(正常或非正常)结束时发出。
	//2、syscall.SIGINT 和 os.Interrupt 是同义词,按下 CTRL+C 时发出。
	//3、SIGTERM（终止）:kill终止进程,允许程序处理问题后退出。
	//4.syscall.SIGHUP,终端控制进程结束(终端连接断开)
	//5、syscall.SIGQUIT，CTRL+\ 退出
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		<-c
		if Config.CleanAllTempFileOnExit {
			fmt.Println("\r" + locale.GetString("start_clear_file"))
			clearTempFilesALL()
		} else {
			//clearTempFilesOne(ReadingBook)
		}
		os.Exit(0)
	}()
}

// setTempDir 设置临时文件夹，退出时会被清理
func setTempDir() {
	//手动设置的临时文件夹
	if Config.TempPATH != "" && tools.ChickExists(Config.TempPATH) && tools.ChickIsDir(Config.TempPATH) {
		Config.CacheFilePath = path.Join(Config.TempPATH)
	} else {
		Config.CacheFilePath = path.Join(os.TempDir(), "comigo_temp_files") //直接使用系统文件夹
	}
	err := os.MkdirAll(Config.CacheFilePath, os.ModePerm)
	if err != nil {
		println(locale.GetString("temp_folder_error"))
	} else {
		fmt.Println(locale.GetString("temp_folder_path") + Config.CacheFilePath)
	}
}

// 清空程序缓存 TODO：生成临时文件，并在退出后清理
func clearTempFilesALL() {
	fmt.Println(locale.GetString("clear_temp_file_start"))
	for _, tempBook := range mapBooks {
		clearTempFilesOne(tempBook)
	}
}

// 清空某一本压缩漫画的解压缓存
func clearTempFilesOne(book *Book) {
	fmt.Println(locale.GetString("clear_temp_file_start"))
	haveThisBook := false
	for _, tempBook := range mapBooks {
		if tempBook.GetBookID() == book.GetBookID() {
			haveThisBook = true
		}
	}
	if haveThisBook {
		extractPath := path.Join(Config.CacheFilePath, book.GetBookID())
		//避免删错文件,解压路径包含UUID，len不可能小于32
		PathLen := len(extractPath)
		if PathLen < 32 {
			return
		}
		err := os.RemoveAll(extractPath)
		if err != nil {
			fmt.Println(locale.GetString("clear_temp_file_error") + extractPath)
		} else {
			fmt.Println(locale.GetString("clear_temp_file_completed") + extractPath)
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
