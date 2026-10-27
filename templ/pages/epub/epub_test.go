package epub

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
)

// EPUB 页面复用原文件接口，支持基础路径，且不载入漫画脚本或暴露磁盘路径。
func TestPageUsesIndependentReader(t *testing.T) {
	oldStore, oldCfg := model.IStore, config.CopyCfg()
	t.Cleanup(func() { model.IStore = oldStore; config.ReplaceConfig(oldCfg) })
	model.IStore = &store.StoreInRam{}
	_, err := config.ChangeConfig(false, func(c *config.Config) error { c.BasePath = "/comics"; return nil })
	if err != nil {
		t.Fatal(err)
	}
	if err := model.IStore.StoreBook(&model.Book{BookInfo: model.BookInfo{BookID: "text", Title: "Text.epub", Type: model.TypeEpub, BookPath: "/private/library/Text.epub"}}); err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	e.GET("/comics/epub", PageHandler)
	e.GET("/comics/epub/:id", PageHandler)
	for _, path := range []string{"/comics/epub", "/comics/epub/text"} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		body := rec.Body.String()
		if rec.Code != http.StatusOK || !strings.Contains(body, "/comics/assets/static/js/epub_reader.js") {
			t.Fatalf("%s: status=%d", path, rec.Code)
		}
		for _, forbidden := range []string{"static/js/reader.js", "wasm_exec.js", "/private/library/"} {
			if strings.Contains(body, forbidden) {
				t.Fatalf("EPUB 页面包含 %s", forbidden)
			}
		}
		if path == "/comics/epub/text" && !strings.Contains(body, `data-source="/comics/api/raw/text/Text.epub"`) {
			t.Fatal("未复用原文件接口")
		}
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/comics/epub/missing", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("不存在的书: %d", rec.Code)
	}
}

// 漫画旧链接跳转时必须保留远程书库身份。
func TestRedirectPreservesRemoteStore(t *testing.T) {
	e := echo.New()
	e.GET("/flip/:id", Redirect)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/flip/book?remote_store=abc&page=2", nil))
	if rec.Code != http.StatusTemporaryRedirect || rec.Header().Get("Location") != config.PrefixPath("/epub/book?page=2&remote_store=abc") {
		t.Fatalf("redirect=%d %s", rec.Code, rec.Header().Get("Location"))
	}
}
