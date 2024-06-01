package handlers

import (
	"github.com/yumenaka/comi/util/file"
	"github.com/yumenaka/comi/util/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/entity"
)

// GetBook 相关参数：
// id：书籍的ID，必须项目       							&id=2b17a130
// author：书籍的作者，未必存在									&author=佚名
// sort_page：按照自然文件名重新排序							&sort_page=true
// 示例 URL： http://127.0.0.1:1234/api/get_book?id=1215a&sort_by=name
// 示例 URL： http://127.0.0.1:1234/api/get_book?&author=Doe&name=book_name
func GetBook(c *gin.Context) {
	author := c.DefaultQuery("author", "")
	sortBy := c.DefaultQuery("sort_by", "default")
	id := c.DefaultQuery("id", "")
	if author != "" {
		bookList, err := entity.GetBookByAuthor(author, sortBy)
		if err != nil {
			logger.Infof("%s", err)
		}
		c.PureJSON(http.StatusOK, bookList)
		return
	}
	if id != "" {
		b, err := entity.GetBookByID(id, sortBy)
		if err != nil {
			logger.Infof("%s", err)
			c.PureJSON(http.StatusBadRequest, "id not found")
			return
		}
		// 如果是epub文件，重新按照Epub信息排序
		if b.Type == entity.TypeEpub && sortBy == "epub_info" {
			imageList, err := file.GetImageListFromEpubFile(b.FilePath)
			if err != nil {
				logger.Infof("%s", err)
				c.PureJSON(http.StatusOK, b)
				return
			}
			b.SortPagesByImageList(imageList)
		}
		c.PureJSON(http.StatusOK, b)
	}
}
