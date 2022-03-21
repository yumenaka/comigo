package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/common"
	"net/http"
)

// 相关参数：
// id：书籍的ID，必须项目       							&id=2b17a130
// author：书籍的作者，未必存在									&author=佚名
// sort_page：按照自然文件名重新排序							&sort_page=true
// 示例 URL： http://127.0.0.1:1234/api/getbook?bid=1215a
// 示例 URL： http://127.0.0.1:1234/api/getbook?&author=Doe&name=book_name
func getBookHandler(c *gin.Context) {
	author := c.DefaultQuery("author", "")
	sort := c.DefaultQuery("sort", "false")
	id := c.DefaultQuery("id", "")
	if author != "" {
		bookList, err := common.GetBookByAuthor(author, sort == "true")
		if err != nil {
			fmt.Println(err)
		} else {
			c.PureJSON(http.StatusOK, bookList)
		}
		return
	}
	if id != "" {
		b, err := common.GetBookByID(id, sort == "true")
		if err != nil {
			fmt.Println(err)
		} else {
			c.PureJSON(http.StatusOK, b)
		}
		return
	}
}
