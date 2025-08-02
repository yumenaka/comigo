package get_data_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

func GetParentBook(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "not set id param")
	}
	book, err := model.GetParentBook(id)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "ParentBookInfo not found")
	}
	return c.JSON(http.StatusOK, book)
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
