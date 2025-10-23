package get_data_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

func GetParentBook(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "not set id param")
	}
	book, err := model.IStore.GetParentBook(id)
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
	model.IStore.ClearBookNotExist()
	bookInfoList, err := model.IStore.TopOfShelfInfo(sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "GetTopOfShelfInfo Failed")
	}
	return c.JSON(http.StatusOK, bookInfoList.BookInfos)
}
