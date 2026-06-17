package data_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/logger"
)

func GetParentBook(c echo.Context) error {
	childID := c.QueryParam("id")
	if childID == "" {
		return apiresp.BadRequest(c, "missing_param", "not set id param", map[string]string{"param": "id"})
	}
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}
	for _, bookGroup := range allBooks {
		if bookGroup.Type != model.TypeBooksGroup {
			continue // 只分析书组类型
		}
		for _, id := range bookGroup.ChildBooksID {
			if id == childID {
				b, err := model.IStore.GetBook(bookGroup.BookID)
				if err != nil {
					return apiresp.Error(c, http.StatusNotFound, "parent_book_not_found", "ParentBookInfo not found", map[string]string{"id": childID})
				}
				return c.JSON(http.StatusOK, b)
			}
		}
	}
	return apiresp.Error(c, http.StatusNotFound, "parent_book_not_found", "ParentBookInfo not found", map[string]string{"id": childID})
}

func GetTopOfShelfInfo(c echo.Context) error {
	// 书籍排列的方式，默认name
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "default"
	}
	topOfShelfInfo, err := store.TopOfShelfInfo(sortBy)
	if err != nil {
		logger.Infof("%s", err)
		return apiresp.BadRequest(c, "top_shelf_failed", err.Error(), nil)
	}
	return c.JSON(http.StatusOK, topOfShelfInfo)
}
