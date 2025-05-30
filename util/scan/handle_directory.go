package scan

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// 全局变量：标记是否正在扫描，避免并发扫描
var (
	scanning  bool       // 标记是否正在扫描，避免并发扫描
	scanMutex sync.Mutex // 保护 scanning 标志的锁
)

// HandleDirectory 扫描目录的核心函数：递归遍历目录，忽略指定名称的文件夹，收集图片文件信息
func HandleDirectory(currentPath string, depth int, option Option) (model.DirNode, []string, []model.MediaFileInfo, error) {
	node := model.DirNode{
		Name: filepath.Base(currentPath), // filepath.Base():返回路径的最后一个元素
		Path: currentPath,
	}
	var foundDirs []string
	var foundFiles []model.MediaFileInfo

	// 如果超过最大深度限制，直接返回空节点
	if option.Cfg.GetMaxScanDepth() >= 0 && depth > option.Cfg.GetMaxScanDepth() {
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
			if option.IsSkipDir(name) {
				continue
			}
			// 递归扫描子目录
			subNode, subDirs, subFiles, subErr := HandleDirectory(fullPath, depth+1, option)
			if subErr != nil {
				// 忽略单个子目录出错，继续扫描其他目录
				logger.Info("扫描子目录出错:", subErr)
				continue
			}
			node.SubDirs = append(node.SubDirs, subNode)
			// 合并子目录扫描结果
			foundDirs = append(foundDirs, subDirs...)
			foundFiles = append(foundFiles, subFiles...)
		} else {
			// 文件：检查扩展名是否为支持的格式
			ext := strings.ToLower(filepath.Ext(name))
			// 非支持媒体或压缩包格式，跳过
			if (!option.IsSupportMedia(ext)) && (!option.IsSupportFile(ext)) {
				continue
			}
			// 获取文件信息
			info, err := entry.Info()
			if err != nil {
				logger.Info("获取文件信息失败:", err)
				continue
			}
			size := info.Size()
			modTime := info.ModTime()
			// 说明: Go 的标准库 os.FileInfo 不直接提供创建时间。如果需要准确创建时间，可以使用第三方库获取。
			mediaFileInfo := model.MediaFileInfo{
				Name:    name,
				Path:    fullPath,
				Size:    size,
				ModTime: modTime,
			}
			node.Files = append(node.Files, mediaFileInfo)
			foundFiles = append(foundFiles, mediaFileInfo)
		}
	}
	return node, foundDirs, foundFiles, nil
}
