package config

import "C"
import (
	"encoding/json"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/locale"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/yumenaka/comigo/config/stores"
	"github.com/yumenaka/comigo/util/logger"
)

// Config Comigo全局配置
type Config struct {
	AutoRescan             bool            `json:"AutoRescan" comment:"刷新页面时，是否自动重新扫描"`
	CachePath              string          `json:"CachePath" comment:"本地图片缓存位置，默认系统临时文件夹"`
	CertFile               string          `json:"CertFile" comment:"TLS/SSL 证书文件路径 (default: ~/.config/.comigo/cert.crt)"`
	ClearCacheExit         bool            `json:"ClearCacheExit" comment:"退出程序的时候，清理web图片缓存"`
	ClearDatabaseWhenExit  bool            `json:"ClearDatabaseWhenExit" comment:"启用本地数据库时，扫描完成后，清除不存在的书籍。"`
	ConfigPath             string          `json:"-" toml:"-" comment:"用户指定的的yaml设置文件路径"`
	Debug                  bool            `json:"Debug" comment:"开启Debug模式"`
	DefaultMode            string          `json:"DefaultMode" comment:"默认阅读模式，默认为空，可以设置为scroll或flip"`
	DisableLAN             bool            `json:"DisableLAN" comment:"只在本机提供阅读服务，不对外共享，此项配置不支持热重载"`
	EnableDatabase         bool            `json:"EnableDatabase" comment:"启用本地数据库，保存扫描到的书籍数据。此项配置不支持热重载。"`
	EnableLogin            bool            `json:"EnableLogin" comment:"是否启用登录。默认不需要登陆。此项配置不支持热重载。"`
	EnableTLS              bool            `json:"EnableTLS" comment:"是否启用HTTPS协议。需要设置证书于key文件。"`
	EnableUpload           bool            `json:"EnableUpload" comment:"启用上传功能"`
	ExcludePath            []string        `json:"ExcludePath" comment:"扫描书籍的时候，需要排除的文件或文件夹的名字"`
	GenerateMetaData       bool            `json:"GenerateMetaData" toml:"GenerateMetaData" comment:"生成书籍元数据"`
	Host                   string          `json:"Host" comment:"自定义二维码显示的主机名"`
	KeyFile                string          `json:"KeyFile" comment:"TLS/SSL key文件路径 (default: ~/.config/.comigo/key.key)"`
	LocalStores            []string        `json:"LocalStores" comment:"本地书库文件夹"`
	LogFileName            string          `json:"LogFileName" comment:"Log文件名"`
	LogFilePath            string          `json:"LogFilePath" comment:"Log文件的保存位置"`
	LogToFile              bool            `json:"LogToFile" comment:"是否保存程序Log到本地文件。默认不保存。"`
	MaxScanDepth           int             `json:"MaxScanDepth" comment:"最大扫描深度"`
	MinImageNum            int             `json:"MinImageNum" comment:"压缩包或文件夹内，至少有几张图片，才算作书籍"`
	OpenBrowser            bool            `json:"OpenBrowser" comment:"是否同时打开浏览器，windows默认true，其他默认false"`
	Password               string          `json:"Password" comment:"启用登陆后，登录界面需要的密码。"`
	Port                   int             `json:"Port" comment:"Comigo设置文件(config.toml)，可保存在：HomeDirectory（$HOME/.config/comigo/config.toml）、WorkingDirectory（当前执行目录）、ProgramDirectory（程序所在目录）下。可用“comi --config-save”生成本文件\n网页服务端口，此项配置不支持热重载"`
	PrintAllPossibleQRCode bool            `json:"PrintAllPossibleQRCode" comment:"扫描完成后，打印所有可能的阅读链接二维码"`
	Stores                 []stores.Store  `json:"BookStores" toml:"-" comment:"书库设置"`
	SupportFileType        []string        `json:"SupportFileType" comment:"支持的书籍压缩包后缀"`
	SupportMediaType       []string        `json:"SupportMediaType" comment:"扫描压缩包时，用于统计图片数量的图片文件后缀"`
	SupportTemplateFile    []string        `json:"SupportTemplateFile" comment:"支持的模板文件类型，默认为html"`
	Timeout                int             `json:"Timeout" comment:"启用登陆后，cookie过期的时间。单位为分钟。默认180分钟后过期。"`
	TimeoutLimitForScan    int             `json:"TimeoutLimitForScan" comment:"扫描文件时，超过几秒钟，就放弃扫描这个文件，避免卡在特殊文件上"`
	UploadDirOption        UploadDirOption `json:"UploadDirOption" comment:"上传目录的位置选项：0-当前执行目录，1-第一个书库目录，2-指定上传路径"`
	UploadPath             string          `json:"UploadPath" comment:"指定上传路径时，上传文件的存储位置"`
	UseCache               bool            `json:"UseCache" comment:"开启本地图片缓存，可以加快二次读取，但会占用硬盘空间"`
	Username               string          `json:"Username" comment:"启用登陆后，登录界面需要的用户名。"`
	ZipFileTextEncoding    string          `json:"ZipFileTextEncoding" comment:"非utf-8编码的ZIP文件，尝试用什么编码解析，默认GBK"`
}

// CheckLocalStores 检查本地书库路径是否已存在
func (c *Config) CheckLocalStores(path string) bool {
	for _, store := range c.Stores {
		if store.Type == stores.Local && store.Local.Path == path {
			return true
		}
	}
	return false
}

// AddLocalStore 添加本地书库(单个路径)
func (c *Config) AddLocalStore(path string) {
	if c.CheckLocalStores(path) {
		return
	}
	c.Stores = append(c.Stores, stores.Store{
		Type: stores.Local,
		Local: stores.LocalOption{
			Path: path,
		},
	})
	c.LocalStores = c.LocalStoresList()
}

// AddLocalStores 添加本地书库（多个路径）
func (c *Config) AddLocalStores(path []string) {
	for _, p := range path {
		if !c.CheckLocalStores(p) {
			c.AddLocalStore(p)
		}
	}
}

// ReplaceLocalStores 替换现有的“本地”类型的书库，保留其他类型的书库
func (c *Config) ReplaceLocalStores(pathList []string) {
	var newStores []stores.Store
	for _, store := range c.Stores {
		if store.Type != stores.Local {
			newStores = append(newStores, store)
		}
	}
	c.Stores = newStores
	c.AddLocalStores(pathList)
	c.LocalStores = c.LocalStoresList()
}

// LocalStoresList 获取本地书库列表
func (c *Config) LocalStoresList() []string {
	var localStoresList []string
	for _, store := range c.Stores {
		if store.Type == stores.Local {
			localStoresList = append(localStoresList, store.Local.Path)
		}
	}
	return localStoresList
}

// UpdateConfig 更新配置。 使用 JSON 反序列化将更新的配置解析为映射，遍历映射并更新配置，减少重复的代码。
func UpdateConfig(config *Config, jsonString string) (*Config, error) {
	oldConfig := *config
	var updates map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &updates); err != nil {
		logger.Infof("Failed to unmarshal JSON: %v", err)
		return &oldConfig, err
	}
	for key, value := range updates {
		switch key {
		case "Port":
			if v, ok := value.(float64); ok {
				config.Port = int(v)
			}
		case "Host":
			if v, ok := value.(string); ok {
				config.Host = v
			}
		case "LocalStores":
			if v, ok := value.([]interface{}); ok {
				var storeList []string
				for _, s := range v {
					if str, ok := s.(string); ok {
						storeList = append(storeList, str)
					}
				}
				config.ReplaceLocalStores(storeList)
			}
		case "UseCache":
			if v, ok := value.(bool); ok {
				config.UseCache = v
			}
		case "CachePath":
			if v, ok := value.(string); ok {
				config.CachePath = v
			}
		case "ClearCacheExit":
			if v, ok := value.(bool); ok {
				config.ClearCacheExit = v
			}
		case "UploadPath":
			if v, ok := value.(string); ok {
				config.UploadPath = v
			}
		case "EnableUpload":
			if v, ok := value.(bool); ok {
				config.EnableUpload = v
			}
		case "EnableDatabase":
			if v, ok := value.(bool); ok {
				config.EnableDatabase = v
			}
		case "ClearDatabaseWhenExit":
			if v, ok := value.(bool); ok {
				config.ClearDatabaseWhenExit = v
			}
		case "OpenBrowser":
			if v, ok := value.(bool); ok {
				config.OpenBrowser = v
			}
		case "DisableLAN":
			if v, ok := value.(bool); ok {
				config.DisableLAN = v
			}
		case "DefaultMode":
			if v, ok := value.(string); ok {
				config.DefaultMode = v
			}
		case "LogToFile":
			if v, ok := value.(bool); ok {
				config.LogToFile = v
			}
		case "MaxScanDepth":
			if v, ok := value.(float64); ok {
				config.MaxScanDepth = int(v)
			}
		case "MinImageNum":
			if v, ok := value.(float64); ok {
				config.MinImageNum = int(v)
			}
		case "ZipFileTextEncoding":
			if v, ok := value.(string); ok {
				config.ZipFileTextEncoding = v
			}
		case "ExcludePath":
			if v, ok := value.([]interface{}); ok {
				var paths []string
				for _, s := range v {
					if str, ok := s.(string); ok {
						paths = append(paths, str)
					}
				}
				config.ExcludePath = paths
			}
		case "SupportMediaType":
			if v, ok := value.([]interface{}); ok {
				var mediaTypes []string
				for _, s := range v {
					if str, ok := s.(string); ok {
						mediaTypes = append(mediaTypes, str)
					}
				}
				config.SupportMediaType = mediaTypes
			}
		case "SupportFileType":
			if v, ok := value.([]interface{}); ok {
				var fileTypes []string
				for _, s := range v {
					if str, ok := s.(string); ok {
						fileTypes = append(fileTypes, str)
					}
				}
				config.SupportFileType = fileTypes
			}
		case "TimeoutLimitForScan":
			if v, ok := value.(float64); ok {
				config.TimeoutLimitForScan = int(v)
			}
		case "PrintAllPossibleQRCode":
			if v, ok := value.(bool); ok {
				config.PrintAllPossibleQRCode = v
			}
		case "Debug":
			if v, ok := value.(bool); ok {
				config.Debug = v
			}
		case "Username":
			if v, ok := value.(string); ok {
				config.Username = v
			}
		case "Password":
			if v, ok := value.(string); ok {
				config.Password = v
			}
		case "Timeout":
			if v, ok := value.(float64); ok {
				config.Timeout = int(v)
			}
		case "GenerateMetaData":
			if v, ok := value.(bool); ok {
				config.GenerateMetaData = v
			}
		default:
			logger.Infof("Unknown config key: %s", key)
		}
	}
	return &oldConfig, nil
}

// SetByExecutableFilename 通过执行文件名设置默认网页模板参数
func (c *Config) SetByExecutableFilename() {
	// 获取可执行文件的名称
	filenameWithSuffix := path.Base(os.Args[0])
	fileSuffix := path.Ext(filenameWithSuffix)
	filenameWithoutSuffix := strings.TrimSuffix(filenameWithSuffix, fileSuffix)

	// 获取可执行文件所在目录
	executablePath, err := os.Executable()
	if err != nil {
		logger.Infof("Error getting executable path: %s", err)
		return
	}
	executableDir := filepath.Dir(executablePath)

	if c.Debug {
		logger.Infof("Executable Name: %s", filenameWithoutSuffix)
		logger.Infof("Executable Path: %s", executableDir)
	}
}

var Cfg = Config{
	ConfigPath:            "",
	CachePath:             "",
	ClearCacheExit:        true,
	ClearDatabaseWhenExit: true,
	DisableLAN:            false,
	DefaultMode:           "scroll",
	EnableUpload:          true,
	EnableDatabase:        false,
	EnableTLS:             false,
	ExcludePath:           []string{".comigo", ".idea", ".vscode", ".git", "node_modules", "flutter_ui", "$RECYCLE.BIN", "System Volume Information", ".cache"},
	Host:                  "",
	LogToFile:             false,
	MaxScanDepth:          4,
	MinImageNum:           3,
	OpenBrowser:           true,
	Port:                  1234,
	Password:              "",
	Stores: []stores.Store{
		{
			Type: stores.SMB,
			Smb: stores.SMBOption{
				Host:      os.Getenv("SMB_HOST"),
				Port:      445,
				Username:  os.Getenv("SMB_USER"),
				Password:  os.Getenv("SMB_PASS"),
				ShareName: os.Getenv("SMB_PATH"),
			},
		},
	},
	SupportFileType:     []string{".zip", ".tar", ".rar", ".cbr", ".cbz", ".epub", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz", ".tar.lz4", ".tlz4", ".tar.sz", ".tsz", ".bz2", ".gz", ".lz4", ".sz", ".xz", ".mp4", ".webm", ".pdf", ".m4v", ".flv", ".avi", ".mp3", ".wav", ".wma", ".ogg"},
	SupportMediaType:    []string{".jpg", ".jpeg", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".gif", ".apng", ".bmp", ".webp", ".ico", ".heic", ".heif", ".avif"},
	SupportTemplateFile: []string{".html"},
	UseCache:            true,
	UploadPath:          "",
	Username:            "comigo",
	ZipFileTextEncoding: "",
}

func GetCfg() *Config {
	return &Cfg
}

func GetConfigPath() string {
	return Cfg.ConfigPath
}

func SetConfigPath(path string) {
	//检查路径是否存在
	if !util.PathExists(path) {
		logger.Info("Invalid config file path.")
		return
	}
	Cfg.ConfigPath = path
}

func GetCachePath() string {
	return Cfg.CachePath
}

func SetCachePath(path string) {
	if !util.PathExists(path) {
		logger.Info("Invalid cache path.")
		return
	}
	Cfg.CachePath = path
}

func AutoSetCachePath() {
	//手动设置的临时文件夹
	if Cfg.CachePath != "" && util.IsExist(Cfg.CachePath) && util.ChickIsDir(Cfg.CachePath) {
		Cfg.CachePath = path.Join(Cfg.CachePath)
	} else {
		Cfg.CachePath = path.Join(os.TempDir(), "comigo_cache") //使用系统文件夹
	}
	err := os.MkdirAll(Cfg.CachePath, os.ModePerm)
	if err != nil {
		logger.Infof("%s", locale.GetString("temp_folder_error"))
	} else {
		logger.Infof("%s", locale.GetString("temp_folder_path")+Cfg.CachePath)
	}
}

func GetClearDatabaseWhenExit() bool {
	return Cfg.ClearDatabaseWhenExit
}
func SetClearDatabaseWhenExit(clearDatabaseWhenExit bool) {
	Cfg.ClearDatabaseWhenExit = clearDatabaseWhenExit
}

func GetDebug() bool {
	return Cfg.Debug
}

func SetDebug(debug bool) {
	Cfg.Debug = debug
}

func GetEnableUpload() bool {
	return Cfg.EnableUpload
}

func GetEnableDatabase() bool {
	return Cfg.EnableDatabase
}

func SetEnableDatabase(enableDatabase bool) {
	Cfg.EnableDatabase = enableDatabase
}

func GetEnableTLS() bool {
	return Cfg.EnableTLS
}

func GetUploadPath() string {
	return Cfg.UploadPath
}

func SetUploadPath(path string) {
	if (!util.IsDir(path)) || (!util.PathExists(path)) {
		logger.Info("Invalid upload path.")
		return
	}
	Cfg.UploadPath = path
}

func GetUseCache() bool {
	return Cfg.UseCache
}

func SetUseCache(useCache bool) {
	Cfg.UseCache = useCache
}

func GetCertFile() string {
	return Cfg.CertFile
}

func GetClearCacheExit() bool {
	return Cfg.ClearCacheExit
}

func GetLocalStoresList() []string {
	return Cfg.LocalStores
}

func GetLogToFile() bool {
	return Cfg.LogToFile
}

func GetLogFilePath() string {
	return Cfg.LogFilePath
}

func GetLogFileName() string {
	return Cfg.LogFileName
}

func GetStores() []stores.Store {
	return Cfg.Stores
}

func GetMaxScanDepth() int {
	return Cfg.MaxScanDepth
}

func SetClearCacheExit(clearCacheExit bool) {
	Cfg.ClearCacheExit = clearCacheExit
}

func GetDefaultMode() string {
	return Cfg.DefaultMode
}

func GetDisableLAN() bool {
	return Cfg.DisableLAN
}

func SetDisableLAN(disableLAN bool) {
	Cfg.DisableLAN = disableLAN
}

func GetMinImageNum() int {
	return Cfg.MinImageNum
}

func GetTimeoutLimitForScan() int {
	return Cfg.TimeoutLimitForScan
}

func GetExcludePath() []string {
	return Cfg.ExcludePath
}

func GetSupportMediaType() []string {
	return Cfg.SupportMediaType
}

func GetSupportFileType() []string {
	return Cfg.SupportFileType
}

func GetSupportTemplateFile() []string {
	return Cfg.SupportTemplateFile
}

func GetZipFileTextEncoding() string {
	return Cfg.ZipFileTextEncoding
}

func GetOpenBrowser() bool {
	return Cfg.OpenBrowser
}

func GetPrintAllPossibleQRCode() bool {
	return Cfg.PrintAllPossibleQRCode
}

func GetPort() int {
	return Cfg.Port
}

func GetUsername() string {
	return Cfg.Username
}

func GetPassword() string {
	return Cfg.Password
}

func GetTimeout() int {
	return Cfg.Timeout
}

func SetPort(port int) {
	if port < 0 || port > 65535 {
		port = 1234
		logger.Infof("Invalid port number. Using default port: %d", port)
	}
	Cfg.Port = port
}

func GetHost() string {
	return Cfg.Host
}

func SetHost(host string) {
	// 如果主机名为空，使用默认主机名
	if host == "" {
		host = ""
		logger.Infof("Invalid host name. Using default host: %s", host)
	}
	Cfg.Host = host
}

func GetKeyFile() string {
	return Cfg.KeyFile
}

var version = "v0.9.13"

func GetVersion() string {
	return version
}
