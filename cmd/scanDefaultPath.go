package cmd

import (
	"fmt"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"os"
	"path"
)

// ScanDefaultPath 2、搜索基本路径，来自程序启动时的参数
func ScanDefaultPath(args []string) {
	//决定如何扫描，扫描哪个路径
	//没有指定路径或文件的情况下
	if len(args) == 0 {
		cmdPath := path.Dir(os.Args[0]) //扫描程序执行的路径
		addList, err := common.ScanAndGetBookList(cmdPath, common.DatabaseBookList)
		if err != nil {
			fmt.Println(locale.GetString("scan_error"), cmdPath, err)
		} else {
			common.AddBooksToStore(addList, cmdPath)
		}
	} else {
		//指定了多个参数的话，都扫描一遍
		for _, p := range args {
			addList, err := common.ScanAndGetBookList(p, common.DatabaseBookList)
			if err != nil {
				fmt.Println(locale.GetString("scan_error"), p, err)
			} else {
				common.AddBooksToStore(addList, p)
			}
		}
	}
}
