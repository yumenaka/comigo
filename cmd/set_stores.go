package cmd

import (
	"os"

	"github.com/spf13/viper"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/upload_api"
	"github.com/yumenaka/comigo/sqlc"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
	"github.com/yumenaka/comigo/util/scan"
)

// ScanStore 解析命令,扫描文件，设置书库等
func ScanStore(args []string) {
	// 1. 初始化数据库
	if config.GetEnableDatabase() {
		// 从数据库中读取书籍信息并持久化
		if err := sqlc.OpenDatabase(viper.ConfigFileUsed()); err != nil {
			logger.Infof("%s", err)
		}
		books, err := sqlc.Repo.GetBooksFromDatabase()
		if err != nil {
			logger.Infof("%s", err)
		} else {
			for _, book := range books {
				err = model.MainStoreGroup.AddBook(book.BookStorePath, book, config.GetMinImageNum())
				if err != nil {
					logger.Infof("AddBook error: %s", err)
				} else {
					logger.Infof("Book %s added from database", book.BookID)
				}
			}
		}
	}
	// 2、设置默认书库路径：扫描CMD指定的路径，或添加当前文件夹为默认路径。
	CreateLocalStores(args)
	// 3、扫描配置文件里面的书库路径
	err := scan.InitAllStore(scan.NewOption(config.GetCfg()))
	if err != nil {
		logger.Infof("Failed to scan store path: %v", err)
	}
	// 4、保存扫描结果到数据库
	if config.GetEnableDatabase() {
		err = scan.SaveResultsToDatabase(viper.ConfigFileUsed(), config.GetClearDatabaseWhenExit())
		if err != nil {
			logger.Infof("Failed SaveResultsToDatabase: %v", err)
			return
		}
	}
}

// CreateLocalStores 添加默认扫描路径 args[1:]是用户指定的扫描路径
func CreateLocalStores(args []string) {
	// 如果用户指定了扫描路径，就把指定的路径都加入到扫描路径里面
	config.GetCfg().InitStoreUrls()
	// 没指定扫描路径,配置文件也没设置书库文件夹的时候，默认把【当前工作目录】作为扫描路径
	if len(args) == 0 && len(config.GetCfg().GetStoreUrls()) == 0 {
		// 获取当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			logger.Infof("Failed to get working directory:%s", err)
		}
		logger.Infof("Working directory:%s", wd)
		config.GetCfg().AddStoreUrl(wd)
	}
	// 指定了书库路径，就都扫描一遍
	for key, arg := range args {
		if config.GetDebug() {
			logger.Infof("args[%d]: %s\n", key, arg)
		}
		config.GetCfg().AddStoreUrl(arg)
	}
	// 如果用户启用上传，且用户指定的上传路径不为空，就把程序预先设定的【默认上传路径】当作书库
	if config.GetEnableUpload() {
		if config.GetUploadPath() != "" {
			// 尝试把上传路径添加为书库里
			config.GetCfg().AddStoreUrl(config.GetUploadPath())
		}
		// 如果用户启用上传，但没有指定上传路径
		if config.GetUploadPath() == "" {
			for _, storeUrl := range config.GetStoreUrls() {
				// 把【本地存储】里面的第一个可用路径作为上传路径
				if util.IsExist(storeUrl) {
					config.SetUploadPath(storeUrl)
					config.GetCfg().AddStoreUrl(config.GetUploadPath())
					break
				}
			}
		}
	}
	// 扫描路径设置，传递给 router
	upload_api.ConfigEnableUpload = &config.GetCfg().EnableUpload
	upload_api.ConfigUploadPath = &config.GetCfg().UploadPath
}
