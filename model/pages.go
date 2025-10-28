package model

import (
	"sort"

	"github.com/yumenaka/comigo/tools"
)

// AllPage 定义页面结构
type AllPage struct {
	PageInfos []PageInfo `json:"page_infos"`
	SortBy    string     `json:"sort_by"`
}

// SortImages 根使用 sortBy参数 与图片信息，排序 AllPage
func (allPage *AllPage) SortImages(sortBy string) {
	if sortBy == "" {
		sortBy = "default"
	}
	var lessFunc func(i, j int) bool
	switch sortBy {
	case "filename":
		lessFunc = func(i, j int) bool {
			return tools.Compare(allPage.PageInfos[i].Name, allPage.PageInfos[j].Name)
		}
	case "filename_reverse":
		lessFunc = func(i, j int) bool {
			return !tools.Compare(allPage.PageInfos[i].Name, allPage.PageInfos[j].Name)
		}
	case "filesize":
		lessFunc = func(i, j int) bool {
			return allPage.PageInfos[i].Size > allPage.PageInfos[j].Size
		}
	case "filesize_reverse":
		lessFunc = func(i, j int) bool {
			return !(allPage.PageInfos[i].Size > allPage.PageInfos[j].Size)
		}
	case "modify_time": // 根据修改时间排序 从新到旧
		lessFunc = func(i, j int) bool {
			return allPage.PageInfos[i].ModTime.After(allPage.PageInfos[j].ModTime)
		}
	case "modify_time_reverse": // 根据修改时间排序 从旧到新
		lessFunc = func(i, j int) bool {
			return allPage.PageInfos[i].ModTime.Before(allPage.PageInfos[j].ModTime)
		}
	default:
		lessFunc = func(i, j int) bool {
			return tools.Compare(allPage.PageInfos[i].Name, allPage.PageInfos[j].Name)
		}
	}
	//  Go 1.8 及以上版本的 sort.Slice 函数。简化排序逻辑，无需再实现 Len、Less 和 Swap 方法。
	sort.Slice(allPage.PageInfos, lessFunc)
}
