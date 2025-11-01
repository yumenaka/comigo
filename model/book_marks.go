package model

import "sort"

type BookMarks []BookMark

// SortBookmarks 根据 sortBy 参数对 BookInfos 进行排序
func (s *BookMarks) SortBookmarks(sortBy string) {
	if sortBy == "" {
		sortBy = "created_at_desc"
	}
	var lessFunc func(i, j int) bool
	switch sortBy {
	case "created_at":
		lessFunc = func(i, j int) bool {
			return (*s)[i].CreatedAt.Before((*s)[j].CreatedAt)
		}
	case "created_at_desc":
		lessFunc = func(i, j int) bool {
			return (*s)[i].CreatedAt.After((*s)[j].CreatedAt)
		}
	case "updated_at":
		lessFunc = func(i, j int) bool {
			return (*s)[i].UpdatedAt.Before((*s)[j].UpdatedAt)
		}
	case "updated_at_desc":
		lessFunc = func(i, j int) bool {
			return (*s)[i].UpdatedAt.After((*s)[j].UpdatedAt)
		}
	case "page_index":
		lessFunc = func(i, j int) bool {
			return (*s)[i].PageIndex < (*s)[j].PageIndex
		}
	case "page_index_desc":
		lessFunc = func(i, j int) bool {
			return (*s)[i].PageIndex > (*s)[j].PageIndex
		}
	default:
		lessFunc = func(i, j int) bool {
			return (*s)[i].CreatedAt.Before((*s)[j].CreatedAt)
		}
	}
	// 对解引用后的切片进行排序
	sort.Slice(*s, lessFunc)
}
