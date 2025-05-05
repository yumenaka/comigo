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
	for _, localPath := range option.Cfg.GetLocalStores() {
		books, err := InitStore(localPath, option)
		if err != nil {
			logger.Infof(locale.GetString("scan_error")+" path:%s %s", localPath, err)
			continue
		}
		AddBooksToStore(books, localPath, option.Cfg.GetMinImageNum())
	}
	// for _, server := range scanStores {
	//	addList, err := Smb(option)
	//	if err != nil {
	//		logger.Infof("smb scan_error"+" path:%s %s", server.ShareName, err)
	//		continue
	//	}
	//	AddBooksToStore(addList, server.ShareName, scanMinImageNum)
	// }
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
