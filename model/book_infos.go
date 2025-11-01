package model

import (
	"sort"

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
			// if (*s)[i].Type == TypeDir || (*s)[j].Type == TypeDir {
			//	logger.Info("!!!!" + (*s)[i].Title + "!!!modify_time:" + (*s)[i].Modified.String())
			//	logger.Info("!!!!" + (*s)[j].Title + "!!!modify_time:" + (*s)[j].Modified.String())
			// }
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
