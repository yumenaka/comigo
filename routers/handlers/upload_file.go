package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
	"net/http"
	"os"
	"path/filepath"
)

var LocalRescanBroadcast *chan string
var EnableUpload *bool
var UploadPath *string

// UploadFile 下载服务器配置
// 除了设置头像以外，也可以做上传文件并阅读功能
// Set a lower memory limit for multipart forms (default is 32 MiB)

// UploadFile engine.MaxMultipartMemory = 60 << 20  // 60 MiB  只限制程序在上传文件时可以使用多少内存，而是不限制上传文件的大小。(default is 32 MiB)
func UploadFile(c *gin.Context) {

	if !*EnableUpload {
		logger.Infof("%s", locale.GetString("UPLOAD_DISABLE_HINT"))
		return
	}

	//默认的上传路径
	uploadDir := "upload"
	//如果设置过上传路径
	if *UploadPath != "" {
		uploadDir = *UploadPath
	}
	//创建上传目录（如果不存在）
	if !util.IsExist(uploadDir) {
		// 创建文件夹
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			// 无法创建上传目录: %s
			logger.Infof("mkdir failed![%s]\n", err)
		} else {
			// 创建上传目录成功: %s
			logger.Infof("%s", "mkdir success!\n")
		}
	}

	// 设置最大上传文件大小（例如 5000 MB）
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5000<<20) // 5000 MB

	// 解析多部分表单，最多允许 5000 MB
	if err := c.Request.ParseMultipartForm(5000 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "上传的文件过大"})
		return
	}

	form := c.Request.MultipartForm
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有上传文件"})
		return
	}

	var uploadedFiles []string
	fmt.Println("上传文件数量:", len(files))

	for _, file := range files {
		// 验证文件大小（例如，不超过 5000 MB）
		if file.Size > 5000<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("文件 %s 超过大小限制", file.Filename)})
			return
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
			//"application/octet-stream":     true, //未知的应用程序文件
		}
		if !allowedTypes[file.Header.Get("Content-Type")] {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("文件类型不允许: %s", file.Filename)})
			return
		}

		// 生成安全的文件名，避免目录遍历攻击
		filename := filepath.Base(file.Filename)

		// 确保文件名唯一（可选）
		destPath := filepath.Join(uploadDir, filename)
		// 如果文件已存在，追加编号
		counter := 1
		ext := filepath.Ext(filename)
		name := filename[:len(filename)-len(ext)]
		for {
			if _, err := os.Stat(destPath); os.IsNotExist(err) {
				break
			}
			filename = fmt.Sprintf("%s_%d%s", name, counter, ext)
			destPath = filepath.Join(uploadDir, filename)
			counter++
		}

		// 保存文件
		if err := c.SaveUploadedFile(file, destPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("无法保存文件 %s", filename)})
			return
		}

		uploadedFiles = append(uploadedFiles, filename)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文件上传成功",
		"files":   uploadedFiles,
	})

	// 通知重新扫描
	*LocalRescanBroadcast <- "upload"
}
