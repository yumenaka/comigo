package model

import (
	"sort"

	"github.com/yumenaka/comigo/tools"
)

// Pages 定义页面结构
type Pages struct {
	Images []MediaFileInfo `json:"images"`
	SortBy string          `json:"sort_by"`
}

// SortImages 根使用 sortBy参数 与图片信息，排序 Pages
func (p *Pages) SortImages(sortBy string) {
	if sortBy == "" {
		sortBy = "default"
	}
	var lessFunc func(i, j int) bool
	switch sortBy {
	case "filename":
		lessFunc = func(i, j int) bool {
			return tools.Compare(p.Images[i].Name, p.Images[j].Name)
		}
	case "filename_reverse":
		lessFunc = func(i, j int) bool {
			return !tools.Compare(p.Images[i].Name, p.Images[j].Name)
		}
	case "filesize":
		lessFunc = func(i, j int) bool {
			return p.Images[i].Size > p.Images[j].Size
		}
	case "filesize_reverse":
		lessFunc = func(i, j int) bool {
			return !(p.Images[i].Size > p.Images[j].Size)
		}
	case "modify_time": // 根据修改时间排序 从新到旧
		lessFunc = func(i, j int) bool {
			return p.Images[i].ModTime.After(p.Images[j].ModTime)
		}
	case "modify_time_reverse": // 根据修改时间排序 从旧到新
		lessFunc = func(i, j int) bool {
			return p.Images[i].ModTime.Before(p.Images[j].ModTime)
		}
	default:
		lessFunc = func(i, j int) bool {
			return tools.Compare(p.Images[i].Name, p.Images[j].Name)
		}
	}
	//  Go 1.8 及以上版本的 sort.Slice 函数。简化排序逻辑，无需再实现 Len、Less 和 Swap 方法。
	sort.Slice(p.Images, lessFunc)
}

// 老写法：三个函数定义好以后，使用sort包排序：sort.Sort(b.Pages)
// // Len 返回页面数量
// func (p Pages) Len() int {
// 	return len(p.Images)
// }
//
// // Less 按照排序方式比较页面
// func (p Pages) Less(i, j int) bool {
// 	switch p.SortBy {
// 	case "filename":
// 		return tools.Compare(p.Images[i].Name, p.Images[j].Name)
// 	case "filesize":
// 		return p.Images[i].Size > p.Images[j].Size
// 	case "modify_time":
// 		return p.Images[i].ModTime.After(p.Images[j].ModTime) // 根据修改时间排序 从新到旧
// 	case "filename_reverse":
// 		return !tools.Compare(p.Images[i].Name, p.Images[j].Name)
// 	case "filesize_reverse":
// 		return p.Images[i].Size < p.Images[j].Size
// 	case "modify_time_reverse":
// 		return p.Images[i].ModTime.Before(p.Images[j].ModTime) // 根据修改时间排序 从旧到新
// 	default:
// 		return tools.Compare(p.Images[i].Name, p.Images[j].Name)
// 	}
// }
//
// // Swap 交换页面
// func (p Pages) Swap(i, j int) {
// 	p.Images[i], p.Images[j] = p.Images[j], p.Images[i]
// }
