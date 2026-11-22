package data_api

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
)

type bookmarkMutationRequest struct {
	BookID      string
	MarkType    model.MarkType
	PageIndex   int
	Description string
}

type storeBookmarkPayload struct {
	Type        string `json:"type"`
	BookID      string `json:"book_id"`
	PageIndex   int    `json:"page_index"`
	Description string `json:"description"`
}

// parseStoreBookmarkRequest 解析新增书签请求，校验失败不写响应。
func parseStoreBookmarkRequest(c echo.Context) (bookmarkMutationRequest, error) {
	var payload storeBookmarkPayload
	if err := c.Bind(&payload); err != nil {
		return bookmarkMutationRequest{}, requestValidationError{"invalid_json", locale.GetString("err_invalid_json_request"), err.Error()}
	}
	markType, err := parseBookmarkMarkType(payload.Type)
	if err != nil {
		return bookmarkMutationRequest{}, err
	}
	req := bookmarkMutationRequest{
		BookID:      payload.BookID,
		MarkType:    markType,
		PageIndex:   payload.PageIndex,
		Description: payload.Description,
	}
	if err := validateBookmarkRequiredFields(req); err != nil {
		return bookmarkMutationRequest{}, err
	}
	return req, nil
}

// parseDeleteBookmarkRequest 解析删除书签请求。
func parseDeleteBookmarkRequest(c echo.Context) (bookmarkMutationRequest, error) {
	pageIndexText := c.QueryParam("page_index")
	if pageIndexText == "" {
		return bookmarkMutationRequest{}, requestValidationError{"missing_page_index", locale.GetString("err_page_index_required"), nil}
	}
	pageIndex, err := strconv.Atoi(pageIndexText)
	if err != nil {
		return bookmarkMutationRequest{}, requestValidationError{"invalid_page_index", locale.GetString("err_page_index_invalid_number"), nil}
	}
	markType, err := parseBookmarkMarkType(c.QueryParam("mark_type"))
	if err != nil {
		return bookmarkMutationRequest{}, err
	}
	req := bookmarkMutationRequest{
		BookID:    c.QueryParam("book_id"),
		MarkType:  markType,
		PageIndex: pageIndex,
	}
	if err := validateBookmarkRequiredFields(req); err != nil {
		return bookmarkMutationRequest{}, err
	}
	return req, nil
}

// validateBookmarkRequiredFields 检查书签必填字段。
func validateBookmarkRequiredFields(req bookmarkMutationRequest) error {
	if req.BookID == "" {
		return requestValidationError{"missing_book_id", locale.GetString("err_book_id_required"), nil}
	}
	if req.PageIndex <= 0 {
		return requestValidationError{"missing_page_index", locale.GetString("err_page_index_required"), nil}
	}
	return nil
}

// parseBookmarkMarkType 解析允许的书签类型。
func parseBookmarkMarkType(markType string) (model.MarkType, error) {
	if markType == "" {
		return "", requestValidationError{"missing_mark_type", locale.GetString("err_mark_type_required"), nil}
	}
	mt := model.MarkType(markType)
	if mt != model.AutoMark && mt != model.UserMark {
		return "", requestValidationError{"invalid_mark_type", locale.GetString("err_mark_type_invalid"), map[string][]string{
			"allowed": []string{string(model.AutoMark), string(model.UserMark)},
		}}
	}
	return mt, nil
}

// requireBookmarkBook 获取书籍并保留不存在的校验错误。
func requireBookmarkBook(bookID string) (*model.Book, error) {
	book, err := model.IStore.GetBook(bookID)
	if err != nil {
		return nil, requestValidationError{"book_not_found", fmt.Sprintf(locale.GetString("err_book_not_found"), bookID), map[string]string{"book_id": bookID}}
	}
	return book, nil
}

// validateBookmarkPageIndex 检查页码是否在书籍范围内。
func validateBookmarkPageIndex(req bookmarkMutationRequest, book *model.Book) error {
	if req.PageIndex <= 0 || req.PageIndex > book.PageCount {
		return requestValidationError{"page_index_out_of_range", locale.GetString("err_page_index_out_of_range"), map[string]int{
			"page_index": req.PageIndex,
			"page_count": book.PageCount,
		}}
	}
	return nil
}
