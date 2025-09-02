package scan

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// InitStore 扫描本地路径，取得路径里的书籍
func InitStore(storePath string, option Option) ([]*model.Book, error) {
	if !tools.PathExists(storePath) {
		return nil, errors.New(locale.GetString("path_not_exist"))
	}
	storePathAbs, err := filepath.Abs(storePath)
	if err != nil {
		logger.Infof("Failed to get absolute path: %s", err)
		storePathAbs = storePath
	}
	logger.Infof(locale.GetString("scan_start_hint")+" %s", storePathAbs)

	// 如果书库URL对应一个文件，返回一本书
	if tools.IsFile(storePathAbs) && option.IsSupportFile(storePathAbs) {
		book, err := scanFileGetBook(storePathAbs, storePathAbs, 0, option)
		if err != nil {
			return nil, err
		}
		// logger.Info("-------------found_in_path:", storePathAbs)
		return []*model.Book{book}, nil
	}

	// 如果书库URL是一个文件夹，使用 HandleDirectory（不支持扫描单个文件）进行扫描
	rootDirectoryNode, foundDirs, foundFiles, err := HandleDirectory(storePathAbs, 0, option)
	if err != nil {
		return nil, err
	}

	var newBookList []*model.Book

	// 处理根目录
	if rootDirectoryNode.Files != nil {
		book, err := scanDirGetBook(storePathAbs, storePathAbs, 0, option)
		if err != nil {
			logger.Infof("Failed to scan root directory: %s, error: %v", storePathAbs, err)
		} else {
			newBookList = append(newBookList, book)
		}
	}

	// 处理子目录
	for _, dir := range foundDirs {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			logger.Infof("Failed to get absolute path: %s", err)
			absDir = storePath
		}

		// 计算路径深度
		relPath, err := filepath.Rel(storePathAbs, absDir)
		if err != nil {
			logger.Infof("Failed to get relative path: %s", err)
			continue
		}
		depth := strings.Count(relPath, string(os.PathSeparator))

		// 检查是否超过最大深度限制
		if option.Cfg.GetMaxScanDepth() >= 0 && depth > option.Cfg.GetMaxScanDepth() {
			logger.Infof(locale.GetString("exceeds_maximum_depth")+" %d, base: %s, scan: %s", option.Cfg.GetMaxScanDepth(), storePathAbs, absDir)
			continue
		}

		// 检查是否在忽略列表中
		if option.IsSkipDir(absDir) {
			logger.Infof(locale.GetString("skip_path")+" %s", absDir)
			continue
		}

		// 扫描目录
		book, err := scanDirGetBook(absDir, storePathAbs, depth, option)
		if err != nil {
			logger.Infof("Failed to scan directory: %s, error: %v", absDir, err)
			continue
		}
		newBookList = append(newBookList, book)
	}

	// 处理文件
	for _, file := range foundFiles {
		if !option.IsSupportFile(file.Name) {
			continue
		}

		// 计算路径深度
		relPath, err := filepath.Rel(storePathAbs, file.Path)
		if err != nil {
			logger.Infof("Failed to get relative path: %s", err)
			continue
		}
		depth := strings.Count(relPath, string(os.PathSeparator))

		// 检查是否超过最大深度限制
		if option.Cfg.GetMaxScanDepth() >= 0 && depth > option.Cfg.GetMaxScanDepth() {
			logger.Infof(locale.GetString("exceeds_maximum_depth")+" %d, base: %s, scan: %s", option.Cfg.GetMaxScanDepth(), storePathAbs, file.Path)
			continue
		}

		// 扫描文件
		book, err := scanFileGetBook(file.Path, storePathAbs, depth, option)
		if err != nil {
			logger.Info(err)
			continue
		}
		newBookList = append(newBookList, book)
	}

	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("found_in_path"), storePathAbs, len(newBookList))
	}
	return newBookList, nil
}
