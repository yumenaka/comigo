package scan

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// 扫描目录，并返回对应书籍
func scanDirGetBook(dirPath string, storePath string, depth int, option Option) (*model.Book, error) {
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
		logger.Infof("Failed to read directory: %s, error: %v", dirPath, err)
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !option.IsSupportMedia(fileName) {
			continue
		}

		fileInfo, err := entry.Info()
		if err != nil {
			logger.Infof("Failed to get file info: %s, error: %v", fileName, err)
			continue
		}

		absPath := filepath.Join(dirPath, fileName)
		tempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(fileName)
		newBook.Images = append(newBook.Images, model.MediaFileInfo{
			Path:    absPath,
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime(),
			Name:    fileName,
			Url:     tempURL,
		})
	}
	newBook.SortPages("default")
	return newBook, nil
}
