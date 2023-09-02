package cmd

import (
	"fmt"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/storage"
)

// initBookStores 解析命令,扫描书库
func initBookStores(args []string) {
	//1. 初始化数据库
	if common.Config.EnableDatabase {
		// 从数据库中读取书籍信息并持久化
		if err := storage.InitDatabase(common.ConfigFilePath); err != nil {
			fmt.Println(err)
			return
		}
		books, err := storage.GetArchiveBookFromDatabase()
		if err != nil {
			fmt.Println(err)
			return
		}
		common.DatabaseBookList = books
	}
	//2、添加CMD路径，默认上传文件夹到书库
	AddPathToStore(args)

	//3、扫描书库
	err := common.ScanStorePathInConfig()
	if err != nil {
		fmt.Println(err)
	}

	//4、保存扫描结果到数据库
	if common.Config.EnableDatabase {
		err = common.SaveResultsToDatabase()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//5、通过“可执行文件名”设置部分默认参数,目前不生效
	common.Config.SetByExecutableFilename()
}
