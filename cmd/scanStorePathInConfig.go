package cmd

import (
	"fmt"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
)

// ScanStorePathInConfig 3、扫描配置文件指定的的书籍库
func ScanStorePathInConfig() {
	if len(common.Config.StoresPath) > 0 {
		for _, p := range common.Config.StoresPath {
			addList, err := common.ScanAndGetBookList(p, databaseBookList)
			if err != nil {
				fmt.Println(locale.GetString("scan_error"), p, err)
			} else {
				common.AddBooksToStore(addList, p)
			}
		}
	}
}
