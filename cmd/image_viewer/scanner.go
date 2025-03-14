package main

import (
	"fmt"
	"github.com/yumenaka/comigo/model"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// 全局变量：标记是否正在扫描，避免并发扫描
var (
	scanning  bool       // 标记是否正在扫描，避免并发扫描
	scanMutex sync.Mutex // 保护 scanning 标志的锁
)

// 扫描目录的核心函数：递归遍历目录，忽略指定名称的文件夹，收集图片文件信息
func scanDirectory(currentPath string, depth, maxDepth int, ignoreDirs map[string]bool) (DirNode, []string, []model.MediaFileInfo, error) {
	node := DirNode{
		Name: filepath.Base(currentPath),
		Path: currentPath,
	}
	var foundDirs []string
	var foundFiles []model.MediaFileInfo

	// 如果超过最大深度限制，直接返回空节点
	if maxDepth >= 0 && depth > maxDepth {
		return node, foundDirs, foundFiles, nil
	}

	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return node, foundDirs, foundFiles, err
	}

	// 当前目录计入 foundDirs（用于记录树状结构）
	foundDirs = append(foundDirs, currentPath)

	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(currentPath, name)
		if entry.IsDir() {
			// 检查是否在忽略列表
			if ignoreDirs[strings.ToLower(name)] {
				continue
			}
			// 递归扫描子目录
			subNode, subDirs, subFiles, subErr := scanDirectory(fullPath, depth+1, maxDepth, ignoreDirs)
			if subErr != nil {
				// 忽略单个子目录出错，继续扫描其他目录
				fmt.Println("扫描子目录出错:", subErr)
				continue
			}
			node.SubDirs = append(node.SubDirs, subNode)
			// 合并子目录扫描结果
			foundDirs = append(foundDirs, subDirs...)
			foundFiles = append(foundFiles, subFiles...)
		} else {
			// 文件：检查扩展名是否为支持的图片格式
			ext := strings.ToLower(filepath.Ext(name))
			if !imageExtensions[ext] {
				continue // 非支持格式，跳过
			}
			// 获取文件信息
			info, err := entry.Info()
			if err != nil {
				fmt.Println("获取文件信息失败:", err)
				continue
			}
			size := info.Size()
			modTime := info.ModTime()
			// 说明: Go 的标准库 os.MediaFileInfo 不直接提供创建时间。如果需要准确创建时间，可以使用第三方库获取。
			// 这里将修改时间先作为创建时间存储。
			imgInfo := model.MediaFileInfo{
				Name:    name,
				Path:    fullPath,
				Size:    size,
				ModTime: modTime,
			}
			node.Files = append(node.Files, imgInfo)
			foundFiles = append(foundFiles, imgInfo)
		}
	}
	return node, foundDirs, foundFiles, nil
}
