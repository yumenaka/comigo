package entity

import (
	"errors"
	"github.com/yumenaka/comi/util/logger"
	"os"
	"path/filepath"
	"sync"
)

// Folder 本地总书库，扫描后生成。可以有多个子书库。
type Folder struct {
	Path       string
	SortBy     string   //新字段   //排序方式
	SubFolders sync.Map //key为路径 存储 *subFolder
}

// 对应某个扫描路径的子书库
type subFolder struct {
	Path         string   //路径
	SortBy       string   //排序方式
	BookMap      sync.Map //拥有的书籍,key是BookID,存储 *BookInfo
	BookGroupMap sync.Map //拥有的书籍组,通过扫描书库生成，key是BookID,存储 *BookInfo。需要通过Depth从深到浅生成
}

// InitFolder 生成书籍组
func (folder *Folder) InitFolder() (e error) {
	//遍历所有子书库
	for _, value := range folder.SubFolders.Range {
		s := value.(*subFolder)
		err := s.AnalyzeFolder()
		if err != nil {
			e = err
		}
	}
	return e
}

func (s *subFolder) AnalyzeFolder() error {
	count := 0
	for _, _ = range s.BookMap.Range {
		count++
	}
	if count == 0 {
		return errors.New("empty Bookstore")
	}
	depthBooksMap := make(map[int][]BookInfo) //key是Depth的临时map
	//定义一个最大深度
	maxDepth := 0
	for _, value := range s.BookMap.Range {
		b := value.(*BookInfo)
		depthBooksMap[b.Depth] = append(depthBooksMap[b.Depth], *b)
		//找到最大深度
		if b.Depth > maxDepth {
			maxDepth = b.Depth
		}
	}

	//从深往浅遍历
	//如果有几本书同时有同一个父文件夹，那么应该【新建]一本书(组)，并加入到depth-1层里面
	for depth := maxDepth; depth >= 0; depth-- {
		//用父文件夹做key的parentMap，后面遍历用
		parentTempMap := make(map[string][]BookInfo)
		////遍历depth等于i的所有book
		for _, b := range depthBooksMap[depth] {
			parentTempMap[b.ParentFolder] = append(parentTempMap[b.ParentFolder], b)
		}
		//循环parentMap，把有相同parent的书创建为一个书组
		for parent, sameParentBookList := range parentTempMap {
			//新建一本书,类型是书籍组
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
			//书名应该设置成parent
			if newBookGroup.Title != parent {
				newBookGroup.Title = parent
			}
			//初始化ChildBook
			//然后把同一parent的书，都加进某个书籍组
			for i, bookInList := range sameParentBookList {
				//顺便设置一下封面，只设置一次
				if i == 0 {
					newBookGroup.SetClover(bookInList.Cover)
				}
				newBookGroup.ChildBook.Store(bookInList.BookID, &sameParentBookList[i])
			}
			newBookGroup.ChildBookNum = len(sameParentBookList)
			//如果书籍组的子书籍数量等于0，那么不需要添加
			if newBookGroup.ChildBookNum == 0 {
				continue
			}
			//检测是否已经生成并添加过
			Added := false
			for _, value := range mapBookGroup.Range {
				group := value.(*BookGroup)
				if group.FilePath == newBookGroup.FilePath {
					Added = true
				}
			}

			//添加过的不需要添加
			if Added {
				continue
			}
			if (depth - 1) < 0 {
				continue
			}
			depthBooksMap[depth-1] = append(depthBooksMap[depth-1], newBookGroup.BookInfo)
			newBookGroup.setAuthor()
			//将这本书加到子书库的BookGroup表（Images.BookGroupMap）里面去
			s.BookGroupMap.Store(newBookGroup.BookID, &newBookGroup.BookInfo)
			//将这本书加到BookGroup总表（mapBookGroup）里面去
			mapBookGroup.Store(newBookGroup.BookID, newBookGroup)
		}
	}
	return nil
}

// AddSubFolder 创建一个新文件夹
func (folder *Folder) AddSubFolder(basePath string) error {
	if _, ok := folder.SubFolders.Load(basePath); ok {
		// 已经有这个key了
		return errors.New("add Bookstore Error： The key already exists [" + basePath + "]")
	}
	s := subFolder{
		Path: basePath,
	}
	folder.SubFolders.Store(basePath, &s)
	return nil
}

// AddBookToSubFolder 将某一本书，放到BookMap里面去
func (folder *Folder) AddBookToSubFolder(searchPath string, b *BookInfo) error {
	if f, ok := folder.SubFolders.Load(searchPath); !ok {
		//创建一个新子书库，并添加一本书
		newStore := subFolder{
			Path: searchPath,
		}
		newStore.BookMap.Store(b.BookID, b)
		folder.SubFolders.Store(searchPath, &newStore)
		return errors.New("add Bookstore Error： The key not found [" + searchPath + "]")
	} else {
		//给已有子书库添加一本书
		temp := f.(*subFolder)
		temp.BookMap.Store(b.BookID, b)
		return nil
	}
}
