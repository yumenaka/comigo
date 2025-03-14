package main

import "github.com/yumenaka/comigo/model"

// 辅助函数：检查字符串是否在 slice 中
func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// 辅助函数：检查路径是否在文件列表中（根据 MediaFileInfo.Path）
func pathInList(path string, files []model.MediaFileInfo) bool {
	for _, f := range files {
		if f.Path == path {
			return true
		}
	}
	return false
}

// 辅助函数：在 DirNode 树中找到指定路径的节点
func findDirNode(root DirNode, targetPath string) *DirNode {
	if root.Path == targetPath {
		return &root
	}
	for i := range root.SubDirs {
		if root.SubDirs[i].Path == targetPath {
			return &root.SubDirs[i]
		}
		// 递归在子目录中查找
		if result := findDirNode(root.SubDirs[i], targetPath); result != nil {
			return result
		}
	}
	return nil
}
