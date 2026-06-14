package data_api

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools/logger"
)

// DeleteBookmark 删除特定书签的 API 处理函数
// DELETE /api/delete-bookmark
// 参数：book_id, mark_type, page_index
func DeleteBookmark(c echo.Context) error {
	request, err := parseDeleteBookmarkRequest(c)
	if err != nil {
		return err
	}
	if localBook, client, _, ok, err := remoteComigoBookFromRequest(c, request.BookID); ok {
		if err != nil {
			logger.Infof("%s", err)
			return writeRemoteComigoError(c, err)
		}
		query := url.Values{}
		query.Set("book_id", localBook.RemoteBookID)
		query.Set("mark_type", string(request.MarkType))
		query.Set("page_index", strconv.Itoa(request.PageIndex))
		if _, _, err := client.Delete("/api/delete-bookmark", query); err != nil {
			logger.Infof(locale.GetString("log_failed_to_delete_bookmark"), err)
			return writeRemoteComigoError(c, err)
		}
		return apiresp.Success(c, "ok", locale.GetString("bookmark_deleted_successfully"), nil)
	}

	// 调用 Store 层删除书签
	err = model.IStore.DeleteBookMark(request.BookID, request.MarkType, request.PageIndex)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_delete_bookmark"), err)
		return apiresp.Error(c, http.StatusInternalServerError, "delete_bookmark_failed", locale.GetString("err_delete_bookmark_failed"), err.Error())
	}

	// 返回成功响应
	return apiresp.Success(c, "ok", locale.GetString("bookmark_deleted_successfully"), nil)
}
