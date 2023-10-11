package cmd

import (
	"fmt"
	"github.com/yumenaka/comi/arch/scan"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/database"
	"github.com/yumenaka/comi/types"
	"log"
	"strconv"
)

// initBookStores 解析命令,扫描书库
func initBookStores(args []string) {
	//1. 初始化数据库
	if config.Config.EnableDatabase {
		// 从数据库中读取书籍信息并持久化
		if err := database.InitDatabase(config.Config.ConfigPath); err != nil {
			fmt.Println(err)
		}
		books, err := database.GetBooksFromDatabase()
		if err != nil {
			fmt.Println(err)
		} else {
			err := types.RestoreDatabaseBooks(books)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("从数据库中读取书籍信息,持久化成功:" + strconv.Itoa(len(books)))
			}
		}
	}
	//2、设置默认书库路径：扫描CMD指定的路径，如果开启上传，额外增加上传文件夹到默认书库路径
	initStorePath(args)

	//3、扫描配置文件里面的书库路径
	option := scan.NewScanOption(
		true,
		config.Config.StoresPath,
		config.Config.MaxScanDepth,
		config.Config.MinImageNum,
		config.Config.TimeoutLimitForScan,
		config.Config.ExcludePath,
		config.Config.SupportMediaType,
		config.Config.SupportFileType,
		config.Config.ZipFileTextEncoding,
		config.Config.EnableDatabase,
		config.Config.ClearDatabaseWhenExit,
		config.Config.Debug,
	)
	err := scan.ScanStorePath(option)
	if err != nil {
		log.Printf("Failed to scan store path: %v", err)
	}

	//4、保存扫描结果到数据库
	if config.Config.EnableDatabase {
		err = scan.SaveResultsToDatabase(config.Config.ConfigPath, config.Config.ClearDatabaseWhenExit)
		if err != nil {
			log.Printf("Failed SaveResultsToDatabase: %v", err)
			return
		}
	}

	//5、通过“可执行文件名”设置部分默认参数,目前不生效
	config.Config.SetByExecutableFilename()
}
