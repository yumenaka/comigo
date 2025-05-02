package config

import (
	"errors"
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

// WriteConfigFile 如果存在本地配置，更新本地配置
func WriteConfigFile() error {
	bytes, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	// HomeDirectory
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

	// 当前执行目录
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

// DefaultConfigLocation 判断当前配置文件应该保存到哪里。
// 逻辑：
//  1. 是否在HomeDirectory已有 config.toml
//  2. 否则是否在WorkingDirectory已有 config.toml
//  3. 否则是否在ProgramDirectory已有 config.toml
//     若都没有，则返回 "HomeDirectory"。
//
// 返回：location字符串
func DefaultConfigLocation() string {
	// 1. 检查HomeDirectory
	home, err := os.UserHomeDir()
	if err != nil {
		// 获取Home失败就先记录一下，但不直接return，继续往后查
		fmt.Println("Warning: failed to get HomeDir:", err)
	} else {
		homePath := path.Join(home, ".config", "comigo", "config.toml")
		if fileExists(homePath) {
			return HomeDirectory
		}
	}

	// 2. 检查WorkingDirectory
	wdPath := "config.toml"
	if fileExists(wdPath) {
		return WorkingDirectory
	}

	// 3. 检查ProgramDirectory
	executable, errExe := os.Executable()
	if errExe != nil {
		// 获取可执行文件路径出错，这种情况也比较少见，先打印日志
		fmt.Println("Warning: failed to get Executable path:", errExe)
	} else {
		progPath := path.Join(path.Dir(executable), "config.toml")
		if fileExists(progPath) {
			return ProgramDirectory
		}
	}

	// 4. 如果都找不到，默认返回 HomeDirectory
	return HomeDirectory
}

// fileExists 封装一个检查文件是否存在的辅助函数
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		// 其它错误比如权限问题，也直接返回false
		// fmt.Println("Warning: cannot access file:", filename, "error:", err)
		return false
	}
	// 确实存在且不是目录
	return !info.IsDir()
}

func SaveConfig(to string) error {
	// 保存配置
	bytes, errMarshal := toml.Marshal(cfg)
	if errMarshal != nil {
		return errMarshal
	}
	logger.Infof("cfg Save To %s", to)
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

func GetWorkingDirectoryConfig() string {
	WorkingDirectoryConfig := ""
	wdPath := "config.toml"
	if fileExists(wdPath) {
		WorkingDirectoryConfig = wdPath
		absPath, errAbs := filepath.Abs(WorkingDirectoryConfig)
		if errAbs == nil {
			WorkingDirectoryConfig = absPath
		}
	}
	return WorkingDirectoryConfig
}

func GetHomeDirectoryConfig() string {
	HomeDirectoryConfig := ""
	home, err := os.UserHomeDir()
	if err == nil {
		homePath := path.Join(home, ".config", "comigo", "config.toml")
		if fileExists(homePath) {
			absPath, errAbs := filepath.Abs(homePath)
			if errAbs == nil {
				HomeDirectoryConfig = absPath
			}
			// 如果取绝对路径出错，就返回相对路径；或者你也可以直接忽略
			HomeDirectoryConfig = homePath
		}
	} else {
		// 获取HomeDir失败，可以做个日志或忽略
		fmt.Println("Warning: failed to get HomeDir:", err)
	}
	return HomeDirectoryConfig
}

func GetProgramDirectoryConfig() string {
	ProgramDirectoryConfig := ""
	exe, errExe := os.Executable()
	if errExe == nil {
		progPath := path.Join(path.Dir(exe), "config.toml")
		if fileExists(progPath) {
			absPath, errAbs := filepath.Abs(progPath)
			if errAbs == nil {
				ProgramDirectoryConfig = absPath
			}
			ProgramDirectoryConfig = progPath
		}
	}
	return ProgramDirectoryConfig
}

func DeleteConfigIn(in string) error {
	logger.Infof("Delete cfg in %s", in)
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
	enableTLS := cfg.CertFile != "" && cfg.KeyFile != ""
	protocol := "http://"
	if enableTLS {
		protocol = "https://"
	}
	// 取得本机的首选出站IP
	OutIP := util.GetOutboundIP().String()
	if cfg.Host == "" {
		return protocol + OutIP + ":" + strconv.Itoa(cfg.Port)
	}
	return protocol + cfg.Host + ":" + strconv.Itoa(cfg.Port)
}
