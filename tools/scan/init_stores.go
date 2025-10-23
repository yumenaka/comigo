package scan

import (
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// InitAllStore 扫描书库路径，取得书籍
func InitAllStore(cfg ConfigInterface) error {
	// 重置所有书籍与书组信息
	model.IStore.ClearAllBook()
	storeUrls := cfg.GetStoreUrls()
	for _, storeUrl := range storeUrls {
		err := InitStore(storeUrl, cfg)
		if err != nil {
			logger.Infof(locale.GetString("scan_error")+" path:%s %s", storeUrl, err)
			continue
		}
	}
	return nil
}

// AddBooksToStore 添加一组书到书库
func AddBooksToStore(bookList []*model.Book) {
	err := model.IStore.AddBooks(bookList, minImageNum)
	if err != nil {
		logger.Infof(locale.GetString("AddBook_error"))
	}
	// 生成虚拟书籍组
	if err := model.IStore.GenerateAllBookGroup(); err != nil {
		logger.Infof("%s", err)
	}
}
