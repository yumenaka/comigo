package config

import (
	"errors"
	"os"
	"path"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

func GetCfg() *Config {
	return &cfg
}

func CopyCfg() Config {
	return cfg
}

// GetConfigDir 获取配置文件所在目录
func GetConfigDir() (dir string, err error) {
	// 如果未设置配置文件路径，返回系统用户配置目录
	if cfg.ConfigFile == "" {
		// On Unix systems, it returns $XDG_CONFIG_HOME else $HOME/.config.
		// On Darwin, it returns $HOME/Library/Application Support.
		// On Windows, it returns %AppData%.
		// On Plan 9, it returns $home/lib.
		userConfigDir, err := os.UserConfigDir()
		if err == nil {
			userConfigDir = path.Join(userConfigDir, "comigo")
			// 创建目录（如果不存在）
			err = os.MkdirAll(userConfigDir, os.ModePerm)
			if err != nil {
				logger.Infof("Failed to create config dir: %s", err)
				return "", err
			}
			return userConfigDir, nil
		}
		// 如果获取用户配置目录失败，使用临时目录 // 该目录既不能保证存在，也不能保证具有可访问权限
		userConfigDir = os.TempDir()
		if tools.PathExists(userConfigDir) {
			userConfigDir = path.Join(userConfigDir, "comigo")
			// 创建目录（如果不存在）
			err = os.MkdirAll(userConfigDir, os.ModePerm)
			if err != nil {
				logger.Infof("Failed to create temp config dir: %s", err)
				return "", err
			}
			return userConfigDir, nil
		}
	}
	// 如果配置文件存在
	if cfg.ConfigFile != "" && tools.PathExists(cfg.ConfigFile) {
		return path.Dir(cfg.ConfigFile), nil
	}
	return "", errors.New("config dir does not exist")
}

func AutoSetCacheDir() {
	// 手动设置的临时文件夹
	if cfg.CacheDir != "" && tools.IsExist(cfg.CacheDir) && tools.ChickIsDir(cfg.CacheDir) {
		cfg.CacheDir = path.Join(cfg.CacheDir)
	} else {
		cfg.CacheDir = path.Join(os.TempDir(), "comigo_cache") // 使用系统文件夹
	}
	err := os.MkdirAll(cfg.CacheDir, os.ModePerm)
	if err != nil {
		logger.Infof("%s", locale.GetString("temp_folder_error"))
	} else {
		logger.Infof("%s", locale.GetString("temp_folder_path")+cfg.CacheDir)
	}
}

func SetUploadPath(path string) {
	if (!tools.IsDir(path)) || (!tools.PathExists(path)) {
		logger.Info("Invalid upload path.")
		return
	}
	cfg.UploadPath = path
}

// GetJwtSigningKey JWT令牌签名key，目前是用户名+密码(如果两者都设置了的话)
func GetJwtSigningKey() string {
	if cfg.Username == "" || cfg.Password == "" {
		logger.Infof("Username or password is empty. Using default Jwt Signing key.")
		tempStr := cfg.Username + cfg.Password + GetVersion()
		for _, store := range cfg.StoreUrls {
			tempStr = tempStr + store
		}
		base62.EncodeToString([]byte(tools.Md5string(tools.Md5string(tempStr))))
	}
	return cfg.Username + cfg.Password
}

func SetPort(port int) {
	if port < 0 || port > 65535 {
		port = 1234
		logger.Infof("Invalid port number. Using default port: %d", port)
	}
	cfg.Port = port
}

func GetKeyFile() string {
	return cfg.KeyFile
}
