package entity

import (
	"errors"
	"github.com/yumenaka/comigo/util"
	"sort"
	"strconv"
)

// BookInfoList Slice
type BookInfoList struct {
	SortBy    string
	BookInfos []BookInfo
}

func GetAllBookInfoList(sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	//首先加上所有真实的书籍
	for _, value := range mapBooks.Range {
		b := value.(*Book)
		info := NewBaseInfo(b)
		infoList.BookInfos = append(infoList.BookInfos, *info)
	}

	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error:can not found bookshelf. GetAllBookInfoList")
}

func (s BookInfoList) Len() int {
	return len(s.BookInfos)
}

// Less 按时间或URL，将图片排序
func (s BookInfoList) Less(i, j int) (less bool) {
	//如何定义 Images[i] < Images[j]
	//根据文件名
	switch s.SortBy {
	case "filename":
		return util.Compare(s.BookInfos[i].Title, s.BookInfos[j].Title) //(比较自然语言字符串)
	case "filesize":
		//两本之中有一本是书籍组。同样是书籍组，比较子书籍数。
		if s.BookInfos[i].Type == TypeBooksGroup || s.BookInfos[j].Type == TypeBooksGroup {
			if s.BookInfos[i].Type == TypeBooksGroup && s.BookInfos[j].Type == TypeBooksGroup {
				return s.BookInfos[i].ChildBookNum > s.BookInfos[j].ChildBookNum
			}
			if s.BookInfos[i].Type != TypeBooksGroup || s.BookInfos[j].Type != TypeBooksGroup {
				return s.BookInfos[i].Type == TypeBooksGroup
			}
		}
		//两本之中有一本是文件夹。同样是文件夹，比较页数。
		if s.BookInfos[i].Type == TypeDir || s.BookInfos[j].Type == TypeDir {
			if s.BookInfos[i].Type == TypeDir && s.BookInfos[j].Type == TypeDir {
				return s.BookInfos[i].PageCount > s.BookInfos[j].PageCount
			}
			if s.BookInfos[i].Type != TypeDir || s.BookInfos[j].Type != TypeDir {
				return s.BookInfos[i].Type == TypeDir
			}
		}
		//一般情况，比较文件大小
		return !util.Compare(strconv.Itoa(int(s.BookInfos[i].FileSize)), strconv.Itoa(int(s.BookInfos[j].FileSize)))
	case "modify_time":
		return !util.Compare(s.BookInfos[i].Modified.String(), s.BookInfos[j].Modified.String())
	case "author":
		return util.Compare(s.BookInfos[i].Author, s.BookInfos[j].Author)
	//如何定义 Images[i] < Images[j] 反向
	case "filename_reverse":
		return !util.Compare(s.BookInfos[i].Title, s.BookInfos[j].Title) //(比较自然语言字符串)
	case "filesize_reverse":
		//两本之中有一本是书籍组。同样是书籍组，比较子书籍数。
		if s.BookInfos[i].Type == TypeBooksGroup || s.BookInfos[j].Type == TypeBooksGroup {
			if s.BookInfos[i].Type == TypeBooksGroup && s.BookInfos[j].Type == TypeBooksGroup {
				return !(s.BookInfos[i].ChildBookNum > s.BookInfos[j].ChildBookNum)
			}
			if s.BookInfos[i].Type != TypeBooksGroup || s.BookInfos[j].Type != TypeBooksGroup {
				return !(s.BookInfos[i].Type == TypeBooksGroup)
			}
		}
		//两本之中有一本是文件夹。同样是文件夹，比较页数。
		if s.BookInfos[i].Type == TypeDir || s.BookInfos[j].Type == TypeDir {
			if s.BookInfos[i].Type == TypeDir && s.BookInfos[j].Type == TypeDir {
				return !(s.BookInfos[i].PageCount > s.BookInfos[j].PageCount)
			}
			if s.BookInfos[i].Type != TypeDir || s.BookInfos[j].Type != TypeDir {
				return !(s.BookInfos[i].Type == TypeDir)
			}
		}
		//一般情况，比较文件大小
		return util.Compare(strconv.Itoa(int(s.BookInfos[i].FileSize)), strconv.Itoa(int(s.BookInfos[j].FileSize)))
	case "modify_time_reverse":
		return util.Compare(s.BookInfos[i].Modified.String(), s.BookInfos[j].Modified.String())
	case "author_reverse":
		return !util.Compare(s.BookInfos[i].Author, s.BookInfos[j].Author)
	default:
		return util.Compare(s.BookInfos[i].Title, s.BookInfos[j].Title)
	}
}

func (s BookInfoList) Swap(i, j int) {
	s.BookInfos[i], s.BookInfos[j] = s.BookInfos[j], s.BookInfos[i]
}

// SortBooks 上面三个函数定义好了，终于可以使用sort包排序了
func (s *BookInfoList) SortBooks(by string) {
	if by != "" {
		s.SortBy = by
		sort.Sort(s)
	}
}
