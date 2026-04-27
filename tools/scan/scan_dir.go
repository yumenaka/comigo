package scan

import (
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/vfs"
)

// scanDirGetBook 扫描本地目录，并返回对应书籍
func scanDirGetBook(dirPath string, storePath string, depth int) (*model.Book, error) {
	// 获取文件夹信息
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}
	newBook, err := model.NewBook(dirPath, dirInfo.ModTime(), dirInfo.Size(), storePath, depth, model.TypeDir)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_read_directory"), dirPath, err)
		return nil, err
	}
	pageNum := 1
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !IsSupportMedia(fileName) {
			continue
		}

		fileInfo, err := entry.Info()
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_get_file_info_scan"), fileName, err)
			continue
		}

		absPath := filepath.Join(dirPath, fileName)
		tempURL := "/api/get-file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(fileName)
		newBook.PageInfos = append(newBook.PageInfos, model.PageInfo{
			Path:    absPath,
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime(),
			Name:    fileName,
			Url:     tempURL,
			PageNum: pageNum,
		})
		pageNum++
	}
	newBook.SortPages("default")
	return newBook, nil
}

// scanRemoteDirGetBook 扫描远程目录，并返回对应书籍
func scanRemoteDirGetBook(fs vfs.FileSystem, dirPath string, storeURL string, depth int) (*model.Book, error) {
	// 获取目录信息
	dirInfo, err := fs.Stat(dirPath)
	if err != nil {
		return nil, err
	}

	// 创建新书籍，使用远程路径
	newBook, err := model.NewBook(dirPath, dirInfo.ModTime(), dirInfo.Size(), storeURL, depth, model.TypeDir)
	if err != nil {
		return nil, err
	}

	// 标记为远程书籍
	newBook.IsRemote = true
	newBook.RemoteURL = storeURL

	entries, err := fs.ReadDir(dirPath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_read_directory"), dirPath, err)
		return nil, err
	}

	pageNum := 1
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !IsSupportMedia(fileName) {
			continue
		}

		fileInfo, err := entry.Info()
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_get_file_info_scan"), fileName, err)
			continue
		}

		// 远程路径使用 path.Join（URL 风格）
		fullPath := path.Join(dirPath, fileName)
		tempURL := "/api/get-file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(fileName)

		newBook.PageInfos = append(newBook.PageInfos, model.PageInfo{
			Path:    fullPath,
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime(),
			Name:    fileName,
			Url:     tempURL,
			PageNum: pageNum,
		})
		pageNum++
	}

	newBook.SortPages("default")
	return newBook, nil
}
