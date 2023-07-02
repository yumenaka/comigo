package cmd

import (
	"fmt"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/storage"
)

// initBookStores 解析命令,扫描书库
func initBookStores(args []string) {
	//初始化数据库
	if common.Config.EnableDatabase {
		//从数据库里面读取书籍信息，持久化
		storage.InitDatabase(common.ConfigFilePath)
		var dataErr error
		databaseBookList, dataErr = storage.GetArchiveBookFromDatabase()
		if dataErr != nil {
			fmt.Println(dataErr)
		}
	}
	//2、搜索基本路径，来自程序启动时的参数
	ScanDefaultPath(args)
	//3、扫描配置文件指定的书籍库
	ScanStorePathInConfig()
	//4、扫描默认上传文件夹
	ReScanUploadPath()
	//5、保存扫描结果到数据库
	SaveResultsToDatabase()
	//6、通过“可执行文件名”设置部分默认参数,目前不生效
	common.Config.SetByExecutableFilename()
}
