package error_page

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

// TestNotFoundCommonReturns404 验证错误页不会被监控或代理误判为成功响应。
func TestNotFoundCommonReturns404(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	if err := NotFoundCommon(e.NewContext(httptest.NewRequest(http.MethodGet, "/missing", nil), rec)); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if got := rec.Header().Get(echo.HeaderContentType); got != echo.MIMETextHTMLCharsetUTF8 {
		t.Fatalf("content type = %q", got)
	}
}
