package config

import (
	"github.com/joho/godotenv"
	"github.com/mitchellh/go-homedir"
	"github.com/yumenaka/comigo/util/logger"
	"os"
)

// home目录 配置
func init() {
	// Find home directory.
	home, err := homedir.Dir()
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
	Config.BookStores[0].Host = os.Getenv("SMB_HOST")
	Config.BookStores[0].Username = os.Getenv("SMB_USER")
	Config.BookStores[0].Password = os.Getenv("SMB_PASS")
	Config.BookStores[0].ShareName = os.Getenv("SMB_SHARE_NAME")
	Config.BookStores[0].Path = os.Getenv("SMB_PATH")
}
