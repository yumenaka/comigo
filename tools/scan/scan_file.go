package scan

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// 扫描本地文件，并返回对应书籍
func scanFileGetBook(filePath string, storePath string, depth int) (*model.Book, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Infof("Failed to open file: %s, error: %v", filePath, err)
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logger.Infof("Failed to get file info: %s, error: %v", filePath, err)
		return nil, err
	}
	//查看书库中是否已经有了这本书，有了就跳过
	for _, book := range model.IStore.ListBooks() {
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			logger.Infof("Error getting absolute path: %v", err)
			continue
		}
		if book.FilePath == absFilePath && book.Type == model.GetBookTypeByFilename(filePath) {
			return nil, errors.New("skip: " + filePath)
		}
	}
	// 创建新书籍
	newBook := model.NewBook(filePath, fileInfo.ModTime(), fileInfo.Size(), storePath, depth, model.GetBookTypeByFilename(filePath))

	switch newBook.Type {
	case model.TypeZip, model.TypeCbz, model.TypeEpub:
		err = handleZipAndEpubFiles(filePath, newBook)
		if err != nil {
			return nil, err
		}
	case model.TypePDF:
		err = handlePdfFiles(filePath, newBook)
		if err != nil {
			return nil, err
		}
	case model.TypeVideo, model.TypeAudio:
		handleMediaFiles(newBook)
	case model.TypeUnknownFile:
		handleOtherFiles(newBook)
	default:
		err = handleOtherArchiveFiles(filePath, newBook)
		if err != nil {
			return nil, err
		}
	}
	newBook.Cover = newBook.GetCover()
	newBook.SortPages("default")
	return newBook, nil
}
