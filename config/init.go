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
	Config.LogFilePath = home
	Config.LogFileName = "comigo.log"
}

// smb配置（TODO:SMB支持）
func init() {
	err := godotenv.Load()
	if err != nil {
		if Config.Debug {
			logger.Infof("Not found .env file")
		}
	}
	Config.Stores[0].Host = os.Getenv("SMB_HOST")
	Config.Stores[0].Username = os.Getenv("SMB_USER")
	Config.Stores[0].Password = os.Getenv("SMB_PASS")
	Config.Stores[0].ShareName = os.Getenv("SMB_SHARE_NAME")
	Config.Stores[0].Path = os.Getenv("SMB_PATH")
}
