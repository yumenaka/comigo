package data_api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
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

func parseStoreBookmarkRequest(c echo.Context) (bookmarkMutationRequest, error) {
	var payload storeBookmarkPayload
	if err := c.Bind(&payload); err != nil {
		return bookmarkMutationRequest{}, apiresp.BadRequest(c, "invalid_json", locale.GetString("err_invalid_json_request"), err.Error())
	}
	markType, err := parseBookmarkMarkType(c, payload.Type)
	if err != nil {
		return bookmarkMutationRequest{}, err
	}
	req := bookmarkMutationRequest{
		BookID:      payload.BookID,
		MarkType:    markType,
		PageIndex:   payload.PageIndex,
		Description: payload.Description,
	}
	if err := validateBookmarkRequiredFields(c, req, false); err != nil {
		return bookmarkMutationRequest{}, err
	}
	return req, nil
}

func parseDeleteBookmarkRequest(c echo.Context) (bookmarkMutationRequest, error) {
	pageIndexText := c.QueryParam("page_index")
	if pageIndexText == "" {
		return bookmarkMutationRequest{}, apiresp.BadRequest(c, "missing_page_index", locale.GetString("err_page_index_required"), nil)
	}
	pageIndex, err := strconv.Atoi(pageIndexText)
	if err != nil {
		return bookmarkMutationRequest{}, apiresp.BadRequest(c, "invalid_page_index", locale.GetString("err_page_index_invalid_number"), nil)
	}
	markType, err := parseBookmarkMarkType(c, c.QueryParam("mark_type"))
	if err != nil {
		return bookmarkMutationRequest{}, err
	}
	req := bookmarkMutationRequest{
		BookID:    c.QueryParam("book_id"),
		MarkType:  markType,
		PageIndex: pageIndex,
	}
	if err := validateBookmarkRequiredFields(c, req, true); err != nil {
		return bookmarkMutationRequest{}, err
	}
	return req, nil
}

func validateBookmarkRequiredFields(c echo.Context, req bookmarkMutationRequest, markTypeAlreadyChecked bool) error {
	if req.BookID == "" {
		return apiresp.BadRequest(c, "missing_book_id", locale.GetString("err_book_id_required"), nil)
	}
	if !markTypeAlreadyChecked && req.MarkType == "" {
		return apiresp.BadRequest(c, "missing_mark_type", locale.GetString("err_mark_type_required"), nil)
	}
	if req.PageIndex <= 0 {
		return apiresp.BadRequest(c, "missing_page_index", locale.GetString("err_page_index_required"), nil)
	}
	return nil
}

func parseBookmarkMarkType(c echo.Context, markType string) (model.MarkType, error) {
	if markType == "" {
		return "", apiresp.BadRequest(c, "missing_mark_type", locale.GetString("err_mark_type_required"), nil)
	}
	mt := model.MarkType(markType)
	if mt != model.AutoMark && mt != model.UserMark {
		return "", apiresp.BadRequest(c, "invalid_mark_type", locale.GetString("err_mark_type_invalid"), map[string][]string{
			"allowed": []string{string(model.AutoMark), string(model.UserMark)},
		})
	}
	return mt, nil
}

func requireBookmarkBook(c echo.Context, bookID string) (*model.Book, error) {
	book, err := model.IStore.GetBook(bookID)
	if err != nil {
		return nil, apiresp.Error(c, http.StatusBadRequest, "book_not_found", fmt.Sprintf(locale.GetString("err_book_not_found"), bookID), map[string]string{"book_id": bookID})
	}
	return book, nil
}

func validateBookmarkPageIndex(c echo.Context, req bookmarkMutationRequest, book *model.Book) error {
	if req.PageIndex <= 0 || req.PageIndex > book.PageCount {
		return apiresp.BadRequest(c, "page_index_out_of_range", locale.GetString("err_page_index_out_of_range"), map[string]int{
			"page_index": req.PageIndex,
			"page_count": book.PageCount,
		})
	}
	return nil
}
