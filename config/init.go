package config

import (
	"os"
	"runtime"

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
		EnableTLS:             false,
		EnablePlugin:          false,
		ExcludePath:           []string{"$RECYCLE.BIN", "System Volume Information", "node_modules"},
		Host:                  "",
		LoginProtection:       false,
		LogToFile:             false,
		MaxScanDepth:          4,
		MinImageNum:           3,
		EnableOAuthLogin:      false,
		OpenBrowser:           true,
		OAuthProviderType:     OAuthProviderTypeGitHub,
		OAuthProviderName:     "",
		OAuthClientID:         "",
		OAuthClientSecret:     "",
		OAuthAuthURL:          "",
		OAuthTokenURL:         "",
		OAuthUserInfoURL:      "",
		OAuthRedirectURL:      "",
		OAuthScopes:           []string{},
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

// cfg 为全局配置,全局单实例
var cfg = newDefaultConfig()

// ResetConfigForRuntime 在嵌入模式重复启动时重置全局配置，避免保留上一次运行的状态。
func ResetConfigForRuntime() {
	cfg = newDefaultConfig()
}
