package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/types"
)

func GetRawFile(c *gin.Context) {
	bookID := c.Param("book_id")
	b, err := types.GetBookByID(bookID, "")
	// 打印文件名
	if err != nil {
		c.String(http.StatusNotFound, "404 page not found")
		return
	}
	fileName := c.Param("file_name")
	logger.Infof("下载文件：%s", fileName)

	// 获取文件信息
	fileInfo, err := os.Stat(b.FilePath)
	if err != nil {
		c.String(http.StatusNotFound, "404 page not found")
		return
	}
	// 如果是目录，返回目录列表
	if fileInfo.IsDir() {
		c.String(http.StatusNotFound, "404 page not found")
		return
	}
	// 如果是文件，返回文件
	c.File(b.FilePath)
}
