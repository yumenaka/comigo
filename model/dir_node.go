package model

// DirNode 表示目录树节点，用于 JSON 存储模式
type DirNode struct {
	Name    string          `json:"name"`
	Path    string          `json:"path"`
	SubDirs []DirNode       `json:"sub_dirs"` // 子目录列表
	Files   []MediaFileInfo `json:"files"`    // 本目录下的图片文件列表
}
