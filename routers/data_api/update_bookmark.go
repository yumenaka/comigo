package data_api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools/logger"
)

func StoreBookmark(c echo.Context) error {
	request, err := parseStoreBookmarkRequest(c)
	if err != nil {
		return err
	}
	if localBook, client, _, ok, err := remoteComigoBookFromRequest(c, request.BookID); ok {
		if err != nil {
			logger.Infof("%s", err)
			return writeRemoteComigoError(c, err)
		}
		payload, err := json.Marshal(storeBookmarkPayload{
			Type:        string(request.MarkType),
			BookID:      localBook.RemoteBookID,
			PageIndex:   request.PageIndex,
			Description: request.Description,
		})
		if err != nil {
			return err
		}
		if _, _, err := client.PostJSON("/api/store-bookmark", payload); err != nil {
			logger.Infof(locale.GetString("log_failed_to_store_bookmark"), err)
			return writeRemoteComigoError(c, err)
		}
		return apiresp.Success(c, "ok", locale.GetString("bookmark_updated_successfully"), nil)
	}
	book, err := requireBookmarkBook(c, request.BookID)
	if err != nil {
		logger.Infof("%s", err)
		return err
	}
	if err := validateBookmarkPageIndex(c, request, book); err != nil {
		return err
	}
	// 创建或更新书签
	bookMark := model.NewBookMark(request.MarkType, book.BookID, book.GetStoreID(), request.PageIndex, request.Description)
	err = model.IStore.StoreBookMark(bookMark)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_store_bookmark"), err)
		return apiresp.Error(c, http.StatusInternalServerError, "store_bookmark_failed", locale.GetString("err_store_bookmark_failed"), err.Error())
	}
	// 返回成功响应
	return apiresp.Success(c, "ok", locale.GetString("bookmark_updated_successfully"), nil)
}

// GetAllBookmarks 获取所有书签的 API 处理函数
func GetAllBookmarks(c echo.Context) error {
	// 获取所有书籍
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}
	// 收集所有书签
	allMarks := []model.BookinfoWithBookMark{}
	for _, book := range allBooks {
		for _, mark := range book.BookMarks {
			book.BookInfo.Cover = book.GetCover()
			bookinfoWithMark := model.BookinfoWithBookMark{
				BookInfo: book.BookInfo,
				BookMark: mark,
			}
			allMarks = append(allMarks, bookinfoWithMark)
		}
	}
	return c.JSON(http.StatusOK, allMarks)
}
