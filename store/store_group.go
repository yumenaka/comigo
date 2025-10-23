package store

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// StoreGroup 内存书库，扫描后生成。可以有多个子书库。
type StoreGroup struct {
	StoreInfo
	ChildStores sync.Map // key为路径 存储 *Store
}

// AddStore 创建一个新书库
func (storeGroup *StoreGroup) AddStore(storeURL string) error {
	if _, ok := storeGroup.ChildStores.Load(storeURL); ok {
		// 已经有这个key了
		return errors.New("add Bookstore Error： The key already exists [" + storeURL + "]")
	}
	s := Store{
		StoreInfo: StoreInfo{
			BackendURL:  storeURL,
			Name:        filepath.Base(storeURL), // 使用路径的最后一部分作为名称
			Description: "Comigo StoreInfo for " + storeURL,
		},
	}
	storeGroup.ChildStores.Store(storeURL, &s)
	return nil
}

// GenerateAllBookGroup 分析所有子书库，并并生成书籍组
func (storeGroup *StoreGroup) GenerateAllBookGroup() (e error) {
	// 遍历所有子书库
	for _, value := range storeGroup.ChildStores.Range {
		s := value.(*Store)
		err := s.GenerateBookGroup()
		if err != nil {
			e = err
		}
	}
	return e
}

func (storeGroup *StoreGroup) ListBooks() []*model.Book {
	var books []*model.Book
	// 遍历 ChildStores 中的所有书籍
	for _, value := range storeGroup.ChildStores.Range {
		childStore := value.(*Store)
		for _, value := range childStore.BookMap.Range {
			book := value.(*model.Book)
			books = append(books, book)
		}
	}
	return books
}

// AddBook 添加一本书
func (storeGroup *StoreGroup) AddBook(b *model.Book, minPageNum int) error {
	if b.BookID == "" {
		return errors.New("add book error: empty BookID")
	}
	if b.GetPageCount() < minPageNum {
		return errors.New("add book error: minPageNum = " + strconv.Itoa(b.GetPageCount()))
	}
	if _, ok := storeGroup.ChildStores.Load(b.BookStorePath); !ok {
		if err := storeGroup.AddStore(b.BookStorePath); err != nil {
			logger.Infof("Error adding subfolder: %s", err)
		}
	}
	return storeGroup.AddBookToSubStore(b.BookStorePath, b)
}

// AddBookToSubStore 将某一本书，放到BookMap里面去
func (storeGroup *StoreGroup) AddBookToSubStore(storeURL string, b *model.Book) error {
	if f, ok := storeGroup.ChildStores.Load(storeURL); !ok {
		// 创建一个新子书库，并添加一本书
		newSubStore := Store{
			StoreInfo: StoreInfo{
				BackendURL:  storeURL,
				Name:        filepath.Base(storeURL), // 使用路径的最后一部分作为名称
				Description: "Comigo StoreInfo for " + storeURL,
			},
		}
		newSubStore.BookMap.Store(b.BookID, b)
		storeGroup.ChildStores.Store(storeURL, &newSubStore)
		return errors.New("add Bookstore Error： The key not found [" + storeURL + "]")
	} else {
		// 给已有子书库添加一本书
		temp := f.(*Store)
		temp.BookMap.Store(b.BookID, b)
		return nil
	}
}

// AddBooks 添加一组书
func (storeGroup *StoreGroup) AddBooks(list []*model.Book, minPageNum int) error {
	for _, b := range list {
		if b.GetPageCount() < minPageNum {
			fmt.Println("Skip book, page number too small:", b.Title, "Pages:", b.GetPageCount())
			continue
		}
		if err := storeGroup.AddBook(b, minPageNum); err != nil {
			return err
		}
	}
	return nil
}

// CheckAllNotExistBooks 检查内存中的书的源文件是否存在，不存在就删掉
func (storeGroup *StoreGroup) ClearBookNotExist() {
	logger.Infof("Checking book files exist...")
	var deletedBooks []string
	// 遍历所有书籍
	for _, book := range storeGroup.ListBooks() {
		// 如果父文件夹存在，但书籍文件不存在，也说明这本书被删除了
		if _, err := os.Stat(book.FilePath); os.IsNotExist(err) {
			deletedBooks = append(deletedBooks, book.FilePath)
			storeGroup.DeleteBook(book.BookID)
		}
	}
	// 重新生成书组
	if len(deletedBooks) > 0 {
		if err := storeGroup.GenerateAllBookGroup(); err != nil {
			logger.Infof("Error initializing main folder: %s", err)
		}
	}
}

// CheckRawFileExist 查看内存中是否已经有相同路径的书，有了就返回 true，让调用者跳过
func (storeGroup *StoreGroup) CheckBookFileExist(filePath string, bookType model.SupportFileType) bool {
	if bookType == model.TypeDir || bookType == model.TypeBooksGroup {
		return false
	}
	for _, realBook := range storeGroup.ListBooks() {
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			logger.Infof("Error getting absolute path: %v", err)
			continue
		}
		if realBook.FilePath == absFilePath && realBook.Type == bookType {
			return true
		}
	}
	return false
}

// GetParentBookID 通过子书籍 ID 获取所属书组 ID
func (storeGroup *StoreGroup) GetParentBookID(childID string) (string, error) {
	for _, bookGroup := range storeGroup.ListBooks() {
		if bookGroup.Type != model.TypeBooksGroup {
			continue // 只处理书组类型
		}
		for _, id := range bookGroup.ChildBooksID {
			if id == childID {
				fmt.Println("Found group for book ID:", childID, "Group ID:", bookGroup.BookID)
				return bookGroup.BookID, nil
			}
		}
	}
	return "", errors.New("cannot find group, id=" + childID)
}

// ClearAll 清空所有书籍与虚拟书组数据
func (storeGroup *StoreGroup) ClearAllBook() {
	storeGroup.ChildStores.Clear() // Clear会删除所有条目，产生一个空 Map。
}

// DeleteBook 删除一本书
func (storeGroup *StoreGroup) DeleteBook(bookID string) {
	// 遍历 ChildStores ，删除指定 ID 的书籍
	for _, value := range storeGroup.ChildStores.Range {
		childStore := value.(*Store)
		childStore.BookMap.Delete(bookID)
	}
}

// GetBooksNumber 获取书籍总数，不包括 BookGroup 类型
func (storeGroup *StoreGroup) GetAllBooksNumber() int {
	// 用于计数的变量
	var count int
	// 遍历 map 并递增计数器
	for _, b := range storeGroup.ListBooks() {
		if b.Type == model.TypeBooksGroup {
			continue // 跳过书组类型
		}
		count++
	}
	return count
}

// GetAllBookList 获取所有书籍列表
func (storeGroup *StoreGroup) ListBookSkipBookGroup() []*model.Book {
	var list []*model.Book
	for _, b := range storeGroup.ListBooks() {
		if b.Type == model.TypeBooksGroup {
			continue // 跳过书组类型
		}
		list = append(list, b)
	}
	return list
}

// GetArchiveBooks 获取所有压缩包格式的书籍
func (storeGroup *StoreGroup) GetArchiveBooks() []*model.Book {
	var list []*model.Book
	for _, b := range storeGroup.ListBooks() {
		if b.Type == model.TypeZip || b.Type == model.TypeRar || b.Type == model.TypeCbz || b.Type == model.TypeCbr || b.Type == model.TypeTar || b.Type == model.TypeEpub {
			list = append(list, b)
		}
	}
	return list
}

// GetBook 根据 BookID 获取书籍
// GetBookByID 根据 BookID 获取书籍
func (storeGroup *StoreGroup) GetBook(id string) (*model.Book, error) {
	// 遍历 ChildStores ，删除指定 ID 的书籍
	for _, value := range storeGroup.ChildStores.Range {
		childStore := value.(*Store)
		if value, ok := childStore.BookMap.Load(id); ok {
			b := value.(*model.Book)
			return b, nil
		}
	}
	return nil, errors.New("GetBook：cannot find book, id=" + id)
}

func (storeGroup *StoreGroup) UpdateBook(b *model.Book) error {
	for _, value := range storeGroup.ChildStores.Range {
		childStore := value.(*Store)
		if _, ok := childStore.BookMap.Load(b.BookID); ok {
			childStore.BookMap.Store(b.BookID, b)
			return nil
		}
	}
	return errors.New("UpdateBook：cannot find book, id=" + b.BookID)
}

// GetParentBook 通过子书籍的 BookID 获取父书组
func (storeGroup *StoreGroup) GetParentBook(childID string) (*model.Book, error) {
	for _, bookGroup := range storeGroup.ListBooks() {
		if bookGroup.Type != model.TypeBooksGroup {
			continue // 只分析书组类型
		}
		for _, id := range bookGroup.ChildBooksID {
			if id == childID {
				b, err := storeGroup.GetBook(bookGroup.BookID)
				if err != nil {
					return nil, errors.New("GetParentBook: cannot find book by childID=" + id)
				}
				return b, nil
			}
		}
	}
	return nil, errors.New("cannot find group, child book ID=" + childID)
}

// GetBookByAuthor 获取同一作者的书籍
func (storeGroup *StoreGroup) GetBookByAuthor(author string, sortBy string) ([]*model.Book, error) {
	var bookList []*model.Book
	for _, b := range storeGroup.ListBooks() {
		if b.Author == author {
			b.SortPages(sortBy)
			bookList = append(bookList, b)
		}
	}

	if len(bookList) > 0 {
		return bookList, nil
	}
	return nil, errors.New("cannot find book, author=" + author)
}

// ClearTempFilesALL 清理所有缓存图片
func (storeGroup *StoreGroup) ClearTempFilesALL(debug bool, cacheFilePath string) {
	for _, book := range storeGroup.ListBooks() {
		//清理某一本书的缓存
		cachePath := path.Join(cacheFilePath, book.GetBookID())
		err := os.RemoveAll(cachePath)
		if err != nil {
			logger.Infof("Error clearing temp files: %s", cachePath)
		} else if debug {
			logger.Infof("Cleared temp files: %s", cachePath)
		}
	}
}

// GetShortBookID 生成短的 BookID，避免冲突
func (storeGroup *StoreGroup) GetShortBookID(fullID string, minLength int) string {
	if len(fullID) <= minLength {
		logger.Infof("Cannot shorten ID: %s", fullID)
		return fullID
	}
	shortID := fullID[:minLength]
	add := 0
	for {
		conflict := false
		for _, b := range storeGroup.ListBooks() {
			if b.BookID == shortID {
				conflict = true
				break
			}
		}
		if !conflict {
			break
		}
		add++
		if minLength+add > len(fullID) {
			break
		}
		shortID = fullID[:minLength+add]
	}
	return shortID
}

// TopOfShelfInfo 获取顶层书架信息
func (storeGroup *StoreGroup) TopOfShelfInfo(sortBy string) (*model.BookInfoList, error) {
	// 显示顶层书库的书籍
	var infoList model.BookInfoList
	for _, b := range storeGroup.ListBooks() {
		if b.Depth == 0 {
			infoList.BookInfos = append(infoList.BookInfos, b.BookInfo)
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	// 没找到任何书
	return nil, errors.New("error: cannot find book in TopOfShelfInfo")
}

// GetChildBooksInfo 根据 ID 获取书籍列表
func (storeGroup *StoreGroup) GetChildBooksInfo(BookID string, sortBy string) (*model.BookInfoList, error) {
	var infoList model.BookInfoList
	tempGroup, err := storeGroup.GetBook(BookID)
	if err != nil {
		return nil, errors.New("cannot find child books info，BookID：" + BookID)
	}
	tempGroup.SortPages(sortBy)
	for _, childID := range tempGroup.ChildBooksID {
		b, err := RamStore.GetBook(childID)
		if err != nil {
			return nil, errors.New("GetParentBook: cannot find book by childID=" + childID)
		}
		infoList.BookInfos = append(infoList.BookInfos, b.BookInfo)
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	} else {
		return nil, errors.New("cannot find child books info，BookID：" + BookID)
	}
}

// GetBookInfoListByParentFolder 根据父文件夹获取书籍列表
func (storeGroup *StoreGroup) GetBookInfoListByParentFolder(parentFolder string, sortBy string) (*model.BookInfoList, error) {
	var infoList model.BookInfoList
	for _, b := range storeGroup.ListBooks() {
		if b.ParentFolder == parentFolder {
			infoList.BookInfos = append(infoList.BookInfos, b.BookInfo)
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("cannot find book, parentFolder=" + parentFolder)
}

// // GetAllBookInfoList 获取所有 BookInfo，并根据 sortBy 参数进行排序
// func (storeGroup *StoreGroup) GetAllBookInfoList(sortBy string) *BookInfoList {
// 	var infoList BookInfoList
// 	// 添加所有真实的书籍
// 	for _, b := range storeGroup.ListBooks() {
// 		info := b.GetBookInfo()
// 		infoList.BookInfos = append(infoList.BookInfos, *info)
// 	}
// 	infoList.SortBooks(sortBy)
// 	return &infoList
// }
//
// func (storeGroup *StoreGroup) GetBookInfoListByMaxDepth(depth int, sortBy string) (*BookInfoList, error) {
// 	var infoList BookInfoList
// 	// 首先加上所有真实的书籍
// 	for _, b := range storeGroup.ListBooks() {
// 		if b.Depth <= depth {
// 			info := b.GetBookInfo()
// 			infoList.BookInfos = append(infoList.BookInfos, *info)
// 		}
// 	}
// 	if len(infoList.BookInfos) > 0 {
// 		infoList.SortBooks(sortBy)
// 		return &infoList, nil
// 	}
// 	return nil, errors.New("error: cannot find bookshelf in GetBookInfoListByMaxDepth")
// }
