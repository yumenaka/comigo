package config

import (
	"fmt"
	"github.com/yumenaka/comigo/config/stores"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
)

var (
	Version = "v0.9.12"
	Srv     *http.Server
	Status  = entity.ConfigStatus{}
	Config  = entity.ComigoConfig{
		Port: 1234,
		Host: "DefaultHost",
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
		SupportFileType:       []string{".zip", ".tar", ".rar", ".cbr", ".cbz", ".epub", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz", ".tar.lz4", ".tlz4", ".tar.sz", ".tsz", ".bz2", ".gz", ".lz4", ".sz", ".xz", ".mp4", ".webm", ".pdf", ".m4v", ".flv", ".avi", ".mp3", ".wav", ".wma", ".ogg"},
		SupportMediaType:      []string{".jpg", ".jpeg", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".gif", ".apng", ".bmp", ".webp", ".ico", ".heic", ".heif", ".avif"},
		SupportTemplateFile:   []string{".html"},
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
	confDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	if util.FileExist(filepath.Join(confDir, "comigo", "config.toml")) {
		err = os.WriteFile(filepath.Join(confDir, "comigo", "config.toml"), bytes, 0644)
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
	executable, err := os.Executable()
	if err != nil {
		fmt.Println(executable)
		return err
	}
	p := path.Join(path.Dir(executable), "config.toml")
	if util.FileExist(p) {
		err = os.WriteFile(p, bytes, 0644)
		if err != nil {
			fmt.Println(path.Join(executable, "config.toml"))
			return err
		}
	}
	return nil
}

func SaveConfig(to string) error {
	//保存配置
	bytes, errMarshal := toml.Marshal(Config)
	if errMarshal != nil {
		return errMarshal
	}
	logger.Infof("Config Save To %s", to)
	switch to {
	case HomeDirectory:
		home, err := os.UserHomeDir()
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
	case WorkingDirectory:
		err := os.WriteFile("config.toml", bytes, 0644)
		if err != nil {
			return err
		}
	case ProgramDirectory:
		executable, err := os.Executable()
		if err != nil {
			fmt.Println(executable)
			return err
		}
		p := path.Join(path.Dir(executable), "config.toml")
		err = os.WriteFile(p, bytes, 0644)
		if err != nil {
			fmt.Println(path.Join(executable, "config.toml"))
			return err
		}
	}
	return nil
}

func DeleteConfigIn(in string) error {
	logger.Infof("Delete Config in %s", in)
	var configFile string
	switch in {
	case HomeDirectory:
		home, errHomeDirectory := os.UserHomeDir()
		if errHomeDirectory != nil {
			return errHomeDirectory
		}
		configFile = path.Join(home, ".config/comigo/config.toml")
	case WorkingDirectory:
		configFile = "config.toml"
	case ProgramDirectory:
		executable, errExecutable := os.Executable()
		if errExecutable != nil {
			return errExecutable
		}
		configFile = path.Join(path.Dir(executable), "config.toml")
	}
	return util.DeleteFileIfExist(configFile)
}

func GetQrcodeURL() string {
	enableTLS := Config.CertFile != "" && Config.KeyFile != ""
	protocol := "http://"
	if enableTLS {
		protocol = "https://"
	}
	//取得本机的首选出站IP
	OutIP := util.GetOutboundIP().String()
	if Config.Host == "DefaultHost" {
		return protocol + OutIP + ":" + strconv.Itoa(Config.Port)
	}
	return protocol + Config.Host + ":" + strconv.Itoa(Config.Port)
}
