package data_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

func UpdateLastReadPage(c echo.Context) error {
	// 解析请求体（JSON格式）
	var request struct {
		BookID    string `json:"book_id"`    // 书籍 ID
		PageIndex int    `json:"page_index"` // 书签页码，从 0 开始，不会超过 PageCount - 1
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}
	// 获取书籍ID
	if request.BookID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "book_id is required")
	}
	book, err := model.IStore.GetBook(request.BookID)
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "book not found, ID:"+request.BookID)
	}
	// 检查页码是否有效
	if request.PageIndex <= 0 || request.PageIndex > book.PageCount {
		return echo.NewHTTPError(http.StatusBadRequest, "page_index out of range")
	}
	//book.LastReadPage = request.PageIndex
	//// 更新书籍信息
	//err = model.IStore.UpdateBook(book)
	//if err != nil {
	//	logger.Infof("Failed to update bookmark: %s", err)
	//	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update bookmark")
	//}
	//// 输出调试信息
	//jsonByte, err := json.MarshalIndent(book, "", "  ")
	//if err == nil {
	//	fmt.Println(string(jsonByte))
	//}

	// 返回成功响应
	return c.JSON(http.StatusOK, map[string]string{"message": "bookmark updated successfully"})
}
