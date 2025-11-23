package cmd

import (
	"os"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
)

// CreateStoreUrls  解析命令,生成StoreUrls
func CreateStoreUrls(args []string) {
	// 如果用户指定了扫描路径，就把指定的路径都加入到扫描路径里面
	config.GetCfg().InitStoreUrls()
	// 没指定扫描路径,配置文件也没设置书库文件夹的时候，默认把【当前工作目录】作为扫描路径
	if len(args) == 0 && len(config.GetCfg().StoreUrls) == 0 {
		// 获取当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			logger.Infof("Failed to get working directory:%s", err)
		}
		logger.Infof("Working directory:%s", wd)
		err = config.GetCfg().AddStoreUrl(wd)
		if err != nil {
			logger.Infof("Failed to add working directory to store urls:%s", err)
		}
	}
	// 指定了书库路径，就都扫描一遍
	for key, arg := range args {
		if config.GetCfg().Debug {
			logger.Infof("args[%d]: %s\n", key, arg)
		}
		err := config.GetCfg().AddStoreUrl(arg)
		if err != nil {
			logger.Infof("Failed to add store url from args:%s", err)
		}
	}
	// 如果用户启用上传，且用户指定的上传路径不为空，就把程序预先设定的【默认上传路径】当作书库
	if config.GetCfg().EnableUpload {
		if config.GetCfg().UploadPath != "" {
			// 尝试把上传路径添加为书库里
			err := config.GetCfg().AddStoreUrl(config.GetCfg().UploadPath)
			if err != nil {
				logger.Infof("Failed to add upload path to store urls:%s", err)
			}
		}
		// 如果用户启用上传，但没有指定上传路径
		if config.GetCfg().UploadPath == "" {
			for _, storeUrl := range config.GetCfg().StoreUrls {
				// 把【本地存储】里面的第一个可用路径作为上传路径
				if tools.IsExist(storeUrl) {
					config.SetUploadPath(storeUrl)
					err := config.GetCfg().AddStoreUrl(config.GetCfg().UploadPath)
					if err != nil {
						logger.Infof("Failed to add upload path to store urls:%s", err)
					}
					break
				}
			}
		}
	}
}

// ScanStore 扫描所有书库，取得书籍
func ScanStore() {
	err := scan.InitAllStore(config.GetCfg())
	if err != nil {
		logger.Infof("Failed to scan store path: %v", err)
	}
}
