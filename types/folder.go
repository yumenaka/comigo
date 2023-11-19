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
	BookMap    map[string]*Book      ///新字段   //拥有的书籍,key是BookID
	SubFolders map[string]*subFolder //key为路径
}

// 对应某个扫描路径的子书库
type subFolder struct {
	Path         string           //路径
	SortBy       string           //排序方式
	BookMap      map[string]*Book //拥有的书籍,key是BookID
	BookGroupMap map[string]*Book //拥有的书籍组,通过扫描书库生成，key是BookID。需要通过Depth从深到浅生成
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
		//用父文件夹做key的parentMap，后面遍历用
		parentTempMap := make(map[string][]Book)
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
			newBook, err := New(filepath.Dir(sameParentBookList[0].FilePath), modTime, 0, s.Path, depth-1, TypeBooksGroup)
			if err != nil {
				logger.Info(err)
				continue
			}
			//名字应该设置成parent
			if newBook.Name != parent {
				newBook.Name = parent
			}
			//初始化ChildBook
			//然后把同一parent的书，都加进某个书籍组
			newBook.ChildBook = make(map[string]*Book)
			for i, bookInList := range sameParentBookList {
				//顺便设置一下封面，只设置一次
				if i == 0 {
					newBook.Cover = bookInList.Cover //
				}
				newBook.ChildBook[bookInList.BookID] = &sameParentBookList[i]
			}
			newBook.ChildBookNum = len(newBook.ChildBook)
			//如果书籍组的子书籍数量大于等于1，且从来没有添加过，将书籍组加到上一层
			if newBook.ChildBookNum == 0 {
				continue
			}
			//检测是否已经生成并添加过
			Added := false
			for _, checkB := range mapBookGroup {
				if checkB.FilePath == newBook.FilePath {
					Added = true
				}
			}
			//添加过的不需要添加
			if Added {
				continue
			}
			depthBooksMap[depth-1] = append(depthBooksMap[depth-1], *newBook)
			newBook.Author, _ = util.GetAuthor(newBook.Name)
			//将这本书加到子书库的BookGroup表（Images.BookGroupMap）里面去
			s.BookGroupMap[newBook.BookID] = newBook
			//将这本书加到BookGroup总表（mapBookGroup）里面去
			mapBookGroup[newBook.BookID] = newBook
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
		BookMap:      make(map[string]*Book),
		BookGroupMap: make(map[string]*Book),
	}
	folder.SubFolders[basePath] = &s
	return nil
}

// AddBookToSubFolder 将某一本书，放到basePath做key的某子书库中
func (folder *Folder) AddBookToSubFolder(searchPath string, b *Book) error {
	if _, ok := folder.SubFolders[searchPath]; !ok {
		//创建一个新子书库，并添加一本书
		folder.SubFolders[searchPath].BookMap[b.BookID] = b
		return errors.New("add Bookstore Error： The key not found [" + searchPath + "]")
	}
	//给已有子书库添加一本书
	folder.SubFolders[searchPath].BookMap[b.BookID] = b
	return nil
}
