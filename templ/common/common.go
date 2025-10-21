package common

import (
	"encoding/base64"
	"mime"
	"path/filepath"
	"strconv"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/state"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// GetPageTitle 获取页面标题
func GetPageTitle(bookID string) string {
	if bookID == "" {
		return state.ServerConfig.GetTopStoreName()
	}
	groupBook, err := model.IStore.GetBookByID(bookID, "")
	if err != nil {
		logger.Info("GetBookByID: %v", err)
		return "Comigo " + state.Version
	}
	return groupBook.Title
}

// GetImageAlt 获取图片的 alt 属性
func GetImageAlt(key int) string {
	return strconv.Itoa(key)
}

// GetReturnUrl 阅读或书架页面，返回按钮实际使用的链接
func GetReturnUrl(BookID string) string {
	// 如果是书籍组，就跳转到父书架
	ParentBook, err := model.IStore.GetParentBook(BookID)
	if err != nil {
		// logger.Infof("ParentBook not found by BookID: %s, error: %v", BookID, err)
		return "/"
	}

	return "/shelf/" + ParentBook.BookID
}

func ShowQuickJumpBar(b *model.Book) (QuickJumpBar bool) {
	_, err := model.IStore.GetBookInfoListByParentFolder(b.ParentFolder, "")
	if err != nil {
		logger.Infof("%s", err)
		return false
	}
	return true
}

func QuickJumpBarBooks(b *model.Book) (list *model.BookInfoList) {
	list, err := model.IStore.GetBookInfoListByParentFolder(b.ParentFolder, "")
	if err != nil {
		logger.Infof("%s", err)
		return nil
	}
	return list
}

func GetFileBase64Text(bookID string, fileName string) string {
	// 获取书籍信息
	bookByID, err := model.IStore.GetBookByID(bookID, "")
	if err != nil {
		logger.Infof("GetBookByID error: %s", err)
		return ""
	}

	// 获取图片数据的选项
	option := fileutil.GetPictureDataOption{
		PictureName:      fileName,
		BookIsPDF:        bookByID.Type == model.TypePDF,
		BookIsDir:        bookByID.Type == model.TypeDir,
		BookIsNonUTF8Zip: bookByID.NonUTF8Zip,
		BookFilePath:     bookByID.FilePath,
		Debug:            config.GetDebug(),
		UseCache:         config.GetUseCache(),
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
