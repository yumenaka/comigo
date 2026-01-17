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
	model.GenerateBookGroup()
	return nil
}

// AddBooksToStore 添加一组书到书库
func AddBooksToStore(books []*model.Book) {
	for _, book := range books {
		if book.Type == model.TypeDir ||
			book.Type == model.TypeZip || book.Type == model.TypeRar ||
			book.Type == model.TypeCbz || book.Type == model.TypeCbr ||
			book.Type == model.TypeTar || book.Type == model.TypeEpub {
			if len(book.PageInfos) < config.GetCfg().MinImageNum {
				continue
			}
		}
		book.PageCount = len(book.PageInfos)
		if err := model.IStore.StoreBook(book); err != nil {
			logger.Infof(locale.GetString("log_add_book_error"), book.BookID, err)
		}
	}
}
