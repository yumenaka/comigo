package cmd

import (
	"os"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/upload_api"
	"github.com/yumenaka/comigo/sqlc"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
)

// ScanStore 解析命令,扫描文件，设置书库等
func ScanStore(args []string) {
	// 从数据库中读取书籍信息并持久化
	// 启动或重启 Tailscale 服务
	configDir, err := config.GetConfigDir()
	if err != nil {
		logger.Errorf("Failed to get config dir: %v", err)
		configDir = ""
	}
	// 1. 初始化数据库
	// 切换到DbStore会导致的已知问题：
	// 书组相关跳转异常
	if config.GetCfg().EnableDatabase {
		if err := sqlc.OpenDatabase(configDir); err != nil {
			logger.Infof("OpenDatabase Error: %s", err)
			model.IStore = store.RamStore
		} else {
			model.IStore = sqlc.DbStore
		}
	}
	//model.IStore = store.RamStore
	// 2、设置默认书库路径：扫描CMD指定的路径，或添加当前文件夹为默认路径。
	CreateStoreUrls(args)
	// 3、扫描配置文件里面的书库路径
	err = scan.InitAllStore(config.GetCfg())
	if err != nil {
		logger.Infof("Failed to scan store path: %v", err)
	}

	// 4、生成虚拟书籍组
	if config.GetCfg().EnableDatabase {
		allBooks, err := sqlc.DbStore.ListBooks()
		if err != nil {
			logger.Infof("Error listing books from database: %s", err)
		} else {
			// 拿到的书加回RamStore
			err = store.RamStore.AddBooks(allBooks)
			if err != nil {
				return
			}
		}
	}
	if err := model.IStore.GenerateBookGroup(); err != nil {
		logger.Infof("%s", err)
	}
	// 5、保存扫描结果到数据库
	if config.GetCfg().EnableDatabase {
		err = scan.SaveBooksToDatabase(config.GetCfg())
		if err != nil {
			logger.Infof("Failed SaveBooksToDatabase: %v", err)
			return
		}
	}
}

// CreateStoreUrls 添加默认扫描路径 args[1:]是用户指定的扫描路径
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
		config.GetCfg().AddStoreUrl(wd)
	}
	// 指定了书库路径，就都扫描一遍
	for key, arg := range args {
		if config.GetCfg().Debug {
			logger.Infof("args[%d]: %s\n", key, arg)
		}
		config.GetCfg().AddStoreUrl(arg)
	}
	// 如果用户启用上传，且用户指定的上传路径不为空，就把程序预先设定的【默认上传路径】当作书库
	if config.GetCfg().EnableUpload {
		if config.GetCfg().UploadPath != "" {
			// 尝试把上传路径添加为书库里
			config.GetCfg().AddStoreUrl(config.GetCfg().UploadPath)
		}
		// 如果用户启用上传，但没有指定上传路径
		if config.GetCfg().UploadPath == "" {
			for _, storeUrl := range config.GetCfg().StoreUrls {
				// 把【本地存储】里面的第一个可用路径作为上传路径
				if tools.IsExist(storeUrl) {
					config.SetUploadPath(storeUrl)
					config.GetCfg().AddStoreUrl(config.GetCfg().UploadPath)
					break
				}
			}
		}
	}
	// 扫描路径设置，传递给 router
	upload_api.ConfigEnableUpload = &config.GetCfg().EnableUpload
	upload_api.ConfigLocked = &config.GetCfg().ConfigLocked
	upload_api.ConfigUploadPath = &config.GetCfg().UploadPath
}
