package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yumenaka/comi/book"
)

// GetBookListHandler
// 可选参数，三择一：
// max_depth：书籍的最大深度									&max_depth=1
// book_group_book_id：按照书籍组的BookID 						&book_group_book_id=abc321
// depth：书籍的深度，      									&depth=0
// 示例 URL： http://127.0.0.1:1234/api/getshelf?sort=name&depth=0
// 示例 URL： http://127.0.0.1:1234/api/getshelf?book_group_id=aedxl
func GetBookListHandler(c *gin.Context) {
	//书籍排列的方式，默认name，TODO:按照修改时间、作者、文件大小等排序书籍
	sortBy := c.DefaultQuery("sort_by", "default")
	//按照书籍所在深度获取书籍信息，0是顶层，即为执行文件夹本身
	maxDepth, err := strconv.Atoi(c.DefaultQuery("max_depth", ""))
	//如果传了maxDepth这个参数
	if err == nil {
		bookInfoList, err := book.GetBookInfoListByMaxDepth(maxDepth, sortBy)
		if err != nil {
			fmt.Println(err)
			return
		}
		bookInfoList.SortBooks(sortBy)
		c.PureJSON(http.StatusOK, bookInfoList.BookInfos)
	}
	//bookGroup的BookId获取
	bookGroupId := c.DefaultQuery("book_group_book_id", "")
	//如果传了bookGroupId这个参数
	if bookGroupId != "" {
		bookInfoList, err := book.GetBookInfoListByBookGroupBookID(bookGroupId, sortBy)
		if err != nil {
			fmt.Println(err)
		}
		bookInfoList.SortBooks(sortBy)
		c.PureJSON(http.StatusOK, bookInfoList.BookInfos)
		return
	}

	//按照书籍所在深度获取书籍信息，0是顶层，即为执行文件夹本身
	depth, err := strconv.Atoi(c.DefaultQuery("depth", ""))
	//如果传了depth这个参数
	if err == nil {
		bookInfoList, err := book.GetBookInfoListByDepth(depth, sortBy)
		if err != nil {
			fmt.Println(err)
		} else {
			bookInfoList.SortBooks(sortBy)
			c.PureJSON(http.StatusOK, bookInfoList.BookInfos)
		}
	}
}
