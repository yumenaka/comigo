package cmd

import (
	"os"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
)

// SetCwdAsScanPath  当没有指定扫描路径时，把当前工作目录作为扫描路径
func SetCwdAsScanPathIfNeed() {
	if len(config.GetCfg().StoreUrls) == 0 {
		// 获取当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_get_working_directory"), err)
		}
		logger.Infof(locale.GetString("log_working_directory"), wd)
		err = config.GetCfg().AddStoreUrl(wd)
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_add_working_directory_to_store_urls"), err)
		}
	}
}

// AddStoreUrls  解析命令行参数,作为路径添加到StoreUrls里
func AddStoreUrls(urls []string) {
	for key, url := range urls {
		if config.GetCfg().Debug {
			logger.Infof(locale.GetString("log_args_index")+"\n", key, url)
		}
		err := config.GetCfg().AddStoreUrl(url)
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_add_store_url_from_args"), err)
		}
	}
}

// SetUploadPath 设置上传路径
func SetUploadPath(args []string) {
	// 如果用户启用上传，且用户指定的上传路径不为空，就把程序预先设定的【默认上传路径】当作书库
	if config.GetCfg().EnableUpload {
		if config.GetCfg().UploadPath != "" && !config.GetCfg().StoreUrlIsExits(config.GetCfg().UploadPath) {
			// 把上传路径添加到书库列表
			err := config.GetCfg().AddStoreUrl(config.GetCfg().UploadPath)
			if err != nil {
				logger.Infof(locale.GetString("log_failed_to_add_upload_path_to_store_urls"), err)
			}
		}
		// 如果用户启用上传，但没有指定上传路径
		if config.GetCfg().UploadPath == "" {
			for _, storeUrl := range config.GetCfg().StoreUrls {
				// 把【本地存储】里面的第一个可用路径作为上传路径
				if tools.IsExist(storeUrl) {
					config.SetUploadPath(storeUrl)
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
		logger.Infof(locale.GetString("log_failed_to_scan_store_path"), err)
	}
}
