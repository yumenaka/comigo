package model

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/yumenaka/comigo/util/logger"
)

// 对应某个扫描路径的子书库
type ChildStore struct {
	Path       string   // 扫描路径
	SubBookMap sync.Map // 拥有的书籍,key是BookID,存储 *Book
}

// GenerateAllBookGroup 分析并生成书籍组
func (storeGroup *StoreGroup) GenerateAllBookGroup() (e error) {
	// 遍历所有子书库
	for _, value := range storeGroup.ChildStores.Range {
		s := value.(*ChildStore)
		err := s.GenerateBookGroup()
		if err != nil {
			e = err
		}
	}
	return e
}

func (subStore *ChildStore) GenerateBookGroup() error {
	// 计没有添加过任何书籍，这不需要生成书组信息
	count := 0
	for range subStore.SubBookMap.Range {
		count++
	}
	if count == 0 {
		return errors.New("empty Bookstore,skipping analysis")
	}
	depthBooksMap := make(map[int][]*Book) // key是Depth的临时map
	// 计算最大深度
	maxDepth := 0
	for _, value := range subStore.SubBookMap.Range {
		b := value.(*Book)
		depthBooksMap[b.Depth] = append(depthBooksMap[b.Depth], b)
		if b.Depth > maxDepth {
			maxDepth = b.Depth
		}
	}

	// 从深往浅遍历
	// 如果有几本书同时有同一个父文件夹，那么应该【新建]一本书(组)，并加入到depth-1层里面
	for depth := maxDepth; depth >= 0; depth-- {
		// 用父文件夹做key的parentMap，后面遍历用
		parentTempMap := make(map[string][]*Book)
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
			tempBookInfo, err := NewBookInfo(filepath.Dir(sameParentBookList[0].FilePath), modTime, 0, subStore.Path, depth-1, TypeBooksGroup)
			if err != nil {
				logger.Infof("%s", err)
				continue
			}
			newBookGroup := &Book{
				BookInfo: *tempBookInfo,
			}
			// 书名应该设置成parent
			if newBookGroup.Title != parent {
				newBookGroup.Title = parent
			}
			// 初始化ChildBook
			// 然后把同一parent的书，都加进某个书籍组
			for _, bookInList := range sameParentBookList {
				newBookGroup.ChildBooksID = append(newBookGroup.ChildBooksID, bookInList.BookID)
			}
			newBookGroup.ChildBooksNum = len(sameParentBookList)
			// 如果书籍组的子书籍数量等于0，那么不需要添加
			if newBookGroup.ChildBooksNum == 0 {
				continue
			}
			// 检测是否已经生成并添加过
			Added := false
			for _, bookGroup := range MainStores.List() {
				if bookGroup.Type == TypeBooksGroup {
					continue // 只处理书籍组类型
				}
				if bookGroup.FilePath == newBookGroup.FilePath {
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
			depthBooksMap[depth-1] = append(depthBooksMap[depth-1], newBookGroup)
			newBookGroup.SetAuthor()
			// 将这本书加到子书库的SubBookMap表里面去
			subStore.SubBookMap.Store(newBookGroup.BookID, newBookGroup)
		}
	}
	return nil
}

// AddSubStore 创建一个新文件夹
func (storeGroup *StoreGroup) AddSubStore(basePath string) error {
	if _, ok := storeGroup.ChildStores.Load(basePath); ok {
		// 已经有这个key了
		return errors.New("add Bookstore Error： The key already exists [" + basePath + "]")
	}
	s := ChildStore{
		Path: basePath,
	}
	storeGroup.ChildStores.Store(basePath, &s)
	return nil
}

// AddBookToSubStore 将某一本书，放到BookMap里面去
func (storeGroup *StoreGroup) AddBookToSubStore(searchPath string, b *Book) error {
	if f, ok := storeGroup.ChildStores.Load(searchPath); !ok {
		// 创建一个新子书库，并添加一本书
		newSubStore := ChildStore{
			Path: searchPath,
		}
		newSubStore.SubBookMap.Store(b.BookID, b)
		storeGroup.ChildStores.Store(searchPath, &newSubStore)
		return errors.New("add Bookstore Error： The key not found [" + searchPath + "]")
	} else {
		// 给已有子书库添加一本书
		temp := f.(*ChildStore)
		temp.SubBookMap.Store(b.BookID, b)
		return nil
	}
}
