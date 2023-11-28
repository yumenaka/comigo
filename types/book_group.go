package types

import (
	"errors"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/util"
	"path/filepath"
	"strings"
	"time"
)

type BookGroup struct {
	BookInfo
	ChildBook map[string]*BookInfo `json:"child_book" ` //key：BookID
}

// New  初始化Book，设置文件路径、书名、BookID等等
func NewBookGroup(filePath string, modified time.Time, fileSize int64, storePath string, depth int, bookType SupportFileType) (*BookGroup, error) {
	//初始化书籍
	var group = BookGroup{
		BookInfo: BookInfo{
			Author:        "",
			Modified:      modified,
			FileSize:      fileSize,
			InitComplete:  false,
			Depth:         depth,
			BookStorePath: storePath,
			Type:          bookType},
	}
	//设置属性：
	//FilePath，转换为绝对路径
	group.setFilePath(filePath)
	group.setTitle(filePath)
	group.Author, _ = util.GetAuthor(group.Title)
	//设置属性：父文件夹
	group.setParentFolder(filePath)
	group.setBookID()
	return &group, nil
}

// 初始化Book时，设置FilePath
func (b *BookInfo) setFilePath(path string) {
	fileAbaPath, err := filepath.Abs(path)
	if err != nil {
		//因为权限问题，无法取得绝对路径的情况下，用相对路径
		logger.Info(err, fileAbaPath)
		b.FilePath = path
	} else {
		b.FilePath = fileAbaPath
	}
}

func (b *BookInfo) setParentFolder(filePath string) {
	//取得文件所在文件夹的路径
	//如果类型是文件夹，同时最后一个字符是路径分隔符的话，就多取一次dir，移除多余的Unix路径分隔符或windows分隔符
	if b.Type == TypeDir {
		if filePath[len(filePath)-1] == '/' || filePath[len(filePath)-1] == '\\' {
			filePath = filepath.Dir(filePath)
		}
	}
	folder := filepath.Dir(filePath)
	post := strings.LastIndex(folder, "/") //Unix路径分隔符
	if post == -1 {
		post = strings.LastIndex(folder, "\\") //windows分隔符
	}
	if post != -1 {
		//FilePath = string([]rune(FilePath)[post:]) //为了防止中文字符被错误截断，先转换成rune，再转回来
		p := folder[post:]
		p = strings.ReplaceAll(p, "\\", "")
		p = strings.ReplaceAll(p, "/", "")
		b.ParentFolder = p
	}
}

func (b *BookInfo) setTitle(filePath string) {
	b.Title = filePath
	//设置属性：书籍名，取文件后缀(可能为 .zip .rar .pdf .mp4等等)。
	if b.Type != TypeBooksGroup { //不是书籍组(book_group)。
		post := strings.LastIndex(filePath, "/") //Unix路径分隔符
		if post == -1 {
			post = strings.LastIndex(filePath, "\\") //windows分隔符
		}
		if post != -1 {
			//FilePath = string([]rune(FilePath)[post:]) //为了防止中文字符被错误截断，先转换成rune，再转回来
			filename := filePath[post:]
			filename = strings.ReplaceAll(filename, "\\", "")
			filename = strings.ReplaceAll(filename, "/", "")
			b.Title = filename
		}
	}
}

func GetBookInfoListByDepth(depth int, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	//首先加上所有真实的书籍
	for _, b := range mapBooks {
		if b.Depth == depth {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	//接下来还要加上扫描生成出来的书籍组
	for _, bs := range MainFolder.SubFolders {
		for _, group := range bs.BookGroupMap {
			if group.Depth == depth {
				infoList.BookInfos = append(infoList.BookInfos, *group)
			}
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error:can not found bookshelf. GetBookInfoListByDepth")
}

func GetBookInfoListByMaxDepth(depth int, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	//首先加上所有真实的书籍
	for _, b := range mapBooks {
		if b.Depth <= depth {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	//扫描生成的书籍组
	for _, bs := range MainFolder.SubFolders {
		for _, group := range bs.BookGroupMap {
			if group.Depth <= depth {
				infoList.BookInfos = append(infoList.BookInfos, *group)
			}
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error:can not found bookshelf. GetBookInfoListByMaxDepth")
}

func GetBookInfoListByID(BookID string, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	group := mapBookGroup[BookID]
	if group != nil {
		//首先加上所有真实的书籍
		for _, g := range group.ChildBook {
			infoList.BookInfos = append(infoList.BookInfos, *g)
		}
		if len(infoList.BookInfos) > 0 {
			infoList.SortBooks(sortBy)
			return &infoList, nil
		}
	}
	return nil, errors.New("can not found bookshelf")
}

func GetBookInfoListByParentFolder(parentFolder string, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	for _, b := range mapBooks {
		if b.ParentFolder == parentFolder {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("can not found book,parentFolder=" + parentFolder)
}
