package model

import (
	"errors"
	"sort"

	"github.com/yumenaka/comigo/util"
)

// BookInfoList 表示 BookInfo 的列表，排序用
type BookInfoList struct {
	BookInfos []BookInfo
}

// GetAllBookInfoList 获取所有 BookInfo，并根据 sortBy 参数进行排序
func GetAllBookInfoList(sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	// 添加所有真实的书籍
	for _, value := range mapBooks.Range {
		b := value.(*Book)
		info := b.GetBookInfo()
		infoList.BookInfos = append(infoList.BookInfos, *info)
	}

	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error: cannot find bookshelf in GetAllBookInfoList")
}

// SortBooks 根据 sortBy 参数对 BookInfos 进行排序
func (s *BookInfoList) SortBooks(sortBy string) {
	if sortBy == "" {
		sortBy = "default"
	}

	var lessFunc func(i, j int) bool

	switch sortBy {
	case "filename":
		lessFunc = func(i, j int) bool {
			return util.Compare(s.BookInfos[i].Title, s.BookInfos[j].Title)
		}
	case "filename_reverse":
		lessFunc = func(i, j int) bool {
			return !util.Compare(s.BookInfos[i].Title, s.BookInfos[j].Title)
		}
	case "filesize":
		lessFunc = func(i, j int) bool {
			return compareByFileSize(s.BookInfos[i], s.BookInfos[j])
		}
	case "filesize_reverse":
		lessFunc = func(i, j int) bool {
			return !compareByFileSize(s.BookInfos[i], s.BookInfos[j])
		}
	case "modify_time": // 根据修改时间排序 从新到旧
		lessFunc = func(i, j int) bool {
			// if s.BookInfos[i].Type == TypeDir || s.BookInfos[j].Type == TypeDir {
			//	logger.Info("!!!!" + s.BookInfos[i].Title + "!!!modify_time:" + s.BookInfos[i].Modified.String())
			//	logger.Info("!!!!" + s.BookInfos[j].Title + "!!!modify_time:" + s.BookInfos[j].Modified.String())
			// }
			return s.BookInfos[i].Modified.After(s.BookInfos[j].Modified)
		}
	case "modify_time_reverse": // 根据修改时间排序 从旧到新
		lessFunc = func(i, j int) bool {
			return s.BookInfos[i].Modified.Before(s.BookInfos[j].Modified)
		}
	case "author":
		lessFunc = func(i, j int) bool {
			return util.Compare(s.BookInfos[i].Author, s.BookInfos[j].Author)
		}
	case "author_reverse":
		lessFunc = func(i, j int) bool {
			return !util.Compare(s.BookInfos[i].Author, s.BookInfos[j].Author)
		}
	default:
		lessFunc = func(i, j int) bool {
			return util.Compare(s.BookInfos[i].Title, s.BookInfos[j].Title)
		}
	}
	//  Go 1.8 及以上版本的 sort.Slice 函数。简化排序逻辑，无需再实现 Len、Less 和 Swap 方法。
	sort.Slice(s.BookInfos, lessFunc)
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
