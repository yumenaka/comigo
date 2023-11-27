package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/types"
)

func HandlerGetBookInfosByMaxDepth(c *gin.Context) {
	//书籍排列的方式，默认name
	sortBy := c.DefaultQuery("sort_by", "default")
	//按照书籍所在深度获取书籍信息，0是顶层，即为执行文件夹本身
	maxDepth, err := strconv.Atoi(c.DefaultQuery("max_depth", "0"))
	if err != nil {
		logger.Info(err)
		c.PureJSON(http.StatusBadRequest, "book_info not found")
		return
	}
	//如果传了maxDepth这个参数
	bookInfoList, err := types.GetBookInfoListByMaxDepth(maxDepth, sortBy)
	if err != nil {
		logger.Info(err)
		c.PureJSON(http.StatusBadRequest, "book_info not found")
		return
	}
	bookInfoList.SortBooks(sortBy)
	c.PureJSON(http.StatusOK, bookInfoList.BookInfos)
}

func HandlerGetBookInfosByDepth(c *gin.Context) {
	//书籍排列的方式，默认name
	sortBy := c.DefaultQuery("sort_by", "default")
	//按照书籍所在深度获取书籍信息，0是顶层，即为执行文件夹本身
	depth, err := strconv.Atoi(c.DefaultQuery("depth", ""))
	if err != nil {
		logger.Info(err)
		c.PureJSON(http.StatusBadRequest, "book_info not found")
		return
	}
	//如果传了depth这个参数
	bookInfoList, err := types.GetBookInfoListByDepth(depth, sortBy)
	if err != nil {
		logger.Info(err)
		c.PureJSON(http.StatusBadRequest, "book_info not found")
		return
	}
	bookInfoList.SortBooks(sortBy)
	c.PureJSON(http.StatusOK, bookInfoList.BookInfos)
}

func HandlerGetBookInfosByGroupID(c *gin.Context) {
	//书籍排列的方式，默认name
	sortBy := c.DefaultQuery("sort_by", "default")
	//bookGroup的BookId获取
	bookGroupId := c.DefaultQuery("book_group_book_id", "")
	if bookGroupId == "" {
		c.PureJSON(http.StatusBadRequest, "book_group_book_id not found")
		return
	}
	//如果传了bookGroupId这个参数
	bookInfoList, err := types.GetBookInfoListByID(bookGroupId, sortBy)
	if err != nil {
		logger.Info(err)
		c.PureJSON(http.StatusBadRequest, "book_info not found")
		return
	}
	bookInfoList.SortBooks(sortBy)
	c.PureJSON(http.StatusOK, bookInfoList.BookInfos)
}

// HandlerSameGroupBookInfo 示例 URL： http://127.0.0.1:1234/api/same_group_book_infos?id=1215a&sort_by=filename
func HandlerSameGroupBookInfo(c *gin.Context) {
	sortBy := c.DefaultQuery("sort_by", "filename")
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.PureJSON(http.StatusBadRequest, "book id not set")
		return
	}
	//sortBy: 根据压缩包原始顺序、时间、文件名排序
	b, err := types.GetBookByID(id, sortBy)
	if err != nil {
		logger.Info(err)
		c.PureJSON(http.StatusBadRequest, "book id not found")
		return
	}
	infoList, err := types.GetBookInfoListByParentFolder(b.ParentFolder, sortBy)
	if err != nil {
		logger.Info(err)
		c.PureJSON(http.StatusBadRequest, "ParentFolder, not found")
		return
	}
	c.PureJSON(http.StatusOK, infoList)
}
