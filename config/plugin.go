package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// CustomPlugin 用户自定义插件结构
type CustomPlugin struct {
	Name     string // 文件名
	Content  string // 文件内容
	FileType string // html/css/js
	Scope    string // global/shelf/flip/scroll
}

// FrpClientConfig frp客户端配置
type FrpClientConfig struct {
	FrpcCommand      string `comment:"手动设定frpc可执行程序的路径,默认为frpc"`
	ServerAddr       string
	ServerPort       int
	Token            string
	FrpType          string // 本地转发端口设置
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

// ScanUserPlugins 扫描并加载用户自定义插件
func ScanUserPlugins() error {
	// 检查是否启用插件
	if !cfg.EnablePlugin {
		if cfg.Debug {
			logger.Infof(locale.GetString("log_plugin_system_disabled_skip_scan"))
		}
		return nil
	}

	// 获取配置目录
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	pluginsPath := filepath.Join(configDir, "plugins")

	// 检查插件目录是否存在
	if _, err := os.Stat(pluginsPath); os.IsNotExist(err) {
		if cfg.Debug {
			logger.Infof(locale.GetString("log_plugin_dir_not_exist_skip_load"), pluginsPath)
		}
		return nil
	}

	// 清空现有的自定义插件列表
	cfg.CustomPlugins = []CustomPlugin{}

	// 扫描不同范围的插件
	scopes := []struct {
		name string
		path string
	}{
		{"global", pluginsPath},                          // 全局插件：直接在 plugins 目录下
		{"shelf", filepath.Join(pluginsPath, "shelf")},   // shelf 页面插件
		{"flip", filepath.Join(pluginsPath, "flip")},     // flip 页面插件
		{"scroll", filepath.Join(pluginsPath, "scroll")}, // scroll 页面插件
	}

	for _, scope := range scopes {
		if err := loadPluginsFromDir(scope.path, scope.name); err != nil {
			logger.Infof(locale.GetString("log_plugin_scope_load_error"), scope.name, err)
		}
	}

	if cfg.Debug {
		logger.Infof(locale.GetString("log_plugin_custom_loaded_count"), len(cfg.CustomPlugins))
		for _, plugin := range cfg.CustomPlugins {
			logger.Infof(locale.GetString("log_plugin_loaded_item"), plugin.Scope, plugin.Name, plugin.FileType)
		}
	}

	return nil
}

// loadPluginsFromDir 从指定目录加载插件文件
func loadPluginsFromDir(dirPath, scope string) error {
	// 检查目录是否存在
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		if cfg.Debug {
			logger.Infof(locale.GetString("log_plugin_dir_not_exist"), dirPath)
		}
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return nil
	}

	// 读取目录内容
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// 跳过子目录（对于 global 范围，跳过 shelf/flip/scroll 子目录）
		if entry.IsDir() {
			if scope == "global" && (entry.Name() == "shelf" || entry.Name() == "flip" || entry.Name() == "scroll") {
				continue
			}
			continue
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))

		// 只处理 .html, .css, .js 文件
		if ext != ".html" && ext != ".css" && ext != ".js" {
			continue
		}

		// 读取文件内容
		filePath := filepath.Join(dirPath, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			logger.Infof(locale.GetString("log_plugin_read_file_failed"), filePath, err)
			continue
		}

		// 创建插件对象
		plugin := CustomPlugin{
			Name:     fileName,
			Content:  string(content),
			FileType: strings.TrimPrefix(ext, "."),
			Scope:    scope,
		}

		cfg.CustomPlugins = append(cfg.CustomPlugins, plugin)
	}

	return nil
}

// GetCustomPluginsByScope 根据范围获取自定义插件
func GetCustomPluginsByScope(scope string) []CustomPlugin {
	var plugins []CustomPlugin
	for _, plugin := range cfg.CustomPlugins {
		if plugin.Scope == scope {
			plugins = append(plugins, plugin)
		}
	}
	return plugins
}

// LoadBookPlugins 加载特定书籍的插件（按需加载）
// bookID: 书籍ID，如 "aBcE4Fz"
// scope: 范围，如 "flip" 或 "scroll"
// 返回该书籍在指定范围下的插件列表
func LoadBookPlugins(bookID, scope string) ([]CustomPlugin, error) {
	if !cfg.EnablePlugin || bookID == "" {
		return nil, nil
	}

	// 获取配置目录
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	// 构建书籍插件路径：plugins/{scope}/{bookID}/
	bookPluginPath := filepath.Join(configDir, "plugins", scope, bookID)

	// 检查目录是否存在
	if _, err := os.Stat(bookPluginPath); os.IsNotExist(err) {
		return nil, nil // 目录不存在，返回空列表
	}

	// 加载该目录下的插件
	var plugins []CustomPlugin
	entries, err := os.ReadDir(bookPluginPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))

		// 只处理 .html, .css, .js 文件
		if ext != ".html" && ext != ".css" && ext != ".js" {
			continue
		}

		// 读取文件内容
		filePath := filepath.Join(bookPluginPath, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			logger.Infof(locale.GetString("log_plugin_read_book_file_failed"), filePath, err)
			continue
		}

		// 创建插件对象
		plugin := CustomPlugin{
			Name:     fileName,
			Content:  string(content),
			FileType: strings.TrimPrefix(ext, "."),
			Scope:    scope + "/" + bookID, // 例如: "flip/aBcE4Fz"
		}

		plugins = append(plugins, plugin)
	}

	if cfg.Debug && len(plugins) > 0 {
		logger.Infof(locale.GetString("log_plugin_loaded_for_book"), bookID, scope, len(plugins))
	}

	return plugins, nil
}
