package common

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/url"
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
		return config.PrefixPath("/")
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
					return config.PrefixPath("/")
				}
				return config.PrefixPath("/shelf/" + b.BookID)
			}
		}
	}
	return config.PrefixPath("/")
}

func QuickJumpBarBooks(b *model.Book) (list *model.BookInfos) {
	list, err := store.GetBookInfoListByBookFolder(b)
	if err != nil {
		logger.Infof("%s", err)
		return nil
	}
	return list
}

// RawBookURL 生成单文件书籍的原始文件访问地址。
// HTML/音视频这类单文件内容不需要阅读模板时，应跳转到该地址交给 raw API 返回源文件。
func RawBookURL(book *model.Book) string {
	if book == nil {
		return config.PrefixPath("/")
	}
	return RawBookInfoURL(book.BookInfo)
}

// RawBookInfoURL 根据 BookInfo 生成原始文件访问地址。
// 前端播放器、HTML 直出等场景只需要稳定的公开 URL，不需要知道本地 BookPath。
func RawBookInfoURL(book model.BookInfo) string {
	if book.BookID == "" {
		return config.PrefixPath("/")
	}
	fileName := book.Title
	if fileName == "" {
		fileName = filepath.Base(book.BookPath)
	}
	rawURL := config.PrefixPath("/api/raw/" + url.PathEscape(book.BookID) + "/" + url.PathEscape(fileName))
	if book.RemoteStoreKey != "" {
		rawURL += "?remote_store=" + url.QueryEscape(book.RemoteStoreKey)
	}
	return rawURL
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
