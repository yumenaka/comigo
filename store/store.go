package store

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// StoreInfo 书库基本信息
type StoreInfo struct {
	BackendURL  string // 本地书库文件夹路径，或远程书库URL
	Name        string
	Description string
	Backend
}

// Store 对应某个扫描路径的子书库，目前只支持本地书库
type Store struct {
	StoreInfo
	BookMap      sync.Map // 拥有的书籍,key是BookID,存储 *Book 与 *BooksGroup
	BookMarksMap sync.Map // 书签,key是BookID,存储 BookMarks
}

// GenerateBookGroup 分析书库中已有书籍的路径，生成书籍组信息
func (store *Store) GenerateBookGroup() error {
	// 遍历 BookMap ，清理所有 BooksGroup 类型的书籍
	for _, value := range store.BookMap.Range {
		b := value.(*model.Book)
		if b.Type == model.TypeBooksGroup {
			store.BookMap.Delete(b.BookID)
		}
	}
	// 然后再重新生成 BooksGroup
	depthBooksMap := make(map[int][]*model.Book) // key是Depth的临时map
	// 计算最大深度
	maxDepth := 0
	for _, value := range store.BookMap.Range {
		b := value.(*model.Book)
		if b.Depth > maxDepth {
			maxDepth = b.Depth
		}
		depthBooksMap[b.Depth] = append(depthBooksMap[b.Depth], b)
	}
	// 从深往浅遍历
	// 如果有几本书同时有同一个父文件夹，那么应该【新建】一本书(组)，并加入到depth-1层里面
	for depth := maxDepth; depth >= 0; depth-- {
		// 用父文件夹做key的parentMap，后面遍历用
		parentTempMap := make(map[string][]*model.Book)
		// //遍历depth等于i的所有book
		for _, b := range depthBooksMap[depth] {
			parentTempMap[b.ParentFolder] = append(parentTempMap[b.ParentFolder], b)
		}
		// 循环parentMap，把有相同parent的书创建为一个书组
		for parent, sameParentBookList := range parentTempMap {
			// 新建一本书,类型是书籍组
			// 获取文件夹信息
			pathInfo, err := os.Stat(sameParentBookList[0].BookPath)
			if err != nil {
				return err
			}
			// 获取修改时间
			modTime := pathInfo.ModTime()
			tempBook, err := model.NewBook(filepath.Dir(sameParentBookList[0].BookPath), modTime, 0, store.BackendURL, depth-1, model.TypeBooksGroup)
			if err != nil {
				if config.GetCfg().Debug {
					logger.Infof("Error creating new book group: %s", err)
				}
				continue
			}
			newBookGroup := tempBook
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
			allBooks, err := RamStore.ListBooks()
			if err != nil {
				logger.Infof("Error listing books: %s", err)
			}
			for _, bookGroup := range allBooks {
				if bookGroup.Type != model.TypeBooksGroup {
					continue // 只关心书籍组类型
				}
				if bookGroup.BookPath == newBookGroup.BookPath {
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
			// 将这本书加到Store的 BookMap 表里面去
			store.BookMap.Store(newBookGroup.BookID, newBookGroup)
		}
	}
	return nil
}
