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
		// 压缩包类型的书籍，页数小于最小图片数，跳过添加
		if book.Type == model.TypeZip || book.Type == model.TypeRar ||
			book.Type == model.TypeCbz || book.Type == model.TypeCbr ||
			book.Type == model.TypeTar || book.Type == model.TypeEpub {
			if len(book.PageInfos) < config.GetCfg().MinImageNum {
				continue
			}
		}
		// 目录类型的书籍，页数小于最小图片数(或1)，跳过添加
		if book.Type == model.TypeDir {
			if len(book.PageInfos) < config.GetCfg().MinImageNum || len(book.PageInfos) == 0 {
				continue
			}
		}
		book.PageCount = len(book.PageInfos)
		if err := model.IStore.StoreBook(book); err != nil {
			logger.Infof(locale.GetString("log_add_book_error"), book.BookID, err)
		}
	}
}
