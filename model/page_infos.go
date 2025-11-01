package model

import (
	"sort"

	"github.com/yumenaka/comigo/tools"
)

// PageInfos 定义页面列表，排序用
// 在 Go 中，方法接收器必须是命名类型，这是为了确保类型具有一个唯一的标识和类型身份，从而可以在包级作用域中明确地定义和调用这些方法。
// 使用类型别名 PageInfos []PageInfo 来提供排序方法，这样可以确保 PageInfos 具有明确的类型身份，并且可以在包级作用域中使用和扩展。
type PageInfos []PageInfo

// SortImages 根据 sortBy 参数对 PageInfos 进行排序
func (s *PageInfos) SortImages(sortBy string) {
	if sortBy == "" {
		sortBy = "default"
	}
	var lessFunc func(i, j int) bool
	switch sortBy {
	case "filename":
		lessFunc = func(i, j int) bool {
			return tools.Compare((*s)[i].Name, (*s)[j].Name)
		}
	case "filename_reverse":
		lessFunc = func(i, j int) bool {
			return !tools.Compare((*s)[i].Name, (*s)[j].Name)
		}
	case "filesize":
		lessFunc = func(i, j int) bool {
			return (*s)[i].Size > (*s)[j].Size
		}
	case "filesize_reverse":
		lessFunc = func(i, j int) bool {
			return !((*s)[i].Size > (*s)[j].Size)
		}
	case "modify_time": // 根据修改时间排序 从新到旧
		lessFunc = func(i, j int) bool {
			return (*s)[i].ModTime.After((*s)[j].ModTime)
		}
	case "modify_time_reverse": // 根据修改时间排序 从旧到新
		lessFunc = func(i, j int) bool {
			return (*s)[i].ModTime.Before((*s)[j].ModTime)
		}
	default:
		lessFunc = func(i, j int) bool {
			return tools.Compare((*s)[i].Name, (*s)[j].Name)
		}
	}
	//  Go 1.8 及以上版本的 sort.Slice 函数。简化排序逻辑，无需再实现 Len、Less 和 Swap 方法。
	sort.Slice(*s, lessFunc)
}
