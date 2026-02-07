package data_api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// DeleteBookmark 删除特定书签的 API 处理函数
// DELETE /api/delete-bookmark
// 参数：book_id, mark_type, page_index
func DeleteBookmark(c echo.Context) error {
	// 获取查询参数
	bookID := c.QueryParam("book_id")
	markType := c.QueryParam("mark_type")
	pageIndexStr := c.QueryParam("page_index")

	// 验证必需参数
	if bookID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "book_id is required")
	}
	if markType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "mark_type is required")
	}
	if pageIndexStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "page_index is required")
	}

	// 解析页码
	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid page_index: must be a number")
	}

	// 转换书签类型
	mt := model.MarkType(markType)
	if mt != model.AutoMark && mt != model.UserMark {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mark_type: must be 'auto' or 'user'")
	}

	// 调用 Store 层删除书签
	err = model.IStore.DeleteBookMark(bookID, mt, pageIndex)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_delete_bookmark"), err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete bookmark: "+err.Error())
	}

	// 返回成功响应
	return c.JSON(http.StatusOK, map[string]string{"message": "bookmark deleted successfully"})
}
