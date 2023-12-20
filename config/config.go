package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/util"
	"net/http"
	"os"
	"path"

	"github.com/yumenaka/comi/types"
)

var (
	Version = "v0.9.7"
	Srv     *http.Server
	Status  = types.ConfigStatus{}
	Config  = types.ComigoConfig{
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

const (
	HomeDirectory    = "HomeDirectory"
	WorkingDirectory = "WorkingDirectory"
	ProgramDirectory = "ProgramDirectory"
)

// UpdateLocalConfig 如果存在本地配置，更新本地配置
func UpdateLocalConfig() error {
	bytes, err := toml.Marshal(Config)
	if err != nil {
		return err
	}
	//HomeDirectory
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	if util.FileExist(path.Join(home, ".config/comigo/config.toml")) {
		err = os.WriteFile(path.Join(home, ".config/comigo/config.toml"), bytes, 0644)
		if err != nil {
			return err
		}
	}

	//当前执行目录
	if util.FileExist("config.toml") {
		err = os.WriteFile("config.toml", bytes, 0644)
		if err != nil {
			return err
		}
	}

	// 可执行程序自身的路径
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println(executablePath)
		return err
	}
	congigPath := path.Join(path.Dir(executablePath), "config.toml")
	if util.FileExist(congigPath) {
		err = os.WriteFile(congigPath, bytes, 0644)
		if err != nil {
			fmt.Println(path.Join(executablePath, "config.toml"))
			return err
		}
	}
	return nil
}

func SaveConfig(Directory string) error {
	//保存配置
	bytes, err := toml.Marshal(Config)
	if err != nil {
		return err
	}
	logger.Info("Config Save To " + Directory)
	// HomeDirectory 目录
	if Directory == HomeDirectory {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		err = os.MkdirAll(path.Join(home, ".config/comigo/"), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.WriteFile(path.Join(home, ".config/comigo/config.toml"), bytes, 0644)
		if err != nil {
			return err
		}
	}
	//当前执行目录
	if Directory == WorkingDirectory {
		err = os.WriteFile("config.toml", bytes, 0644)
		if err != nil {
			return err
		}
	}
	// 可执行程序自身的文件路径
	if Directory == ProgramDirectory {
		executablePath, err := os.Executable()
		if err != nil {
			fmt.Println(executablePath)
			return err
		}
		congigPath := path.Join(path.Dir(executablePath), "config.toml")
		err = os.WriteFile(congigPath, bytes, 0644)
		if err != nil {
			fmt.Println(path.Join(executablePath, "config.toml"))
			return err
		}
	}
	return nil
}

func DeleteConfigIn(Directory string) (err error) {
	logger.Info("Delete Config in " + Directory)
	var filePath string

	switch Directory {
	case HomeDirectory:
		home, err := homedir.Dir()
		if err == nil {
			filePath = path.Join(home, ".config/comigo/config.toml")
		}
	case WorkingDirectory:
		filePath = "config.toml"
	case ProgramDirectory:
		executablePath, err := os.Executable()
		if err != nil {
			return err
		}
		filePath = path.Join(path.Dir(executablePath), "config.toml")
	}
	if err != nil {
		return err
	}
	return util.DeleteFileIfExist(filePath)
}
