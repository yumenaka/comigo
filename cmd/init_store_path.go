package cmd

import (
	"os"
	"path"

	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/util/logger"
)

// initStorePath 添加默认扫描路径
func initStorePath(args []string) {
	wd, err := os.Getwd()
	if err != nil {
		logger.Infof("Failed to get working directory:%s", err)
	}
	logger.Infof("Working directory:%s", wd)

	//没指定路径或文件,同时也配置文件也没设定书库文件夹
	if len(args) == 0 && len(config.Config.LocalStores) == 0 {
		config.Config.LocalStores = append(config.Config.LocalStores, wd)
	}
	//指定了多个路径，就都扫描一遍
	for _, arg := range args {
		config.Config.LocalStores = append(config.Config.LocalStores, arg)
	}

	//启用上传，则添加upload目录
	if config.Config.EnableUpload {
		if config.Config.UploadPath != "" {
			config.Config.LocalStores = append(config.Config.LocalStores, config.Config.UploadPath)
		}
		if config.Config.UploadPath == "" && len(config.Config.LocalStores) > 0 {
			createUploadFolder := true
			for _, checkPath := range config.Config.LocalStores {
				if checkPath == path.Join(config.Config.LocalStores[0], "upload") {
					createUploadFolder = false
				}
			}
			if createUploadFolder {
				config.Config.LocalStores = append(config.Config.LocalStores, path.Join(config.Config.LocalStores[0], "upload"))
			}
		}
	}
}
