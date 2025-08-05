package model

import "github.com/yumenaka/comigo/util"

// Pages 定义页面结构
type Pages struct {
	Images []MediaFileInfo `json:"images"`
	SortBy string          `json:"sort_by"`
}

// Len 返回页面数量
func (s Pages) Len() int {
	return len(s.Images)
}

// Less 按照排序方式比较页面
func (s Pages) Less(i, j int) bool {
	switch s.SortBy {
	case "filename":
		return util.Compare(s.Images[i].Name, s.Images[j].Name)
	case "filesize":
		return s.Images[i].Size > s.Images[j].Size
	case "modify_time":
		return s.Images[i].ModTime.After(s.Images[j].ModTime) // 根据修改时间排序 从新到旧
	case "filename_reverse":
		return !util.Compare(s.Images[i].Name, s.Images[j].Name)
	case "filesize_reverse":
		return s.Images[i].Size < s.Images[j].Size
	case "modify_time_reverse":
		return s.Images[i].ModTime.Before(s.Images[j].ModTime) // 根据修改时间排序 从旧到新
	default:
		return util.Compare(s.Images[i].Name, s.Images[j].Name)
	}
}

// Swap 交换页面
func (s Pages) Swap(i, j int) {
	s.Images[i], s.Images[j] = s.Images[j], s.Images[i]
}
