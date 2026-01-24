package data_api

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/epub"
	"github.com/yumenaka/comigo/model"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// DownloadEpub 将书籍转换为 EPUB 格式下载
// 示例 URL： http://127.0.0.1:1234/api/download-epub?id=2b17a13
// 相关参数：
// id：书籍的ID，必须参数  &id=2b17a13
func DownloadEpub(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	// 获取书籍信息
	book, err := model.IStore.GetBook(id)
	if err != nil {
		logger.Infof("GetBook error: %s", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}

	// 检查书籍类型，不支持 TypeBooksGroup、TypeVideo、TypeAudio、TypeHTML
	switch book.Type {
	case model.TypeBooksGroup:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot convert book group to EPUB"})
	case model.TypeVideo:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot convert video to EPUB"})
	case model.TypeAudio:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot convert audio to EPUB"})
	case model.TypeHTML:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot convert HTML to EPUB"})
	case model.TypeEpub:
		// 如果已经是 EPUB，直接返回原文件
		return c.File(book.BookPath)
	}

	// 检查页面数量
	if len(book.PageInfos) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Book has no pages"})
	}

	// 创建 EPUB 生成器
	generator, err := epub.NewGenerator()
	if err != nil {
		logger.Infof("Failed to create EPUB generator: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create EPUB generator"})
	}

	// 收集图片数据
	imageFiles := make([]epub.ImageFile, 0, len(book.PageInfos))
	for _, page := range book.PageInfos {
		var imgData []byte
		var imgErr error

		// 根据书籍类型获取图片数据
		option := fileutil.GetPictureDataOption{
			PictureName:      page.Name,
			BookIsDir:        book.Type == model.TypeDir,
			BookIsPDF:        book.Type == model.TypePDF,
			BookIsNonUTF8Zip: book.NonUTF8Zip,
			BookPath:         book.BookPath,
		}

		imgData, _, imgErr = fileutil.GetPictureData(option)
		if imgErr != nil {
			logger.Infof("Failed to get image %s: %v", page.Name, imgErr)
			continue
		}

		imageFiles = append(imageFiles, epub.ImageFile{
			Name: page.Name,
			Data: imgData,
		})
	}

	if len(imageFiles) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to extract any images from the book"})
	}

	// 创建书籍数据
	bookData := epub.CreateBookData(book.BookID, book.Title, book.Author, imageFiles)

	// 设置响应头
	epubFileName := filepath.Base(book.Title)
	// 移除原有扩展名
	if ext := filepath.Ext(epubFileName); ext != "" {
		epubFileName = epubFileName[:len(epubFileName)-len(ext)]
	}
	epubFileName += ".epub"
	encodedFileName := url.PathEscape(epubFileName)

	c.Response().Header().Set(echo.HeaderContentType, "application/epub+zip")
	// 使用 RFC 5987 格式支持中文文件名
	c.Response().Header().Set(echo.HeaderContentDisposition,
		fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", encodedFileName, encodedFileName))
	c.Response().WriteHeader(http.StatusOK)

	// 生成 EPUB 并写入响应
	if err := generator.Generate(c.Response().Writer, bookData, imageFiles); err != nil {
		logger.Infof("Failed to generate EPUB: %s", err)
		// 由于已经开始写入响应，无法返回 JSON 错误
		return nil
	}

	return nil
}
