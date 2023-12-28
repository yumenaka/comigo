package cmd

import (
	"os"
	"path"

	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/logger"
)

// initStorePath 添加默认扫描路径
func initStorePath(args []string) {
	wd, err := os.Getwd()
	if err != nil {
		logger.Info("Failed to get working directory:", err)
	}
	logger.Info("Working directory:", wd)

	//没指定路径或文件,同时也配置文件也没设定书库文件夹
	if len(args) == 0 && len(config.Config.StoresPath) == 0 {
		config.Config.StoresPath = append(config.Config.StoresPath, wd)
	}
	//指定了多个路径，就都扫描一遍
	for _, arg := range args {
		config.Config.StoresPath = append(config.Config.StoresPath, arg)
	}

	//启用上传，则添加upload目录
	if config.Config.EnableUpload {
		if config.Config.UploadPath != "" {
			config.Config.StoresPath = append(config.Config.StoresPath, config.Config.UploadPath)
		}
		if config.Config.UploadPath == "" && len(config.Config.StoresPath) > 0 {
			createUploadFolder := true
			for _, checkPath := range config.Config.StoresPath {
				if checkPath == path.Join(config.Config.StoresPath[0], "upload") {
					createUploadFolder = false
				}
			}
			if createUploadFolder {
				config.Config.StoresPath = append(config.Config.StoresPath, path.Join(config.Config.StoresPath[0], "upload"))
			}
		}
	}
}
