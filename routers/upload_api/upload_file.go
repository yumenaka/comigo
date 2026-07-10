package upload_api

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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

const (
	maxUploadFileSize = int64(5000 << 20)
	maxUploadBodySize = maxUploadFileSize + 1<<20 // 为 multipart 边界与字段预留 1 MiB。
)

func uploadError(key string, args ...interface{}) map[string]interface{} {
	return map[string]interface{}{
		"error": fmt.Sprintf(locale.GetString(key), args...),
	}
}

// UploadFile 处理上传入口；Echo 的 MaxMultipartMemory 只限制解析时内存占用，不限制文件总大小。
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
	form, err := parseUploadMultipartForm(c, maxUploadBodySize)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return c.JSON(http.StatusRequestEntityTooLarge, uploadError("upload_file_too_large", "multipart"))
		}
		return c.JSON(http.StatusBadRequest, uploadError("upload_parse_form_failed"))
	}
	defer form.RemoveAll()

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, uploadError("upload_no_files"))
	}

	var uploadedFiles []string
	logger.Infof(locale.GetString("log_upload_file_count"), len(files))
	// 遍历所有上传的文件
	for _, file := range files {
		// 验证文件大小（例如，不超过 5000 MB）
		if file.Size > maxUploadFileSize {
			return c.JSON(http.StatusBadRequest, uploadError("upload_file_too_large", file.Filename))
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
			return c.JSON(http.StatusBadRequest, uploadError("upload_file_type_not_allowed", file.Filename, contentType))
		}

		// 生成安全的文件名，避免目录遍历攻击
		filename := filepath.Base(file.Filename)

		filename, err = saveUploadedFile(storeUrl, filename, file)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, uploadError("upload_save_file_failed", filename))
		}
		logger.Infof(locale.GetString("log_file_upload_success"), filename)
		uploadedFiles = append(uploadedFiles, filename)
	}

	// 同步执行扫描（等待完成）
	// 扫描上传目录的文件
	err = scan.InitStore(storeUrl, config.GetCfg())
	if err != nil {
		logger.Infof(locale.GetString("scan_error")+"path:%s  %s", storeUrl, err)
		return c.JSON(http.StatusInternalServerError, uploadError("upload_scan_failed", err))
	}
	// 保存扫描结果到数据库（如果开启）
	if config.GetCfg().EnableDatabase {
		if err := scan.SaveBooksToDatabase(config.GetCfg()); err != nil {
			return c.JSON(http.StatusInternalServerError, uploadError("upload_save_database_failed", err))
		}
	}
	model.GenerateBookGroup()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": locale.GetString("file_uploaded_successfully"),
		"files":   uploadedFiles,
	})
}

// parseUploadMultipartForm 在 multipart 写入临时磁盘前限制整个请求大小。
func parseUploadMultipartForm(c echo.Context, maxBytes int64) (*multipart.Form, error) {
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, maxBytes)
	return c.MultipartForm()
}

// saveUploadedFile 用 O_EXCL 原子选择文件名，并在复制结束后立即关闭文件。
func saveUploadedFile(storeURL, filename string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return filename, err
	}
	defer src.Close()

	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	for counter := 0; ; counter++ {
		candidate := filename
		if counter > 0 {
			candidate = fmt.Sprintf("%s_%d%s", name, counter, ext)
		}
		destPath := filepath.Join(storeURL, candidate)
		dst, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if errors.Is(err, os.ErrExist) {
			continue
		}
		if err != nil {
			return candidate, err
		}
		_, copyErr := io.Copy(dst, src)
		closeErr := dst.Close()
		if copyErr != nil {
			_ = os.Remove(destPath)
			return candidate, copyErr
		}
		if closeErr != nil {
			_ = os.Remove(destPath)
			return candidate, closeErr
		}
		return candidate, nil
	}
}
