package store

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// StoreInRam 内存书库，扫描后生成。可以有多个子书库。
type StoreInRam struct {
	StoreInfo
	ChildStores sync.Map // key为路径 存储 *Store
}

// AddStore 创建一个新书库
func (ramStore *StoreInRam) AddStore(storeURL string) error {
	if _, ok := ramStore.ChildStores.Load(storeURL); ok {
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
	ramStore.ChildStores.Store(storeURL, &s)
	return nil
}

// GenerateAllBookGroup 分析所有子书库，并并生成书籍组
func (ramStore *StoreInRam) GenerateAllBookGroup() (e error) {
	// 遍历所有子书库
	for _, value := range ramStore.ChildStores.Range {
		s := value.(*Store)
		err := s.GenerateBookGroup()
		if err != nil {
			e = err
		}
	}
	return e
}

func (ramStore *StoreInRam) ListBooks() ([]*model.Book, error) {
	var books []*model.Book
	// 遍历 ChildStores 中的所有书籍
	for _, value := range ramStore.ChildStores.Range {
		childStore := value.(*Store)
		for _, value := range childStore.BookMap.Range {
			book := value.(*model.Book)
			books = append(books, book)
		}
	}
	return books, nil
}

// AddBook 添加一本书
func (ramStore *StoreInRam) AddBook(b *model.Book) error {
	if b.BookID == "" {
		return errors.New("add book error: empty BookID")
	}
	if _, ok := ramStore.ChildStores.Load(b.StoreUrl); !ok {
		if err := ramStore.AddStore(b.StoreUrl); err != nil {
			logger.Infof("Error adding subfolder: %s", err)
		}
	}
	return ramStore.AddBookToSubStore(b.StoreUrl, b)
}

// AddBookToSubStore 将某一本书，放到BookMap里面去
func (ramStore *StoreInRam) AddBookToSubStore(storeURL string, b *model.Book) error {
	if f, ok := ramStore.ChildStores.Load(storeURL); !ok {
		// 创建一个新子书库，并添加一本书
		newSubStore := Store{
			StoreInfo: StoreInfo{
				BackendURL:  storeURL,
				Name:        filepath.Base(storeURL), // 使用路径的最后一部分作为名称
				Description: "Comigo StoreInfo for " + storeURL,
			},
		}
		newSubStore.BookMap.Store(b.BookID, b)
		ramStore.ChildStores.Store(storeURL, &newSubStore)
		return errors.New("add Bookstore Error： The key not found [" + storeURL + "]")
	} else {
		// 给已有子书库添加一本书
		temp := f.(*Store)
		temp.BookMap.Store(b.BookID, b)
		return nil
	}
}

// AddBooks 添加多本书
func (ramStore *StoreInRam) AddBooks(books []*model.Book) error {
	for _, b := range books {
		if err := ramStore.AddBook(b); err != nil {
			logger.Infof("Error adding book %s: %s", b.BookID, err)
		}
	}
	return nil
}

// GetParentBookID 通过子书籍 ID 获取所属书组 ID
func (ramStore *StoreInRam) GetParentBookID(childID string) (string, error) {
	allBooks, err := ramStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, bookGroup := range allBooks {
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

// GetArchiveBooks 获取所有压缩包格式的书籍
func (ramStore *StoreInRam) GetArchiveBooks() []*model.Book {
	var list []*model.Book
	allBooks, err := ramStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
		if b.Type == model.TypeZip || b.Type == model.TypeRar || b.Type == model.TypeCbz || b.Type == model.TypeCbr || b.Type == model.TypeTar || b.Type == model.TypeEpub {
			list = append(list, b)
		}
	}
	return list
}

// GetBook 根据 BookID 获取书籍
// GetBookByID 根据 BookID 获取书籍
func (ramStore *StoreInRam) GetBook(id string) (*model.Book, error) {
	// 遍历 ChildStores ，删除指定 ID 的书籍
	for _, value := range ramStore.ChildStores.Range {
		childStore := value.(*Store)
		if value, ok := childStore.BookMap.Load(id); ok {
			b := value.(*model.Book)
			return b, nil
		}
	}
	return nil, errors.New("GetBook：cannot find book, id=" + id)
}

func (ramStore *StoreInRam) DeleteBook(id string) error {
	for _, value := range ramStore.ChildStores.Range {
		childStore := value.(*Store)
		if _, ok := childStore.BookMap.Load(id); ok {
			childStore.BookMap.Delete(id)
			return nil
		}
	}
	return errors.New("DeleteBook：cannot find book, id=" + id)
}

func (ramStore *StoreInRam) UpdateBook(b *model.Book) error {
	for _, value := range ramStore.ChildStores.Range {
		childStore := value.(*Store)
		if _, ok := childStore.BookMap.Load(b.BookID); ok {
			childStore.BookMap.Store(b.BookID, b)
			return nil
		}
	}
	return errors.New("UpdateBook：cannot find book, id=" + b.BookID)
}

// GetBookByAuthor 获取同一作者的书籍
func (ramStore *StoreInRam) GetBookByAuthor(author string, sortBy string) ([]*model.Book, error) {
	var bookList []*model.Book
	allBooks, err := ramStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
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

// TopOfShelfInfo 获取顶层书架信息
func TopOfShelfInfo(sortBy string) (*model.BookInfoList, error) {
	// 显示顶层书库的书籍
	var infoList model.BookInfoList
	allBooks, err := model.IStore.ListBooks()
	//allBooksB, err := RamStore.ListBooks()
	//// 比较两个书库的数量是否一致
	//if err != nil {
	//	logger.Infof("Error listing books: %s", err)
	//}
	//if len(allBooks) != len(allBooksB) {
	//	logger.Infof("Warning: TopOfShelfInfo: the number of books in RamStore (%d) and IStore (%d) are not equal.", len(allBooksB), len(allBooks))
	//}
	//// 打印不一致的书籍ID
	//bookIDMap := make(map[string]bool)
	//for _, b := range allBooksB {
	//	bookIDMap[b.BookID] = true
	//}
	//for _, b := range allBooks {
	//	if _, ok := bookIDMap[b.BookID]; !ok {
	//		logger.Infof("Warning: TopOfShelfInfo: BookID %s is in IStore but not in RamStore.", b.BookID)
	//	}
	//}

	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
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
func GetChildBooksInfo(BookID string) (*model.BookInfoList, error) {
	var infoList model.BookInfoList
	parentBook, err := model.IStore.GetBook(BookID)
	if err != nil {
		return nil, errors.New("cannot find child books info，BookID：" + BookID)
	}
	for _, childID := range parentBook.ChildBooksID {
		b, err := model.IStore.GetBook(childID)
		if err != nil {
			return nil, errors.New("GetParentBook: cannot find book by childID=" + childID)
		}
		infoList.BookInfos = append(infoList.BookInfos, b.BookInfo)
	}
	if len(infoList.BookInfos) > 0 {
		return &infoList, nil
	} else {
		return nil, errors.New("cannot find child books info，BookID：" + BookID)
	}
}

// GetBookInfoListByParentFolder 根据父文件夹获取书籍列表
func GetBookInfoListByParentFolder(parentFolder string) (*model.BookInfoList, error) {
	var infoList model.BookInfoList
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
		if b.ParentFolder == parentFolder {
			infoList.BookInfos = append(infoList.BookInfos, b.BookInfo)
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks("filename")
		return &infoList, nil
	}
	return nil, errors.New("cannot find book, parentFolder=" + parentFolder)
}

// // GetAllBookInfoList 获取所有 BookInfo，并根据 sortBy 参数进行排序
// func (storeGroup *StoreInRam) GetAllBookInfoList(sortBy string) *BookInfoList {
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
// func (storeGroup *StoreInRam) GetBookInfoListByMaxDepth(depth int, sortBy string) (*BookInfoList, error) {
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
