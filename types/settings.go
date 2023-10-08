package types

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ServerConfig 服务器设置(config.toml)
type ServerConfig struct {
	Port                   int      `json:"Port" comment:"Comigo设置文件(config.toml)，放在默认保存目录中。可选值：RAM（不保存）、HomeDir（$HOME/.config/comigo/config.toml）、NowDir（当前执行目录）、ProgramDir（程序所在目录）下。可用“comi --config-save”生成本文件\n网页服务端口，此项配置不支持热重载"`
	ConfigSaveTo           string   `json:"ConfigSaveTo" toml:"-" comment:"配置文件的默认保存位置，可选值：RAM（不保存）、HomeDir（$HOME/.config/comigo/）、NowDir（当前执行目录）、ProgramDir（程序所在目录）"`
	ConfigFileUsed         string   `json:"-" toml:"-" comment:"当前生效的yaml设置文件路径，数据库文件(comigo.db)在同一个文件夹"`
	Host                   string   `json:"Host" comment:"自定义二维码显示的主机名"`
	StoresPath             []string `json:"StoresPath" comment:"默认扫描的书库文件夹"`
	ExcludePath            []string `json:"ExcludePath" comment:"扫描书籍的时候，需要排除的文件或文件夹的名字"`
	SupportMediaType       []string `json:"SupportMediaType" comment:"扫描压缩包时，用于统计图片数量的图片文件后缀"`
	SupportFileType        []string `json:"SupportFileType" comment:"支持的书籍压缩包后缀"`
	MinImageNum            int      `json:"MinImageNum" comment:"压缩包或文件夹内，至少有几张图片，才算作书籍"`
	TimeoutLimitForScan    int      `json:"TimeoutLimitForScan" comment:"扫描文件时，超过几秒钟，就放弃扫描这个文件，避免卡在特殊文件上"`
	EnableUpload           bool     `json:"EnableUpload" comment:"启用上传功能"`
	UploadPath             string   `json:"UploadPath" comment:"上传文件存储位置，默认在当前执行目录下创建 upload 文件夹"`
	MaxScanDepth           int      `json:"MaxScanDepth" comment:"最大扫描深度"`
	ZipFileTextEncoding    string   `json:"ZipFileTextEncoding" comment:"非utf-8编码的ZIP文件，尝试用什么编码解析，默认GBK"`
	PrintAllPossibleQRCode bool     `json:"PrintAllPossibleQRCode" comment:"扫描完成后，打印所有可能的阅读链接二维码"`
	Debug                  bool     `json:"Debug" comment:"开启Debug模式"`
	OpenBrowser            bool     `json:"OpenBrowser" comment:"是否同时打开浏览器，windows默认true，其他默认false"`
	DisableLAN             bool     `json:"DisableLAN" comment:"只在本机提供阅读服务，不对外共享，此项配置不支持热重载"`
	DefaultMode            string   `json:"DefaultMode" comment:"默认阅读模式，默认为空，可以设置为scroll或flip"`
	EnableLogin            bool     `json:"EnableLogin" comment:"是否启用登录。默认不需要登陆。此项配置不支持热重载。"`
	Username               string   `json:"Username" comment:"启用登陆后，登录界面需要的用户名。"`
	Password               string   `json:"Password" comment:"启用登陆后，登录界面需要的密码。"`
	Timeout                int      `json:"Timeout" comment:"启用登陆后，cookie过期的时间。单位为分钟。默认180分钟后过期。"`
	EnableTLS              bool     `json:"EnableTLS" comment:"是否启用HTTPS协议。需要设置证书于key文件。"`
	CertFile               string   `json:"CertFile" comment:"TLS/SSL 证书文件路径 (default: "~/.config/.comigo/cert.crt")"`
	KeyFile                string   `json:"KeyFile" comment:"TLS/SSL key文件路径 (default: "~/.config/.comigo/key.key")"`
	UseCache               bool     `json:"UseCache" comment:"开启本地图片缓存，可以加快二次读取，但会占用硬盘空间"`
	CachePath              string   `json:"CachePath" comment:"本地图片缓存位置，默认系统临时文件夹"`
	ClearCacheExit         bool     `json:"ClearCacheExit" comment:"退出程序的时候，清理web图片缓存"`
	EnableDatabase         bool     `json:"EnableDatabase" comment:"启用本地数据库，保存扫描到的书籍数据。此项配置不支持热重载。"`
	ClearDatabaseWhenExit  bool     `json:"ClearDatabaseWhenExit" comment:"启用本地数据库时，扫描完成后，清除不存在的书籍。"`
	LogToFile              bool     `json:"LogToFile" comment:"是否保存程序Log到本地文件。默认不保存。"`
	LogFilePath            string   `json:"LogFilePath" comment:"Log文件的保存位置"`
	LogFileName            string   `json:"LogFileName" comment:"Log文件名"`
	GenerateMetaData       bool     `json:"GenerateMetaData" toml:"GenerateMetaData" comment:"生成书籍元数据"`

	//EnableWebpServer     bool             `json:"enable_webp_server"`
	//WebpConfig           WebPServerConfig `json:"-"  comment:" WebPServer设置"`
	//DatabaseFilePath     string           `json:"-" comment:"数据库文件存储位置，默认config目录"`
}

func UpdateConfig(oldConfig ServerConfig, jsonString string) (newConfig ServerConfig, err error) {
	newConfig = oldConfig
	Port := gjson.Get(jsonString, "Port")
	if Port.Exists() {
		newConfig.Port = int(Port.Int())
	}
	Host := gjson.Get(jsonString, "Host")
	if Host.Exists() {
		newConfig.Host = Host.String()
	}
	StoresPath := gjson.Get(jsonString, "StoresPath")
	if StoresPath.Exists() {
		// 将字符串解析为字符串切片
		arr, err := parseString(StoresPath.String())
		if err != nil {
			fmt.Println("Failed to parse string:", err)
			return newConfig, err
		}
		newConfig.StoresPath = arr
	}
	UseCache := gjson.Get(jsonString, "UseCache")
	if UseCache.Exists() {
		newConfig.UseCache = UseCache.Bool()
	}
	CachePath := gjson.Get(jsonString, "CachePath")
	if CachePath.Exists() {
		newConfig.CachePath = CachePath.String()
	}
	ClearCacheExit := gjson.Get(jsonString, "ClearCacheExit")
	if ClearCacheExit.Exists() {
		newConfig.ClearCacheExit = ClearCacheExit.Bool()
	}
	UploadPath := gjson.Get(jsonString, "UploadPath")
	if UploadPath.Exists() {
		newConfig.UploadPath = UploadPath.String()
	}
	EnableUpload := gjson.Get(jsonString, "EnableUpload")
	if EnableUpload.Exists() {
		newConfig.EnableUpload = EnableUpload.Bool()
	}
	EnableDatabase := gjson.Get(jsonString, "EnableDatabase")
	if EnableDatabase.Exists() {
		newConfig.EnableDatabase = EnableDatabase.Bool()
	}
	ClearDatabaseWhenExit := gjson.Get(jsonString, "ClearDatabaseWhenExit")
	if ClearDatabaseWhenExit.Exists() {
		newConfig.ClearDatabaseWhenExit = ClearDatabaseWhenExit.Bool()
	}
	OpenBrowser := gjson.Get(jsonString, "OpenBrowser")
	if OpenBrowser.Exists() {
		newConfig.OpenBrowser = OpenBrowser.Bool()
	}
	DisableLAN := gjson.Get(jsonString, "DisableLAN")
	if DisableLAN.Exists() {
		newConfig.DisableLAN = DisableLAN.Bool()
	}
	DefaultMode := gjson.Get(jsonString, "DefaultMode")
	if DefaultMode.Exists() {
		newConfig.DefaultMode = DefaultMode.String()
	}
	LogToFile := gjson.Get(jsonString, "LogToFile")
	if LogToFile.Exists() {
		newConfig.LogToFile = LogToFile.Bool()
	}
	MaxScanDepth := gjson.Get(jsonString, "MaxScanDepth")
	if MaxScanDepth.Exists() {
		newConfig.MaxScanDepth = int(MaxScanDepth.Int())
	}
	MinImageNum := gjson.Get(jsonString, "MinImageNum")
	if MinImageNum.Exists() {
		newConfig.MinImageNum = int(MinImageNum.Int())
	}
	ZipFileTextEncoding := gjson.Get(jsonString, "ZipFileTextEncoding")
	if ZipFileTextEncoding.Exists() {
		newConfig.ZipFileTextEncoding = ZipFileTextEncoding.String()
	}
	ExcludePath := gjson.Get(jsonString, "ExcludePath")
	if ExcludePath.Exists() {
		// 将字符串解析为字符串切片
		arr, err := parseString(ExcludePath.String())
		if err != nil {
			fmt.Println("Failed to parse string:", err)
			return newConfig, err
		}
		newConfig.ExcludePath = arr
	}
	SupportMediaType := gjson.Get(jsonString, "SupportMediaType")
	if SupportMediaType.Exists() {
		// 将字符串解析为字符串切片
		arr, err := parseString(SupportMediaType.String())
		if err != nil {
			fmt.Println("Failed to parse string:", err)
			return newConfig, err
		}
		newConfig.SupportMediaType = arr
	}
	SupportFileType := gjson.Get(jsonString, "SupportFileType")
	if SupportFileType.Exists() {
		// 将字符串解析为字符串切片
		arr, err := parseString(SupportFileType.String())
		if err != nil {
			fmt.Println("Failed to parse string:", err)
			return newConfig, err
		}
		newConfig.SupportFileType = arr
	}
	TimeoutLimitForScan := gjson.Get(jsonString, "TimeoutLimitForScan")
	if TimeoutLimitForScan.Exists() {
		newConfig.TimeoutLimitForScan = int(TimeoutLimitForScan.Int())
	}
	PrintAllPossibleQRCode := gjson.Get(jsonString, "PrintAllPossibleQRCode")
	if PrintAllPossibleQRCode.Exists() {
		newConfig.PrintAllPossibleQRCode = PrintAllPossibleQRCode.Bool()
	}
	Debug := gjson.Get(jsonString, "Debug")
	if Debug.Exists() {
		newConfig.Debug = Debug.Bool()
	}
	Username := gjson.Get(jsonString, "Username")
	if Username.Exists() {
		newConfig.Username = Username.String()
	}
	Password := gjson.Get(jsonString, "Password")
	if Password.Exists() {
		newConfig.Password = Password.String()
	}
	Timeout := gjson.Get(jsonString, "Timeout")
	if Timeout.Exists() {
		newConfig.Timeout = int(Timeout.Int())
	}
	GenerateMetaData := gjson.Get(jsonString, "GenerateMetaData")
	if GenerateMetaData.Exists() {
		newConfig.GenerateMetaData = GenerateMetaData.Bool()
	}
	ConfigSaveTo := gjson.Get(jsonString, "ConfigSaveTo")
	if ConfigSaveTo.Exists() {
		newConfig.ConfigSaveTo = ConfigSaveTo.String()
	}
	return newConfig, nil
}

// 将字符串解析为字符串切片
func parseString(str string) ([]string, error) {
	var arr []string
	//fmt.Println("str =", str)
	err := json.Unmarshal([]byte(str), &arr)
	if err != nil {
		return nil, err
	}
	//arr := strings.Split(str, ",")
	return arr, nil
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
func (config *ServerConfig) IsSupportMedia(checkPath string) bool {
	for _, ex := range config.SupportMediaType {
		suffix := strings.ToLower(path.Ext(checkPath)) //strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportArchiver 是否是支持的压缩文件
func (config *ServerConfig) IsSupportArchiver(checkPath string) bool {
	for _, ex := range config.SupportFileType {
		suffix := path.Ext(checkPath)
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSkipDir  检查路径是否应该跳过（排除文件，文件夹列表）。
func (config *ServerConfig) IsSkipDir(path string) bool {
	for _, substr := range config.ExcludePath {
		if strings.HasSuffix(path, substr) {
			return true
		}
	}
	return false
}

// SetByExecutableFilename 通过执行文件名设置默认网页模板参数
func (config *ServerConfig) SetByExecutableFilename() {
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
