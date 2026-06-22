package data_api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	fileutil "github.com/yumenaka/comigo/tools/file"
)

// 验证文件请求在禁用缓存或带变换参数时会跳过缓存。
func TestParseGetFileRequestDisablesCacheForNoCacheAndTransforms(t *testing.T) {
	e := echo.New()

	cases := []struct {
		name  string
		query string
	}{
		{
			name:  "no-cache query",
			query: "/api/get-file?id=book1&filename=page.jpg&no-cache=true",
		},
		{
			name:  "image transform query",
			query: "/api/get-file?id=book1&filename=page.jpg&resize_width=320",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			rec := httptest.NewRecorder()
			parsed, err := parseGetFileRequest(e.NewContext(req, rec))
			if err != nil {
				t.Fatalf("parseGetFileRequest error: %v", err)
			}
			if !parsed.disableCache {
				t.Fatalf("%s 应禁用图片缓存", tc.name)
			}
		})
	}
}

// 验证文件请求会拒绝非法图片处理参数。
func TestParseGetFileRequestRejectsInvalidImageParams(t *testing.T) {
	e := echo.New()
	cases := []string{
		"/api/get-file?id=book1&filename=page.jpg&resize_width=abc",
		"/api/get-file?id=book1&filename=page.jpg&resize_width=-1",
		"/api/get-file?id=book1&filename=page.jpg&resize_width=4097",
		"/api/get-file?id=book1&filename=page.jpg&auto_crop=101",
		"/api/get-file?id=book1&filename=page.jpg&blurhash=3",
	}

	for _, query := range cases {
		t.Run(query, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, query, nil)
			rec := httptest.NewRecorder()
			if _, err := parseGetFileRequest(e.NewContext(req, rec)); err == nil {
				t.Fatalf("invalid image query should return validation error")
			}
		})
	}
}

// 验证图片参数非法时会先返回错误，不继续查书。
func TestGetFileRejectsInvalidImageParamBeforeBookLookup(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/get-file?id=book1&filename=page.jpg&resize_width=4097", nil)
	rec := httptest.NewRecorder()
	if err := GetFile(e.NewContext(req, rec)); err != nil {
		t.Fatalf("GetFile returned transport error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid image param should return 400 before book lookup, got %d", rec.Code)
	}
}

// 验证封面请求会拒绝危险的缩放高度。
func TestParseCoverRequestRejectsUnsafeResizeHeight(t *testing.T) {
	e := echo.New()
	for _, query := range []string{
		"/api/get-cover?id=book1&resize_height=0",
		"/api/get-cover?id=book1&resize_height=4097",
	} {
		t.Run(query, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, query, nil)
			rec := httptest.NewRecorder()
			if _, err := parseCoverRequest(e.NewContext(req, rec)); err == nil {
				t.Fatalf("invalid cover resize_height should return validation error")
			}
		})
	}
}

// 验证带禁用缓存参数时图片响应不会读取缓存。
func TestServeCachedPictureSkipsCacheWhenNoCacheQuerySet(t *testing.T) {
	cfg := config.GetCfg()
	oldUseCache := cfg.UseCache
	oldCacheDir := cfg.CacheDir
	oldDebug := cfg.Debug
	defer func() {
		cfg.UseCache = oldUseCache
		cfg.CacheDir = oldCacheDir
		cfg.Debug = oldDebug
	}()

	cfg.UseCache = true
	cfg.CacheDir = t.TempDir()
	cfg.Debug = false

	e := echo.New()
	bookID := "book1"
	filename := "page.jpg"
	req := httptest.NewRequest(http.MethodGet, "/api/get-file?id="+bookID+"&filename="+filename+"&no-cache=true", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 即使同一查询串已经存在缓存，no-cache 请求也必须跳过读取，强制走源文件路径。
	if err := fileutil.SaveFileToCache(bookID, filename, []byte("cached"), fileutil.GetQueryString(req.URL.Query()), "image/jpeg", cfg.CacheDir, false); err != nil {
		t.Fatalf("SaveFileToCache error: %v", err)
	}

	parsed, err := parseGetFileRequest(c)
	if err != nil {
		t.Fatalf("parseGetFileRequest error: %v", err)
	}
	handled, err := serveCachedPicture(c, parsed)
	if err != nil {
		t.Fatalf("serveCachedPicture error: %v", err)
	}
	if handled {
		t.Fatalf("no-cache 请求不应命中已有图片缓存")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("跳过缓存时不应提前写响应，got status %d", rec.Code)
	}
}
