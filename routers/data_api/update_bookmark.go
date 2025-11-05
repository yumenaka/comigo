package data_api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

func StoreBookmark(c echo.Context) error {
	// 解析请求体（JSON格式）
	var request struct {
		Type        string `json:"type"`        // 书签类型，例如 "auto" 表示自动书签
		BookID      string `json:"book_id"`     // 书籍 ID
		PageIndex   int    `json:"page_index"`  // 书签页码，从 0 开始，不会超过 PageCount - 1
		Description string `json:"description"` // 书签描述，自动书签的话，是浏览器+系统信息（like：Firefox in Linux）
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
	markType := model.MarkType(request.Type)
	// 创建或更新书签
	bookMark := model.NewBookMark(markType, book.BookID, book.GetStoreID(), request.PageIndex, request.Description)
	err = model.IStore.StoreBookMark(bookMark)
	if err != nil {
		logger.Infof("Failed to store bookmark: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to store bookmark")
	}
	// 输出调试信息
	jsonByte, err := json.MarshalIndent(book.BookMarks, "", "  ")
	if err == nil {
		if config.GetCfg().Debug {
			logger.Infof("Updated bookmarks for book ID %s: %s", book.BookID, string(jsonByte))
		}
	}
	// 返回成功响应
	return c.JSON(http.StatusOK, map[string]string{"message": "bookmark updated successfully"})
}
