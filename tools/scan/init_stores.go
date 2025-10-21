package scan

import (
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// InitAllStore 扫描书库路径，取得书籍
func InitAllStore(option Option) error {
	// 重置所有书籍与书组信息
	model.IStore.ClearAll()
	storeUrls := option.Cfg.GetStoreUrls()
	for _, storeUrl := range storeUrls {
		books, err := InitStore(storeUrl, option)
		if err != nil {
			logger.Infof(locale.GetString("scan_error")+" path:%s %s", storeUrl, err)
			continue
		}
		AddBooksToStore(storeUrl, books, option.Cfg.GetMinImageNum())
	}
	return nil
}

// AddBooksToStore 添加一组书到书库
func AddBooksToStore(storeUrl string, bookList []*model.Book, MinImageNum int) {
	err := model.IStore.AddBooks(storeUrl, bookList, MinImageNum)
	if err != nil {
		logger.Infof(locale.GetString("AddBook_error")+"%s", storeUrl)
	}
	// 生成虚拟书籍组
	if err := model.IStore.GenerateAllBookGroup(); err != nil {
		logger.Infof("%s", err)
	}
}
