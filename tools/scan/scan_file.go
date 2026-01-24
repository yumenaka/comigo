package scan

import (
	"os"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// 扫描本地文件，并返回对应书籍
func scanFileGetBook(filePath string, storePath string, depth int) (*model.Book, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_open_file"), filePath, err)
		return nil, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_get_file_info"), filePath, err)
		return nil, err
	}
	// 创建新书籍
	newBook, err := model.NewBook(filePath, fileInfo.ModTime(), fileInfo.Size(), storePath, depth, model.GetBookTypeByFilename(filePath))
	if err != nil {
		return nil, err
	}
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
	case model.TypeVideo, model.TypeAudio, model.TypeHTML:
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
