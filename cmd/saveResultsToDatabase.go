package cmd

import (
	"fmt"
	"github.com/yumenaka/comi/book"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/storage"
)

// SaveResultsToDatabase 4，保存扫描结果到数据库，并清理不存在的书籍
func SaveResultsToDatabase() {
	if common.Config.EnableDatabase {
		AllBook := book.GetAllBookList()
		//设置清理数据库的时候，是否清理没扫描到的书籍信息
		if common.Config.ClearDatabase {
			for _, checkBook := range databaseBookList {
				needClear := true //这条数据是否需要清理
				for _, b := range AllBook {
					if b.BookID == checkBook.BookID {
						needClear = false //如果扫到了这本书,就不清理相关数据
					}
				}
				if needClear {
					storage.ClearBookData(checkBook, common.Config.Debug)
				}
			}
		}
		saveErr := storage.SaveBookListToDatabase(AllBook)
		if saveErr != nil {
			fmt.Println(saveErr)
		}
	}
}
