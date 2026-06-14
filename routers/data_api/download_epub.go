package data_api

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/epub"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// DownloadEpub 将书籍转换为 EPUB 格式下载
// 示例 URL： http://127.0.0.1:1234/api/download-epub?id=2b17a13
// 相关参数：
// id：书籍的ID，必须参数  &id=2b17a13
func DownloadEpub(c echo.Context) error {
	if localBook, client, _, ok, err := remoteComigoBookFromRequest(c, c.QueryParam("id")); ok {
		if err != nil {
			logger.Infof("%s", err)
			return writeRemoteComigoError(c, err)
		}
		resp, err := client.GetResponse("/api/download-epub", remoteComigoQuery(c, localBook.RemoteBookID), nil)
		if err != nil {
			logger.Infof("%s", err)
			return writeRemoteComigoError(c, err)
		}
		defer resp.Body.Close()
		return streamRemoteDownload(c, resp)
	}

	book, handled, err := requireBookByQueryID(c)
	if handled || err != nil {
		return err
	}

	// 检查书籍类型，不支持 TypeBooksGroup、TypeVideo、TypeAudio、TypeHTML
	switch book.Type {
	case model.TypeBooksGroup:
		return apiresp.BadRequest(c, "unsupported_book_type", "Cannot convert book group to EPUB", map[string]string{"type": string(book.Type)})
	case model.TypeVideo:
		return apiresp.BadRequest(c, "unsupported_book_type", "Cannot convert video to EPUB", map[string]string{"type": string(book.Type)})
	case model.TypeAudio:
		return apiresp.BadRequest(c, "unsupported_book_type", "Cannot convert audio to EPUB", map[string]string{"type": string(book.Type)})
	case model.TypeHTML:
		return apiresp.BadRequest(c, "unsupported_book_type", "Cannot convert HTML to EPUB", map[string]string{"type": string(book.Type)})
	case model.TypeEpub:
		// 如果已经是 EPUB，直接返回原文件
		return c.File(book.BookPath)
	}

	// 检查页面数量
	if len(book.PageInfos) == 0 {
		return apiresp.BadRequest(c, "book_has_no_pages", "Book has no pages", map[string]string{"id": book.BookID})
	}

	// 创建 EPUB 生成器
	generator, err := epub.NewGenerator()
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_create_epub_generator"), err)
		return apiresp.Error(c, http.StatusInternalServerError, "create_epub_generator_failed", "Failed to create EPUB generator", err.Error())
	}

	// 收集图片数据
	imageFiles := make([]epub.ImageFile, 0, len(book.PageInfos))
	for _, page := range book.PageInfos {
		var imgData []byte
		var imgErr error

		// 根据书籍类型获取图片数据
		option := fileutil.GetPictureDataOption{
			PictureName:      page.Name,
			BookID:           book.BookID,
			BookIsDir:        book.Type == model.TypeDir,
			BookIsPDF:        book.Type == model.TypePDF,
			BookIsNonUTF8Zip: book.NonUTF8Zip,
			BookPath:         book.BookPath,
			// 远程书籍支持
			IsRemote:  book.IsRemote,
			RemoteURL: book.RemoteURL,
		}

		imgData, _, imgErr = fileutil.GetPictureData(option)
		if imgErr != nil {
			logger.Infof(locale.GetString("log_failed_to_get_image_epub"), page.Name, imgErr)
			continue
		}

		imageFiles = append(imageFiles, epub.ImageFile{
			Name: page.Name,
			Data: imgData,
		})
	}

	if len(imageFiles) == 0 {
		return apiresp.BadRequest(c, "extract_epub_images_failed", "Failed to extract any images from the book", map[string]string{"id": book.BookID})
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
	setAttachmentHeaders(c, "application/epub+zip", epubFileName)
	c.Response().WriteHeader(http.StatusOK)

	// 生成 EPUB 并写入响应
	if err := generator.Generate(c.Response().Writer, bookData, imageFiles); err != nil {
		logger.Infof(locale.GetString("log_failed_to_generate_epub"), err)
		// 由于已经开始写入响应，无法返回 JSON 错误
		return nil
	}

	return nil
}

// streamRemoteDownload 透传远端下载响应头和内容，让本地只做代理不保留原始文件。
func streamRemoteDownload(c echo.Context, resp *http.Response) error {
	contentType := resp.Header.Get(echo.HeaderContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	copyRemoteRawHeader(c, resp, echo.HeaderContentDisposition)
	copyRemoteRawHeader(c, resp, echo.HeaderContentLength)
	return c.Stream(resp.StatusCode, contentType, resp.Body)
}
