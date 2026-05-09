package data_api

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools/logger"
)

// DownloadZip 将 TypeDir 类型的书籍打包为 zip 文件下载
// 示例 URL： http://127.0.0.1:1234/api/download-zip?id=2b17a13
// 相关参数：
// id：书籍的ID，必须参数  &id=2b17a13
func DownloadZip(c echo.Context) error {
	book, handled, err := requireBookByQueryID(c)
	if handled || err != nil {
		return err
	}

	// 检查书籍类型是否为 TypeDir
	if book.Type != model.TypeDir {
		return apiresp.BadRequest(c, "unsupported_book_type", fmt.Sprintf("Only TypeDir books can be downloaded as zip, current type: %s", book.Type), map[string]string{"type": string(book.Type)})
	}

	// 检查目录是否存在
	if _, err := os.Stat(book.BookPath); os.IsNotExist(err) {
		return apiresp.Error(c, http.StatusNotFound, "book_directory_not_found", "Book directory not found", map[string]string{"id": book.BookID})
	}

	// 设置响应头
	// 使用 URL 编码处理文件名中的特殊字符
	zipFileName := filepath.Base(book.BookPath) + ".zip"
	setAttachmentHeaders(c, "application/zip", zipFileName)
	c.Response().WriteHeader(http.StatusOK)

	// 创建 zip writer，直接写入响应流
	zipWriter := zip.NewWriter(c.Response().Writer)
	defer func() {
		if err := zipWriter.Close(); err != nil {
			logger.Infof(locale.GetString("log_error_closing_zip_writer"), err)
		}
	}()

	// 遍历书籍的所有页面，添加到 zip
	for _, page := range book.PageInfos {
		// 获取文件的实际路径
		filePath := page.Path
		if filePath == "" {
			// 如果 Path 为空，尝试使用 BookPath + Name 构建路径
			filePath = filepath.Join(book.BookPath, page.Name)
		}

		// 检查文件是否存在
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			logger.Infof(locale.GetString("log_file_not_found_skipping"), filePath)
			continue
		}

		// 跳过目录
		if fileInfo.IsDir() {
			continue
		}

		// 打开源文件
		srcFile, err := os.Open(filePath)
		if err != nil {
			logger.Infof(locale.GetString("log_error_opening_file"), filePath, err)
			continue
		}

		// 在 zip 中创建文件，使用相对路径作为文件名
		// 这里使用 page.Name 保持文件名一致性
		header := &zip.FileHeader{
			Name:     page.Name,
			Method:   zip.Deflate,
			Modified: fileInfo.ModTime(),
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			srcFile.Close()
			logger.Infof(locale.GetString("log_error_creating_zip_entry"), page.Name, err)
			continue
		}

		// 复制文件内容到 zip
		_, err = io.Copy(writer, srcFile)
		srcFile.Close()
		if err != nil {
			logger.Infof(locale.GetString("log_error_writing_file_to_zip"), page.Name, err)
			continue
		}
	}

	return nil
}
