package scan

import (
	"os"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// 扫描本地文件，并返回对应书籍
func scanFileGetBook(filePath string, storePath string, depth int, scanOption Option) (*model.Book, error) {
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

	newBook, err := model.NewBook(filePath, fileInfo.ModTime(), fileInfo.Size(), storePath, depth, model.GetBookTypeByFilename(filePath))
	if err != nil {
		return nil, err
	}

	switch newBook.Type {
	case model.TypeZip, model.TypeCbz, model.TypeEpub:
		err = handleZipAndEpubFiles(filePath, newBook, scanOption)
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
		handleMediaFiles(newBook)
	default:
		err = handleOtherArchiveFiles(filePath, newBook, scanOption)
		if err != nil {
			return nil, err
		}
	}

	newBook.SortPages("default")
	return newBook, nil
}
