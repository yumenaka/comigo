package get_data_api

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

func GetRawFile(c echo.Context) error {
	bookID := c.Param("book_id")
	b, err := model.GetBookByID(bookID, "")
	// 打印文件名
	if err != nil {
		return c.String(http.StatusNotFound, "404 page not found")
	}
	fileName := c.Param("file_name")
	logger.Infof("下载文件：%s", fileName)

	// 获取文件信息
	fileInfo, err := os.Stat(b.FilePath)
	if err != nil {
		return c.String(http.StatusNotFound, "404 page not found")
	}
	// 如果是目录，返回目录列表
	if fileInfo.IsDir() {
		return c.String(http.StatusNotFound, "404 page not found")
	}
	// 如果是文件，返回文件
	return c.File(b.FilePath)
}
