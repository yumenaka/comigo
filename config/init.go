package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/yumenaka/comigo/util/logger"
)

// home目录 配置
func init() {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Infof("%s", err)
	}
	cfg.LogFilePath = home
	cfg.LogFileName = "comigo.log"
}

// smb配置（TODO:SMB支持）
func init() {
	err := godotenv.Load()
	if err != nil {
		if cfg.Debug {
			logger.Infof("Not found .env file")
		}
	}
	cfg.Stores[0].Smb.Host = os.Getenv("SMB_HOST")
	cfg.Stores[0].Smb.Username = os.Getenv("SMB_USER")
	cfg.Stores[0].Smb.Password = os.Getenv("SMB_PASS")
	cfg.Stores[0].Smb.ShareName = os.Getenv("SMB_SHARE_NAME")
	cfg.Stores[0].Smb.Path = os.Getenv("SMB_PATH")
}
