package scan

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/vfs"
)

// InitStore 扫描书库路径，取得路径里的书籍
// 支持本地路径和远程 URL（WebDAV 等）
func InitStore(storePath string, cfg ConfigInterface) error {
	InitConfig(cfg)

	// 判断是否为远程 URL
	isRemote := store.IsRemoteURL(storePath)

	if isRemote {
		return initRemoteStore(storePath, cfg)
	}
	return initLocalStore(storePath, cfg)
}

// initLocalStore 扫描本地路径
func initLocalStore(storePath string, cfg ConfigInterface) error {
	if !tools.PathExists(storePath) {
		return errors.New(locale.GetString("path_not_exist"))
	}
	storePathAbs, err := filepath.Abs(storePath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_get_absolute_path_scan"), err)
		storePathAbs = storePath
	}
	logger.Infof(locale.GetString("scan_start_hint")+" %s", storePathAbs)

	// 创建本地文件系统实例
	localFS, err := vfs.NewLocalFS(storePathAbs)
	if err != nil {
		return fmt.Errorf("创建本地文件系统失败: %w", err)
	}
	SetCurrentFS(localFS)
	defer func() {
		SetCurrentFS(nil)
		_ = localFS.Close()
	}()

	// 如果书库URL对应一个文件，返回一本书
	if tools.IsFile(storePathAbs) && IsSupportFile(storePathAbs) {
		// 如果书库里已经有这本书，跳过
		if checkBookInStore(storePathAbs, storePathAbs) {
			return nil
		}
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
		// 如果是以 . 开头的隐藏文件，跳过
		if strings.HasPrefix(filepath.Base(file.Name), ".") {
			continue
		}
		if !IsSupportFile(file.Name) {
			logger.Infof(locale.GetString("log_skip_unsupported_file_type")+" (路径: %s)", file.Name, file.Path)
			continue
		}
		// logger.Infof(locale.GetString("log_processing_file"), file.Name, file.Path)
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
		// 如果书库里已经有这本书(文件类型)，跳过。文件夹类型不能跳过。
		if tools.IsFile(file.Path) && checkBookInStore(storePathAbs, file.Path) {
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
		logger.Infof(locale.GetString("how_many_books_update"), storePathAbs, len(newBookList))
	}
	AddBooksToStore(newBookList)
	return nil
}

// initRemoteStore 扫描远程书库（WebDAV 等）
func initRemoteStore(storeURL string, cfg ConfigInterface) error {
	logger.Infof("开始扫描远程书库: %s", storeURL)

	// 创建 VFS 实例
	opts := vfs.Options{
		Debug:        cfg.GetDebug(),
		Timeout:      cfg.GetTimeoutLimitForScan(),
		CacheEnabled: true, // 启用 VFS 层面的文件读取缓存（用于优化性能，压缩包仍为流式读取）
	}

	fs, err := vfs.GetOrCreate(storeURL, opts)
	if err != nil {
		return fmt.Errorf("连接远程书库失败: %w", err)
	}
	SetCurrentFS(fs)
	defer func() {
		SetCurrentFS(nil)
		// 注意：不关闭 fs，因为它被注册在 vfs.registry 中供后续使用
	}()

	// 获取基础路径
	basePath := ""
	if webdavFS, ok := fs.(*vfs.WebDAVFS); ok {
		basePath = webdavFS.GetBasePath()
	}
	if basePath == "" {
		basePath = "/"
	}

	logger.Infof(locale.GetString("scan_start_hint")+" %s (远程路径: %s)", storeURL, basePath)

	// 检查路径是否存在
	exists, err := fs.Exists(basePath)
	if err != nil || !exists {
		return fmt.Errorf("远程路径不存在或无法访问: %s", basePath)
	}

	// 检查是否为目录
	isDir, err := fs.IsDir(basePath)
	if err != nil {
		return fmt.Errorf("无法获取远程路径信息: %w", err)
	}

	var newBookList []*model.Book

	if !isDir {
		// 单个文件
		if IsSupportFile(basePath) {
			if checkBookInStore(storeURL, basePath) {
				return nil
			}
			book, err := scanRemoteFileGetBook(fs, basePath, storeURL, 0)
			if err != nil {
				return err
			}
			AddBooksToStore([]*model.Book{book})
		}
		return nil
	}

	// 目录扫描
	rootNode, foundDirs, foundFiles, err := HandleDirectoryVFS(fs, basePath, 0)
	if err != nil {
		return err
	}

	// 处理根目录
	if rootNode.Files != nil && len(rootNode.Files) >= cfg.GetMinImageNum() {
		book, err := scanRemoteDirGetBook(fs, basePath, storeURL, 0)
		if err != nil {
			logger.Infof(locale.GetString("log_skip_to_scan_root_directory"), basePath, err)
		} else {
			newBookList = append(newBookList, book)
		}
	}

	// 处理子目录
	for _, dir := range foundDirs {
		// 跳过根目录（已处理）
		if dir == basePath {
			continue
		}

		// 计算路径深度
		relPath, err := fs.RelPath(basePath, dir)
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_get_relative_path"), err)
			continue
		}
		depth := strings.Count(relPath, "/")

		// 检查是否超过最大深度限制
		if cfg.GetMaxScanDepth() >= 0 && depth > cfg.GetMaxScanDepth() {
			logger.Infof(locale.GetString("exceeds_maximum_depth")+" %d, base: %s, scan: %s", cfg.GetMaxScanDepth(), basePath, dir)
			continue
		}

		// 检查是否在忽略列表中
		if IsSkipDir(dir) {
			logger.Infof(locale.GetString("skip_path")+" %s", dir)
			continue
		}

		// 扫描目录
		book, err := scanRemoteDirGetBook(fs, dir, storeURL, depth)
		if err != nil {
			logger.Infof(locale.GetString("log_skip_to_scan_directory"), dir, err)
			continue
		}
		newBookList = append(newBookList, book)
	}

	// 处理文件
	for _, file := range foundFiles {
		// 如果是以 . 开头的隐藏文件，跳过
		if strings.HasPrefix(getBaseName(file.Name), ".") {
			continue
		}
		if !IsSupportFile(file.Name) {
			continue
		}

		// 计算路径深度
		relPath, err := fs.RelPath(basePath, file.Path)
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_get_relative_path"), err)
			continue
		}
		depth := strings.Count(relPath, "/")

		// 检查是否超过最大深度限制
		if cfg.GetMaxScanDepth() >= 0 && depth > cfg.GetMaxScanDepth() {
			logger.Infof(locale.GetString("exceeds_maximum_depth")+" %d, base: %s, scan: %s", cfg.GetMaxScanDepth(), basePath, file.Path)
			continue
		}

		// 如果书库里已经有这本书，跳过
		if checkBookInStore(storeURL, file.Path) {
			continue
		}

		// 扫描文件
		book, err := scanRemoteFileGetBook(fs, file.Path, storeURL, depth)
		if err != nil {
			logger.Info(err)
			continue
		}
		newBookList = append(newBookList, book)
	}

	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("how_many_books_update"), storeURL, len(newBookList))
	}
	AddBooksToStore(newBookList)
	return nil
}

func checkBookInStore(storePathAbs string, filePathAbs string) bool {
	bookInStore, err := getBookByPath(storePathAbs, filePathAbs)
	if err != nil || bookInStore == nil {
		return false
	}
	logger.Infof(locale.GetString("log_book_data_already_exists"), bookInStore.BookID, filePathAbs)
	return true
}

// getBookByPath 获取指定路径的书籍
func getBookByPath(storePath string, filePath string) (*model.Book, error) {
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}
	for _, b := range allBooks {
		if b.StoreUrl == storePath && b.BookPath == filePath {
			return b, nil
		}
	}
	return nil, fmt.Errorf("err_getbook_cannot_find_by_path")
}
