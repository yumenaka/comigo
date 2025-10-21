package get_data_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

func UpdateBookmark(c echo.Context) error {
	// 解析请求体（JSON格式）
	var request struct {
		BookID      string `json:"book_id"`     // 书籍 ID
		PageIndex   int    `json:"page_index"`  // 书签页码，从 0 开始，不会超过 PageCount - 1
		Description string `json:"description"` // 用户添加的备注
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}
	// 获取书籍ID
	if request.BookID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "book_id is required")
	}
	model.IStore.CheckAllNotExistBooks()
	book, err := model.IStore.GetBookByID(request.BookID, "")
	if err != nil {
		logger.Infof("%s", err)
		return c.JSON(http.StatusBadRequest, "book not found, ID:"+request.BookID)
	}
	// 检查页码是否有效
	if request.PageIndex <= 0 || request.PageIndex > book.PageCount {
		return echo.NewHTTPError(http.StatusBadRequest, "page_index out of range")
	}
	// 查找是否已有该页码的书签
	existingBookmark := book.BookMarks.FindByPageIndex(request.PageIndex)
	if existingBookmark != nil {
		// 更新已有书签
		existingBookmark.Description = request.Description
		existingBookmark.UpdatedAt = time.Now()
	} else {
		// 创建新书签
		newBookmark := model.NewBookMark(request.BookID, request.PageIndex, book.PageCount, request.Description)
		book.BookMarks.Add(*newBookmark)
	}
	b, err := json.MarshalIndent(book.BookMarks, "", "  ")
	if err != nil {
		fmt.Println("unexpected error: %w", err)
	}
	fmt.Println(string(b))
	// 返回成功响应
	return c.JSON(http.StatusOK, map[string]string{"message": "bookmark updated successfully"})
}
