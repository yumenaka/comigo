package scan

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/comigo_remote"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/vfs"
)

// InitStore 扫描书库路径，取得路径里的书籍
// 支持本地路径和远程 URL（WebDAV 等）
func InitStore(storePath string, cfg ConfigInterface) error {
	InitConfig(cfg)

	if tools.DetectStoreURL(storePath).Type == tools.StoreBackendComigo {
		return initComigoStore(storePath, cfg)
	}

	// 判断是否为远程 URL
	isRemote := tools.IsRemoteStoreURL(storePath)

	if isRemote {
		return initRemoteStore(storePath, cfg)
	}
	return initLocalStore(storePath, cfg)
}

// initComigoStore 扫描另一个 Comigo 服务，只保存远端元数据和代理定位信息。
func initComigoStore(storeURL string, cfg ConfigInterface) error {
	logger.Infof(locale.GetString("log_scan_remote_store_start"), storeURL)
	client, err := comigo_remote.NewClient(storeURL, cfg.GetTimeoutLimitForScan())
	if err != nil {
		return err
	}
	shelves, err := client.GetTopShelf("default")
	if err != nil {
		return err
	}

	fetched := map[string]bool{}
	var books []*model.Book
	scanStart := time.Now()
	lastProgressLog := scanStart.Add(-2 * time.Second)
	// 每 2 秒输出一次远端扫描进度；大量小请求不会触发 timeout 等待日志，也要让用户看到进展。
	logProgress := func(stage string) {
		if time.Since(lastProgressLog) < 2*time.Second {
			return
		}
		logger.Infof(locale.GetString("log_scan_remote_comigo_progress"), storeURL, len(fetched), len(books), stage, time.Since(scanStart).Round(time.Second))
		lastProgressLog = time.Now()
	}
	logProgress("top-shelf")

	// skipped 记录远端里本身就是远程书库的条目；它们不进入 fetched，重扫时会清掉旧的嵌套代理数据。
	skipped := map[string]bool{}
	// fetchBook 返回该远端书是否实际导入；父书组据此过滤 ChildBooksID，避免留下指向被跳过书籍的子项。
	var fetchBook func(remoteBookID string, shelfKey string, shelfName string, topLevel bool) (bool, error)
	fetchBook = func(remoteBookID string, shelfKey string, shelfName string, topLevel bool) (bool, error) {
		localBookID := comigo_remote.LocalBookID(storeURL, remoteBookID)
		if fetched[localBookID] {
			return true, nil
		}
		if skipped[localBookID] {
			return false, nil
		}
		remoteBook, err := client.GetBook(remoteBookID, "default")
		if err != nil {
			return false, err
		}
		// 远端 Comigo 中的远程书籍只属于远端服务本身，本机不再套一层远程代理。
		if !shouldImportRemoteComigoBook(remoteBook) {
			skipped[localBookID] = true
			logProgress(remoteBook.Title)
			return false, nil
		}
		localBook := comigo_remote.LocalizeBookInShelf(storeURL, remoteBook, shelfKey, shelfName)
		if topLevel {
			// 远端顶层列表才知道它属于哪个顶级书库；旧版详情 API 不暴露远端 StoreUrl。
			localBook.Depth = 0
		}
		fetched[localBook.BookID] = true
		books = append(books, localBook)
		logProgress(remoteBook.Title)
		// 重新生成子书列表，只保留本轮已经导入的非远程子书。
		childBookIDs := make([]string, 0, len(remoteBook.ChildBooksID))
		for _, childRemoteID := range remoteBook.ChildBooksID {
			imported, err := fetchBook(childRemoteID, shelfKey, shelfName, false)
			if err != nil {
				logger.Infof(locale.GetString("log_get_book_error"), err)
				continue
			}
			if imported {
				childBookIDs = append(childBookIDs, comigo_remote.LocalBookID(storeURL, childRemoteID))
			}
		}
		localBook.ChildBooksID = childBookIDs
		return true, nil
	}

	for shelfIndex, shelf := range shelves {
		shelfKey := comigo_remote.ShelfKey(storeURL, shelf, shelfIndex)
		shelfName := shelf.DisplayName
		for _, bookInfo := range shelf.BookInfos {
			if _, err := fetchBook(bookInfo.BookID, shelfKey, shelfName, true); err != nil {
				logger.Infof(locale.GetString("log_get_book_error"), err)
			}
		}
	}

	deleteStaleComigoRemoteBooks(storeURL, fetched)
	AddBooksToStore(books)
	model.GenerateBookGroup()
	return nil
}

// shouldImportRemoteComigoBook 避免导入远端服务中的远程书籍，防止 WebDAV/另一台 Comigo 等书库被二次嵌套代理。
func shouldImportRemoteComigoBook(remoteBook *model.Book) bool {
	return remoteBook != nil && !remoteBook.IsRemote
}

// deleteStaleComigoRemoteBooks 删除本次远端 Comigo 扫描后已经过期的本地条目。
// 本地生成的远程书组没有 RemoteBookID，必须在重新生成前整批清掉，否则远端书库删除后会残留旧书组。
func deleteStaleComigoRemoteBooks(storeURL string, current map[string]bool) {
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
		return
	}
	for _, book := range allBooks {
		if book.StoreUrl != storeURL {
			continue
		}
		if book.Type == model.TypeBooksGroup && book.RemoteBookID == "" && book.RemoteStoreKey != "" {
			if err := model.IStore.DeleteBook(book.BookID); err != nil {
				logger.Infof(locale.GetString("log_error_deleting_book"), book.BookID, err)
			}
			continue
		}
		if book.RemoteBookID == "" || current[book.BookID] {
			continue
		}
		if err := model.IStore.DeleteBook(book.BookID); err != nil {
			logger.Infof(locale.GetString("log_error_deleting_book"), book.BookID, err)
		}
	}
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
		fileInfo, statErr := os.Stat(storePathAbs)
		if statErr != nil {
			return statErr
		}
		if shouldSkipFailedArchiveFile(storePathAbs, storePathAbs, fileInfo.Size(), fileInfo.ModTime(), false) {
			return nil
		}
		previousBook, skip := prepareBookPathForScan(storePathAbs, storePathAbs, fileInfo.Size(), fileInfo.ModTime())
		if skip {
			return nil
		}
		book, err := scanFileGetBook(storePathAbs, storePathAbs, 0)
		if err != nil {
			restorePreviousBookAfterScanFailure(previousBook)
			recordArchiveScanFailure(storePathAbs, storePathAbs, fileInfo.Size(), fileInfo.ModTime(), false, err)
			return err
		}
		mergePreviousBookState(book, previousBook)
		clearArchiveScanFailure(storePathAbs, storePathAbs, false)
		AddBooksToStore([]*model.Book{book})
		return nil
	}

	// 如果书库URL是一个文件夹，使用 HandleDirectory（不支持扫描单个文件）进行扫描
	rootDirectoryNode, _, foundFiles, err := HandleDirectory(storePathAbs, 0)
	if err != nil {
		return err
	}

	var newBookList []*model.Book

	appendLocalDirBooks(&newBookList, rootDirectoryNode, storePathAbs)

	// 处理文件
	for _, file := range foundFiles {
		// 如果是以 . 开头的隐藏文件，跳过
		if strings.HasPrefix(filepath.Base(file.Name), ".") {
			continue
		}
		if !IsSupportFile(file.Name) {
			if IsSupportMedia(file.Name) {
				continue
			}
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
		if shouldSkipFailedArchiveFile(storePathAbs, file.Path, file.Size, file.ModTime, false) {
			continue
		}
		previousBook, skip := prepareBookPathForScan(storePathAbs, file.Path, file.Size, file.ModTime)
		if skip {
			continue
		}
		// 扫描文件
		book, err := scanFileGetBook(file.Path, storePathAbs, depth)
		if err != nil {
			restorePreviousBookAfterScanFailure(previousBook)
			recordArchiveScanFailure(storePathAbs, file.Path, file.Size, file.ModTime, false, err)
			logger.Info(err)
			continue
		}
		mergePreviousBookState(book, previousBook)
		clearArchiveScanFailure(storePathAbs, file.Path, false)
		newBookList = append(newBookList, book)
	}

	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("how_many_books_update"), storePathAbs, len(newBookList))
	}
	AddBooksToStore(newBookList)
	return nil
}

// appendLocalDirBooks 复用 HandleDirectory 已经读到的文件列表，避免本地目录扫描阶段重复 ReadDir。
func appendLocalDirBooks(newBookList *[]*model.Book, node DirNode, storePathAbs string) {
	appendLocalDirBook(newBookList, node, storePathAbs)
	for _, subNode := range node.SubDirs {
		appendLocalDirBooks(newBookList, subNode, storePathAbs)
	}
}

// appendLocalDirBook 处理单个本地目录节点，并保留旧书签迁移与失败回滚逻辑。
func appendLocalDirBook(newBookList *[]*model.Book, node DirNode, storePathAbs string) {
	depth, err := localScanDepth(storePathAbs, node.Path)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_get_relative_path"), err)
		return
	}
	if cfg.GetMaxScanDepth() >= 0 && depth > cfg.GetMaxScanDepth() {
		logger.Infof(locale.GetString("exceeds_maximum_depth")+" %d, base: %s, scan: %s", cfg.GetMaxScanDepth(), storePathAbs, node.Path)
		return
	}
	dirInfo, statErr := os.Stat(node.Path)
	if statErr != nil {
		logger.Infof(locale.GetString("log_failed_to_get_file_info_scan"), node.Path, statErr)
		return
	}
	previousBook, skip := prepareBookPathForScan(storePathAbs, node.Path, dirInfo.Size(), dirInfo.ModTime())
	if skip {
		return
	}
	book, err := bookFromLocalDirNode(node, storePathAbs, depth)
	if err != nil {
		restorePreviousBookAfterScanFailure(previousBook)
		logger.Infof(locale.GetString("log_skip_to_scan_directory"), node.Path, err)
		return
	}
	mergePreviousBookState(book, previousBook)
	*newBookList = append(*newBookList, book)
}

// localScanDepth 计算与旧逻辑一致的扫描深度：顶层文件夹为 0，每多一级子目录加 1。
func localScanDepth(storePathAbs string, path string) (int, error) {
	relPath, err := filepath.Rel(storePathAbs, path)
	if err != nil {
		return 0, err
	}
	return strings.Count(relPath, string(os.PathSeparator)), nil
}

// initRemoteStore 扫描远程书库（WebDAV 等）
func initRemoteStore(storeURL string, cfg ConfigInterface) error {
	logger.Infof(locale.GetString("log_scan_remote_store_start"), storeURL)

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
	} else if sftpFS, ok := fs.(*vfs.SFTPFS); ok {
		basePath = sftpFS.GetBasePath()
	} else if smbFS, ok := fs.(*vfs.SMBFS); ok {
		basePath = smbFS.GetBasePath()
	} else if ftpFS, ok := fs.(*vfs.FTPFS); ok {
		basePath = ftpFS.GetBasePath()
	} else if s3FS, ok := fs.(*vfs.S3FS); ok {
		basePath = s3FS.GetBasePath()
	}
	if basePath == "" {
		basePath = "/"
	}

	logger.Infof(locale.GetString("log_scan_start_hint_remote"), storeURL, basePath)

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
			fileInfo, statErr := fs.Stat(basePath)
			if statErr != nil {
				return statErr
			}
			if shouldSkipFailedArchiveFile(storeURL, basePath, fileInfo.Size(), fileInfo.ModTime(), true) {
				return nil
			}
			previousBook, skip := prepareBookPathForScan(storeURL, basePath, fileInfo.Size(), fileInfo.ModTime())
			if skip {
				return nil
			}
			book, err := scanRemoteFileGetBook(fs, basePath, storeURL, 0)
			if err != nil {
				restorePreviousBookAfterScanFailure(previousBook)
				recordArchiveScanFailure(storeURL, basePath, fileInfo.Size(), fileInfo.ModTime(), true, err)
				return err
			}
			mergePreviousBookState(book, previousBook)
			clearArchiveScanFailure(storeURL, basePath, true)
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
		dirInfo, statErr := fs.Stat(basePath)
		if statErr != nil {
			return statErr
		}
		previousBook, skip := prepareBookPathForScan(storeURL, basePath, dirInfo.Size(), dirInfo.ModTime())
		if !skip {
			book, err := scanRemoteDirGetBook(fs, basePath, storeURL, 0)
			if err != nil {
				restorePreviousBookAfterScanFailure(previousBook)
				logger.Infof(locale.GetString("log_skip_to_scan_root_directory"), basePath, err)
			} else {
				mergePreviousBookState(book, previousBook)
				newBookList = append(newBookList, book)
			}
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

		dirInfo, statErr := fs.Stat(dir)
		if statErr != nil {
			logger.Infof(locale.GetString("log_failed_to_get_file_info_scan"), dir, statErr)
			continue
		}
		previousBook, skip := prepareBookPathForScan(storeURL, dir, dirInfo.Size(), dirInfo.ModTime())
		if skip {
			continue
		}

		// 扫描目录
		book, err := scanRemoteDirGetBook(fs, dir, storeURL, depth)
		if err != nil {
			restorePreviousBookAfterScanFailure(previousBook)
			logger.Infof(locale.GetString("log_skip_to_scan_directory"), dir, err)
			continue
		}
		mergePreviousBookState(book, previousBook)
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

		if shouldSkipFailedArchiveFile(storeURL, file.Path, file.Size, file.ModTime, true) {
			continue
		}
		previousBook, skip := prepareBookPathForScan(storeURL, file.Path, file.Size, file.ModTime)
		if skip {
			continue
		}

		// 扫描文件
		book, err := scanRemoteFileGetBook(fs, file.Path, storeURL, depth)
		if err != nil {
			restorePreviousBookAfterScanFailure(previousBook)
			recordArchiveScanFailure(storeURL, file.Path, file.Size, file.ModTime, true, err)
			logger.Info(err)
			continue
		}
		mergePreviousBookState(book, previousBook)
		clearArchiveScanFailure(storeURL, file.Path, true)
		newBookList = append(newBookList, book)
	}

	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("how_many_books_update"), storeURL, len(newBookList))
	}
	AddBooksToStore(newBookList)
	return nil
}

// prepareBookPathForScan 在扫描前处理同路径旧书籍：
// - 指纹未变化时跳过扫描，沿用旧 metadata；
// - 指纹变化时先移除旧条目，避免 NewBook 的同路径判重阻止重扫。
// 返回 previousBook 用于扫描失败回滚或扫描成功后迁移用户状态。
func prepareBookPathForScan(storePath string, filePath string, size int64, modTime time.Time) (previousBook *model.Book, skip bool) {
	previousBook, err := getBookByPath(storePath, filePath)
	if err != nil || previousBook == nil {
		return nil, false
	}
	if previousBook.FileSize == size && previousBook.Modified.Equal(modTime) {
		logger.Infof(locale.GetString("log_book_data_already_exists"), previousBook.BookID, filePath)
		return previousBook, true
	}
	if err := model.IStore.DeleteBook(previousBook.BookID); err != nil {
		logger.Infof(locale.GetString("log_error_deleting_book"), previousBook.BookID, err)
		return previousBook, true
	}
	return previousBook, false
}

func mergePreviousBookState(newBook *model.Book, previousBook *model.Book) {
	if newBook == nil || previousBook == nil {
		return
	}
	// 重扫后 BookID 可能因文件大小变化而改变；迁移书签时同步改到新书 ID。
	for _, mark := range previousBook.BookMarks {
		mark.BookID = newBook.BookID
		mark.BookStoreID = newBook.GetStoreID()
		newBook.BookMarks = append(newBook.BookMarks, mark)
	}
	newBook.BookComplete = previousBook.BookComplete
}

func restorePreviousBookAfterScanFailure(previousBook *model.Book) {
	if previousBook == nil {
		return
	}
	if err := model.IStore.StoreBook(previousBook); err != nil {
		logger.Infof(locale.GetString("log_error_adding_book"), previousBook.BookID, err)
	}
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
