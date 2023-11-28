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
	config.Config.StoresPath = append(config.Config.StoresPath, wd)

	//没指定路径或文件的情况下
	if len(args) != 0 {
		//指定了多个参数的话，都扫描一遍
		for _, arg := range args {
			config.Config.StoresPath = append(config.Config.StoresPath, arg)
		}
	}

	//启用上传，则添加upload目录
	if config.Config.EnableUpload {
		if config.Config.UploadPath != "" {
			config.Config.StoresPath = append(config.Config.StoresPath, config.Config.UploadPath)
		} else {
			config.Config.StoresPath = append(config.Config.StoresPath, path.Join(wd, "upload"))
		}
	}
}
