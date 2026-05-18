package data_api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	fileutil "github.com/yumenaka/comigo/tools/file"
)

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
