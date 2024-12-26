package config

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
)

var Srv *http.Server

// UpdateLocalConfig 如果存在本地配置，更新本地配置
func UpdateLocalConfig() error {
	bytes, err := toml.Marshal(Cfg)
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

const (
	HomeDirectory    = "HomeDirectory"
	WorkingDirectory = "WorkingDirectory"
	ProgramDirectory = "ProgramDirectory"
)

func SaveConfig(to string) error {
	//保存配置
	bytes, errMarshal := toml.Marshal(Cfg)
	if errMarshal != nil {
		return errMarshal
	}
	logger.Infof("Cfg Save To %s", to)
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
	logger.Infof("Delete Cfg in %s", in)
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
	enableTLS := Cfg.CertFile != "" && Cfg.KeyFile != ""
	protocol := "http://"
	if enableTLS {
		protocol = "https://"
	}
	//取得本机的首选出站IP
	OutIP := util.GetOutboundIP().String()
	if Cfg.Host == "" {
		return protocol + OutIP + ":" + strconv.Itoa(Cfg.Port)
	}
	return protocol + Cfg.Host + ":" + strconv.Itoa(Cfg.Port)
}
