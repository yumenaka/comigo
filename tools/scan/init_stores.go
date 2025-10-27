package scan

import (
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// InitAllStore 扫描书库路径，取得书籍
func InitAllStore(cfg ConfigInterface) error {
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
func AddBooksToStore(books []*model.Book) {
	for _, book := range books {
		if len(book.Images) < config.GetCfg().MinImageNum {
			continue
		}
		book.PageCount = len(book.Images)
		if err := model.IStore.AddBook(book); err != nil {
			logger.Infof("AddBook_error"+" bookID:%s %s", book.BookID, err)
		}
	}
}
