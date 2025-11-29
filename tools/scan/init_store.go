package scan

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// InitStore 扫描本地路径，取得路径里的书籍
func InitStore(storePath string, cfg ConfigInterface) error {
	InitConfig(cfg)
	if !tools.PathExists(storePath) {
		return errors.New(locale.GetString("path_not_exist"))
	}
	storePathAbs, err := filepath.Abs(storePath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_get_absolute_path_scan"), err)
		storePathAbs = storePath
	}
	logger.Infof(locale.GetString("scan_start_hint")+" %s", storePathAbs)

	// 如果书库URL对应一个文件，返回一本书
	if tools.IsFile(storePathAbs) && IsSupportFile(storePathAbs) {
		book, err := scanFileGetBook(storePathAbs, storePathAbs, 0)
		if err != nil {
			return err
		}
		AddBooksToStore([]*model.Book{book})
		return nil
	}

	// 如果书库URL是一个文件夹，使用 HandleDirectory（不支持扫描单个文件）进行扫描
	rootDirectoryNode, foundDirs, foundFiles, err := HandleDirectory(storePathAbs, 0)
	if err != nil {
		return err
	}

	var newBookList []*model.Book

	// 处理根目录
	if rootDirectoryNode.Files != nil {
		book, err := scanDirGetBook(storePathAbs, storePathAbs, 0)
		if err != nil {
			logger.Infof(locale.GetString("log_skip_to_scan_root_directory"), storePathAbs, err)
		} else {
			newBookList = append(newBookList, book)
		}
	}

	// 处理子目录
	for _, dir := range foundDirs {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_get_absolute_path_scan"), err)
			absDir = storePath
		}

		// 计算路径深度
		relPath, err := filepath.Rel(storePathAbs, absDir)
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_get_relative_path"), err)
			continue
		}
		depth := strings.Count(relPath, string(os.PathSeparator))

		// 检查是否超过最大深度限制
		if cfg.GetMaxScanDepth() >= 0 && depth > cfg.GetMaxScanDepth() {
			logger.Infof(locale.GetString("exceeds_maximum_depth")+" %d, base: %s, scan: %s", cfg.GetMaxScanDepth(), storePathAbs, absDir)
			continue
		}

		// 检查是否在忽略列表中
		if IsSkipDir(absDir) {
			logger.Infof(locale.GetString("skip_path")+" %s", absDir)
			continue
		}

		// 扫描目录
		book, err := scanDirGetBook(absDir, storePathAbs, depth)
		if err != nil {
			logger.Infof(locale.GetString("log_skip_to_scan_directory"), absDir, err)
			continue
		}
		newBookList = append(newBookList, book)
	}

	// 处理文件
	for _, file := range foundFiles {
		////如果是以 . 开头的隐藏文件，跳过
		if strings.HasPrefix(path.Base(file.Name), ".") {
			continue
		}
		if !IsSupportFile(file.Name) {
			continue
		}
		// 计算路径深度
		relPath, err := filepath.Rel(storePathAbs, file.Path)
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_get_relative_path"), err)
			continue
		}
		depth := strings.Count(relPath, string(os.PathSeparator))

		// 检查是否超过最大深度限制
		if cfg.GetMaxScanDepth() >= 0 && depth > cfg.GetMaxScanDepth() {
			logger.Infof(locale.GetString("exceeds_maximum_depth")+" %d, base: %s, scan: %s", cfg.GetMaxScanDepth(), storePathAbs, file.Path)
			continue
		}

		// 扫描文件
		book, err := scanFileGetBook(file.Path, storePathAbs, depth)
		if err != nil {
			logger.Info(err)
			continue
		}
		newBookList = append(newBookList, book)
	}

	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("found_in_path"), storePathAbs, len(newBookList))
	}
	AddBooksToStore(newBookList)
	return nil
}
