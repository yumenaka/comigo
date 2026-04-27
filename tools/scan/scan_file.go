package scan

import (
	"os"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/vfs"
)

// scanFileGetBook 扫描本地文件，并返回对应书籍
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

// scanRemoteFileGetBook 扫描远程文件，并返回对应书籍
// 压缩包类型（ZIP/RAR/EPUB等）直接从 WebDAV 流式读取，不下载到本地
// PDF 文件需要下载到本地缓存（因为 PDF 处理库需要随机访问）
func scanRemoteFileGetBook(fs vfs.FileSystem, filePath string, storeURL string, depth int) (*model.Book, error) {
	// 获取文件信息
	fileInfo, err := fs.Stat(filePath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_get_file_info"), filePath, err)
		return nil, err
	}

	// 创建新书籍
	newBook, err := model.NewBook(filePath, fileInfo.ModTime(), fileInfo.Size(), storeURL, depth, model.GetBookTypeByFilename(filePath))
	if err != nil {
		return nil, err
	}

	// 标记为远程书籍
	newBook.IsRemote = true
	newBook.RemoteURL = storeURL

	// 根据书籍类型处理
	switch newBook.Type {
	case model.TypeZip, model.TypeCbz, model.TypeEpub:
		// 远程压缩包：直接从 WebDAV 流式读取，不下载到本地
		err = handleRemoteZipFile(fs, filePath, newBook)
		if err != nil {
			return nil, err
		}
	case model.TypePDF:
		// 远程 PDF：需要下载到本地缓存（PDF 处理库需要随机访问）
		err = handleRemotePdfFile(fs, filePath, newBook)
		if err != nil {
			return nil, err
		}
	case model.TypeVideo, model.TypeAudio, model.TypeHTML:
		handleMediaFiles(newBook)
	case model.TypeUnknownFile:
		handleOtherFiles(newBook)
	default:
		// 其他压缩包格式：直接从 WebDAV 流式读取，不下载到本地
		err = handleRemoteOtherArchiveFile(fs, filePath, newBook)
		if err != nil {
			return nil, err
		}
	}

	newBook.Cover = newBook.GetCover()
	newBook.SortPages("default")
	return newBook, nil
}
