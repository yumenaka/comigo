package common

import (
	"encoding/base64"
	"fmt"
	"mime"
	"path/filepath"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// GetPageTitle 获取页面标题
func GetPageTitle(bookID string, nowBookNum int, storeBookInfos []model.StoreBookInfo, childBookInfos []model.BookInfo) string {
	if bookID == "" {
		return fmt.Sprintf("%v(x%v) ", config.GetCfg().GetTopStoreName(), nowBookNum)
	}
	groupBook, err := model.IStore.GetBook(bookID)
	if err != nil {
		logger.Infof(locale.GetString("log_getbook_error_scroll"), err)
		return "Comigo " + config.GetVersion()
	}
	return groupBook.Title
}

// GetBookTitle 获取页面标题
func GetBookTitle(bookID string) string {
	if bookID == "" {
		return "Comigo " + config.GetVersion()
	}
	groupBook, err := model.IStore.GetBook(bookID)
	if err != nil {
		logger.Infof(locale.GetString("log_getbook_error_scroll"), err)
		return "Comigo " + config.GetVersion()
	}
	return groupBook.Title
}

// GetReturnUrl 阅读或书架页面，返回按钮实际使用的链接
func GetReturnUrl(BookID string) string {
	childID := BookID
	if childID == "" {
		return "/"
	}
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}
	for _, bookGroup := range allBooks {
		if bookGroup.Type != model.TypeBooksGroup {
			continue // 只分析书组类型
		}
		for _, id := range bookGroup.ChildBooksID {
			if id == childID {
				b, err := model.IStore.GetBook(bookGroup.BookID)
				if err != nil {
					return "/"
				}
				return "/shelf/" + b.BookID
			}
		}
	}
	return "/"
}

func QuickJumpBarBooks(b *model.Book) (list *model.BookInfos) {
	list, err := store.GetBookInfoListByParentFolder(b.ParentFolder)
	if err != nil {
		logger.Infof("%s", err)
		return nil
	}
	return list
}

func GetFileBase64Text(bookID string, fileName string) string {
	// 获取书籍信息
	bookByID, err := model.IStore.GetBook(bookID)
	if err != nil {
		logger.Infof(locale.GetString("log_getbook_error_common"), err)
		return ""
	}

	// 获取图片数据的选项
	option := fileutil.GetPictureDataOption{
		PictureName:      fileName,
		BookIsPDF:        bookByID.Type == model.TypePDF,
		BookIsDir:        bookByID.Type == model.TypeDir,
		BookIsNonUTF8Zip: bookByID.NonUTF8Zip,
		BookPath:         bookByID.BookPath,
		Debug:            config.GetCfg().Debug,
		UseCache:         config.GetCfg().UseCache,
	}

	// 获取图片数据
	imgData, _, err := fileutil.GetPictureData(option)
	if err != nil {
		logger.Infof(locale.GetString("log_getpicturedata_error"), err)
		return ""
	}

	// 转换为Base64字符串
	mimeType := mime.TypeByExtension(filepath.Ext(fileName))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	dataURI := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(imgData)

	return dataURI
}
