package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comi/types"
	"net/http"
	"os"
	"path"
)

var (
	Version = "v0.9.5"
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
		ConfigSaveTo:          "RAM",
		ConfigFileUsed:        "",
	}
)

func SaveConfig() {
	//保存配置
	if Config.ConfigSaveTo == "RAM" {
		fmt.Println("Config Save To RAM")
		return
	}
	bytes, err := toml.Marshal(Config)
	if err != nil {
		fmt.Println("toml.Marshal Error")
		return
	}
	//在命令行打印
	fmt.Println("Config Save To " + Config.ConfigSaveTo)
	//fmt.Printf("Config: %s \n", string(bytes))
	//保存到文件
	if Config.ConfigSaveTo == "HomeDir" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Printf("homedir.Dir Error: %s \n", err)
		}
		// 创建目录
		err = os.MkdirAll(path.Join(home, ".config/comigo/"), os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(path.Join(home, ".config/comigo/config.toml"), bytes, 0644)
		if err != nil {
			fmt.Printf("os.WriteFile Error: %s \n", err)
		}
	}
	if Config.ConfigSaveTo == "NowDir" {
		err = os.WriteFile("config.toml", bytes, 0644)
		if err != nil {
			fmt.Printf("os.WriteFile Error: %s \n", err)
		}
	}
	if Config.ConfigSaveTo == "ProgramDir" {
		// 获取可执行程序自身的文件路径
		executablePath, err := os.Executable()
		if err != nil {
			fmt.Printf("os.Executable Error: %s \n", err)
		}
		err = os.WriteFile(path.Join(executablePath, "config.toml"), bytes, 0644)
		if err != nil {
			fmt.Printf("os.WriteFile Error: %s \n", err)
		}
	}
}
