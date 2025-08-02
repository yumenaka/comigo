package model

import (
	"errors"
	"strconv"
	"sync"

	"github.com/yumenaka/comigo/util/logger"
)

// 使用并发安全的 sync.Map 存储书籍和书组
var (
	mapBooks sync.Map // 实际存在的书 key: string (BookID), value: *Book
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

// DeleteBook 删除一本书
func DeleteBook(bookID string) {
	mapBooks.Delete(bookID)
}

// GetBooksNumber 获取书籍总数，不包括 BookGroup 类型
func GetBooksNumber() int {
	// 用于计数的变量
	var count int
	// 遍历 map 并递增计数器
	for _, b := range mapBooks.Range {
		if b.(*Book).Type == TypeBooksGroup {
			continue // 跳过书组类型
		}
		count++
	}
	return count
}

// GetAllBookList 获取所有书籍列表
func GetAllBookList() []*Book {
	var list []*Book
	for _, value := range mapBooks.Range {
		b := value.(*Book)
		if b.Type == TypeBooksGroup {
			continue // 跳过书组类型
		}
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
	return nil, errors.New("GetBookByID：cannot find book, id=" + id)
}

// GetParentBook 通过子书籍 ID 获取所属书组信息
func GetParentBook(childID string) (*Book, error) {
	for _, value := range mapBooks.Range {
		group := value.(*Book)
		if group.Type != TypeBooksGroup {
			continue // 只处理书组类型
		}
		for _, id := range group.ChildBooksID {
			if id == childID {
				b, err := GetBookByID(group.BookID, "")
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
