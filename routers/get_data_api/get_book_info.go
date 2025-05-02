package get_data_api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

func GetParentBookInfo(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "not set id param")
	}
	info, err := model.GetBookGroupInfoByChildBookID(id)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "ParentBookInfo not found")
	}
	return c.JSON(http.StatusOK, info)
}

func GetBookInfos(c echo.Context) error {
	// 书籍排列的方式，默认name
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "default"
	}

	// 参数未提供
	if c.QueryParam("book_group_id") == "" &&
		c.QueryParam("depth") == "" &&
		c.QueryParam("max_depth") == "" {
		return c.JSON(http.StatusBadRequest, "need book_group_id or depth or max_depth")
	}

	model.CheckAllBookFileExist()

	// 按照最大书籍所在深度获取书籍信息
	if c.QueryParam("max_depth") != "" {
		return GetBookInfosByMaxDepth(c, sortBy)
	}

	// 按照书籍组ID获取书籍信息
	if c.QueryParam("book_group_id") != "" {
		return GetBookInfosByGroupID(c, sortBy)
	}

	// 按照书籍所在深度获取书籍信息
	if c.QueryParam("depth") != "" {
		return GetBookInfosByDepth(c, sortBy)
	}

	return nil
}

func GetTopOfShelfInfo(c echo.Context) error {
	// 书籍排列的方式，默认name
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "default"
	}

	model.CheckAllBookFileExist()
	bookInfoList, err := model.TopOfShelfInfo(sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "GetTopOfShelfInfo Failed")
	}
	return c.JSON(http.StatusOK, bookInfoList.BookInfos)
}

func GetBookInfosByMaxDepth(c echo.Context, sortBy string) error {
	model.CheckAllBookFileExist()
	// 按照书籍所在深度获取书籍信息，0是顶层，即为执行文件夹本身
	maxDepth, err := strconv.Atoi(c.QueryParam("max_depth"))
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "book_info not found")
	}
	// 如果传了maxDepth这个参数
	bookInfoList, err := model.GetBookInfoListByMaxDepth(maxDepth, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "book_info not found")
	}
	return c.JSON(http.StatusOK, bookInfoList.BookInfos)
}

func GetBookInfosByDepth(c echo.Context, sortBy string) error {
	model.CheckAllBookFileExist()
	// 按照书籍所在深度获取书籍信息，0是顶层，即为执行文件夹本身
	depth, err := strconv.Atoi(c.QueryParam("depth"))
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "book_info not found")
	}
	// 如果传了depth这个参数
	bookInfoList, err := model.GetBookInfoListByDepth(depth, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "book_info not found")
	}
	bookInfoList.SortBooks(sortBy)
	return c.JSON(http.StatusOK, bookInfoList.BookInfos)
}

func GetBookInfosByGroupID(c echo.Context, sortBy string) error {
	model.CheckAllBookFileExist()
	// 按照 BookGroupID获取
	bookGroupID := c.QueryParam("book_group_id")
	if bookGroupID == "" {
		return c.JSON(http.StatusBadRequest, "book_group_id not found")
	}
	// 如果传了bookGroupId这个参数
	bookInfoList, err := model.GetBookInfoListByID(bookGroupID, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "book_info not found")
	}
	bookInfoList.SortBooks(sortBy)
	return c.JSON(http.StatusOK, bookInfoList.BookInfos)
}

// GroupInfo 示例 URL： http://127.0.0.1:1234/api/group_info?id=1215a&sort_by=filename
func GroupInfo(c echo.Context) error {
	model.CheckAllBookFileExist()
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "default"
	}
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "book id not set")
	}
	// sortBy: 根据压缩包原始顺序、时间、文件名排序
	b, err := model.GetBookByID(id, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "book id not found")
	}
	infoList, err := model.GetBookInfoListByParentFolder(b.ParentFolder, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "ParentFolder, not found")
	}
	return c.JSON(http.StatusOK, infoList)
}

// GroupInfoFilter 示例 URL： http://127.0.0.1:1234/api/group_info_filter?id=1215a&sort_by=filename
func GroupInfoFilter(c echo.Context) error {
	model.CheckAllBookFileExist()
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "default"
	}
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "book id not set")
	}
	// sortBy: 根据压缩包原始顺序、时间、文件名排序
	b, err := model.GetBookByID(id, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "book id not found")
	}
	infoList, err := model.GetBookInfoListByParentFolder(b.ParentFolder, sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "ParentFolder, not found")
	}
	// 过滤掉不需要的类型
	filterList := model.BookInfoList{}
	for _, info := range infoList.BookInfos {
		if info.Type == model.TypeZip ||
			info.Type == model.TypeRar ||
			info.Type == model.TypeDir ||
			info.Type == model.TypeCbz ||
			info.Type == model.TypeCbr ||
			info.Type == model.TypePDF ||
			info.Type == model.TypeEpub {
			filterList.BookInfos = append(filterList.BookInfos, info)
		}
	}
	filterList.SortBooks(sortBy)
	return c.JSON(http.StatusOK, filterList)
}
