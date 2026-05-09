package data_api

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/vfs"
)

func GetRawFile(c echo.Context) error {
	bookID := c.Param("book_id")
	b, err := model.IStore.GetBook(bookID)
	// 打印文件名
	if err != nil {
		return rawFileNotFound(c)
	}
	fileName := c.Param("file_name")
	logger.Infof(locale.GetString("log_download_file"), fileName)

	resource, err := openRawBookFile(b)
	if err != nil {
		return rawFileNotFound(c)
	}
	defer func() {
		_ = resource.file.Close()
	}()

	setRawFileHeaders(c, resource.contentType)

	// 使用标准库 ServeContent：支持 Range（206），适合媒体播放/拖动进度。
	http.ServeContent(c.Response().Writer, c.Request(), fileName, resource.modTime, resource.file)
	return nil
}

type rawFileResource struct {
	file        http.File
	modTime     time.Time
	contentType string
}

func openRawBookFile(book *model.Book) (*rawFileResource, error) {
	if book.IsRemote {
		return openRemoteRawBookFile(book)
	}
	return openLocalRawBookFile(book)
}

func openRemoteRawBookFile(book *model.Book) (*rawFileResource, error) {
	vfsInstance, err := vfs.GetOrCreate(book.RemoteURL, vfs.Options{
		CacheEnabled: false,
		Timeout:      30,
	})
	if err != nil {
		logger.Infof(locale.GetString("log_remote_store_connect_failed"), book.RemoteURL, err)
		return nil, err
	}

	fileInfo, err := vfsInstance.Stat(book.BookPath)
	if err != nil {
		logger.Infof(locale.GetString("log_remote_file_stat_failed"), book.BookPath, err)
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, os.ErrNotExist
	}

	openedFile, err := vfsInstance.Open(book.BookPath)
	if err != nil {
		logger.Infof(locale.GetString("log_remote_file_open_failed"), book.BookPath, err)
		return nil, err
	}

	return &rawFileResource{
		file:        &vfsFileAdapter{file: openedFile, fileInfo: fileInfo},
		modTime:     fileInfo.ModTime(),
		contentType: rawFileContentType(book.BookPath),
	}, nil
}

func openLocalRawBookFile(book *model.Book) (*rawFileResource, error) {
	fileInfo, err := os.Stat(book.BookPath)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, os.ErrNotExist
	}

	openedFile, err := os.Open(book.BookPath)
	if err != nil {
		return nil, err
	}
	return &rawFileResource{
		file:        openedFile,
		modTime:     fileInfo.ModTime(),
		contentType: rawFileContentType(book.BookPath),
	}, nil
}

func rawFileContentType(bookPath string) string {
	return mime.TypeByExtension(strings.ToLower(filepath.Ext(bookPath)))
}

func setRawFileHeaders(c echo.Context, contentType string) {
	if contentType == "" {
		return
	}
	c.Response().Header().Set(echo.HeaderContentType, contentType)
	// 对音视频尽量 inline，避免被当成附件下载导致无法播放。
	if strings.HasPrefix(contentType, "audio/") || strings.HasPrefix(contentType, "video/") {
		c.Response().Header().Set(echo.HeaderContentDisposition, "inline")
	}
}

func rawFileNotFound(c echo.Context) error {
	return c.String(http.StatusNotFound, "404 page not found")
}

// vfsFileAdapter 将 vfs.File 适配为 http.File（io.ReadSeeker）接口
// http.File 接口要求实现 io.Reader, io.Closer, io.Seeker, 以及 Stat() 方法
type vfsFileAdapter struct {
	file     vfs.File
	fileInfo vfs.FileInfo
}

func (a *vfsFileAdapter) Read(p []byte) (n int, err error) {
	return a.file.Read(p)
}

func (a *vfsFileAdapter) Close() error {
	return a.file.Close()
}

func (a *vfsFileAdapter) Seek(offset int64, whence int) (int64, error) {
	return a.file.Seek(offset, whence)
}

func (a *vfsFileAdapter) Stat() (os.FileInfo, error) {
	return a.fileInfo, nil
}

func (a *vfsFileAdapter) Readdir(count int) ([]os.FileInfo, error) {
	// http.File 接口要求实现 Readdir，但对于文件来说不需要
	return nil, nil
}
