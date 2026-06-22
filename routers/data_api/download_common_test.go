package data_api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

// 验证缺少书籍查询参数时会写入统一接口错误。
func TestRequireBookByQueryIDMissingIDWritesAPIError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/download-zip", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	book, handled, err := requireBookByQueryID(c)
	if err != nil {
		t.Fatalf("requireBookByQueryID returned error: %v", err)
	}
	if !handled {
		t.Fatal("missing id should be handled by helper")
	}
	if book != nil {
		t.Fatalf("expected nil book, got %#v", book)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var body struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if body.Code != "missing_param" || body.Message != "id is required" {
		t.Fatalf("expected API error body, got %#v", body)
	}
}

// 验证下载附件文件名会正确转义，避免响应头异常。
func TestSetAttachmentHeadersEscapesFileName(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/download-zip", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	setAttachmentHeaders(c, "application/zip", "测试 01.zip")

	if got := rec.Header().Get(echo.HeaderContentType); got != "application/zip" {
		t.Fatalf("content type = %q, want application/zip", got)
	}
	disposition := rec.Header().Get(echo.HeaderContentDisposition)
	if !strings.Contains(disposition, "attachment;") {
		t.Fatalf("content disposition missing attachment: %q", disposition)
	}
	if !strings.Contains(disposition, "filename*=UTF-8''%E6%B5%8B%E8%AF%95%2001.zip") {
		t.Fatalf("content disposition missing encoded filename: %q", disposition)
	}
}
