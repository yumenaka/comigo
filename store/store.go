package store

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yumenaka/comigo/assets/locale"
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
	// 记录已创建过的书组路径，避免重复创建（仅在当前子书库内判重）
	addedGroupPath := make(map[string]bool)
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
		// 用父目录“绝对路径”做 key 的 parentMap（仅用 ParentFolder 的名字会发生同名目录误合并）
		parentTempMap := make(map[string][]*model.Book)
		// //遍历depth等于i的所有book
		for _, b := range depthBooksMap[depth] {
			// 计算父目录路径：
			// - 文件型书籍：BookPath 是文件路径，父目录为 Dir(BookPath)
			// - 目录型书籍/书组：BookPath 是目录路径，先去掉结尾分隔符再取父目录
			parentPath := b.BookPath
			if b.Type == model.TypeDir || b.Type == model.TypeBooksGroup {
				parentPath = strings.TrimRight(parentPath, "/\\")
			}
			parentPath = filepath.Dir(parentPath)
			parentTempMap[parentPath] = append(parentTempMap[parentPath], b)
		}
		// 循环parentMap，把有相同parent的书创建为一个书组
		for parentPath, sameParentBookList := range parentTempMap {
			// depth-1 小于 0 说明已经到达子书库根目录以上，不生成书组
			if (depth - 1) < 0 {
				continue
			}
			// 书组路径判重（同一子书库内）
			if addedGroupPath[parentPath] {
				continue
			}
			addedGroupPath[parentPath] = true
			// 新建一本书,类型是书籍组
			// 获取父目录信息（作为书组的时间信息来源）
			pathInfo, err := os.Stat(parentPath)
			if err != nil {
				// 父目录可能暂时不可访问，退化为使用子项目的时间信息
				pathInfo, err = os.Stat(sameParentBookList[0].BookPath)
				if err != nil {
					return err
				}
			}
			// 获取修改时间
			modTime := pathInfo.ModTime()
			tempBook, err := model.NewBook(parentPath, modTime, 0, store.BackendURL, depth-1, model.TypeBooksGroup)
			if err != nil {
				if config.GetCfg().Debug {
					logger.Infof(locale.GetString("log_error_creating_new_book_group"), err)
				}
				continue
			}
			newBookGroup := tempBook
			// 书名设置为目录名（更符合“文件夹/书组”语义）
			parentName := filepath.Base(parentPath)
			if newBookGroup.Title != parentName {
				newBookGroup.Title = parentName
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
			depthBooksMap[depth-1] = append(depthBooksMap[depth-1], newBookGroup)
			// 将这本书加到Store的 BookMap 表里面去
			store.BookMap.Store(newBookGroup.BookID, newBookGroup)
		}
	}
	return nil
}
