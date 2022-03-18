package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/common"
	"net/http"
)

// 相关参数：
// uuid：书籍的UUID，必须项目       							&uuid=2b17a130
// author：书籍的作者，未必存在									&author=佚名
// sort_page：按照自然文件名重新排序							&sort_page=true
// 示例 URL： http://127.0.0.1:1234/api/getbook?uuid=1215a
// 示例 URL： http://127.0.0.1:1234/api/getbook?&author=Doe&name=book_name
func getBookHandler(c *gin.Context) {
	author := c.DefaultQuery("author", "")
	sort := c.DefaultQuery("sort", "false")
	uuid := c.DefaultQuery("uuid", "")
	if author != "" {
		bookList, err := common.GetBookByAuthor(author, sort == "true")
		if err != nil {
			fmt.Println(err)
		} else {
			c.PureJSON(http.StatusOK, bookList)
		}
		return
	}

	if uuid != "" {
		b, err := common.GetBookByUUID(uuid, sort == "true")
		if err != nil {
			fmt.Println(err)
		} else {
			c.PureJSON(http.StatusOK, b)
		}
		return
	}
}
