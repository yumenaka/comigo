package config

import (
	"os"
	"runtime"

	"github.com/yumenaka/comigo/util/logger"
)

// home目录 配置
func init() {
	// 在非js环境下
	if runtime.GOOS == "js" {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Infof("%s", err)
		}
		cfg.LogFilePath = home
		cfg.LogFileName = "comigo.log"
	}
}
