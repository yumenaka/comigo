package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/yumenaka/comigo/config/stores"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
)

// Config Comigo全局配置
type Config struct {
	AutoRescan             bool            `json:"AutoRescan" comment:"刷新页面时，是否自动重新扫描书籍是否存在"`
	CachePath              string          `json:"CachePath" comment:"本地图片缓存位置，默认系统临时文件夹"`
	CertFile               string          `json:"CertFile" comment:"TLS/SSL 证书文件路径 (default: ~/.config/.comigo/cert.crt)"`
	ClearCacheExit         bool            `json:"ClearCacheExit" comment:"退出程序的时候，清理web图片缓存"`
	ClearDatabaseWhenExit  bool            `json:"ClearDatabaseWhenExit" comment:"启用本地数据库时，扫描完成后，清除不存在的书籍。"`
	ConfigPath             string          `json:"-" toml:"-" comment:"用户指定的的yaml设置文件路径"`
	Debug                  bool            `json:"Debug" comment:"开启Debug模式"`
	DefaultMode            string          `json:"DefaultMode" comment:"默认阅读模式，默认为空，可以设置为scroll或flip"`
	DisableLAN             bool            `json:"DisableLAN" comment:"只在本机提供阅读服务，不对外共享"`
	EnableDatabase         bool            `json:"EnableDatabase" comment:"启用本地数据库，保存扫描到的书籍数据。"`
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
	Password               string          `json:"Password" comment:"登录界面需要的密码。"`
	Port                   int             `json:"Port" comment:"Comigo设置文件(config.toml)，可保存在：HomeDirectory（$HOME/.config/comigo/config.toml）、WorkingDirectory（当前执行目录）、ProgramDirectory（程序所在目录）下。可用“comi --config-save”生成本文件\n网页服务端口"`
	PrintAllPossibleQRCode bool            `json:"PrintAllPossibleQRCode" comment:"扫描完成后，打印所有可能的阅读链接二维码"`
	Stores                 []stores.Store  `json:"BookStores" toml:"-" comment:"书库设置"`
	SupportFileType        []string        `json:"SupportFileType" comment:"支持的书籍压缩包后缀"`
	SupportMediaType       []string        `json:"SupportMediaType" comment:"扫描压缩包时，用于统计图片数量的图片文件后缀"`
	SupportTemplateFile    []string        `json:"SupportTemplateFile" comment:"支持的模板文件类型，默认为html"`
	Timeout                int             `json:"Timeout" comment:"cookie过期的时间。单位为分钟。默认60*24*30分钟后过期。"`
	TimeoutLimitForScan    int             `json:"TimeoutLimitForScan" comment:"扫描文件时，超过几秒钟，就放弃扫描这个文件，避免卡在特殊文件上"`
	UploadDirOption        UploadDirOption `json:"UploadDirOption" comment:"上传目录的位置选项：0-当前执行目录，1-第一个书库目录，2-指定上传路径"`
	UploadPath             string          `json:"UploadPath" comment:"指定上传路径时，上传文件的存储位置"`
	UseCache               bool            `json:"UseCache" comment:"开启本地图片缓存，可以加快二次读取，但会占用硬盘空间"`
	Username               string          `json:"Username" comment:"登录界用的用户名。"`
	ZipFileTextEncoding    string          `json:"ZipFileTextEncoding" comment:"非utf-8编码的ZIP文件，尝试用什么编码解析，默认GBK"`
}

func (c *Config) GetLocalStores() []string {
	return c.LocalStores
}

func (c *Config) GetStores() []stores.Store {
	return c.Stores
}

func (c *Config) GetMaxScanDepth() int {
	return c.MaxScanDepth
}

func (c *Config) GetMinImageNum() int {
	return c.MinImageNum
}

func (c *Config) GetTimeoutLimitForScan() int {
	return c.TimeoutLimitForScan
}

func (c *Config) GetExcludePath() []string {
	return c.ExcludePath
}

func (c *Config) GetSupportMediaType() []string {
	return c.SupportMediaType
}

func (c *Config) GetSupportFileType() []string {
	return c.SupportFileType
}

func (c *Config) GetSupportTemplateFile() []string {
	return c.SupportTemplateFile
}

func (c *Config) GetZipFileTextEncoding() string {
	return c.ZipFileTextEncoding
}

func (c *Config) GetEnableDatabase() bool {
	return c.EnableDatabase
}

func (c *Config) GetClearDatabaseWhenExit() bool {
	return c.ClearDatabaseWhenExit
}

func (c *Config) GetDebug() bool {
	return c.Debug
}

func (c *Config) GetTopStoreName() string {
	if len(c.LocalStores) == 0 {
		return "未设置书库"
	}
	return strings.Join(c.LocalStores, ", ")
}

// SetConfigValue 根据字段名和字符串形式的值，来更新 Config 的相应字段。
// 如果字段名不存在，则返回错误。如果字段类型不在支持范围内，也返回错误。
func (c *Config) SetConfigValue(fieldName, fieldValue string) error {
	// 使用反射获得指向结构体的 Value
	v := reflect.ValueOf(c).Elem()

	// 根据 fieldName 获取对应字段的 reflect.Value
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		return fmt.Errorf("不存在名为 '%s' 的字段", fieldName)
	}
	if !f.CanSet() {
		return fmt.Errorf("无法对字段 '%s' 进行设置", fieldName)
	}

	// 根据字段的类型，进行解析并赋值
	switch f.Kind() {
	case reflect.Bool:
		// ParseBool 返回字符串表示的布尔值。它接受 1、t、T、TRUE、true、True、0、f、F、FALSE、false、False。任何其他值都会返回错误。
		boolVal, err := strconv.ParseBool(fieldValue)
		if err != nil {
			return fmt.Errorf("无法将 '%s' 解析为 bool: %v", fieldValue, err)
		}
		f.SetBool(boolVal)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// ParseInt 以给定基数（0、2 到 36）和位大小（0 到 64）解释字符串 s 并返回相应的值 i。该字符串可以以前导符号开头：“+”或“-”。
		intVal, err := strconv.ParseInt(fieldValue, 10, 64)
		if err != nil {
			return fmt.Errorf("无法将 '%s' 解析为 int: %v", fieldValue, err)
		}
		// 这里为了简单，直接设定 int64。如果结构体字段是 int（非 int64），
		// Go 会自动做一次范围检查和转换；若超范围将报错。
		f.SetInt(intVal)

	case reflect.String:
		f.SetString(fieldValue)

	case reflect.Slice:
		// 针对常见的 []string 做一个示例：
		if f.Type().Elem().Kind() == reflect.String {
			// 假设传进来的字符串用逗号分割
			sliceVal := strings.Split(fieldValue, ",")
			f.Set(reflect.ValueOf(sliceVal))
		} else {
			return errors.New("暂不支持此 slice 的设置(仅支持 []string)")
		}

	default:
		return fmt.Errorf("暂不支持设置字段 '%s' 的类型: %s", fieldName, f.Type().String())
	}

	return nil
}

// getStringSliceField 是一个帮助函数，用于根据字段名，获取可设置的 []string 字段。
func getStringSliceField(c *Config, fieldName string) (reflect.Value, []string, error) {
	// 取得 *Config 的 Value
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return reflect.Value{}, nil, errors.New("必须是一个非空的 *Config 指针")
	}
	// 取得实际元素
	v = v.Elem()

	// 根据 fieldName 获取对应字段
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		return reflect.Value{}, nil, fmt.Errorf("不存在名为 '%s' 的字段", fieldName)
	}
	if !f.CanSet() {
		return reflect.Value{}, nil, fmt.Errorf("无法对字段 '%s' 进行设置", fieldName)
	}
	// 检查字段是否是切片类型
	if f.Kind() != reflect.Slice {
		return reflect.Value{}, nil, fmt.Errorf("字段 '%s' 不是切片类型", fieldName)
	}
	// 检查切片元素类型是否是string
	if f.Type().Elem().Kind() != reflect.String {
		return reflect.Value{}, nil, fmt.Errorf("字段 '%s' 的元素类型不是 string", fieldName)
	}
	// 转换为 []string
	oldSlice := f.Interface().([]string)

	return f, oldSlice, nil
}

// AddStringArrayConfig 往指定的 []string 字段中添加一个新字符串
func (c *Config) AddStringArrayConfig(fieldName, addValue string) ([]string, error) {
	f, oldSlice, err := getStringSliceField(c, fieldName)
	if err != nil {
		return nil, err
	}
	// 检查新元素是否已存在
	for _, v := range oldSlice {
		if v == addValue {
			logger.Infof("AddStringArrayConfig: 字符串 '%s' 已存在", addValue)
			return oldSlice, nil // 已存在则直接返回
		}
	}
	// 若不存在则添加
	newSlice := append(oldSlice, addValue)
	// 赋值给字段
	f.Set(reflect.ValueOf(newSlice))
	return newSlice, nil
}

// DeleteStringArrayConfig 从指定的 []string 字段中删除某个字符串
func (c *Config) DeleteStringArrayConfig(fieldName, deleteValue string) ([]string, error) {
	f, oldSlice, err := getStringSliceField(c, fieldName)
	if err != nil {
		return nil, err
	}

	found := false
	// 过滤掉需要删除的字符串
	newSlice := make([]string, 0, len(oldSlice))
	for _, v := range oldSlice {
		if v == deleteValue {
			found = true
			continue
		}
		newSlice = append(newSlice, v)
	}

	// 如果没找到，返回旧切片即可
	if !found {
		return oldSlice, nil
	}

	// 找到了才赋值给字段
	f.Set(reflect.ValueOf(newSlice))
	return newSlice, nil
}

// UpdateConfigByJson 使用 JSON 字符串反序列化将更新的配置解析为映射，遍历映射并更新配置，减少重复的代码。
func UpdateConfigByJson(jsonString string) error {
	var updates map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &updates); err != nil {
		logger.Infof("Failed to unmarshal JSON: %v", err)
		return err
	}
	for key, value := range updates {
		switch key {
		case "Port":
			if v, ok := value.(float64); ok {
				cfg.Port = int(v)
			}
		case "Host":
			if v, ok := value.(string); ok {
				cfg.Host = v
			}
		case "LocalStores":
			if v, ok := value.([]interface{}); ok {
				var storeList []string
				for _, s := range v {
					if str, ok := s.(string); ok {
						storeList = append(storeList, str)
					}
				}
				ReplaceLocalStores(storeList)
			}
		case "UseCache":
			if v, ok := value.(bool); ok {
				cfg.UseCache = v
			}
		case "CachePath":
			if v, ok := value.(string); ok {
				cfg.CachePath = v
			}
		case "ClearCacheExit":
			if v, ok := value.(bool); ok {
				cfg.ClearCacheExit = v
			}
		case "UploadPath":
			if v, ok := value.(string); ok {
				cfg.UploadPath = v
			}
		case "EnableUpload":
			if v, ok := value.(bool); ok {
				cfg.EnableUpload = v
			}
		case "EnableDatabase":
			if v, ok := value.(bool); ok {
				cfg.EnableDatabase = v
			}
		case "ClearDatabaseWhenExit":
			if v, ok := value.(bool); ok {
				cfg.ClearDatabaseWhenExit = v
			}
		case "OpenBrowser":
			if v, ok := value.(bool); ok {
				cfg.OpenBrowser = v
			}
		case "DisableLAN":
			if v, ok := value.(bool); ok {
				cfg.DisableLAN = v
			}
		case "DefaultMode":
			if v, ok := value.(string); ok {
				cfg.DefaultMode = v
			}
		case "LogToFile":
			if v, ok := value.(bool); ok {
				cfg.LogToFile = v
			}
		case "MaxScanDepth":
			if v, ok := value.(float64); ok {
				cfg.MaxScanDepth = int(v)
			}
		case "MinImageNum":
			if v, ok := value.(float64); ok {
				cfg.MinImageNum = int(v)
			}
		case "ZipFileTextEncoding":
			if v, ok := value.(string); ok {
				cfg.ZipFileTextEncoding = v
			}
		case "ExcludePath":
			if v, ok := value.([]interface{}); ok {
				var paths []string
				for _, s := range v {
					if str, ok := s.(string); ok {
						paths = append(paths, str)
					}
				}
				cfg.ExcludePath = paths
			}
		case "SupportMediaType":
			if v, ok := value.([]interface{}); ok {
				var mediaTypes []string
				for _, s := range v {
					if str, ok := s.(string); ok {
						mediaTypes = append(mediaTypes, str)
					}
				}
				cfg.SupportMediaType = mediaTypes
			}
		case "SupportFileType":
			if v, ok := value.([]interface{}); ok {
				var fileTypes []string
				for _, s := range v {
					if str, ok := s.(string); ok {
						fileTypes = append(fileTypes, str)
					}
				}
				cfg.SupportFileType = fileTypes
			}
		case "TimeoutLimitForScan":
			if v, ok := value.(float64); ok {
				cfg.TimeoutLimitForScan = int(v)
			}
		case "PrintAllPossibleQRCode":
			if v, ok := value.(bool); ok {
				cfg.PrintAllPossibleQRCode = v
			}
		case "Debug":
			if v, ok := value.(bool); ok {
				cfg.Debug = v
			}
		case "Username":
			if v, ok := value.(string); ok {
				cfg.Username = v
			}
		case "Password":
			if v, ok := value.(string); ok {
				cfg.Password = v
			}
		case "Timeout":
			if v, ok := value.(float64); ok {
				cfg.Timeout = int(v)
			}
		case "GenerateMetaData":
			if v, ok := value.(bool); ok {
				cfg.GenerateMetaData = v
			}
		default:
			logger.Infof("Unknown config key: %s", key)
		}
	}
	return nil
}

// SetByExecutableFilename 通过执行文件名设置默认网页模板参数
func SetByExecutableFilename() {
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

	if cfg.Debug {
		logger.Infof("Executable Name: %s", filenameWithoutSuffix)
		logger.Infof("Executable Path: %s", executableDir)
	}
}

// cfg 为全局配置,全局单实例
var cfg = Config{
	ConfigPath:            "",
	CachePath:             "",
	ClearCacheExit:        true,
	ClearDatabaseWhenExit: true,
	DisableLAN:            false,
	DefaultMode:           "scroll",
	EnableUpload:          true,
	EnableDatabase:        false,
	EnableTLS:             false,
	ExcludePath:           []string{"$RECYCLE.BIN", "System Volume Information", ".cache", ".comigo", ".idea", ".vscode", ".git", "node_modules"},
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

// CanAddLocalStores 检查本地书库路径是否可添加
func CanAddLocalStores(path string) bool {
	// 如果文件/文件夹不存在
	if !util.PathExists(path) {
		logger.Infof("Path not exists: %v", path)
		return false
	}

	// 检查本地书库路径是否已存在
	for _, store := range cfg.Stores {
		if store.Type == stores.Local && store.Local.Path == path {
			logger.Infof("Local store already exists: %s", path)
			return false
		}
	}
	return true
}

// AddLocalStore 添加本地书库(单个路径)
func AddLocalStore(path string) {
	if !CanAddLocalStores(path) {
		return
	}
	cfg.Stores = append(cfg.Stores, stores.Store{
		Type: stores.Local,
		Local: stores.LocalOption{
			Path: path,
		},
	})
	cfg.LocalStores = GetLocalStoresList()
}

// AddLocalStores 添加本地书库（多个路径）
func AddLocalStores(path []string) {
	for _, p := range path {
		if CanAddLocalStores(p) {
			AddLocalStore(p)
		}
	}
}

// InitCfgStores 初始化配置文件中的书库
func InitCfgStores() {
	AddLocalStores(cfg.LocalStores)
}

// ReplaceLocalStores 替换现有的“本地”类型的书库，保留其他类型的书库
func ReplaceLocalStores(pathList []string) {
	var newStores []stores.Store
	for _, store := range cfg.Stores {
		if store.Type != stores.Local {
			newStores = append(newStores, store)
		}
	}
	cfg.Stores = newStores
	AddLocalStores(pathList)
	cfg.LocalStores = GetLocalStoresList()
}

func GetLocalStoresList() []string {
	var localStoresList []string
	for _, store := range cfg.Stores {
		if store.Type == stores.Local {
			localStoresList = append(localStoresList, store.Local.Path)
		}
	}
	return localStoresList
}
