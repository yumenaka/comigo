package common

import (
	"encoding/base64"
	"mime"
	"path/filepath"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/templ/state"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// GetPageTitle 获取页面标题
func GetPageTitle(bookID string) string {
	if bookID == "" {
		return state.ServerConfig.GetTopStoreName()
	}
	groupBook, err := model.IStore.GetBook(bookID)
	if err != nil {
		logger.Info("GetBook: %v", err)
		return "Comigo " + state.Version
	}
	return groupBook.Title
}

// GetReturnUrl 阅读或书架页面，返回按钮实际使用的链接
func GetReturnUrl(BookID string) string {
	childID := BookID
	if childID == "" {
		return "/"
	}
	for _, bookGroup := range model.IStore.ListBooks() {
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

func ShowQuickJumpBar(b *model.Book) (QuickJumpBar bool) {
	_, err := store.GetBookInfoListByParentFolder(b.ParentFolder)
	if err != nil {
		logger.Infof("%s", err)
		return false
	}
	return true
}

func QuickJumpBarBooks(b *model.Book) (list *model.BookInfoList) {
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
		logger.Infof("GetBook error: %s", err)
		return ""
	}

	// 获取图片数据的选项
	option := fileutil.GetPictureDataOption{
		PictureName:      fileName,
		BookIsPDF:        bookByID.Type == model.TypePDF,
		BookIsDir:        bookByID.Type == model.TypeDir,
		BookIsNonUTF8Zip: bookByID.NonUTF8Zip,
		BookFilePath:     bookByID.FilePath,
		Debug:            config.GetCfg().Debug,
		UseCache:         config.GetCfg().UseCache,
	}

	// 获取图片数据
	imgData, _, err := fileutil.GetPictureData(option)
	if err != nil {
		logger.Infof("GetPictureData error: %s", err)
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
