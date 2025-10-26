package data_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/logger"
)

func GetParentBook(c echo.Context) error {
	childID := c.QueryParam("id")
	if childID == "" {
		return c.JSON(http.StatusBadRequest, "not set id param")
	}
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, bookGroup := range allBooks {
		if bookGroup.Type != model.TypeBooksGroup {
			continue // 只分析书组类型
		}
		for _, id := range bookGroup.ChildBooksID {
			if id == childID {
				b, err := model.IStore.GetBook(bookGroup.BookID)
				if err != nil {
					return c.JSON(http.StatusBadRequest, "ParentBookInfo not found")
				}
				return c.JSON(http.StatusOK, b)
			}
		}
	}
	return c.JSON(http.StatusBadRequest, "ParentBookInfo not found")
}

func GetTopOfShelfInfo(c echo.Context) error {
	// 书籍排列的方式，默认name
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "default"
	}
	ClearBookNotExist()
	bookInfoList, err := store.TopOfShelfInfo(sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "GetTopOfShelfInfo Failed")
	}
	return c.JSON(http.StatusOK, bookInfoList.BookInfos)
}
