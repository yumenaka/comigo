package common

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"
)

// Bookstores 本地总书库，扫描后生成。可以有多个子书库。
type Bookstores struct {
	mapBookstores map[string]*singleBookstore //key为子书库搜索路径（）
	SortBy        string
}

//GenerateBookGroup 生成书籍组
func (bs *Bookstores) GenerateBookGroup() error {
	for _, single := range bs.mapBookstores {
		err := single.initBookGroupMap()
		if err != nil {
			return err
		}
	}
	return nil
}

// 对应某个扫描路径的子书库
type singleBookstore struct {
	StorePath    string           //书库ID，从0开始？用 mapBookstores的大小定义
	BasePath     string           //扫描路径，可能是相对路径。Bookstores的Key。
	Name         string           //书库名，不指定就默认等于Path
	BookMap      map[string]*Book //拥有的书籍,key是BookID
	BookGroupMap map[string]*Book //拥有的书籍组,通过扫描书库生成，key是BookID。需要通过Depth从深到浅生成
}

func (s *singleBookstore) initBookGroupMap() error {
	if len(s.BookMap) == 0 {
		return errors.New("empty Bookstore")
	}
	depthBooksMap := make(map[int][]Book) //key是Depth的临时map
	maxDepth := 0
	for _, b := range s.BookMap {
		depthBooksMap[b.Depth] = append(depthBooksMap[b.Depth], *b)
		if b.Depth > maxDepth {
			maxDepth = b.Depth
		}
	}
	//从深往浅遍历
	//如果有几本书同时有同一个父文件夹，那么应该【新建]一本书(组)，并加入到depth-1层里面
	for depth := maxDepth; depth >= 0; depth-- {
		//最上层（0）与第一层（1）的书籍直接展示，不需要创建书籍组
		if depth < 2 {
			continue
		}
		//用父文件夹做key的parentMap，后面遍历用
		parentTempMap := make(map[string][]Book)
		////遍历depth等于i的所有book
		for _, b := range depthBooksMap[depth] {
			parentTempMap[b.ParentFolder] = append(parentTempMap[b.ParentFolder], b)
		}
		//循环parentMap，把有相同parent的书创建为一个书组
		for parent, list := range parentTempMap {
			//新建一本书籍
			newBook := NewBook(filepath.Dir(list[0].FilePath), time.Now(), 0, s.StorePath, depth-1)
			//类型是书籍组
			newBook.BookType = BookTypeBooksGroup
			//名字应该设置成Name
			if newBook.Name != parent {
				fmt.Printf("newBook.Name!=parent!?")
			}
			//初始化ChildBook
			newBook.ChildBook = make(map[string]*Book)

			//然后把同一parent的书，都加进某个书籍组
			setCover := true
			for _, b := range list {
				//顺便设置一下封面，只设置一次
				if setCover {
					setCover = false
					newBook.Cover = b.Cover //
				}
				newBook.ChildBook[b.BookID] = &b
			}
			//如果书籍组的子书籍数量大于等于1，将书籍组加到上一层
			if len(newBook.ChildBook) >= 1 {
				depthBooksMap[depth-1] = append(depthBooksMap[depth-1], *newBook)
				s.BookGroupMap[newBook.BookID] = newBook
				fmt.Print("生成book_group：")
				fmt.Println(newBook)
			}
		}
	}
	return nil
}

// NewSingleBookstore 创建一个新书库
func (bs *Bookstores) NewSingleBookstore(basePath string) error {
	if _, ok := bs.mapBookstores[basePath]; ok {
		// 已经有这个key了
		return errors.New("add Bookstore Error： The key already exists [" + basePath + "]")
	}
	s := singleBookstore{
		StorePath:    basePath,
		BasePath:     basePath, //暂时与路径同名 TODO：自定义书库名
		Name:         basePath,
		BookMap:      make(map[string]*Book),
		BookGroupMap: make(map[string]*Book),
	}
	bs.mapBookstores[basePath] = &s
	return nil
}

// AddBookToStores 将某一本书，放到basePath做key的某子书库中
func (bs *Bookstores) AddBookToStores(searchPath string, b *Book) error {
	if _, ok := bs.mapBookstores[searchPath]; !ok {
		//创建一个新子书库，并添加一本书
		bs.mapBookstores[searchPath].BookMap[b.BookID] = b
		return errors.New("add Bookstore Error： The key not found [" + searchPath + "]")
	}
	//给已有子书库添加一本书
	bs.mapBookstores[searchPath].BookMap[b.BookID] = b
	return nil
}
