package config

import (
	"os"
	"runtime"
	"sync/atomic"

	"github.com/yumenaka/comigo/tools/logger"
)

func newDefaultConfig() Config {
	defaultCfg := Config{
		ConfigFile:            "",
		CacheDir:              "",
		ClearCacheExit:        true,
		ClearDatabaseWhenExit: true,
		DisableLAN:            false,
		EnableUpload:          true,
		EnableDatabase:        false,
		DBType:                "sqlite",
		DBDSN:                 "",
		EnableTLS:             false,
		EnablePlugin:          false,
		ExcludePath:           []string{"$RECYCLE.BIN", "System Volume Information", "node_modules"},
		BasePath:              "",
		Host:                  "",
		LogToFile:             false,
		MaxScanDepth:          4,
		MinImageNum:           3,
		TimeoutLimitForScan:   20,
		OpenBrowser:           true,
		Port:                  1234,
		Password:              "",
		SupportFileType:       []string{".zip", ".tar", ".rar", ".cbr", ".cbz", ".epub", ".mp4", ".m4a", ".webm", ".mov", ".pdf", ".flv", ".avi", ".mp3", ".aac", ".ogg", ".wav", ".wma", ".html", ".htm"},
		SupportMediaType:      []string{".jpg", ".jpeg", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".gif", ".apng", ".bmp", ".webp", ".ico", ".heic", ".heif", ".avif"},
		UseCache:              true,
		Username:              "comigo",
		ZipFileTextEncoding:   "",
		Language:              "auto",
		BuildInPluginList:     []string{"auto_flip", "auto_scroll", "clock", "comigo_xyz", "sample", "sketch_practice"},
		UserPluginList:        []string{},
		EnabledPluginList:     []string{},
	}
	// 在嵌入模式下仍复用默认日志文件位置，避免后续路径判断为空。
	if runtime.GOOS != "js" {
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Infof("%s", err)
		} else {
			defaultCfg.LogFilePath = home
			defaultCfg.LogFileName = "comigo.log"
		}
	}
	return defaultCfg
}

// 配置以不可变快照发布，读取无需持有跨请求的锁。
var currentConfig = func() *atomic.Pointer[Config] {
	value := newDefaultConfig()
	pointer := new(atomic.Pointer[Config])
	pointer.Store(&value)
	return pointer
}()

// ResetConfigForRuntime 在嵌入模式重复启动时重置全局配置，避免保留上一次运行的状态。
func ResetConfigForRuntime() {
	ReplaceConfig(newDefaultConfig())
}
