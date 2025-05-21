package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util"

	"github.com/yumenaka/comigo/cmd/image_viewer/ent"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/directory"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/image"
)

// 全局变量（根据存储模式不同含义不同）:
var (
	rootData  DirNode      // JSON 模式下的内存目录树根节点
	dataMutex sync.RWMutex // 保护 rootData 读写的锁
	entClient *ent.Client  // SQLite 模式下的 ent 客户端
)

// 更新JSON存储
func updateJSONStorage(newRoot DirNode) error {
	// 用互斥锁保护数据更新
	dataMutex.Lock()
	rootData = newRoot // 更新全局的内存数据为最新扫描结果
	dataMutex.Unlock()

	// 将结果写入 JSON 文件（覆盖保存）
	dataBytes, err := json.MarshalIndent(rootData, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 JSON 时出错: %v", err)
	}

	err = os.WriteFile("scan_results.json", dataBytes, 0o644)
	if err != nil {
		return fmt.Errorf("写入 JSON 文件失败: %v", err)
	}

	return nil
}

// 加载JSON存储
func loadJSONStorage() error {
	dataBytes, err := os.ReadFile("scan_results.json")
	if err != nil {
		return err
	}

	var existingRoot DirNode
	if err := json.Unmarshal(dataBytes, &existingRoot); err != nil {
		return err
	}

	dataMutex.Lock()
	rootData = existingRoot
	dataMutex.Unlock()

	return nil
}

// 更新SQLite存储 - 简化版本，不使用事务
func updateSQLiteStorage(foundDirs []string, foundFiles []model.MediaFileInfo) error {
	// 检查entClient是否已初始化
	if entClient == nil {
		return fmt.Errorf("数据库客户端未初始化")
	}

	ctx := context.Background()

	// 预处理目录路径
	validDirs := make([]string, 0, len(foundDirs))
	for _, dir := range foundDirs {
		// 规范化路径
		normalizedPath := util.NormalizePath(dir)
		if util.IsValidPath(normalizedPath) {
			validDirs = append(validDirs, normalizedPath)
		}
	}

	// 按路径长度排序目录，确保先创建父目录再创建子目录
	sort.Slice(validDirs, func(i, j int) bool {
		return len(validDirs[i]) < len(validDirs[j])
	})

	// 将扫描结果更新到数据库（插入新数据或更新已有数据）
	// 先确保所有目录在数据库中存在
	dirMap := make(map[string]*ent.Directory)

	// 查询扫描结果中所有目录对应的数据库记录（提高批量操作效率）
	for _, dirPath := range validDirs {
		// 确保目录名不为空
		dirName := filepath.Base(dirPath)
		if dirName == "" || dirName == "." || dirName == "/" {
			dirName = "root" // 使用"root"作为默认名称
		}

		// 查找数据库是否已有该目录
		d, err := entClient.Directory.Query().Where(directory.Path(dirPath)).Only(ctx)
		if err != nil && ent.IsNotFound(err) {
			// 数据库中不存在则创建
			parentPath := filepath.Dir(dirPath)
			parentPath = util.NormalizePath(parentPath)

			// 创建目录记录
			builder := entClient.Directory.Create().
				SetName(dirName).
				SetPath(dirPath)

			// 只有当父目录路径有效且不等于当前目录路径时才设置父目录
			if parentPath != "." && parentPath != dirPath && parentPath != "" && parentPath != "/" {
				// 查找父目录的记录
				parent, err := entClient.Directory.Query().Where(directory.Path(parentPath)).Only(ctx)
				if err == nil && parent != nil {
					builder = builder.SetParent(parent)
				}
			}

			// 保存目录记录
			newDir, err := builder.Save(ctx)
			if err == nil && newDir != nil {
				d = newDir
			}
		}

		// 将目录记录存入 map
		if d != nil {
			dirMap[dirPath] = d
		}
	}

	// 预处理图片路径
	validImages := make([]model.MediaFileInfo, 0, len(foundFiles))
	for _, img := range foundFiles {
		// 规范化路径
		normalizedPath := util.NormalizePath(img.Path)
		if util.IsValidPath(normalizedPath) {
			img.Path = normalizedPath
			validImages = append(validImages, img)
		}
	}

	// 更新或插入图片记录
	for _, img := range validImages {
		// 确保图片名称不为空
		imgName := img.Name
		if imgName == "" {
			imgName = filepath.Base(img.Path)
		}
		if imgName == "" {
			imgName = "unnamed" // 默认名称
		}

		// 获取父目录路径
		parentDirPath := filepath.Dir(img.Path)
		parentDirPath = util.NormalizePath(parentDirPath)

		// 查找图片记录
		rec, err := entClient.Image.Query().Where(image.Path(img.Path)).Only(ctx)
		if err != nil && ent.IsNotFound(err) {
			// 数据库中无此图片，创建新记录
			parentDir := dirMap[parentDirPath]
			if parentDir == nil {
				continue
			}

			entClient.Image.Create().
				SetName(imgName).
				SetPath(img.Path).
				SetSize(img.Size).
				SetModTime(img.ModTime).
				SetDirectory(parentDir).
				Save(ctx)
		} else if err == nil && rec != nil {
			// 已存在此图片记录，更新其元数据（大小或时间可能变动）
			entClient.Image.UpdateOneID(rec.ID).
				SetSize(img.Size).
				SetModTime(img.ModTime).
				Save(ctx)
		}
	}

	return nil
}
