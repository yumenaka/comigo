package model

import (
	"errors"
	"strconv"
	"sync"

	"github.com/yumenaka/comigo/util/logger"
)

// 使用并发安全的 sync.Map 存储书籍和书组
var (
	mapBooks     sync.Map // 实际存在的书 key: string (BookID), value: *Book
	mapBookGroup sync.Map // 虚拟书组    key: string (BookID), value: *BookGroup
	// MainStore 带有层级关系的总书组，用于前端展示
	MainStore = Store{}
)

// AddBook 添加一本书
func AddBook(b *Book, basePath string, minPageNum int) error {
	if b.BookID == "" {
		return errors.New("add book error: empty BookID")
	}
	if b.GetPageCount() < minPageNum {
		return errors.New("add book error: minPageNum = " + strconv.Itoa(b.GetPageCount()))
	}

	if _, ok := MainStore.SubStores.Load(basePath); !ok {
		if err := MainStore.AddSubStore(basePath); err != nil {
			logger.Infof("Error adding subfolder: %s", err)
		}
	}
	mapBooks.Store(b.BookID, b)
	return MainStore.AddBookToSubStore(basePath, &b.BookInfo)
}

// AddBooks 添加一组书
func AddBooks(list []*Book, basePath string, minPageNum int) error {
	for _, b := range list {
		if b.GetPageCount() < minPageNum {
			continue
		}
		if err := AddBook(b, basePath, minPageNum); err != nil {
			return err
		}
	}
	return nil
}

// DeleteBookByID 删除一本书
func DeleteBookByID(bookID string) {
	mapBooks.Delete(bookID)
}

// GetBooksNumber 获取书籍总数，不包括 BookGroup
func GetBooksNumber() int {
	// 用于计数的变量
	var count int
	// 遍历 map 并递增计数器
	for range mapBooks.Range {
		count++
	}
	return count
}

// GetAllBookList 获取所有书籍列表
func GetAllBookList() []*Book {
	var list []*Book
	for _, value := range mapBooks.Range {
		b := value.(*Book)
		list = append(list, b)
	}
	return list
}

// GetArchiveBooks 获取所有压缩包格式的书籍
func GetArchiveBooks() []*Book {
	var list []*Book
	for _, value := range mapBooks.Range {
		b := value.(*Book)
		if b.Type == TypeZip || b.Type == TypeRar || b.Type == TypeCbz || b.Type == TypeCbr || b.Type == TypeTar || b.Type == TypeEpub {
			list = append(list, b)
		}
	}
	return list
}

// GetBookByID 根据 BookID 获取书籍
func GetBookByID(id string, sortBy string) (*Book, error) {
	if value, ok := mapBooks.Load(id); ok {
		b := value.(*Book)
		b.SortPages(sortBy)
		return b, nil
	}
	if value, ok := mapBookGroup.Load(id); ok {
		group := value.(*BookGroup)
		temp := Book{
			BookInfo: group.BookInfo,
		}
		return &temp, nil
	}
	return nil, errors.New("GetBookByID：cannot find book, id=" + id)
}

// GetBookGroupInfoByChildBookID 通过子书籍 ID 获取所属书组信息
func GetBookGroupInfoByChildBookID(id string) (*BookGroup, error) {
	for _, value := range mapBookGroup.Range {
		group := value.(*BookGroup)
		for _, v := range group.ChildBook.Range {
			b := v.(*BookInfo)
			if b.BookID == id {
				return group, nil
			}
		}
	}
	return nil, errors.New("cannot find group, id=" + id)
}

// GetBookByAuthor 获取同一作者的书籍
func GetBookByAuthor(author string, sortBy string) ([]*Book, error) {
	var bookList []*Book
	for _, value := range mapBooks.Range {
		b := value.(*Book)
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
