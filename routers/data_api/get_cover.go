package data_api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools/comigo_remote"
	coverutil "github.com/yumenaka/comigo/tools/cover"
	"github.com/yumenaka/comigo/tools/logger"
)

type coverRequest struct {
	bookID       string
	resizeHeight int
}

// GetCover 获取书籍封面
// 相关参数：
// id：书籍的ID，必须参数 &id=2B17a
// resize_height：可选参数，指定封面高度，默认值为352 &resize_height=500
// 示例 URL： http://127.0.0.1:1234/api/get-cover?id=2b17a13
// 示例 URL（自定义高度）： http://127.0.0.1:1234/api/get-cover?id=2b17a13&resize_height=500
func GetCover(c echo.Context) error {
	req, err := parseCoverRequest(c)
	if err != nil {
		return writeValidationError(c, err)
	}
	if localBook, client, ok, err := remoteComigoCoverTarget(c, req.bookID); ok {
		if err != nil {
			return writeRemoteComigoError(c, err)
		}
		data, contentType, err := client.GetBytes("/api/get-cover", remoteComigoQuery(c, localBook.RemoteBookID))
		if err != nil {
			return writeRemoteComigoError(c, err)
		}
		return serveCoverBytes(c, contentType, data)
	}
	result, err := coverutil.GetBookCover(coverutil.Request{BookID: req.bookID, ResizeHeight: req.resizeHeight})
	if err != nil {
		if errors.Is(err, coverutil.ErrBookNotFound) {
			return apiresp.Error(c, http.StatusNotFound, "book_not_found", "Book not found", map[string]string{"id": req.bookID})
		}
		logger.Infof(locale.GetString("log_get_file_error"), err)
		return apiresp.BadRequest(c, "get_cover_failed", "Get file error: "+err.Error(), nil)
	}

	return serveCoverBytes(c, result.ContentType, result.Data)
}

// remoteComigoCoverTarget 解析远端 Comigo 封面目标。
// 本地生成的远程书组没有 RemoteBookID，需要改用实际提供封面的子书远端 ID。
func remoteComigoCoverTarget(c echo.Context, localBookID string) (*model.Book, *comigo_remote.Client, bool, error) {
	remoteStoreKey := c.QueryParam(comigo_remote.RemoteStoreQuery)
	if remoteStoreKey == "" {
		return nil, nil, false, nil
	}
	client, _, err := remoteComigoClientByKey(remoteStoreKey)
	if err != nil {
		return nil, nil, true, err
	}
	book, err := model.IStore.GetBook(localBookID)
	if err != nil {
		return nil, nil, true, err
	}
	if book.RemoteBookID != "" {
		return book, client, true, nil
	}
	if book.Type == model.TypeBooksGroup {
		coverBook, err := remoteComigoBookGroupCoverBook(book)
		if err != nil {
			return nil, nil, true, err
		}
		return coverBook, client, true, nil
	}
	return nil, nil, true, errors.New("remote Comigo book id is empty")
}

func remoteComigoBookGroupCoverBook(book *model.Book) (*model.Book, error) {
	for _, childID := range book.ChildBooksID {
		childBook, err := model.IStore.GetBook(childID)
		if err != nil {
			continue
		}
		if childBook.RemoteBookID != "" {
			return childBook, nil
		}
	}
	return nil, errors.New("remote Comigo book group cover child id is empty")
}

// parseCoverRequest 解析 HTTP 查询参数，并保持 get-cover 原有默认高度。
func parseCoverRequest(c echo.Context) (coverRequest, error) {
	resizeHeight, err := parseOptionalBoundedInt(c, "resize_height", 352, 1, imageQueryMaxDimension)
	if err != nil {
		return coverRequest{}, err
	}
	req := coverRequest{
		bookID:       c.QueryParam("id"),
		resizeHeight: resizeHeight,
	}
	if req.bookID == "" {
		return req, requestValidationError{
			code:    "missing_param",
			message: "id is required",
			details: map[string]string{"param": "id"},
		}
	}
	return req, nil
}

// serveCoverBytes 统一输出封面字节，Content-Type 由复用的封面解析逻辑返回。
func serveCoverBytes(c echo.Context, contentType string, imgData []byte) error {
	return c.Blob(http.StatusOK, contentType, imgData)
}
