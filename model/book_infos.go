package model

import (
	"sort"
	"time"

	"github.com/yumenaka/comigo/tools"
)

// BookInfos 表示 BookInfo 的列表，排序用
// 在 Go 中，方法接收器必须是命名类型，这是为了确保类型具有一个唯一的标识和类型身份，从而可以在包级作用域中明确地定义和调用这些方法。
// 使用类型别名 BookInfos []BookInfo 来提供排序方法，这样可以确保 BookInfos 具有明确的类型身份，并且可以在包级作用域中使用和扩展。
type BookInfos []BookInfo

// SortBooks 根据 sortBy 参数对 BookInfos 进行排序
func (s *BookInfos) SortBooks(sortBy string) {
	if sortBy == "" {
		sortBy = "default"
	}

	var lessFunc func(i, j int) bool
	lastReadTimes := map[string]time.Time{}
	if sortBy == "last_read" {
		lastReadTimes = s.collectLastReadTimes()
	}

	switch sortBy {
	case "filename":
		lessFunc = func(i, j int) bool {
			return tools.Compare((*s)[i].Title, (*s)[j].Title)
		}
	case "filename_reverse":
		lessFunc = func(i, j int) bool {
			return !tools.Compare((*s)[i].Title, (*s)[j].Title)
		}
	case "filesize":
		lessFunc = func(i, j int) bool {
			return compareByFileSize((*s)[i], (*s)[j])
		}
	case "filesize_reverse":
		lessFunc = func(i, j int) bool {
			return !compareByFileSize((*s)[i], (*s)[j])
		}
	case "modify_time": // 根据修改时间排序 从新到旧
		lessFunc = func(i, j int) bool {
			return (*s)[i].Modified.After((*s)[j].Modified)
		}
	case "last_read": // 根据最后阅读时间排序 从新到旧。没有阅读记录的书籍按照修改时间排序

		lessFunc = func(i, j int) bool {
			iLastReadTime := lastReadTimes[(*s)[i].BookID]
			jLastReadTime := lastReadTimes[(*s)[j].BookID]
			iHasRead := !iLastReadTime.IsZero()
			jHasRead := !jLastReadTime.IsZero()
			if iHasRead != jHasRead {
				return iHasRead
			}
			if iHasRead && !iLastReadTime.Equal(jLastReadTime) {
				return iLastReadTime.After(jLastReadTime)
			}
			return (*s)[i].Modified.After((*s)[j].Modified)
		}
	case "modify_time_reverse": // 根据修改时间排序 从旧到新
		lessFunc = func(i, j int) bool {
			return (*s)[i].Modified.Before((*s)[j].Modified)
		}
	case "author":
		lessFunc = func(i, j int) bool {
			return tools.Compare((*s)[i].Author, (*s)[j].Author)
		}
	case "author_reverse":
		lessFunc = func(i, j int) bool {
			return !tools.Compare((*s)[i].Author, (*s)[j].Author)
		}
	default:
		lessFunc = func(i, j int) bool {
			return tools.Compare((*s)[i].Title, (*s)[j].Title)
		}
	}
	//  Go 1.8 及以上版本的 sort.Slice 函数。简化排序逻辑，无需再实现 Len、Less 和 Swap 方法。
	sort.Slice(*s, lessFunc)
}

// collectLastReadTimes 预先读取自动书签的更新时间，避免 sort 比较函数反复访问 Store。
func (s *BookInfos) collectLastReadTimes() map[string]time.Time {
	lastReadTimes := make(map[string]time.Time, len(*s))
	if IStore == nil {
		return lastReadTimes
	}
	for _, book := range *s {
		bookMarks, err := IStore.GetBookMarks(book.BookID)
		if err != nil || bookMarks == nil {
			continue
		}
		lastReadTimes[book.BookID] = bookMarks.GetLastReadTime()
	}
	return lastReadTimes
}

// compareByFileSize 按文件大小比较两个 BookInfo
func compareByFileSize(a, b BookInfo) bool {
	// 如果其中一本是书籍组，比较子书籍数量
	if a.Type == TypeBooksGroup || b.Type == TypeBooksGroup {
		if a.Type == TypeBooksGroup && b.Type == TypeBooksGroup {
			return a.ChildBooksNum > b.ChildBooksNum
		}
		return a.Type == TypeBooksGroup
	}
	// 如果其中一本是文件夹，比较页数
	if a.Type == TypeDir || b.Type == TypeDir {
		if a.Type == TypeDir && b.Type == TypeDir {
			return a.PageCount > b.PageCount
		}
		return a.Type == TypeDir
	}
	// 一般情况下，直接比较文件大小
	return a.FileSize > b.FileSize
}

type StoreBookInfo struct {
	StoreUrl     string    `json:"-"` // 书库真实路径只用于服务端分组匹配，普通 JSON 不输出
	DisplayName  string    `json:"display_name"`
	IsRemote     bool      `json:"is_remote"`
	ChildBookNum int       `json:"child_book_num"`
	BookInfos    BookInfos `json:"book_infos"`
}
