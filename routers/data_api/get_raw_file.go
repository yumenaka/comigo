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
		return c.String(http.StatusNotFound, "404 page not found")
	}
	fileName := c.Param("file_name")
	logger.Infof("下载文件：%s", fileName)

	var file http.File
	var modTime time.Time

	// 判断是否为远程书籍
	if b.IsRemote {
		// 远程书库：使用 VFS 获取文件信息
		vfsInstance, err := vfs.GetOrCreate(b.RemoteURL, vfs.Options{
			CacheEnabled: false,
			Timeout:      30,
		})
		if err != nil {
			logger.Infof(locale.GetString("log_remote_store_connect_failed"), b.RemoteURL, err)
			return c.String(http.StatusNotFound, "404 page not found")
		}

		// 获取文件信息
		vfsFileInfo, err := vfsInstance.Stat(b.BookPath)
		if err != nil {
			logger.Infof(locale.GetString("log_remote_file_stat_failed"), b.BookPath, err)
			return c.String(http.StatusNotFound, "404 page not found")
		}

		// 如果是目录，返回 404
		if vfsFileInfo.IsDir() {
			return c.String(http.StatusNotFound, "404 page not found")
		}

		// 打开文件
		vfsFile, err := vfsInstance.Open(b.BookPath)
		if err != nil {
			logger.Infof(locale.GetString("log_remote_file_open_failed"), b.BookPath, err)
			return c.String(http.StatusNotFound, "404 page not found")
		}
		defer func() {
			_ = vfsFile.Close()
		}()

		// 将 vfs.File 适配为 http.File（io.ReadSeeker）
		file = &vfsFileAdapter{file: vfsFile, fileInfo: vfsFileInfo}
		modTime = vfsFileInfo.ModTime()
	} else {
		// 本地书库：使用 os.Stat 和 os.Open
		osFileInfo, err := os.Stat(b.BookPath)
		if err != nil {
			return c.String(http.StatusNotFound, "404 page not found")
		}
		// 如果是目录，返回 404
		if osFileInfo.IsDir() {
			return c.String(http.StatusNotFound, "404 page not found")
		}

		// 打开文件（os.File 实现了 io.ReadSeeker）
		osFile, err := os.Open(b.BookPath)
		if err != nil {
			return c.String(http.StatusNotFound, "404 page not found")
		}
		defer func() {
			_ = osFile.Close()
		}()

		file = osFile
		modTime = osFileInfo.ModTime()
	}

	// 显式设置 Content-Type，避免浏览器无法识别媒体类型
	ext := strings.ToLower(filepath.Ext(b.BookPath))
	contentType := mime.TypeByExtension(ext)
	if contentType != "" {
		c.Response().Header().Set(echo.HeaderContentType, contentType)
	}

	// 对音视频尽量 inline，避免被当成附件下载导致无法播放
	if strings.HasPrefix(contentType, "audio/") || strings.HasPrefix(contentType, "video/") {
		c.Response().Header().Set(echo.HeaderContentDisposition, "inline")
	}

	// 使用标准库 ServeContent：支持 Range（206），适合媒体播放/拖动进度
	// http.ServeContent 需要一个实现了 io.ReadSeeker 的对象
	// 对于本地文件，os.File 已经实现了 io.ReadSeeker
	// 对于远程文件，vfsFileAdapter 将 vfs.File 适配为 io.ReadSeeker
	http.ServeContent(c.Response().Writer, c.Request(), fileName, modTime, file)
	return nil
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
