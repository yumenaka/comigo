package data_api

import (
	"os"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/logger"
)

// ClearBookNotExist  检查内存中的书的源文件是否存在，不存在就删掉
func ClearBookNotExist() {
	logger.Infof("Checking book files exist...")
	var deletedBooks []string
	// 遍历所有书籍
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, book := range allBooks {
		// 如果父文件夹存在，但书籍文件不存在，也说明这本书被删除了
		if _, err := os.Stat(book.BookPath); os.IsNotExist(err) {
			deletedBooks = append(deletedBooks, book.BookPath)
			err := model.IStore.DeleteBook(book.BookID)
			if err != nil {
				logger.Infof("Error deleting book %s: %s", book.BookID, err)
			}
		}
	}
	// 重新生成书组
	if len(deletedBooks) > 0 {
		if config.GetCfg().EnableDatabase {
			// 先把剩下的书加回RamStore
			err := store.RamStore.AddBooks(allBooks)
			if err != nil {
				return
			}
		}
		if err := store.RamStore.GenerateAllBookGroup(); err != nil {
			logger.Infof("Error initializing main folder: %s", err)
		}
	}
}
