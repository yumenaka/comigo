package data_api

import (
	"net/http"

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

	// 调用 Store 层删除书签
	err = model.IStore.DeleteBookMark(request.BookID, request.MarkType, request.PageIndex)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_delete_bookmark"), err)
		return apiresp.Error(c, http.StatusInternalServerError, "delete_bookmark_failed", locale.GetString("err_delete_bookmark_failed"), err.Error())
	}

	// 返回成功响应
	return apiresp.Success(c, "ok", locale.GetString("bookmark_deleted_successfully"), nil)
}
