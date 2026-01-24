package config

import (
	"os"
	"path/filepath"

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
// 优先级：1. 用户指定的配置文件路径 2. COMIGO_CONFIG_DIR 环境变量 3. 用户主目录 4. 临时目录
func GetConfigDir() (dir string, err error) {
	// 如果配置文件存在
	if cfg.ConfigFile != "" && tools.PathExists(cfg.ConfigFile) {
		return filepath.Dir(cfg.ConfigFile), nil
	}

	// 检查 COMIGO_CONFIG_DIR 环境变量（适用于 Docker 等容器环境）
	if envConfigDir := os.Getenv("COMIGO_CONFIG_DIR"); envConfigDir != "" {
		// 创建目录（如果不存在）
		err = os.MkdirAll(envConfigDir, os.ModePerm)
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_create_config_dir"), err)
			return "", err
		}
		return envConfigDir, nil
	}

	// 获取用户主目录
	home, err := os.UserHomeDir()
	if err != nil {
		// 如果获取Home失败，使用临时目录
		// 不过该目录既不能保证存在，也不能保证具有可访问权限
		tempDir := os.TempDir()
		if tools.PathExists(tempDir) {
			configDir := filepath.Join(tempDir, "comigo")
			// 创建目录（如果不存在）
			err = os.MkdirAll(configDir, os.ModePerm)
			if err != nil {
				logger.Infof(locale.GetString("log_failed_to_create_temp_config_dir"), err)
				return "", err
			}
			return configDir, nil
		}
	}
	configDir := filepath.Join(home, ".config", "comigo")
	// 创建目录（如果不存在）
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_create_config_dir"), err)
		return "", err
	}
	return configDir, nil
}

func AutoSetCacheDir() {
	// 手动设置的临时文件夹
	if cfg.CacheDir != "" && tools.IsExist(cfg.CacheDir) && tools.ChickIsDir(cfg.CacheDir) {
		cfg.CacheDir = filepath.Join(cfg.CacheDir)
	} else {
		cfg.CacheDir = filepath.Join(os.TempDir(), "comigo_cache") // 使用系统文件夹
	}
	err := os.MkdirAll(cfg.CacheDir, os.ModePerm)
	if err != nil {
		logger.Infof("%s", locale.GetString("temp_folder_error"))
	} else {
		logger.Infof("%s", locale.GetString("temp_folder_path")+cfg.CacheDir)
	}
}

// GetJwtSigningKey JWT令牌签名key，目前是用户名+密码(如果两者都设置了的话)
func GetJwtSigningKey() string {
	if cfg.Username == "" || cfg.Password == "" {
		logger.Infof(locale.GetString("log_username_or_password_empty"))
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
		logger.Infof(locale.GetString("log_invalid_port_number"), port)
	}
	cfg.Port = port
}
