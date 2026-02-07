package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// isRemoteStoreURL 判断是否为远程存储 URL（WebDAV、SMB、SFTP等）
// 这是一个简单的检查，避免循环导入 store 包
func isRemoteStoreURL(storeURL string) bool {
	// 检查是否为本地绝对路径
	if strings.HasPrefix(storeURL, "/") ||
		(len(storeURL) > 2 && storeURL[1] == ':' && (storeURL[2] == '\\' || storeURL[2] == '/')) ||
		strings.HasPrefix(storeURL, "\\\\") {
		return false
	}

	// 检查是否为 file:// 协议
	if strings.HasPrefix(storeURL, "file://") {
		return false
	}

	// 尝试解析为 URL
	u, err := url.Parse(storeURL)
	if err != nil {
		return false
	}

	// 检查 scheme 是否为远程协议
	switch strings.ToLower(u.Scheme) {
	case "webdav", "dav", "davs", "http", "https", "smb", "sftp", "ftp", "ftps", "s3":
		return true
	default:
		return false
	}
}

// getRemoteStoreHost 获取远程存储的主机名
func getRemoteStoreHost(storeURL string) string {
	u, err := url.Parse(storeURL)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// Config Comigo全局配置
type Config struct {
	AutoRescanIntervalMinutes int            `json:"AutoRescanIntervalMinutes" comment:"定期扫描书库间隔。单位为分钟。默认为 0，表示禁用自动定期扫描。"`
	CacheDir                  string         `json:"CacheDir" comment:"本地图片缓存位置，默认系统临时文件夹"`
	ClearCacheExit            bool           `json:"ClearCacheExit" comment:"退出程序的时候，清理web图片缓存"`
	ClearDatabaseWhenExit     bool           `json:"ClearDatabaseWhenExit" comment:"启用本地数据库时，扫描完成后，清除不存在的书籍。"`
	ConfigFile                string         `json:"-" toml:"-" comment:"用户指定的的yaml设置文件路径"`
	ReadOnlyMode              bool           `json:"ReadOnlyMode" comment:"只读模式。禁止网页端更改配置或上传文件。"`
	Debug                     bool           `json:"Debug" comment:"开启Debug模式"`
	EnablePlugin              bool           `json:"EnablePlugin" comment:"启用插件系统"`
	PluginDirectory           string         `json:"-" toml:"-"  comment:"插件存放目录"`
	BuildInPluginList         []string       `json:"-" toml:"-"  comment:"内置插件列表"`
	UserPluginList            []string       `json:"-" toml:"-"  comment:"用户自定义插件列表"`
	EnabledPluginList         []string       `json:"-" toml:"-"  comment:"已启用插件列表"`
	CustomPlugins             []CustomPlugin `json:"-" toml:"-"  comment:"用户自定义插件内容列表"`
	DisableLAN                bool           `json:"DisableLAN" comment:"只在本机提供阅读服务，不对外共享"`
	EnableDatabase            bool           `json:"EnableDatabase" comment:"启用本地数据库，保存扫描到的书籍数据。"`
	EnableTLS                 bool           `json:"EnableTLS" comment:"是否启用HTTPS协议。需要设置证书于key文件。"`
	AutoTLSCertificate        bool           `json:"AutoTLSCertificate" comment:"自动申请、签发 HTTPS 证书（Let's Encrypt）"`
	Host                      string         `json:"Host" comment:"自定义二维码显示的主机名，如果为空，则使用自动检测到的局域网IP地址。自动申请HTTPS证书时，必须设置为公网可访问的域名"`
	KeyFile                   string         `json:"KeyFile" comment:"TLS/SSL key文件路径 (default: ~/.config/.comigo/key.key)"`
	CertFile                  string         `json:"CertFile" comment:"TLS/SSL 证书文件路径 (default: ~/.config/.comigo/cert.crt)"`
	EnableUpload              bool           `json:"EnableUpload" comment:"启用上传功能"`
	ExcludePath               []string       `json:"ExcludePath" comment:"扫描书籍的时候，需要排除的文件或文件夹的名字"`
	GenerateMetaData          bool           `json:"GenerateMetaData" toml:"GenerateMetaData" comment:"生成书籍元数据"`
	StoreUrls                 []string       `json:"StoreUrls" comment:"本地书库路径列表，支持多个路径。可以是本地文件夹或网络书库地址。"` // 书库地址列表
	LogFileName               string         `json:"LogFileName" comment:"Log文件名"`
	LogFilePath               string         `json:"LogFilePath" comment:"Log文件的保存位置"`
	LogToFile                 bool           `json:"LogToFile" comment:"是否保存程序Log到本地文件。默认不保存。"`
	MaxScanDepth              int            `json:"MaxScanDepth" comment:"最大扫描深度"`
	MinImageNum               int            `json:"MinImageNum" comment:"压缩包或文件夹内，至少有几张图片，才算作书籍"`
	OpenBrowser               bool           `json:"OpenBrowser" comment:"是否同时打开浏览器，windows默认true，其他默认false"`
	Password                  string         `json:"Password" comment:"登录界面需要的密码。"`
	Port                      int            `json:"Port" comment:"Comigo设置文件(config.toml)，可保存在：HomeDirectory（$HOME/.config/comigo/config.toml）、WorkingDirectory（当前执行目录）、ProgramDirectory（程序所在目录）下。可用“comi --config-save”生成本文件\n网页服务端口，启用auto TLS时强制使用443端口"`
	PrintAllPossibleQRCode    bool           `json:"PrintAllPossibleQRCode" comment:"扫描完成后，打印所有可能的阅读链接二维码"`
	SupportFileType           []string       `json:"SupportFileType" comment:"支持的书籍压缩包后缀"`
	SupportMediaType          []string       `json:"SupportMediaType" comment:"扫描压缩包时，用于统计图片数量的图片文件后缀"`
	SupportTemplateFile       []string       `json:"SupportTemplateFile" comment:"支持的模板文件类型，默认为html"`
	Timeout                   int            `json:"Timeout" comment:"cookie过期的时间。单位为分钟。默认60*24*30分钟后过期。"`
	TimeoutLimitForScan       int            `json:"TimeoutLimitForScan" comment:"扫描文件时，超过几秒钟，就放弃扫描这个文件，避免卡在特殊文件上"`
	UseCache                  bool           `json:"UseCache" comment:"开启本地图片缓存，可以加快二次读取，但会占用硬盘空间"`
	Username                  string         `json:"Username" comment:"登录界用的用户名。"`
	EnableTailscale           bool           `json:"EnableTailscale" comment:"启用Tailscale网络支持"`
	TailscaleHostname         string         `json:"TailscaleHostname" comment:"Tailscale网络的主机名，默认为comigo"`
	FunnelTunnel              bool           `json:"FunnelTunnel" comment:"启用Tailscale的Funnel模式，允许通过Tailscale公开comigo服务到公网。"`
	FunnelLoginCheck          bool           `json:"funnel_enforce_password" comment:"启用Funnel模式时，强制要求使用密码登录comigo服务。"`
	TailscalePort             int            `json:"TailscalePort" comment:"Tailscale网络的端口，默认为443"`
	TailscaleAuthKey          string         `json:"TailscaleAuthKey" comment:"Tailscale身份验证密钥。另外，也可以将本地环境变量 TS_AUTHKEY 设置为身份验证密钥"`
	ZipFileTextEncoding       string         `json:"ZipFileTextEncoding" comment:"非utf-8编码的ZIP文件，尝试用什么编码解析，默认GBK"`
	EnableSingleInstance      bool           `json:"EnableSingleInstance" comment:"启用单实例模式，确保同一时间只有一个程序实例运行"`
	Language                  string         `json:"Language" comment:"界面语言设置，可选值：auto（自动检测）、zh（中文）、en（英文）、ja（日文），默认为auto"`
	RegisterContextMenu       bool           `json:"RegisterContextMenu" comment:"在 Windows 上注册资源管理器文件夹右键菜单：使用Comigo打开"`
	UnregisterContextMenu     bool           `json:"UnregisterContextMenu" comment:"在 Windows 上卸载资源管理器文件夹右键菜单：使用Comigo打开"`
}

func (c *Config) GetHost() string {
	return c.Host
}

func (c *Config) GetPort() int {
	return c.Port
}

func (c *Config) GetEnableUpload() bool {
	return c.EnableUpload
}

func (c *Config) GetDebug() bool {
	return c.Debug
}

func (c *Config) GetStoreUrls() []string {
	return c.StoreUrls
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

// IsPathOverlapping 检查新路径是否与已有路径重合（完全相同、父子关系）
// 返回值：
// - overlapping: true 表示有重合
// - conflictPath: 冲突的已有路径
// - message: 错误消息
func (c *Config) IsPathOverlapping(newPath string) (overlapping bool, conflictPath string, message string) {
	// 判断新路径是否为远程 URL
	isNewRemote := isRemoteStoreURL(newPath)

	var newPathNormalized string
	if isNewRemote {
		// 远程 URL 直接使用原始 URL（已在 AddStoreUrl 中标准化）
		newPathNormalized = newPath
	} else {
		// 本地路径使用统一的路径标准化函数
		var err error
		newPathNormalized, err = tools.NormalizeAbsPath(newPath)
		if err != nil {
			return true, "", fmt.Sprintf(locale.GetString("err_invalid_store_path"), newPath)
		}
	}

	// 检查是否与已有路径冲突
	for _, existingUrl := range c.StoreUrls {
		isExistingRemote := isRemoteStoreURL(existingUrl)

		var existingNormalized string
		if isExistingRemote {
			existingNormalized = existingUrl
		} else {
			var err error
			existingNormalized, err = tools.NormalizeAbsPath(existingUrl)
			if err != nil {
				// 如果现有路径无法转换，跳过检查
				continue
			}
		}

		// 1. 检查是否完全相同
		if newPathNormalized == existingNormalized {
			return true, existingUrl, fmt.Sprintf(locale.GetString("err_store_url_already_exists_error"), existingUrl)
		}

		// 2. 对于本地路径，检查父子目录关系
		// 远程 URL 之间或远程与本地之间不需要检查父子关系
		if !isNewRemote && !isExistingRemote {
			// 检查新路径是否是已有路径的子目录
			if isSubPath(existingNormalized, newPathNormalized) {
				return true, existingUrl, fmt.Sprintf(locale.GetString("err_store_path_is_subdir_of_existing"), newPath, existingUrl)
			}

			// 检查新路径是否是已有路径的父目录
			if isSubPath(newPathNormalized, existingNormalized) {
				return true, existingUrl, fmt.Sprintf(locale.GetString("err_store_path_is_parent_of_existing"), newPath, existingUrl)
			}
		}
	}

	return false, "", ""
}

// isSubPath 检查 child 是否是 parent 的子路径
// parent 和 child 应该都是已经清理过的绝对路径
func isSubPath(parent, child string) bool {
	// 使用 filepath.Rel 计算相对路径
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	// 如果相对路径是 "."，说明两个路径相同（这种情况在调用前已经检查过）
	if rel == "." {
		return false
	}

	// 如果相对路径以 ".." 开头，说明 child 不在 parent 内
	// 如果相对路径不以 ".." 开头，说明 child 在 parent 内
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".."
}

// StoreUrlIsExits 检查书库URL是否可添加（已弃用，使用 IsPathOverlapping 代替）
func (c *Config) StoreUrlIsExits(url string) bool {
	// 检查书库URL是否已存在
	for _, storeUrl := range c.StoreUrls {
		if storeUrl == url {
			if c.Debug {
				logger.Infof(locale.GetString("log_store_url_already_exists"), storeUrl)
			}
			return true
		}
	}
	return false
}

// AddStoreUrl 添加书库（支持本地路径和远程 URL）
// 本地路径会将相对路径转换为绝对路径，并检查路径冲突
// 远程 URL（WebDAV 等）会在扫描时验证连接
func (c *Config) AddStoreUrl(storeURL string) error {
	var normalizedURL string

	if isRemoteStoreURL(storeURL) {
		// 远程 URL 处理（WebDAV, SMB, SFTP 等）
		// 验证 URL 格式
		u, err := url.Parse(storeURL)
		if err != nil {
			return fmt.Errorf(locale.GetString("err_invalid_store_path")+": %w", err)
		}

		if u.Host == "" {
			return fmt.Errorf(locale.GetString("err_invalid_store_path")+": 缺少主机名", storeURL)
		}

		if c.Debug {
			logger.Infof(locale.GetString("log_add_remote_store"), storeURL, u.Scheme, u.Host)
		}

		// 使用原始 URL（保留认证信息）
		normalizedURL = storeURL
	} else {
		// 本地路径处理
		absPath, err := tools.NormalizeAbsPath(storeURL)
		if err != nil {
			return fmt.Errorf(locale.GetString("err_invalid_store_path"), storeURL)
		}

		// 检查路径是否存在（可选，不强制要求路径必须存在）
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			// 路径不存在，记录警告但不阻止添加
			if c.Debug {
				logger.Infof(locale.GetString("path_not_exist")+": %s", absPath)
			}
		}

		normalizedURL = absPath
	}

	// 使用路径重合检测
	overlapping, conflictPath, message := c.IsPathOverlapping(normalizedURL)
	if overlapping {
		if c.Debug {
			logger.Infof("%s: %s <-> %s", locale.GetString("err_store_path_conflict"), normalizedURL, conflictPath)
		}
		return errors.New(message)
	}

	// 添加到配置
	c.StoreUrls = append(c.StoreUrls, normalizedURL)
	return nil
}

// InitConfigStoreUrls 初始化配置文件中的书库
// 本地路径会将相对路径转换为绝对路径
// 远程 URL（WebDAV 等）保持原样
func (c *Config) InitConfigStoreUrls() {
	// 保存原始的 StoreUrls
	originalUrls := make([]string, len(c.StoreUrls))
	copy(originalUrls, c.StoreUrls)

	// 清空 StoreUrls，然后重新添加（这样可以触发路径标准化）
	c.StoreUrls = []string{}

	for _, storeUrl := range originalUrls {
		// 使用 AddStoreUrl 添加，会自动处理本地路径和远程 URL
		err := c.AddStoreUrl(storeUrl)
		if err != nil {
			// 如果添加失败，记录错误但继续处理其他路径
			logger.Infof(locale.GetString("log_failed_to_add_store_url"), err)
		}
	}
}

// RequiresAuth 是否需要登录
func (c *Config) RequiresAuth() bool {
	return c.Username != "" && c.Password != ""
}

func (c *Config) GetTopStoreName() string {
	if len(c.StoreUrls) == 0 {
		return "No Found Store"
	}
	storeUrls := make([]string, len(c.StoreUrls))
	for i, storeUrl := range c.StoreUrls {
		// 判断是否为远程 URL
		if isRemoteStoreURL(storeUrl) {
			// 远程 URL：尝试提取主机名
			host := getRemoteStoreHost(storeUrl)
			if host != "" {
				storeUrls[i] = host
			} else {
				// 回退：使用 URL 的最后路径部分
				storeUrls[i] = filepath.Base(storeUrl)
			}
		} else {
			// 本地路径：只取路径的最后一部分作为书库名称
			storeUrls[i] = filepath.Base(storeUrl)
		}
	}
	return strings.Join(storeUrls, ", ")
}

// SetConfigValue 更新 Config 的相应字段，如果【fieldName】不存在、或【fieldValue】类型有问题，都返回错误。
func (c *Config) SetConfigValue(fieldName, fieldValue string) error {
	// 使用反射获得指向结构体的 Value
	v := reflect.ValueOf(c).Elem()

	// 根据 fieldName 获取对应字段的 reflect.Value
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		return fmt.Errorf(locale.GetString("err_field_not_exists"), fieldName)
	}
	if !f.CanSet() {
		return fmt.Errorf(locale.GetString("err_field_cannot_set"), fieldName)
	}

	// 根据字段的类型，进行解析并赋值
	switch f.Kind() {
	case reflect.Bool:
		// ParseBool 返回字符串表示的布尔值。它接受 1、t、T、TRUE、true、True、0、f、F、FALSE、false、False。任何其他值都会返回错误。
		boolVal, err := strconv.ParseBool(fieldValue)
		if err != nil {
			return fmt.Errorf(locale.GetString("err_failed_to_parse_bool"), fieldValue, err)
		}
		f.SetBool(boolVal)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// ParseInt 以给定基数（0、2 到 36）和位大小（0 到 64）解释字符串 s 并返回相应的值 i。该字符串可以以前导符号开头：“+”或“-”。
		intVal, err := strconv.ParseInt(fieldValue, 10, 64)
		if err != nil {
			return fmt.Errorf(locale.GetString("err_failed_to_parse_int"), fieldValue, err)
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
			return errors.New(locale.GetString("err_slice_not_supported"))
		}

	default:
		return fmt.Errorf(locale.GetString("err_field_type_not_supported"), fieldName, f.Type().String())
	}

	return nil
}

// getStringSliceField 是一个帮助函数，用于根据字段名，获取可设置的 []string 字段。
func getStringSliceField(c *Config, fieldName string) (reflect.Value, []string, error) {
	// 取得 *Config 的 Value
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return reflect.Value{}, nil, errors.New(locale.GetString("err_must_be_nonempty_config_pointer"))
	}
	// 取得实际元素
	v = v.Elem()

	// 根据 fieldName 获取对应字段
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		return reflect.Value{}, nil, fmt.Errorf(locale.GetString("err_field_not_exists"), fieldName)
	}
	if !f.CanSet() {
		return reflect.Value{}, nil, fmt.Errorf(locale.GetString("err_field_cannot_set"), fieldName)
	}
	// 检查字段是否是切片类型
	if f.Kind() != reflect.Slice {
		return reflect.Value{}, nil, fmt.Errorf(locale.GetString("err_field_not_slice_type"), fieldName)
	}
	// 检查切片元素类型是否是string
	if f.Type().Elem().Kind() != reflect.String {
		return reflect.Value{}, nil, fmt.Errorf(locale.GetString("err_field_element_not_string"), fieldName)
	}
	// 转换为 []string
	oldSlice := f.Interface().([]string)

	return f, oldSlice, nil
}

// AddStringArrayConfig 往指定的 []string 字段中添加一个新字符串
// 对于 StoreUrls 字段，会进行路径验证和标准化
func (c *Config) AddStringArrayConfig(fieldName, addValue string) ([]string, error) {
	// 针对 StoreUrls 做特殊处理：使用 AddStoreUrl 进行路径验证和转换
	if fieldName == "StoreUrls" {
		err := c.AddStoreUrl(addValue)
		if err != nil {
			return c.StoreUrls, err
		}
		return c.StoreUrls, nil
	}

	// 其他字段使用通用逻辑
	f, oldSlice, err := getStringSliceField(c, fieldName)
	if err != nil {
		return nil, err
	}
	// 检查新元素是否已存在
	for _, v := range oldSlice {
		if v == addValue {
			logger.Infof(locale.GetString("log_string_already_exists"), addValue)
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
// 对于 StoreUrls 字段，删除时也需要使用绝对路径进行匹配
func (c *Config) DeleteStringArrayConfig(fieldName, deleteValue string) ([]string, error) {
	// 针对 StoreUrls 做特殊处理：使用统一的路径标准化函数
	if fieldName == "StoreUrls" {
		// 将要删除的路径转换为绝对路径
		deleteAbs, err := tools.NormalizeAbsPath(deleteValue)
		if err == nil {
			// 过滤掉需要删除的路径
			newSlice := make([]string, 0, len(c.StoreUrls))
			found := false
			for _, v := range c.StoreUrls {
				vAbs, err := tools.NormalizeAbsPath(v)
				if err == nil && vAbs == deleteAbs {
					found = true
					continue
				}
				newSlice = append(newSlice, v)
			}

			if found {
				c.StoreUrls = newSlice
			}
			return c.StoreUrls, nil
		}
	}

	// 其他字段使用通用逻辑
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

// IsPluginEnabled 检查指定插件是否已启用
func (c *Config) IsPluginEnabled(pluginName string) bool {
	for _, enabledPlugin := range c.EnabledPluginList {
		if enabledPlugin == pluginName {
			return true
		}
	}
	return false
}

// AddPlugin 启用指定插件（添加到 EnabledPluginList）
func (c *Config) AddPlugin(pluginName string) error {
	// 检查插件是否已经在启用列表中
	if c.IsPluginEnabled(pluginName) {
		return nil // 已经启用，直接返回
	}
	// 添加到启用列表
	c.EnabledPluginList = append(c.EnabledPluginList, pluginName)
	return nil
}

// DisablePlugin 禁用指定插件（从 EnabledPluginList 移除）
func (c *Config) DisablePlugin(pluginName string) error {
	// 使用 DeleteStringArrayConfig 来删除
	_, err := c.DeleteStringArrayConfig("EnabledPluginList", pluginName)
	return err
}

// UpdateConfigByJson 使用 JSON 字符串反序列化将更新的配置解析为映射，遍历映射并更新配置
// 使用反射自动处理字段赋值，减少重复代码
func UpdateConfigByJson(jsonString string) error {
	var updates map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &updates); err != nil {
		logger.Infof(locale.GetString("log_failed_to_unmarshal_json"), err)
		return err
	}

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for key, value := range updates {
		// 查找字段
		field := v.FieldByName(key)
		if !field.IsValid() || !field.CanSet() {
			logger.Infof(locale.GetString("log_unknown_config_key"), key)
			continue
		}

		// 根据字段类型进行赋值
		if err := setFieldValue(field, value); err != nil {
			logger.Infof(locale.GetString("log_failed_to_set_field"), key, err)
			continue
		}

		// 特殊处理：Language 字段更新后的附加逻辑可在此处理
		_ = t // 避免未使用警告
	}
	return nil
}

// setFieldValue 根据字段类型设置值
func setFieldValue(field reflect.Value, value interface{}) error {
	if value == nil {
		return nil
	}

	switch field.Kind() {
	case reflect.Bool:
		if v, ok := value.(bool); ok {
			field.SetBool(v)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// JSON 数字默认解析为 float64
		if v, ok := value.(float64); ok {
			field.SetInt(int64(v))
		}
	case reflect.String:
		if v, ok := value.(string); ok {
			field.SetString(v)
		}
	case reflect.Slice:
		// 处理 []string 类型
		if field.Type().Elem().Kind() == reflect.String {
			if arr, ok := value.([]interface{}); ok {
				strSlice := make([]string, 0, len(arr))
				for _, item := range arr {
					if str, ok := item.(string); ok {
						strSlice = append(strSlice, str)
					}
				}
				field.Set(reflect.ValueOf(strSlice))
			}
		}
	}
	return nil
}
