package data_api

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
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

	// 获取文件信息
	fileInfo, err := os.Stat(b.BookPath)
	if err != nil {
		return c.String(http.StatusNotFound, "404 page not found")
	}
	// 如果是目录，返回目录列表
	if fileInfo.IsDir() {
		return c.String(http.StatusNotFound, "404 page not found")
	}
	// 如果是文件，返回文件（支持 Range）
	f, err := os.Open(b.BookPath)
	if err != nil {
		return c.String(http.StatusNotFound, "404 page not found")
	}
	defer func() {
		_ = f.Close()
	}()

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
	http.ServeContent(c.Response().Writer, c.Request(), fileName, fileInfo.ModTime(), f)
	return nil
}
