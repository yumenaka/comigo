package types

import (
	"errors"
	"github.com/yumenaka/comi/logger"
	"os"
	"path/filepath"

	"github.com/yumenaka/comi/util"
)

// Folder 本地总书库，扫描后生成。可以有多个子书库。
type Folder struct {
	Path       string
	SortBy     string                //新字段   //排序方式
	BookMap    map[string]*BookInfo  ///新字段   //拥有的书籍,key是BookID
	SubFolders map[string]*subFolder //key为路径
}

// 对应某个扫描路径的子书库
type subFolder struct {
	Path         string               //路径
	SortBy       string               //排序方式
	BookMap      map[string]*BookInfo //拥有的书籍,key是BookID
	BookGroupMap map[string]*BookInfo //拥有的书籍组,通过扫描书库生成，key是BookID。需要通过Depth从深到浅生成
}

// InitFolder 生成书籍组
func (folder *Folder) InitFolder() error {
	for _, single := range folder.SubFolders {
		err := single.AnalyzeFolder()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *subFolder) AnalyzeFolder() error {
	if len(s.BookMap) == 0 {
		return errors.New("empty Bookstore")
	}
	depthBooksMap := make(map[int][]BookInfo) //key是Depth的临时map
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
				logger.Info(err)
				continue
			}
			//书名应该设置成parent
			if newBookGroup.Title != parent {
				newBookGroup.Title = parent
			}
			//初始化ChildBook
			//然后把同一parent的书，都加进某个书籍组
			newBookGroup.ChildBook = make(map[string]*BookInfo)
			for i, bookInList := range sameParentBookList {
				//顺便设置一下封面，只设置一次
				if i == 0 {
					newBookGroup.Cover = bookInList.Cover //
				}
				newBookGroup.ChildBook[bookInList.BookID] = &sameParentBookList[i]
			}
			newBookGroup.ChildBookNum = len(newBookGroup.ChildBook)
			//如果书籍组的子书籍数量大于等于1，且从来没有添加过，将书籍组加到上一层
			if newBookGroup.ChildBookNum == 0 {
				continue
			}
			//检测是否已经生成并添加过
			Added := false
			for _, group := range mapBookGroup {
				if group.FilePath == newBookGroup.FilePath {
					Added = true
				}
			}
			//添加过的不需要添加
			if Added {
				continue
			}
			depthBooksMap[depth-1] = append(depthBooksMap[depth-1], newBookGroup.BookInfo)
			newBookGroup.Author, _ = util.GetAuthor(newBookGroup.Title)
			//将这本书加到子书库的BookGroup表（Images.BookGroupMap）里面去
			s.BookGroupMap[newBookGroup.BookID] = &newBookGroup.BookInfo
			//将这本书加到BookGroup总表（mapBookGroup）里面去
			mapBookGroup[newBookGroup.BookID] = newBookGroup
		}
	}
	return nil
}

// AddSubFolder 创建一个新文件夹
func (folder *Folder) AddSubFolder(basePath string) error {
	if _, ok := folder.SubFolders[basePath]; ok {
		// 已经有这个key了
		return errors.New("add Bookstore Error： The key already exists [" + basePath + "]")
	}
	s := subFolder{
		Path:         basePath,
		BookMap:      make(map[string]*BookInfo),
		BookGroupMap: make(map[string]*BookInfo),
	}
	folder.SubFolders[basePath] = &s
	return nil
}

// AddBookToSubFolder 将某一本书，放到BookMap里面去
func (folder *Folder) AddBookToSubFolder(searchPath string, b *BookInfo) error {
	if _, ok := folder.SubFolders[searchPath]; !ok {
		//创建一个新子书库，并添加一本书
		folder.SubFolders[searchPath].BookMap[b.BookID] = b
		return errors.New("add Bookstore Error： The key not found [" + searchPath + "]")
	}
	//给已有子书库添加一本书
	folder.SubFolders[searchPath].BookMap[b.BookID] = b
	return nil
}
