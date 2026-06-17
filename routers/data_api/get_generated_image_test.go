package data_api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestGetGeneratedImageRejectsUnsafeParams(t *testing.T) {
	e := echo.New()
	cases := []string{
		"/api/get-generated-image?width=4097&height=1&font_size=12&text=x",
		"/api/get-generated-image?width=4096&height=4096&font_size=513&text=x",
		"/api/get-generated-image?width=100&height=100&font_size=NaN&text=x",
		"/api/get-generated-image?width=4096&height=4096&font_size=12&text=" + url.QueryEscape(strings.Repeat("字", generatedImageMaxTextRunes+1)),
	}

	for _, query := range cases {
		t.Run(query, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, query, nil)
			rec := httptest.NewRecorder()
			if err := GetGeneratedImage(e.NewContext(req, rec)); err != nil {
				t.Fatalf("GetGeneratedImage returned transport error: %v", err)
			}
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("unsafe generated image query should return 400, got %d", rec.Code)
			}
		})
	}
}

func TestParseGeneratedImageRequestAcceptsBoundaryValues(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/get-generated-image?width=4096&height=4096&font_size=512&text=ok", nil)
	rec := httptest.NewRecorder()
	parsed, err := parseGeneratedImageRequest(e.NewContext(req, rec))
	if err != nil {
		t.Fatalf("parseGeneratedImageRequest should accept boundary values: %v", err)
	}
	if parsed.width != 4096 || parsed.height != 4096 || parsed.fontSize != 512 {
		t.Fatalf("unexpected parsed request: %+v", parsed)
	}
}
