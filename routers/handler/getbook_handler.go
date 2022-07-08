package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yumenaka/comi/book"
)

// GetBookHandler 相关参数：
// id：书籍的ID，必须项目       							&id=2b17a130
// author：书籍的作者，未必存在									&author=佚名
// sort_page：按照自然文件名重新排序							&sort_page=true
// 示例 URL： http://127.0.0.1:1234/api/getbook?bid=1215a
// 示例 URL： http://127.0.0.1:1234/api/getbook?&author=Doe&name=book_name
func GetBookHandler(c *gin.Context) {
	author := c.DefaultQuery("author", "")
	sortBy := c.DefaultQuery("sort_by", "default")
	id := c.DefaultQuery("id", "")
	if author != "" {
		//TODO: 根据压缩包原始顺序、时间、文件名排序
		bookList, err := book.GetBookByAuthor(author, sortBy)
		if err != nil {
			fmt.Println(err)
		} else {
			c.PureJSON(http.StatusOK, bookList)
		}
		return
	}
	if id != "" {
		//TODO: 根据压缩包原始顺序、时间、文件名排序
		b, err := book.GetBookByID(id, sortBy)
		if err != nil {
			fmt.Println(err)
		} else {
			c.PureJSON(http.StatusOK, b)
		}
		return
	}
}
