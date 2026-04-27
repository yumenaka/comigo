package config

import (
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

var (
	Server         *http.Server
	Mutex          sync.Mutex
	ConfigFileLock sync.Mutex
)

func configFilePathForLocation(location string) (string, error) {
	switch location {
	case HomeDirectory:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return path.Join(home, ".config/comigo/config.toml"), nil
	case WorkingDirectory:
		return "config.toml", nil
	case ProgramDirectory:
		executable, err := os.Executable()
		if err != nil {
			return "", err
		}
		return path.Join(path.Dir(executable), "config.toml"), nil
	default:
		return "", errors.New("unknown config location")
	}
}

func writeConfigBytes(filePath string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}
	return writeFileAtomically(filePath, data, 0o644)
}

// UpdateConfigFile 更新当前正在使用的配置文件。
// 如果尚未存在配置文件，则会自动在默认位置创建一份。
func UpdateConfigFile() error {
	if runtime.GOOS == "js" {
		return nil
	}
	ConfigFileLock.Lock()
	defer ConfigFileLock.Unlock()

	bytes, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	targetPath := cfg.ConfigFile
	if targetPath == "" {
		targetPath, err = configFilePathForLocation(DefaultConfigLocation())
		if err != nil {
			return err
		}
	}

	if err := writeConfigBytes(targetPath, bytes); err != nil {
		return err
	}
	cfg.ConfigFile = targetPath
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
	// 只有在非js环境下才需要检查配置文件位置
	if runtime.GOOS == "js" {
		return ""
	}
	// 1. 检查HomeDirectory
	home, err := os.UserHomeDir()
	if err != nil {
		// 获取Home失败就先记录一下，但不直接return，继续往后查
		logger.Infof(locale.GetString("log_warning_failed_to_get_homedir"), err)
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
		logger.Infof(locale.GetString("log_warning_failed_to_get_executable_path"), errExe)
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
		// logger.Info("Warning: cannot access file:", filename, "error:", err)
		return false
	}
	// 确实存在且不是目录
	return !info.IsDir()
}

func SaveConfig(to string) error {
	// 在js环境下
	if runtime.GOOS == "js" {
		return nil
	}
	ConfigFileLock.Lock()
	defer ConfigFileLock.Unlock()

	// 保存配置
	bytes, errMarshal := toml.Marshal(cfg)
	if errMarshal != nil {
		return errMarshal
	}
	logger.Infof(locale.GetString("log_cfg_save_to"), to)
	configPath, err := configFilePathForLocation(to)
	if err != nil {
		return err
	}
	if err := writeConfigBytes(configPath, bytes); err != nil {
		return err
	}
	cfg.ConfigFile = configPath
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
	// 在js环境下
	if runtime.GOOS == "js" {
		return ""
	}
	HomeDirectoryConfig := ""
	home, err := os.UserHomeDir()
	if err == nil {
		homePath := path.Join(home, ".config", "comigo", "config.toml")
		if fileExists(homePath) {
			absPath, errAbs := filepath.Abs(homePath)
			if errAbs == nil {
				HomeDirectoryConfig = absPath
			} else {
				HomeDirectoryConfig = homePath
			}
		}
	} else {
		// 获取HomeDir失败，可以做个日志或忽略
		logger.Info("Warning: failed to get HomeDir:", err)
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
			} else {
				ProgramDirectoryConfig = progPath
			}
		}
	}
	return ProgramDirectoryConfig
}

func DeleteConfigIn(in string) error {
	// 在非js环境下
	if runtime.GOOS == "js" {
		return nil
	}
	ConfigFileLock.Lock()
	defer ConfigFileLock.Unlock()

	logger.Infof(locale.GetString("log_try_delete_cfg_in"), in)
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
	return tools.DeleteFileIfExist(configFile)
}

func writeFileAtomically(filePath string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(filePath)
	tempFile, err := os.CreateTemp(dir, ".config-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tempFile.Name()
	cleanup := func() {
		_ = os.Remove(tmpName)
	}
	if _, err = tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		cleanup()
		return err
	}
	if err = tempFile.Chmod(perm); err != nil {
		_ = tempFile.Close()
		cleanup()
		return err
	}
	if err = tempFile.Close(); err != nil {
		cleanup()
		return err
	}
	if err = os.Rename(tmpName, filePath); err != nil {
		cleanup()
		return err
	}
	return nil
}

func GetQrcodeURL() string {
	enableTLS := cfg.CertFile != "" && cfg.KeyFile != ""
	protocol := "http://"
	if enableTLS {
		protocol = "https://"
	}
	if cfg.Host != "" {
		return protocol + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) + PrefixPath("/")
	}
	// 取得本机的首选出站IP
	OutIP := tools.GetOutboundIP().String()
	return protocol + OutIP + ":" + strconv.Itoa(int(cfg.Port)) + PrefixPath("/")
}

func OpenBrowserIfNeeded() {
	if cfg.OpenBrowser == true {
		protocol := "http://"
		if cfg.EnableTLS {
			protocol = "https://"
		}
		go tools.OpenBrowserByURL(protocol + "127.0.0.1:" + strconv.Itoa(cfg.Port) + PrefixPath("/"))
	}
}
