package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/yumenaka/comigo/model"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/directory"
)

//go:embed static_files/index.html
var indexHTML embed.FS

// -------------- 路由 Handler ----------------
// 1) 首页，返回嵌入的 index.html
func indexHandler(c echo.Context) error {
	data, err := indexHTML.ReadFile("static_files/index.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, "index.html not found")
	}
	return c.HTML(http.StatusOK, string(data))
}

// 2) 获取文件列表 API
func listHandler(c echo.Context) error {
	pathParam := c.QueryParam("path")
	pageParam := c.QueryParam("page")
	pageSizeParam := c.QueryParam("pageSize")

	if pathParam == "" {
		pathParam = rootPath
	}
	page := 1
	pageSize := 20
	fmt.Sscanf(pageParam, "%d", &page)
	fmt.Sscanf(pageSizeParam, "%d", &pageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	if useSQLite {
		ctx := context.Background()
		// 获取当前目录记录
		var dirRecord *ent.Directory
		var err error
		// 如果 pathParam == rootPath, 则找没有parent的那个(根)
		if pathParam == rootPath || pathParam == "" || pathParam == "/" {
			dirRecord, err = entClient.Directory.Query().
				Where(directory.Not(directory.HasParent())).
				Only(ctx)
		} else {
			dirRecord, err = entClient.Directory.Query().
				Where(directory.Path(pathParam)).
				Only(ctx)
		}
		if err != nil || dirRecord == nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "目录不存在"})
		}
		// 查子目录
		subDirs, _ := dirRecord.QueryChildren().All(ctx)
		// 查图片
		images, _ := dirRecord.QueryImages().All(ctx)
		total := len(images)

		// 分页
		if pageSize > 0 {
			start := (page - 1) * pageSize
			end := start + pageSize
			if start >= total {
				images = []*ent.Image{}
			} else {
				if end > total {
					end = total
				}
				images = images[start:end]
			}
		}

		// 构造返回数据
		resp := ListResponse{
			Directories: []DirectoryInfo{},
			Images:      []model.MediaFileInfo{},
			TotalImages: total,
			Page:        page,
			PageSize:    pageSize,
		}
		for _, d := range subDirs {
			resp.Directories = append(resp.Directories, DirectoryInfo{
				Name: d.Name,
				Path: d.Path,
			})
		}
		for _, img := range images {
			resp.Images = append(resp.Images, model.MediaFileInfo{
				Name:    img.Name,
				Path:    img.Path,
				Size:    img.Size,
				ModTime: img.ModTime,
			})
		}
		return c.JSON(http.StatusOK, resp)
	} else {
		dataMutex.RLock()
		defer dataMutex.RUnlock()
		// 在内存rootData中搜索指定路径的目录节点
		node := findDirNode(rootData, pathParam)
		if node == nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "目录不存在"})
		}
		total := len(node.Files)

		// 分页
		files := node.Files
		if pageSize > 0 {
			start := (page - 1) * pageSize
			end := start + pageSize
			if start >= total {
				files = []model.MediaFileInfo{}
			} else {
				if end > total {
					end = total
				}
				files = files[start:end]
			}
		}
		// 返回数据
		resp := ListResponse{
			Directories: []DirectoryInfo{},
			Images:      files,
			TotalImages: total,
			Page:        page,
			PageSize:    pageSize,
		}
		for _, sd := range node.SubDirs {
			resp.Directories = append(resp.Directories, DirectoryInfo{
				Name: sd.Name,
				Path: sd.Path,
			})
		}
		return c.JSON(http.StatusOK, resp)
	}
}

// 3) 手动触发重新扫描的端点
func rescanHandler(c echo.Context) error {
	go performScan() // 异步执行扫描，不阻塞请求
	return c.JSON(200, echo.Map{"message": "扫描任务已启动"})
}

// 4) 提供图片原始数据
func rawImageHandler(c echo.Context) error {
	file := c.QueryParam("file")
	if file == "" {
		return c.String(http.StatusBadRequest, "参数file不能为空")
	}
	// 如果要做安全控制，这里可校验file是否在根目录下
	return c.File(file)
}

// 执行扫描并更新存储
func performScan() {
	// 防止并发触发扫描
	scanMutex.Lock()
	if scanning {
		scanMutex.Unlock()
		return // 已经有扫描在进行中
	}
	scanning = true
	scanMutex.Unlock()

	// 执行扫描
	start := time.Now()
	newRoot, foundDirs, foundFiles, err := scanDirectory(rootPath, 0, maxDepth, ignoreDirs)
	if err != nil {
		fmt.Println("扫描过程中发生错误:", err)
	}
	scanDuration := time.Since(start)
	fmt.Printf("扫描完成，用时 %.2fs，发现 %d 个目录，%d 张图片。\n",
		scanDuration.Seconds(), len(foundDirs), len(foundFiles))

	if useSQLite {
		// 使用 SQLite 存储
		err := updateSQLiteStorage(foundDirs, foundFiles)
		if err != nil {
			fmt.Println("更新SQLite数据库失败:", err)
		}
	} else {
		// 使用 JSON 文件存储
		err := updateJSONStorage(newRoot)
		if err != nil {
			fmt.Println("更新JSON存储失败:", err)
		}
	}

	// 扫描完成，重置 scanning 标志
	scanMutex.Lock()
	scanning = false
	scanMutex.Unlock()
}
