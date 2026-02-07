package upload_api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
)

var RescanBroadcast *chan string

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

	// 获取用户选择的书库路径
	storeUrl := c.QueryParam("store_url")
	if storeUrl == "" {
		logger.Infof(locale.GetString("log_upload_no_store_selected"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": locale.GetString("store_validation_failed"),
		})
	}

	// 验证书库路径是否在配置的书库列表中
	storeUrls := config.GetCfg().StoreUrls
	isValidStore := false
	for _, url := range storeUrls {
		absUrl, err := filepath.Abs(url)
		if err != nil {
			absUrl = url
		}
		if absUrl == storeUrl {
			isValidStore = true
			break
		}
	}
	if !isValidStore {
		logger.Infof(locale.GetString("log_upload_invalid_store_path"), storeUrl)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": locale.GetString("store_validation_failed"),
		})
	}

	// 检查书库路径是否存在
	if !tools.IsExist(storeUrl) {
		logger.Infof(locale.GetString("log_upload_store_path_not_exist"), storeUrl)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": locale.GetString("store_not_exists"),
		})
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

		// 验证文件类型
		// 同时检查 Content-Type 和 文件扩展名，以支持更多格式
		allowedContentTypes := map[string]bool{
			"image/jpeg":                   true,
			"image/png":                    true,
			"image/gif":                    true,
			"image/webp":                   true,
			"application/pdf":              true,
			"application/zip":              true,
			"application/x-rar":            true,
			"application/x-rar-compressed": true,
			"application/x-zip-compressed": true,
			"application/x-7z-compressed":  true,
			"application/x-tar":            true,
			"application/gzip":             true,
			"application/epub+zip":         true,
		}

		// 允许的文件扩展名（用于处理 application/octet-stream 情况）
		allowedExtensions := map[string]bool{
			".zip":  true,
			".cbz":  true,
			".rar":  true,
			".cbr":  true,
			".7z":   true,
			".cb7":  true,
			".tar":  true,
			".gz":   true,
			".pdf":  true,
			".epub": true,
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
			".webp": true,
		}

		contentType := file.Header.Get("Content-Type")
		ext := strings.ToLower(filepath.Ext(file.Filename))

		// 如果 Content-Type 是 application/octet-stream，则基于扩展名判断
		isAllowed := allowedContentTypes[contentType]
		if !isAllowed && (contentType == "application/octet-stream" || contentType == "") {
			isAllowed = allowedExtensions[ext]
		}
		// 即使 Content-Type 不在列表中，如果扩展名合法也允许上传
		if !isAllowed {
			isAllowed = allowedExtensions[ext]
		}

		if !isAllowed {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": fmt.Sprintf("文件类型不允许: %s (类型: %s)", file.Filename, contentType),
			})
		}

		// 生成安全的文件名，避免目录遍历攻击
		filename := filepath.Base(file.Filename)

		// 确保文件名唯一（可选）
		destPath := filepath.Join(storeUrl, filename)
		// 如果文件已存在，追加编号
		counter := 1
		fileExt := filepath.Ext(filename)
		name := filename[:len(filename)-len(fileExt)]
		for {
			if _, err := os.Stat(destPath); os.IsNotExist(err) {
				break
			}
			filename = fmt.Sprintf("%s_%d%s", name, counter, fileExt)
			destPath = filepath.Join(storeUrl, filename)
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
	err = scan.InitStore(storeUrl, config.GetCfg())
	if err != nil {
		logger.Infof(locale.GetString("scan_error")+"path:%s  %s", storeUrl, err)
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
	model.GenerateBookGroup()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "文件上传成功",
		"files":   uploadedFiles,
	})
}
