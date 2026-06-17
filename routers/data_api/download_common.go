package data_api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools/logger"
)

// requireBookByQueryID 从 ?id= 读取书籍，统一下载接口在流式写出前的错误响应。
// handled=true 表示已经写入响应，调用方应直接返回 err。
func requireBookByQueryID(c echo.Context) (book *model.Book, handled bool, err error) {
	id := c.QueryParam("id")
	if id == "" {
		return nil, true, apiresp.BadRequest(c, "missing_param", "id is required", map[string]string{"param": "id"})
	}

	book, err = model.IStore.GetBook(id)
	if err != nil {
		logger.Infof(locale.GetString("log_getbook_error_common"), err)
		return nil, true, apiresp.Error(c, http.StatusNotFound, "book_not_found", "Book not found", map[string]string{"id": id})
	}
	return book, false, nil
}

// setAttachmentHeaders 设置附件下载响应头，使用 RFC 5987 支持中文文件名。
func setAttachmentHeaders(c echo.Context, contentType, fileName string) {
	encodedFileName := url.PathEscape(fileName)
	c.Response().Header().Set(echo.HeaderContentType, contentType)
	c.Response().Header().Set(
		echo.HeaderContentDisposition,
		fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", encodedFileName, encodedFileName),
	)
}
