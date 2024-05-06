package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/types"
)

func GetParentBookInfo(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.PureJSON(http.StatusBadRequest, "not set id param")
		return
	}
	info, err := types.GetBookGroupInfoByChildBookID(id)
	if err != nil {
		logger.Infof("%s", err)
		c.PureJSON(http.StatusBadRequest, "ParentBookInfo not found")
		return
	}
	c.PureJSON(http.StatusOK, info)
}

func GetBookInfos(c *gin.Context) {
	//书籍排列的方式，默认name
	sortBy := c.DefaultQuery("sort_by", "default")
	//参数未提供
	if c.Query("book_group_id") == "" &&
		c.Query("depth") == "" &&
		c.Query("max_depth") == "" {
		c.PureJSON(http.StatusBadRequest, "need book_group_id or depth or max_depth")
		return
	}
	//按照最大书籍所在深度获取书籍信息
	if c.Query("max_depth") != "" {
		GetBookInfosByMaxDepth(c, sortBy)
		return
	}
	//
	if c.Query("book_group_id") != "" {
		GetBookInfosByGroupID(c, sortBy)
		return
	}
	//按照书籍所在深度获取书籍信息
	if c.Query("depth") != "" {
		//按照书籍所在深度获取书籍信息，0是顶层，即为执行文件夹本身
		GetBookInfosByDepth(c, sortBy)
		return
	}
}

func GetBookInfosByMaxDepth(c *gin.Context, sortBy string) {
	//按照书籍所在深度获取书籍信息，0是顶层，即为执行文件夹本身
	maxDepth, err := strconv.Atoi(c.DefaultQuery("max_depth", "0"))
	if err != nil {
		logger.Infof("%s", err)
		c.PureJSON(http.StatusBadRequest, "book_info not found")
		return
	}
	//如果传了maxDepth这个参数
	bookInfoList, err := types.GetBookInfoListByMaxDepth(maxDepth, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		c.PureJSON(http.StatusBadRequest, "book_info not found")
		return
	}
	bookInfoList.SortBooks(sortBy)
	c.PureJSON(http.StatusOK, bookInfoList.BookInfos)
}

func GetBookInfosByDepth(c *gin.Context, sortBy string) {
	//按照书籍所在深度获取书籍信息，0是顶层，即为执行文件夹本身
	depth, err := strconv.Atoi(c.DefaultQuery("depth", ""))
	if err != nil {
		logger.Infof("%s", err)
		c.PureJSON(http.StatusBadRequest, "book_info not found")
		return
	}
	//如果传了depth这个参数
	bookInfoList, err := types.GetBookInfoListByDepth(depth, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		c.PureJSON(http.StatusBadRequest, "book_info not found")
		return
	}
	bookInfoList.SortBooks(sortBy)
	c.PureJSON(http.StatusOK, bookInfoList.BookInfos)
}

func GetBookInfosByGroupID(c *gin.Context, sortBy string) {
	//按照 BookGroupID获取
	bookGroupID := c.DefaultQuery("book_group_id", "")
	if bookGroupID == "" {
		c.PureJSON(http.StatusBadRequest, "book_group_id not found")
		return
	}
	//如果传了bookGroupId这个参数
	bookInfoList, err := types.GetBookInfoListByID(bookGroupID, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		c.PureJSON(http.StatusBadRequest, "book_info not found")
		return
	}
	bookInfoList.SortBooks(sortBy)
	c.PureJSON(http.StatusOK, bookInfoList.BookInfos)
}

// GroupInfo 示例 URL： http://127.0.0.1:1234/api/group_info?id=1215a&sort_by=filename
func GroupInfo(c *gin.Context) {
	sortBy := c.DefaultQuery("sort_by", "filename")
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.PureJSON(http.StatusBadRequest, "book id not set")
		return
	}
	//sortBy: 根据压缩包原始顺序、时间、文件名排序
	b, err := types.GetBookByID(id, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		c.PureJSON(http.StatusBadRequest, "book id not found")
		return
	}
	infoList, err := types.GetBookInfoListByParentFolder(b.ParentFolder, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		c.PureJSON(http.StatusBadRequest, "ParentFolder, not found")
		return
	}
	c.PureJSON(http.StatusOK, infoList)
}

// GroupInfoFilter 示例 URL： http://127.0.0.1:1234/api/group_info_filter?id=1215a&sort_by=filename
func GroupInfoFilter(c *gin.Context) {
	sortBy := c.DefaultQuery("sort_by", "filename")
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.PureJSON(http.StatusBadRequest, "book id not set")
		return
	}
	//sortBy: 根据压缩包原始顺序、时间、文件名排序
	b, err := types.GetBookByID(id, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		c.PureJSON(http.StatusBadRequest, "book id not found")
		return
	}
	infoList, err := types.GetBookInfoListByParentFolder(b.ParentFolder, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		c.PureJSON(http.StatusBadRequest, "ParentFolder, not found")
		return
	}
	//过滤掉不需要的类型
	filterList := types.BookInfoList{}
	filterList.SortBy = infoList.SortBy
	for _, info := range infoList.BookInfos {
		if info.Type == types.TypeZip ||
			info.Type == types.TypeRar ||
			info.Type == types.TypeDir ||
			info.Type == types.TypeCbz ||
			info.Type == types.TypeCbr ||
			info.Type == types.TypePDF ||
			info.Type == types.TypeEpub {
			filterList.BookInfos = append(filterList.BookInfos, info)
		}
	}
	c.PureJSON(http.StatusOK, filterList)
}
