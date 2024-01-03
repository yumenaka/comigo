package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/yumenaka/comi/logger"
)

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		logger.Infof("%s", err)
	}
	Config.LogFilePath = home
	Config.LogFileName = "comigo.log"
}
