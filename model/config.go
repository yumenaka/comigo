package model

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/yumenaka/comigo/config/stores"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
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
	logger.Info("Checking Config ShareName")

	// 初始化
	c.In = "None"
	c.Path.WorkingDirectory = ""
	c.Path.HomeDirectory = ""
	c.Path.ProgramDirectory = ""

	// 配置文件搜索路径及顺序
	type configPath struct {
		name string
		path string
	}
	var configPaths []configPath

	// 获取可执行文件所在目录
	executablePath, err := os.Executable()
	if err != nil {
		return errors.New("error: failed to find executable path")
	}
	programDir := path.Dir(executablePath)

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.New("error: failed to find home directory")
	}

	// 添加搜索路径
	configPaths = append(configPaths, configPath{
		name: "HomeDirectory",
		path: util.GetAbsPath(path.Join(homeDir, ".config/comigo/config.toml")),
	})
	configPaths = append(configPaths, configPath{
		name: "ProgramDirectory",
		path: util.GetAbsPath(path.Join(programDir, "config.toml")),
	})
	configPaths = append(configPaths, configPath{
		name: "WorkingDirectory",
		path: util.GetAbsPath("config.toml"),
	})

	// 按顺序检查配置文件是否存在
	for _, cp := range configPaths {
		if util.IsExist(cp.path) {
			switch cp.name {
			case "HomeDirectory":
				c.Path.HomeDirectory = cp.path
			case "ProgramDirectory":
				c.Path.ProgramDirectory = cp.path
			case "WorkingDirectory":
				c.Path.WorkingDirectory = cp.path
			}
			c.In = cp.name
			return nil
		}
	}
	return nil
}

type UploadDirOption int

type ComigoConfig struct {
	//AutoRescan             bool            `json:"AutoRescan" comment:"刷新页面时，是否自动重新扫描"`
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
func (c *ComigoConfig) CheckLocalStores(path string) bool {
	for _, store := range c.Stores {
		if store.Type == stores.Local && store.Local.Path == path {
			return true
		}
	}
	return false
}

// AddLocalStore 添加本地书库(单个路径)
func (c *ComigoConfig) AddLocalStore(path string) {
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
func (c *ComigoConfig) AddLocalStores(path []string) {
	for _, p := range path {
		if !c.CheckLocalStores(p) {
			c.AddLocalStore(p)
		}
	}
}

// ReplaceLocalStores 替换现有的“本地”类型的书库，保留其他类型的书库
func (c *ComigoConfig) ReplaceLocalStores(pathList []string) {
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
func (c *ComigoConfig) LocalStoresList() []string {
	var localStoresList []string
	for _, store := range c.Stores {
		if store.Type == stores.Local {
			localStoresList = append(localStoresList, store.Local.Path)
		}
	}
	return localStoresList
}

// UpdateConfig 更新配置。 使用 JSON 反序列化将更新的配置解析为映射，遍历映射并更新配置，减少重复的代码。
func UpdateConfig(config *ComigoConfig, jsonString string) (*ComigoConfig, error) {
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
