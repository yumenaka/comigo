package types

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/tidwall/gjson"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/util"
)

type ConfigStatus struct {
	// 当前生效的配置文件路径 None、HomeDirectory、WorkingDirectory、ProgramDirectory
	// 设置读取顺序：None（默认值） -> HomeDirectory -> ProgramDirectory -> WorkingDirectory
	In   string
	Path struct {
		// 对应配置文件的绝对路径
		WorkingDirectory string
		HomeDirectory    string
		ProgramDirectory string
	}
}

func (c *ConfigStatus) SetConfigStatus() error {

	c.In = ""
	c.Path.WorkingDirectory = ""
	c.Path.HomeDirectory = ""
	c.Path.ProgramDirectory = ""
	logger.Info("Check Config Path\n")

	// 可执行程序自身的文件路径
	executablePath, err := os.Executable()
	if err != nil {
		return errors.New("error: Failed find executable path")
	}
	if util.IsExist(path.Join(path.Dir(executablePath), "config.toml")) {
		c.Path.ProgramDirectory = util.GetAbsPath(path.Join(path.Dir(executablePath), "config.toml"))
		c.In = "ProgramDirectory"
	}

	// HomeDirectory 目录
	home, err := homedir.Dir()
	if err != nil {
		return errors.New("error: Failed find home directory")
	}
	if util.IsExist(path.Join(home, ".config/comigo/config.toml")) {
		c.Path.HomeDirectory = util.GetAbsPath(path.Join(home, ".config/comigo/config.toml"))
		c.In = "HomeDirectory"
	}

	//当前执行目录
	if util.IsExist("config.toml") {
		c.Path.WorkingDirectory = util.GetAbsPath("config.toml")
		c.In = "WorkingDirectory"
	}
	return nil
}

type RemoteStore struct {
	// 远程书库的类型,支持ftp、sftp、webdav、smb（2或3）
	Type string
	// 远程书库的地址
	Host string
	// 远程书库的端口
	Port int
	// 远程书库的用户名
	Username string
	// 远程书库的密码
	Password string
}

// ComigoConfig 服务器设置(config.toml)
type ComigoConfig struct {
	Port                   int           `json:"Port" comment:"Comigo设置文件(config.toml)，可保存在：HomeDirectory（$HOME/.config/comigo/config.toml）、WorkingDirectory（当前执行目录）、ProgramDirectory（程序所在目录）下。可用“comi --config-save”生成本文件\n网页服务端口，此项配置不支持热重载"`
	ConfigPath             string        `json:"-" toml:"-" comment:"用户指定的的yaml设置文件路径"`
	Host                   string        `json:"Host" comment:"自定义二维码显示的主机名"`
	LocalStores            []string      `json:"LocalStores" comment:"本地书库文件夹"`
	RemoteStores           []RemoteStore `json:"RemoteStores" toml:"-" comment:"用户指定的的yaml设置文件路径"`
	ExcludePath            []string      `json:"ExcludePath" comment:"扫描书籍的时候，需要排除的文件或文件夹的名字"`
	SupportMediaType       []string      `json:"SupportMediaType" comment:"扫描压缩包时，用于统计图片数量的图片文件后缀"`
	SupportFileType        []string      `json:"SupportFileType" comment:"支持的书籍压缩包后缀"`
	MinImageNum            int           `json:"MinImageNum" comment:"压缩包或文件夹内，至少有几张图片，才算作书籍"`
	TimeoutLimitForScan    int           `json:"TimeoutLimitForScan" comment:"扫描文件时，超过几秒钟，就放弃扫描这个文件，避免卡在特殊文件上"`
	EnableUpload           bool          `json:"EnableUpload" comment:"启用上传功能"`
	UploadPath             string        `json:"UploadPath" comment:"上传文件存储位置，默认在当前执行目录下创建 upload 文件夹"`
	MaxScanDepth           int           `json:"MaxScanDepth" comment:"最大扫描深度"`
	ZipFileTextEncoding    string        `json:"ZipFileTextEncoding" comment:"非utf-8编码的ZIP文件，尝试用什么编码解析，默认GBK"`
	PrintAllPossibleQRCode bool          `json:"PrintAllPossibleQRCode" comment:"扫描完成后，打印所有可能的阅读链接二维码"`
	Debug                  bool          `json:"Debug" comment:"开启Debug模式"`
	OpenBrowser            bool          `json:"OpenBrowser" comment:"是否同时打开浏览器，windows默认true，其他默认false"`
	DisableLAN             bool          `json:"DisableLAN" comment:"只在本机提供阅读服务，不对外共享，此项配置不支持热重载"`
	DefaultMode            string        `json:"DefaultMode" comment:"默认阅读模式，默认为空，可以设置为scroll或flip"`
	EnableLogin            bool          `json:"EnableLogin" comment:"是否启用登录。默认不需要登陆。此项配置不支持热重载。"`
	Username               string        `json:"Username" comment:"启用登陆后，登录界面需要的用户名。"`
	Password               string        `json:"Password" comment:"启用登陆后，登录界面需要的密码。"`
	Timeout                int           `json:"Timeout" comment:"启用登陆后，cookie过期的时间。单位为分钟。默认180分钟后过期。"`
	EnableTLS              bool          `json:"EnableTLS" comment:"是否启用HTTPS协议。需要设置证书于key文件。"`
	CertFile               string        `json:"CertFile" comment:"TLS/SSL 证书文件路径 (default: "~/.config/.comigo/cert.crt")"`
	KeyFile                string        `json:"KeyFile" comment:"TLS/SSL key文件路径 (default: "~/.config/.comigo/key.key")"`
	UseCache               bool          `json:"UseCache" comment:"开启本地图片缓存，可以加快二次读取，但会占用硬盘空间"`
	CachePath              string        `json:"CachePath" comment:"本地图片缓存位置，默认系统临时文件夹"`
	ClearCacheExit         bool          `json:"ClearCacheExit" comment:"退出程序的时候，清理web图片缓存"`
	EnableDatabase         bool          `json:"EnableDatabase" comment:"启用本地数据库，保存扫描到的书籍数据。此项配置不支持热重载。"`
	ClearDatabaseWhenExit  bool          `json:"ClearDatabaseWhenExit" comment:"启用本地数据库时，扫描完成后，清除不存在的书籍。"`
	LogToFile              bool          `json:"LogToFile" comment:"是否保存程序Log到本地文件。默认不保存。"`
	LogFilePath            string        `json:"LogFilePath" comment:"Log文件的保存位置"`
	LogFileName            string        `json:"LogFileName" comment:"Log文件名"`
	GenerateMetaData       bool          `json:"GenerateMetaData" toml:"GenerateMetaData" comment:"生成书籍元数据"`
}

func UpdateConfig(config *ComigoConfig, jsonString string) (*ComigoConfig, error) {
	oldConfig := *config
	Port := gjson.Get(jsonString, "Port")
	if Port.Exists() {
		config.Port = int(Port.Int())
	}
	Host := gjson.Get(jsonString, "Host")
	if Host.Exists() {
		config.Host = Host.String()
	}
	LocalStores := gjson.Get(jsonString, "LocalStores")
	if LocalStores.Exists() {
		// 将字符串解析为字符串切片
		arr, err := parseString(LocalStores.String())
		if err != nil {
			logger.Infof("Failed to parse string:%s", err)
			return config, err
		}
		config.LocalStores = arr
	}
	UseCache := gjson.Get(jsonString, "UseCache")
	if UseCache.Exists() {
		config.UseCache = UseCache.Bool()
	}
	CachePath := gjson.Get(jsonString, "CachePath")
	if CachePath.Exists() {
		config.CachePath = CachePath.String()
	}
	ClearCacheExit := gjson.Get(jsonString, "ClearCacheExit")
	if ClearCacheExit.Exists() {
		config.ClearCacheExit = ClearCacheExit.Bool()
	}
	UploadPath := gjson.Get(jsonString, "UploadPath")
	if UploadPath.Exists() {
		config.UploadPath = UploadPath.String()
	}
	EnableUpload := gjson.Get(jsonString, "EnableUpload")
	if EnableUpload.Exists() {
		config.EnableUpload = EnableUpload.Bool()
	}
	EnableDatabase := gjson.Get(jsonString, "EnableDatabase")
	if EnableDatabase.Exists() {
		config.EnableDatabase = EnableDatabase.Bool()
	}
	ClearDatabaseWhenExit := gjson.Get(jsonString, "ClearDatabaseWhenExit")
	if ClearDatabaseWhenExit.Exists() {
		config.ClearDatabaseWhenExit = ClearDatabaseWhenExit.Bool()
	}
	OpenBrowser := gjson.Get(jsonString, "OpenBrowser")
	if OpenBrowser.Exists() {
		config.OpenBrowser = OpenBrowser.Bool()
	}
	DisableLAN := gjson.Get(jsonString, "DisableLAN")
	if DisableLAN.Exists() {
		config.DisableLAN = DisableLAN.Bool()
	}
	DefaultMode := gjson.Get(jsonString, "DefaultMode")
	if DefaultMode.Exists() {
		config.DefaultMode = DefaultMode.String()
	}
	LogToFile := gjson.Get(jsonString, "LogToFile")
	if LogToFile.Exists() {
		config.LogToFile = LogToFile.Bool()
	}
	MaxScanDepth := gjson.Get(jsonString, "MaxScanDepth")
	if MaxScanDepth.Exists() {
		config.MaxScanDepth = int(MaxScanDepth.Int())
	}
	MinImageNum := gjson.Get(jsonString, "MinImageNum")
	if MinImageNum.Exists() {
		config.MinImageNum = int(MinImageNum.Int())
	}
	ZipFileTextEncoding := gjson.Get(jsonString, "ZipFileTextEncoding")
	if ZipFileTextEncoding.Exists() {
		config.ZipFileTextEncoding = ZipFileTextEncoding.String()
	}
	ExcludePath := gjson.Get(jsonString, "ExcludePath")
	if ExcludePath.Exists() {
		// 将字符串解析为字符串切片
		arr, err := parseString(ExcludePath.String())
		if err != nil {
			logger.Infof("Failed to parse string:%s", err)
			return config, err
		}
		config.ExcludePath = arr
	}
	SupportMediaType := gjson.Get(jsonString, "SupportMediaType")
	if SupportMediaType.Exists() {
		// 将字符串解析为字符串切片
		arr, err := parseString(SupportMediaType.String())
		if err != nil {
			logger.Infof("Failed to parse string:%s", err)
			return config, err
		}
		config.SupportMediaType = arr
	}
	SupportFileType := gjson.Get(jsonString, "SupportFileType")
	if SupportFileType.Exists() {
		// 将字符串解析为字符串切片
		arr, err := parseString(SupportFileType.String())
		if err != nil {
			logger.Infof("Failed to parse string:%s", err)
			return config, err
		}
		config.SupportFileType = arr
	}
	TimeoutLimitForScan := gjson.Get(jsonString, "TimeoutLimitForScan")
	if TimeoutLimitForScan.Exists() {
		config.TimeoutLimitForScan = int(TimeoutLimitForScan.Int())
	}
	PrintAllPossibleQRCode := gjson.Get(jsonString, "PrintAllPossibleQRCode")
	if PrintAllPossibleQRCode.Exists() {
		config.PrintAllPossibleQRCode = PrintAllPossibleQRCode.Bool()
	}
	Debug := gjson.Get(jsonString, "Debug")
	if Debug.Exists() {
		config.Debug = Debug.Bool()
	}
	Username := gjson.Get(jsonString, "Username")
	if Username.Exists() {
		config.Username = Username.String()
	}
	Password := gjson.Get(jsonString, "Password")
	if Password.Exists() {
		config.Password = Password.String()
	}
	Timeout := gjson.Get(jsonString, "Timeout")
	if Timeout.Exists() {
		config.Timeout = int(Timeout.Int())
	}
	GenerateMetaData := gjson.Get(jsonString, "GenerateMetaData")
	if GenerateMetaData.Exists() {
		config.GenerateMetaData = GenerateMetaData.Bool()
	}
	return &oldConfig, nil
}

// 将字符串解析为字符串切片
func parseString(str string) ([]string, error) {
	var arr []string
	err := json.Unmarshal([]byte(str), &arr)
	if err != nil {
		return nil, err
	}
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

// SetByExecutableFilename 通过执行文件名设置默认网页模板参数
func (c *ComigoConfig) SetByExecutableFilename() {
	// 当前执行目录
	//targetPath, _ := os.Getwd()
	//logger.Infof(locale.GetString("target_path"), targetPath)
	// 带后缀的执行文件名 comi.exe  sketch.exe
	filenameWithSuffix := path.Base(os.Args[0])
	// 执行文件名后缀
	fileSuffix := path.Ext(filenameWithSuffix)
	// 去掉后缀后的执行文件名
	filenameWithOutSuffix := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	//logger.Infof("filenameWithOutSuffix =", filenameWithOutSuffix)
	ex, err := os.Executable()
	if err != nil {
		logger.Infof("%s", err)
	}
	extPath := filepath.Dir(ex)
	//logger.Infof("extPath =",extPath)
	ExtFileName := strings.TrimPrefix(filenameWithOutSuffix, extPath)
	if c.Debug {
		logger.Infof("ExtFileName =%s", ExtFileName)
	}
}
