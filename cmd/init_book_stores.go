package cmd

import (
	"strconv"

	"github.com/spf13/viper"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/internal/database"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/file/scan"
	"github.com/yumenaka/comigo/util/logger"
)

// StartScan 解析命令,扫描书库
func StartScan(args []string) {
	//1. 初始化数据库
	if config.GetEnableDatabase() {
		// 从数据库中读取书籍信息并持久化
		if err := database.InitDatabase(viper.ConfigFileUsed()); err != nil {
			logger.Infof("%s", err)
		}
		books, err := database.GetBooksFromDatabase()
		if err != nil {
			logger.Infof("%s", err)
		} else {
			model.RestoreDatabaseBooks(books)
			logger.Infof("从数据库中读取书籍信息,一共有 %d 本书", strconv.Itoa(len(books)))
		}
	}
	//2、设置默认书库路径：扫描CMD指定的路径，如果开启上传，额外增加上传文件夹到默认书库路径
	initStorePath(args)

	//3、扫描配置文件里面的书库路径
	option := scan.NewScanOption(
		true,
		config.GetLocalStoresList(),
		config.GetStores(),
		config.GetMaxScanDepth(),
		config.GetMinImageNum(),
		config.GetTimeoutLimitForScan(),
		config.GetExcludePath(),
		config.GetSupportMediaType(),
		config.GetSupportFileType(),
		config.GetSupportTemplateFile(),
		config.GetZipFileTextEncoding(),
		config.GetEnableDatabase(),
		config.GetClearDatabaseWhenExit(),
		config.GetDebug(),
	)
	err := scan.InitStore(option)
	if err != nil {
		logger.Infof("Failed to scan store path: %v", err)
	}
	//4、保存扫描结果到数据库
	if config.GetEnableDatabase() {
		err = scan.SaveResultsToDatabase(viper.ConfigFileUsed(), config.GetClearDatabaseWhenExit())
		if err != nil {
			logger.Infof("Failed SaveResultsToDatabase: %v", err)
			return
		}
	}
	//5、通过“可执行文件名”设置部分默认参数,目前不生效
	config.SetByExecutableFilename()
}
