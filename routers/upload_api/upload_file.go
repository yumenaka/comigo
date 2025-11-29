package upload_api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
)

var (
	RescanBroadcast *chan string
)

// UploadFile 上传文件
// engine.MaxMultipartMemory = 60 << 20  // 60 MiB  只限制程序在上传文件时可以使用多少内存，而是不限制上传文件的大小。(default is 32 MiB)
func UploadFile(c echo.Context) error {
	// 是否开启上传功能
	if !config.GetCfg().EnableUpload || config.GetCfg().ReadOnlyMode {
		logger.Infof("%s", locale.GetString("upload_disable_hint"))
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error": locale.GetString("upload_disable_hint"),
		})
	}
	// 默认的上传路径是否已设置
	if config.GetCfg().UploadPath == "" {
		logger.Infof("%s", locale.GetString("log_upload_path_not_set"))
	}
	// 创建上传目录（如果不存在）
	if !tools.IsExist(config.GetCfg().UploadPath) {
		// 创建文件夹
		err := os.MkdirAll(config.GetCfg().UploadPath, os.ModePerm)
		if err != nil {
			// 无法创建上传目录: %s
			logger.Infof(locale.GetString("log_mkdir_failed"), err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": fmt.Sprintf("无法创建上传目录: %s", config.GetCfg().UploadPath),
			})
		}
		// 创建上传目录成功: %s
		logger.Infof(locale.GetString("log_mkdir_upload_folder_success"), config.GetCfg().UploadPath)
	}

	// 获取表单文件
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "解析表单失败",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "没有上传文件",
		})
	}

	var uploadedFiles []string
	logger.Infof(locale.GetString("log_upload_file_count"), len(files))
	// 遍历所有上传的文件
	for _, file := range files {
		// 验证文件大小（例如，不超过 5000 MB）
		if file.Size > 5000<<20 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": fmt.Sprintf("文件 %s 超过大小限制", file.Filename),
			})
		}
		// 验证文件类型（可选）
		// 例如，仅允许图片和PDF文件和压缩包文件
		allowedTypes := map[string]bool{
			"image/jpeg":                   true,
			"image/png":                    true,
			"image/gif":                    true,
			"application/pdf":              true,
			"application/zip":              true,
			"application/x-rar":            true,
			"application/x-rar-compressed": true,
			"application/x-zip-compressed": true,
			// "application/octet-stream":     true, //application/octet-stream 这是二进制文件的默认值。 由于这意味着未知的二进制文件，浏览器一般不会自动执行或询问执行
		}
		if !allowedTypes[file.Header.Get("Content-Type")] {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": fmt.Sprintf("文件类型不允许: %s", file.Filename),
			})
		}

		// 生成安全的文件名，避免目录遍历攻击
		filename := filepath.Base(file.Filename)

		// 确保文件名唯一（可选）
		destPath := filepath.Join(config.GetCfg().UploadPath, filename)
		// 如果文件已存在，追加编号
		counter := 1
		ext := filepath.Ext(filename)
		name := filename[:len(filename)-len(ext)]
		for {
			if _, err := os.Stat(destPath); os.IsNotExist(err) {
				break
			}
			filename = fmt.Sprintf("%s_%d%s", name, counter, ext)
			destPath = filepath.Join(config.GetCfg().UploadPath, filename)
			counter++
		}
		// 保存文件
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": fmt.Sprintf("无法打开文件 %s", filename),
			})
		}
		defer src.Close()

		dst, err := os.Create(destPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": fmt.Sprintf("无法创建文件 %s", filename),
			})
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": fmt.Sprintf("无法保存文件 %s", filename),
			})
		}
		logger.Infof(locale.GetString("log_file_upload_success"), filename)
		uploadedFiles = append(uploadedFiles, filename)
	}

	//// 通知重新扫描（不等待完成）
	//*RescanBroadcast <- "rescan_upload_path"

	// 同步执行扫描（等待完成）
	// 扫描上传目录的文件
	err = scan.InitStore(config.GetCfg().UploadPath, config.GetCfg())
	if err != nil {
		logger.Infof(locale.GetString("scan_error")+"path:%s  %s", config.GetCfg().UploadPath, err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": fmt.Sprintf("扫描上传目录失败: %s", err),
		})
	}
	// 保存扫描结果到数据库（如果开启）
	if config.GetCfg().EnableDatabase {
		if err := scan.SaveBooksToDatabase(config.GetCfg()); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": fmt.Sprintf("保存数据库失败: %s", err),
			})
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "文件上传成功",
		"files":   uploadedFiles,
	})
}
