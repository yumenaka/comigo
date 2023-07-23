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
	Port                 int             `json:"Port" comment:"Comigo设置文件(config.toml)，放在执行目录中，或$HOME/.config/comigo/下。可用“comi --generate-config”生成本文件\n网页服务端口，此项配置不支持热重载"`
	Host                 string          `json:"Host" comment:"自定义二维码显示的主机名"`
	StoresPath           []string        `json:"StoresPath" comment:"设置默认扫描的书库文件夹"`
	MaxDepth             int             `json:"MaxDepth" comment:"最大扫描深度"`
	OpenBrowser          bool            `json:"OpenBrowser" comment:"是否同时打开浏览器，windows默认true，其他默认false"`
	DisableLAN           bool            `json:"DisableLAN" comment:"只在本机localhost提供服务，不对外共享，此项配置不支持热重载"`
	DefaultMode          string          `json:"DefaultMode" comment:"默认阅读模式，默认为空，可以设置为scroll或flip"`
	UserName             string          `json:"UserName" comment:"访问限制：用户名。需要设置密码"`
	Password             string          `json:"Password" comment:"访问限制：密码。需要设置用户名。"`
	Timeout              int             `json:"Timeout" comment:"cookie过期时间。单位为分钟。默认180分钟"`
	CertFile             string          `json:"CertFile" comment:"Https证书，同时设置KeyFile则启用HTTPS协议"`
	KeyFile              string          `json:"KeyFile" comment:"Https证书，同时设置CertFile则启用HTTPS协议"`
	CacheEnable          bool            `json:"CacheEnable" comment:"是否保存web图片缓存，可以加快二次读取，但会占用硬盘空间"`
	CachePath            string          `json:"CachePath" comment:"web图片缓存存储位置，默认系统临时文件夹"`
	CacheClean           bool            `json:"CacheClean" comment:"退出程序的时候，清理web图片缓存"`
	EnableUpload         bool            `json:"EnableUpload" comment:"启用文件上传功能"`
	UploadPath           string          `json:"UploadPath" comment:"上传文件的存储位置，默认在当前执行目录下创建 ComigoUpload 文件夹"`
	EnableDatabase       bool            `json:"EnableDatabase" comment:"启用本地数据库，保存扫描到的书籍数据"`
	ClearDatabase        bool            `json:"ClearDatabase" comment:"启用本地数据库时，扫描完成后，清除不存在的书籍"`
	ExcludeFileOrFolders []string        `json:"ExcludeFileOrFolders" comment:"需要排除的文件或文件夹"`
	SupportMediaType     []string        `json:"SupportMediaType" comment:"需要扫描的图片文件后缀"`
	SupportFileType      []string        `json:"SupportFileType" comment:"需要扫描的图书文件后缀"`
	MinImageNum          int             `json:"MinImageNum" comment:"压缩包或文件夹内，至少有几张图片，才算作书籍"`
	TimeoutLimitForScan  int             `json:"TimeoutLimitForScan" comment:"扫描文件时，超过几秒钟，就放弃扫描这个文件，避免卡在特殊文件上"`
	PrintAllIP           bool            `json:"PrintAllIP" comment:"打印所有可能阅读链接的二维码"`
	Debug                bool            `json:"Debug" comment:"开启Debug模式"`
	LogToFile            bool            `json:"LogToFile" comment:"记录Log到本地文件"`
	LogFilePath          string          `json:"LogFilePath" comment:"Log保存位置"`
	LogFileName          string          `json:"LogFileName" comment:"Log文件名"`
	ZipFileTextEncoding  string          `json:"ZipFileTextEncoding" comment:"非utf-8编码的ZIP文件，尝试用什么编码解析，默认GBK"`
	EnableFrpcServer     bool            `json:"EnableFrpcServer" comment:"后台启动FrpClient"`
	FrpConfig            FrpClientConfig `json:"FrpConfig" comment:"FrpClient设置"`
	GenerateMetaData     bool            `json:"GenerateMetaData" toml:"GenerateMetaData" comment:"生成书籍元数据"`
	//EnableWebpServer     bool             `json:"enable_webp_server"`
	//SketchCountSeconds   int              `json:"sketch_count_seconds"`
	//WebpConfig           WebPServerConfig `json:"-"  comment:" WebPServer设置"`
	//Template             string           `json:"-"`
	//DatabaseFilePath     string           `json:"-" comment:"数据库文件存储位置，默认当前目录"`
}

// FrpClientConfig frp客户端配置
type FrpClientConfig struct {
	FrpcCommand      string `comment:"手动设定frpc可执行程序的路径,默认为frpc"`
	ServerAddr       string
	ServerPort       int
	Token            string
	FrpType          string //本地转发端口设置
	RemotePort       int
	RandomRemotePort bool
}

// WebPServerConfig  WebPServer服务端配置
type WebPServerConfig struct {
	WebpCommand  string
	HOST         string
	PORT         string
	ImgPath      string
	QUALITY      int
	AllowedTypes []string
	ExhaustPath  string
}

// IsSupportMedia 判断压缩包内的文件是否需要展示（包括图片、音频、视频、PDF在内的媒体文件）
func (config *ServerSettings) IsSupportMedia(checkPath string) bool {
	for _, ex := range config.SupportMediaType {
		suffix := strings.ToLower(path.Ext(checkPath)) //strings.ToLower():某些文件会用大写文件名
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
	if config.Debug {
		fmt.Println("ExtFileName =", ExtFileName)
	}
}
