package model

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/yumenaka/comigo/util/logger"
)

// Store 本地总书库，扫描后生成。可以有多个子书库。
type Store struct {
	SubStores sync.Map // key为路径 存储 *subStore
}

// 对应某个扫描路径的子书库
type subStore struct {
	Path         string   // 扫描路径
	SortBy       string   // 排序方式
	BookMap      sync.Map // 拥有的书籍,key是BookID,存储 *BookInfo
	BookGroupMap sync.Map // 拥有的书籍组,通过扫描书库生成，key是BookID,存储 *BookInfo。需要通过Depth从深到浅生成
}

// AnalyzeStore 分析并生成书籍组
func (folder *Store) AnalyzeStore() (e error) {
	// 遍历所有子书库
	for _, value := range folder.SubStores.Range {
		s := value.(*subStore)
		err := s.AnalyzeFolder()
		if err != nil {
			e = err
		}
	}
	return e
}

func (s *subStore) AnalyzeFolder() error {
	count := 0
	for range s.BookMap.Range {
		count++
	}
	if count == 0 {
		return errors.New("empty Bookstore")
	}
	depthBooksMap := make(map[int][]BookInfo) // key是Depth的临时map
	// 定义一个最大深度
	maxDepth := 0
	for _, value := range s.BookMap.Range {
		b := value.(*BookInfo)
		depthBooksMap[b.Depth] = append(depthBooksMap[b.Depth], *b)
		// 找到最大深度
		if b.Depth > maxDepth {
			maxDepth = b.Depth
		}
	}

	// 从深往浅遍历
	// 如果有几本书同时有同一个父文件夹，那么应该【新建]一本书(组)，并加入到depth-1层里面
	for depth := maxDepth; depth >= 0; depth-- {
		// 用父文件夹做key的parentMap，后面遍历用
		parentTempMap := make(map[string][]BookInfo)
		// //遍历depth等于i的所有book
		for _, b := range depthBooksMap[depth] {
			parentTempMap[b.ParentFolder] = append(parentTempMap[b.ParentFolder], b)
		}
		// 循环parentMap，把有相同parent的书创建为一个书组
		for parent, sameParentBookList := range parentTempMap {
			// 新建一本书,类型是书籍组
			// 获取文件夹信息
			pathInfo, err := os.Stat(sameParentBookList[0].FilePath)
			if err != nil {
				return err
			}
			// 获取修改时间
			modTime := pathInfo.ModTime()
			newBookGroup, err := NewBookGroup(filepath.Dir(sameParentBookList[0].FilePath), modTime, 0, s.Path, depth-1, TypeBooksGroup)
			if err != nil {
				logger.Infof("%s", err)
				continue
			}
			// 书名应该设置成parent
			if newBookGroup.Title != parent {
				newBookGroup.Title = parent
			}
			// 初始化ChildBook
			// 然后把同一parent的书，都加进某个书籍组
			for i, bookInList := range sameParentBookList {
				newBookGroup.ChildBook.Store(bookInList.BookID, &sameParentBookList[i])
			}
			newBookGroup.ChildBookNum = len(sameParentBookList)
			// 如果书籍组的子书籍数量等于0，那么不需要添加
			if newBookGroup.ChildBookNum == 0 {
				continue
			}
			// 检测是否已经生成并添加过
			Added := false
			for _, value := range mapBookGroup.Range {
				group := value.(*BookGroup)
				if group.FilePath == newBookGroup.FilePath {
					Added = true
				}
			}

			// 添加过的不需要添加
			if Added {
				continue
			}
			if (depth - 1) < 0 {
				continue
			}
			depthBooksMap[depth-1] = append(depthBooksMap[depth-1], newBookGroup.BookInfo)
			newBookGroup.SetAuthor()
			// 将这本书加到子书库的BookGroup表（Images.BookGroupMap）里面去
			s.BookGroupMap.Store(newBookGroup.BookID, &newBookGroup.BookInfo)
			// 将这本书加到BookGroup总表（mapBookGroup）里面去
			mapBookGroup.Store(newBookGroup.BookID, newBookGroup)
		}
	}
	return nil
}

// AddSubStore 创建一个新文件夹
func (folder *Store) AddSubStore(basePath string) error {
	if _, ok := folder.SubStores.Load(basePath); ok {
		// 已经有这个key了
		return errors.New("add Bookstore Error： The key already exists [" + basePath + "]")
	}
	s := subStore{
		Path: basePath,
	}
	folder.SubStores.Store(basePath, &s)
	return nil
}

// AddBookToSubStore 将某一本书，放到BookMap里面去
func (folder *Store) AddBookToSubStore(searchPath string, b *BookInfo) error {
	if f, ok := folder.SubStores.Load(searchPath); !ok {
		// 创建一个新子书库，并添加一本书
		newStore := subStore{
			Path: searchPath,
		}
		newStore.BookMap.Store(b.BookID, b)
		folder.SubStores.Store(searchPath, &newStore)
		return errors.New("add Bookstore Error： The key not found [" + searchPath + "]")
	} else {
		// 给已有子书库添加一本书
		temp := f.(*subStore)
		temp.BookMap.Store(b.BookID, b)
		return nil
	}
}
