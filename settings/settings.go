package settings

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ServerSettings 服务器设置(config.toml)，配置文件放在被扫描的根目录中，或$HOME/config/comigo.可以用“comi --generate-config”生成本示例文件
type ServerSettings struct {
	Port                 int             `json:"port" comment:"Comigo设置(config.toml)，放在执行目录中，或$HOME/.config/comigo/下。可用“comi --generate-config”生成本文件\n# 网页端口"`
	Host                 string          `json:"host" comment:"自定义二维码显示的主机名"`
	StoresPath           []string        `json:"-"    comment:"设置默认扫描的书库文件夹"`
	MaxDepth             int             `json:"-" comment:"最大扫描深度"`
	OpenBrowser          bool            `json:"-" comment:"是否同时打开浏览器，windows默认true，其他默认false"`
	DisableLAN           bool            `json:"-" comment:"只在本机localhost提供服务，不对外共享"`
	DefaultMode          string          `json:"default_mode" comment:"默认阅读模式，默认为空，可以设置为scroll或flip"`
	UserName             string          `json:"-" comment:"访问限制：用户名。需要设置密码"`
	Password             string          `json:"-" comment:"访问限制：密码。需要设置用户名。"`
	CertFile             string          `json:"-" comment:"Https证书，同时设置KeyFile则启用HTTPS协议"`
	KeyFile              string          `json:"-" comment:"Https证书，同时设置CertFile则启用HTTPS协议"`
	CacheFileEnable      bool            `json:"-" comment:"是否保存web图片缓存，可以加快二次读取，但会占用硬盘空间"`
	CacheFilePath        string          `json:"-" comment:"web图片缓存存储位置，默认系统临时文件夹"`
	CacheFileClean       bool            `json:"-" comment:"退出程序的时候，清理web图片缓存"`
	DatabaseFilePath     string          `json:"-" comment:"数据库文件存储位置，默认当前目录"`
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
	//EnableWebpServer       bool             `json:"enable_webp_server"`
	//SketchCountSeconds     int              `json:"sketch_count_seconds"`
	//WebpConfig             WebPServerConfig `json:"-"  comment:" WebPServer设置"`
	//Template               string           `json:"-"`
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

// 判断压缩包内的文件是否需要展示（是图片或媒体文件）
func (config *ServerSettings) IsSupportMedia(checkPath string) bool {
	for _, ex := range config.SupportMediaType {
		//strings.ToLower():某些文件会用大写文件名
		suffix := strings.ToLower(path.Ext(checkPath))
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportArchiver 是否是支持的压缩文件
func (config *ServerSettings) IsSupportArchiver(checkPath string) bool {
	for _, ex := range config.SupportFileType {
		suffix := path.Ext(checkPath)
		if ex == suffix {
			return true
		}
	}
	return false
}

// CheckPathSkip 检查路径是否应该跳过（排除文件，文件夹列表）。
func (config *ServerSettings) IsSkipDir(path string) bool {
	for _, substr := range config.ExcludeFileOrFolders {
		if strings.HasSuffix(path, substr) {
			return true
		}
	}
	return false
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
