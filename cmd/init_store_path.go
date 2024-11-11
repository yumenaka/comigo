package cmd

import (
	"github.com/yumenaka/comigo/routers/handlers"
	"os"
	"path/filepath"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/logger"
)

// initStorePath 添加默认扫描路径
func initStorePath(args []string) {
	//没指定扫描路径,配置文件也没设置书库文件夹的时候，默认把【当前工作目录】作为扫描路径
	if len(args) == 0 && len(config.Config.LocalStores) == 0 {
		//获取当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			logger.Infof("Failed to get working directory:%s", err)
		}
		logger.Infof("Working directory:%s", wd)
		config.Config.LocalStores = append(config.Config.LocalStores, wd)
	}
	//指定了多个路径，就都扫描一遍
	for _, arg := range args {
		config.Config.LocalStores = append(config.Config.LocalStores, arg)
	}

	//如果用户启用上传，且用户指定的上传路径不为空，就把上传路径也加入到扫描路径
	if config.Config.EnableUpload {
		if config.Config.UploadPath != "" {
			//判断上传路径是否已经在扫描路径里面了
			for _, store := range config.Config.LocalStores {
				//如果用户指定的上传路径，已经在扫描路径里面了，就不需要添加
				if store == config.Config.UploadPath {
					return
				}
			}
			//把上传路径添加到扫描路径里面去
			config.Config.LocalStores = append(config.Config.LocalStores, config.Config.UploadPath)
		}
		if config.Config.UploadPath == "" {
			config.Config.UploadPath = filepath.Join(config.Config.LocalStores[0], "upload")
		}
	}
	//把扫描路径设置，传递给handlers包
	handlers.ConfigEnableUpload = &config.Config.EnableUpload
	handlers.ConfigUploadPath = &config.Config.UploadPath
}
