package main

import "github.com/yumenaka/comigo/model"

// 支持的图片文件扩展名（统一用小写比较）
var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// DirNode 表示目录树节点，用于 JSON 存储模式
type DirNode struct {
	Name    string                `json:"name"`
	Path    string                `json:"path"`
	SubDirs []DirNode             `json:"sub_dirs"` // 子目录列表
	Files   []model.MediaFileInfo `json:"files"`    // 本目录下的图片文件列表
}

// ListResponse 用于 API 返回目录内容（子目录和文件），支持分页
type ListResponse struct {
	Directories []DirectoryInfo       `json:"directories"`
	Images      []model.MediaFileInfo `json:"images"`
	TotalImages int                   `json:"total_images,omitempty"` // 符合条件的图片总数（用于分页）
	Page        int                   `json:"page,omitempty"`
	PageSize    int                   `json:"page_size,omitempty"`
}

// DirectoryInfo 用于 API 输出的子目录基本信息
type DirectoryInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
