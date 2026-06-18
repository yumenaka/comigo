package model

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/vfs"
)

// Book 定义书籍结构
type Book struct {
	BookInfo  // 嵌入 BookInfo 结构体
	BookMarks // 书签信息
	PageInfos // 书籍内所有页面的信息
}

// CloneForView 返回用于渲染和 API 输出的书籍副本，避免排序、URL 改写污染内存书库。
func (b *Book) CloneForView() *Book {
	if b == nil {
		return nil
	}
	clone := *b
	clone.PageInfos = append(PageInfos(nil), b.PageInfos...)
	clone.BookMarks = append(BookMarks(nil), b.BookMarks...)
	clone.ChildBooksID = append([]string(nil), b.ChildBooksID...)
	return &clone
}

// GuessCover 猜测书籍的封面
func (b *Book) GuessCover() (cover PageInfo) {
	// 按 cover/0/1 的常见命名优先猜测封面，找不到时回退到第一页。
	for i := range b.PageInfos {
		filenameLower := strings.ToLower(b.PageInfos[i].Name)
		filenameWithoutExt := strings.TrimSuffix(filenameLower, filepath.Ext(filenameLower))
		filenameTrimmed := strings.TrimLeft(filenameWithoutExt, "0")
		if strings.Contains(filenameWithoutExt, "cover") || filenameTrimmed == "" || filenameTrimmed == "1" {
			cover = b.PageInfos[i]
			return cover
		}
	}
	if len(b.PageInfos) > 0 {
		cover = b.PageInfos[0]
	}
	return cover
}

// NewBook 初始化 Book，设置文件路径、书名、BookID 等
func NewBook(bookPath string, modified time.Time, fileSize int64, storePath string, depth int, bookType SupportFileType) (*Book, error) {
	// 初始化书籍
	book := &Book{
		BookInfo: BookInfo{
			Modified:         modified,
			FileSize:         fileSize,
			InitComplete:     false,
			Depth:            depth,
			StoreUrl:         storePath,
			Type:             bookType,
			CreatedByVersion: config.GetVersion(), // 记录创建时的版本号
		},
	}
	// 设置文件路径、书名、BookID
	_, err := book.setFilePath(bookPath).setParentFolder(bookPath).setTitle(bookPath).setAuthor().initBookID(bookPath)
	return book, err
}

// setPageNum 设置书籍的页数
func (b *Book) setPageNum() {
	b.PageCount = len(b.PageInfos)
}

// SortPages 对页面进行排序
func (b *Book) SortPages(s string) {
	if b.Type == TypeEpub && s == "default" {
		return
	}
	if s != "" {
		b.PageInfos.SortImages(s)
	}
}

// SortPagesByImageList 根据给定的文件列表排序页面（用于 EPUB）
func (b *Book) SortPagesByImageList(imageList []string) {
	if len(imageList) == 0 {
		return
	}
	imageInfos := b.PageInfos
	var reSortList []PageInfo
	for _, imgName := range imageList {
		for _, imgInfo := range imageInfos {
			if imgInfo.Name == imgName {
				reSortList = append(reSortList, imgInfo)
				break
			}
		}
	}
	if len(reSortList) == 0 {
		logger.Infof(locale.GetString("epub_cannot_resort")+"%s", b.BookPath)
		return
	}
	// 添加不在列表中的图片
	for _, imgInfo := range imageInfos {
		found := false
		for _, imgName := range imageList {
			if imgInfo.Name == imgName {
				found = true
				break
			}
		}
		if !found {
			reSortList = append(reSortList, imgInfo)
		}
	}
	b.PageInfos = reSortList
}

// GetStoreID 获取编码后的书库ID，StoreID是书库URL路径的 base62 编码
func (b *Book) GetStoreID() string {
	return base62.EncodeToString([]byte(b.StoreUrl))
}

// GetPageCount 获取页数
func (b *Book) GetPageCount() int {
	if !b.InitComplete {
		b.setPageNum()
		b.InitComplete = true
	}
	return b.PageCount
}

// GetAllBooksNumber  获取书籍总数，不包括 BookGroup 类型
func GetAllBooksNumber() int {
	// 用于计数的变量
	var count int
	// 遍历 map 并递增计数器
	allBooks, err := IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}
	for _, b := range allBooks {
		if b.Type == TypeBooksGroup {
			continue // 跳过书组类型
		}
		count++
	}
	return count
}

// ClearBookNotExist  检查内存中的书的源文件是否存在，不存在就删掉
func ClearBookNotExist() {
	logger.Info(locale.GetString("log_checking_book_files_exist"))
	// 遍历所有书籍
	var deletedBooks []string
	allBooks, err := IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}
	for _, book := range allBooks {
		if book.RemoteBookID != "" || book.RemoteStoreKey != "" {
			// Comigo 远端书籍和本地生成的远端书组由重扫时对比远端书架清理，不通过 VFS 检查。
			continue
		}
		var exists bool
		if book.IsRemote {
			// 远程书籍：使用 VFS 检查文件是否存在
			fs, fsErr := vfs.GetOrCreate(book.RemoteURL, vfs.Options{
				CacheEnabled: false,
				Timeout:      10, // 使用较短的超时时间，因为只是检查存在性
			})
			if fsErr != nil {
				// 无法连接远程服务器，跳过检查（可能是网络问题）
				logger.Infof(locale.GetString("log_remote_store_check_book_existence_failed"), book.RemoteURL, fsErr)
				continue
			}
			exists, err = fs.Exists(book.BookPath)
			if err != nil {
				// 检查出错，跳过这本书（可能是网络问题或路径问题，不删除书籍）
				logger.Infof(locale.GetString("log_remote_book_existence_check_failed"), book.BookPath, err)
				logger.Infof(locale.GetString("log_remote_book_existence_check_failed_detail"),
					book.BookID, book.RemoteURL, book.BookPath, err)
				continue
			}
		} else {
			// 本地书籍：使用 os.Stat 检查文件是否存在
			if _, err := os.Stat(book.BookPath); os.IsNotExist(err) {
				exists = false
			} else if err != nil {
				// 其他错误，跳过这本书
				logger.Infof(locale.GetString("log_local_book_existence_check_failed"), book.BookPath, err)
				continue
			} else {
				exists = true
			}
		}

		// 如果文件不存在，删除这本书
		if !exists {
			deletedBooks = append(deletedBooks, book.BookPath)
			err := IStore.DeleteBook(book.BookID)
			if err != nil {
				logger.Infof(locale.GetString("log_error_deleting_book"), book.BookID, err)
			}
			continue
		}

		if book.Type == TypeDir {
			changed, empty := cleanupMissingDirPages(book)
			if empty {
				deletedBooks = append(deletedBooks, book.BookPath)
				err := IStore.DeleteBook(book.BookID)
				if err != nil {
					logger.Infof(locale.GetString("log_error_deleting_book"), book.BookID, err)
				}
				continue
			}
			if changed {
				book.setPageNum()
				book.Cover = book.GuessCover()
				if err := IStore.StoreBook(book); err != nil {
					logger.Infof(locale.GetString("log_error_adding_book"), book.BookID, err)
				}
			}
		}
	}
	// 重新生成书组
	if len(deletedBooks) > 0 {
		if err := IStore.GenerateBookGroup(); err != nil {
			logger.Infof(locale.GetString("log_error_initializing_main_folder"), err)
		}
	}
}

// cleanupMissingDirPages 清理目录书籍中已经不存在的页面文件。
// 返回值 changed 表示 PageInfos 有变化；empty 表示页面已全部失效，需要删除整本目录书。
func cleanupMissingDirPages(book *Book) (changed bool, empty bool) {
	if book == nil || book.Type != TypeDir {
		return false, false
	}
	if len(book.PageInfos) == 0 {
		return false, true
	}
	remaining := make(PageInfos, 0, len(book.PageInfos))
	for _, page := range book.PageInfos {
		exists, err := dirPageExists(book, page)
		if err != nil {
			// 权限或网络错误不等同于文件不存在，保守保留页面，避免误删 metadata。
			if book.IsRemote {
				logger.Infof(locale.GetString("log_remote_book_existence_check_failed"), pagePathForExistCheck(book, page), err)
			} else {
				logger.Infof(locale.GetString("log_local_book_existence_check_failed"), pagePathForExistCheck(book, page), err)
			}
			remaining = append(remaining, page)
			continue
		}
		if exists {
			remaining = append(remaining, page)
		}
	}
	if len(remaining) == len(book.PageInfos) {
		return false, false
	}
	book.PageInfos = remaining
	return true, len(remaining) == 0
}

// dirPageExists 按本地/远程目录书籍分别检查页面文件是否存在。
func dirPageExists(book *Book, page PageInfo) (bool, error) {
	pagePath := pagePathForExistCheck(book, page)
	if book.IsRemote {
		fs, err := vfs.GetOrCreate(book.RemoteURL, vfs.Options{
			CacheEnabled: false,
			Timeout:      10,
		})
		if err != nil {
			return false, err
		}
		return fs.Exists(pagePath)
	}
	_, err := os.Stat(pagePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// pagePathForExistCheck 兼容旧 metadata：优先使用 PageInfo.Path，缺失时用书籍目录和页面名拼出路径。
func pagePathForExistCheck(book *Book, page PageInfo) string {
	if page.Path != "" {
		return page.Path
	}
	if book.IsRemote {
		return path.Join(book.BookPath, page.Name)
	}
	return filepath.Join(book.BookPath, page.Name)
}

// ClearBookWhenStoreUrlNotExist 清理书库中不存在的书籍源对应的书籍，传入的是当前存在的书籍源
func ClearBookWhenStoreUrlNotExist(nowStoreUrls []string) {
	logger.Info(locale.GetString("log_checking_store_exist"))
	// 遍历所有书籍
	var deletedBooks []string
	allBooks, err := IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
		return
	}
	// 如果传入的书库URL列表为空，清空所有书籍
	if len(nowStoreUrls) == 0 {
		for _, book := range allBooks {
			deletedBooks = append(deletedBooks, book.BookPath)
			err := IStore.DeleteBook(book.BookID)
			if err != nil {
				logger.Infof(locale.GetString("log_error_deleting_book"), book.BookID, err)
			}
		}
	} else {
		// 将传入的书库URL转换为绝对路径，用于后续比较
		normalizedStoreUrls := make(map[string]bool, len(nowStoreUrls))
		for _, storeUrl := range nowStoreUrls {
			storePathAbs, err := filepath.Abs(storeUrl)
			if err != nil {
				logger.Infof(locale.GetString("log_error_getting_absolute_path"), err)
				storePathAbs = storeUrl
			}
			normalizedStoreUrls[storePathAbs] = true
		}
		// 遍历所有书籍，删除不在当前书库列表中的书籍
		for _, book := range allBooks {
			// 将书籍的书库URL转换为绝对路径
			bookStoreUrlAbs, err := filepath.Abs(book.StoreUrl)
			if err != nil {
				logger.Infof(locale.GetString("log_error_getting_absolute_path"), err)
				bookStoreUrlAbs = book.StoreUrl
			}
			// 如果书籍的书库URL不在当前存在的书库列表中，删除该书籍
			if !normalizedStoreUrls[bookStoreUrlAbs] {
				deletedBooks = append(deletedBooks, book.BookPath)
				err := IStore.DeleteBook(book.BookID)
				if err != nil {
					logger.Infof(locale.GetString("log_error_deleting_book"), book.BookID, err)
				}
			}
		}
	}
	// 重新生成书组
	if len(deletedBooks) > 0 {
		if err := IStore.GenerateBookGroup(); err != nil {
			logger.Infof(locale.GetString("log_error_initializing_main_folder"), err)
		}
	}
}
