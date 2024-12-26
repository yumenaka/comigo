package cmd

import (
	"os"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/handlers"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
)

// initStorePath 添加默认扫描路径
func initStorePath(args []string) {
	//如果用户指定了扫描路径，就把指定的路径都加入到扫描路径里面
	config.Cfg.AddLocalStores(config.Cfg.LocalStores)
	//没指定扫描路径,配置文件也没设置书库文件夹的时候，默认把【当前工作目录】作为扫描路径
	if len(args) == 0 && len(config.Cfg.LocalStoresList()) == 0 {
		//获取当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			logger.Infof("Failed to get working directory:%s", err)
		}
		logger.Infof("Working directory:%s", wd)
		config.Cfg.AddLocalStore(wd)
	}
	//指定了多个路径，就都扫描一遍
	for _, arg := range args {
		config.Cfg.AddLocalStore(arg)
	}

	//如果用户启用上传，且用户指定的上传路径不为空，就把上传路径也加入到扫描路径
	if config.Cfg.EnableUpload {
		if config.Cfg.UploadPath != "" {
			//判断上传路径是否已经在扫描路径里面了
			for _, store := range config.Cfg.LocalStoresList() {
				//如果用户指定的上传路径，已经在扫描路径里面了，就不需要添加
				if store == config.Cfg.UploadPath {
					return
				}
			}
			//把上传路径添加到扫描路径里面去
			config.Cfg.AddLocalStore(config.Cfg.UploadPath)
		}
		//如果用户启用上传，但是用户没有指定上传路径，就把【本地存储】里面的第一个路径作为上传路径
		if config.Cfg.UploadPath == "" {
			for _, store := range config.Cfg.LocalStoresList() {
				if util.IsExist(store) {
					config.Cfg.UploadPath = store
					config.Cfg.AddLocalStore(config.Cfg.UploadPath)
					break
				}
			}
		}
	}
	//把扫描路径设置，传递给handlers包
	handlers.ConfigEnableUpload = &config.Cfg.EnableUpload
	handlers.ConfigUploadPath = &config.Cfg.UploadPath
}
