package scan

import (
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// InitAllStore 扫描书库路径，取得书籍
func InitAllStore(option Option) error {
	// 重置所有书籍与书组信息
	model.ClearAllBookData()
	stores := option.Cfg.GetLocalStores()
	// logger.Info("--------------------new stores------------------------------")
	// logger.Info(stores)
	// logger.Info("--------------------new stores------------------------------")
	for _, localPath := range stores {
		books, err := InitStore(localPath, option)
		if err != nil {
			logger.Infof(locale.GetString("scan_error")+" path:%s %s", localPath, err)
			continue
		}
		AddBooksToStore(books, localPath, option.Cfg.GetMinImageNum())
	}
	return nil
}

// AddBooksToStore 添加一组书到书库
func AddBooksToStore(bookList []*model.Book, basePath string, MinImageNum int) {
	err := model.AddBooks(bookList, basePath, MinImageNum)
	if err != nil {
		logger.Infof(locale.GetString("AddBook_error")+"%s", basePath)
	}
	// 生成虚拟书籍组
	if err := model.MainStore.AnalyzeStore(); err != nil {
		logger.Infof("%s", err)
	}
}
