package config

import (
	"net/http"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/types"
)

var (
	Version = "v0.9.6"
	Srv     *http.Server
	Config  = types.ServerConfig{
		Port:                  1234,
		Host:                  "DefaultHost",
		StoresPath:            []string{},
		SupportFileType:       []string{".zip", ".tar", ".rar", ".cbr", ".cbz", ".epub", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz", ".tar.lz4", ".tlz4", ".tar.sz", ".tsz", ".bz2", ".gz", ".lz4", ".sz", ".xz", ".mp4", ".webm", ".pdf", ".m4v", ".flv", ".avi", ".mp3", ".wav", ".wma", ".ogg"},
		SupportMediaType:      []string{".jpg", ".jpeg", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".gif", ".apng", ".bmp", ".webp", ".ico", ".heic", ".heif", ".avif"},
		ExcludePath:           []string{".comigo", ".idea", ".vscode", ".git", "node_modules", "flutter_ui", "$RECYCLE.BIN", "System Volume Information", ".cache"},
		MaxScanDepth:          4,
		MinImageNum:           3,
		ZipFileTextEncoding:   "",
		OpenBrowser:           true,
		UseCache:              true,
		CachePath:             "",
		ClearCacheExit:        true,
		UploadPath:            "",
		EnableUpload:          true,
		EnableDatabase:        false,
		ClearDatabaseWhenExit: true,
		EnableTLS:             false,
		Username:              "comigo",
		Password:              "",
		DisableLAN:            false,
		DefaultMode:           "scroll",
		LogToFile:             false,
		ConfigPath:            "",
	}
)

func SaveConfig(SaveTo string) {
	//保存配置
	bytes, err := toml.Marshal(Config)
	if err != nil {
		logger.Info("toml.Marshal Error")
		return
	}
	//在命令行打印
	logger.Info("Config Save To " + SaveTo)
	//保存到文件
	if SaveTo == "HomeDir" {
		home, err := homedir.Dir()
		if err != nil {
			logger.Infof("homedir.Dir Error: %s \n", err)
		}
		// 创建目录
		err = os.MkdirAll(path.Join(home, ".config/comigo/"), os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(path.Join(home, ".config/comigo/config.toml"), bytes, 0644)
		if err != nil {
			logger.Infof("os.WriteFile Error: %s \n", err)
		}
	}
	if SaveTo == "NowDir" {
		err = os.WriteFile("config.toml", bytes, 0644)
		if err != nil {
			logger.Infof("os.WriteFile Error: %s \n", err)
		}
	}
	if SaveTo == "ProgramDir" {
		// 获取可执行程序自身的文件路径
		executablePath, err := os.Executable()
		if err != nil {
			logger.Infof("os.Executable Error: %s \n", err)
		}
		err = os.WriteFile(path.Join(executablePath, "config.toml"), bytes, 0644)
		if err != nil {
			logger.Infof("os.WriteFile Error: %s \n", err)
		}
	}
}
