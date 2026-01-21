package config

import (
	"os"
	"runtime"

	"github.com/yumenaka/comigo/tools/logger"
)

// home目录 配置
func init() {
	// 在非js环境下
	if runtime.GOOS != "js" {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Infof("%s", err)
		}
		cfg.LogFilePath = home
		cfg.LogFileName = "comigo.log"
	}
}

// cfg 为全局配置,全局单实例
var cfg = Config{
	ConfigFile:            "",
	CacheDir:              "",
	ClearCacheExit:        true,
	ClearDatabaseWhenExit: true,
	DisableLAN:            false,
	EnableUpload:          true,
	EnableDatabase:        false,
	EnableTLS:             false,
	EnablePlugin:          false,
	ExcludePath:           []string{"$RECYCLE.BIN", "System Volume Information", "node_modules"},
	Host:                  "",
	LogToFile:             false,
	MaxScanDepth:          4,
	MinImageNum:           3,
	OpenBrowser:           true,
	Port:                  1234,
	Password:              "",
	SupportFileType:       []string{".zip", ".tar", ".rar", ".cbr", ".cbz", ".epub", ".mp4", ".m4a", ".webm", ".mov", ".pdf", ".flv", ".avi", ".mp3", ".aac", ".ogg", ".wav", ".wma", ".html", ".htm"},
	SupportMediaType:      []string{".jpg", ".jpeg", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".gif", ".apng", ".bmp", ".webp", ".ico", ".heic", ".heif", ".avif"},
	SupportTemplateFile:   []string{".html"},
	UseCache:              true,
	Username:              "comigo",
	ZipFileTextEncoding:   "",
	EnableSingleInstance:  false,
	Language:              "auto",
	// 初始化内置插件列表
	BuildInPluginList: []string{"auto_flip", "auto_scroll", "clock", "comigo_xyz", "sample"},
	// 用户自定义插件插（TODO）
	UserPluginList:    []string{},
	EnabledPluginList: []string{}, // 初始状态，不启用任何插件
}
