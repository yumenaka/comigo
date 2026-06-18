package shelf

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/model"
)

func TestBookCardRendersWailsContextMenuHook(t *testing.T) {
	var html bytes.Buffer
	marks := model.BookMarks{}
	book := model.BookInfo{
		BookID:    "book-id",
		Title:     "book.zip",
		Type:      model.TypeZip,
		PageCount: 1,
	}

	if err := BookCard(nil, book, marks).Render(context.Background(), &html); err != nil {
		t.Fatalf("render BookCard: %v", err)
	}
	if !strings.Contains(html.String(), "data-wails-book-card") {
		t.Fatalf("BookCard missing Wails context menu hook in %s", html.String())
	}
	if !strings.Contains(html.String(), `data-wails-book-id="book-id"`) {
		t.Fatalf("BookCard missing Wails book id in %s", html.String())
	}
}

func TestShelfHeaderTitleRendersRescanButton(t *testing.T) {
	var html bytes.Buffer

	if err := ShelfHeaderTitle("Library").Render(context.Background(), &html); err != nil {
		t.Fatalf("render ShelfHeaderTitle: %v", err)
	}
	rendered := html.String()
	for _, want := range []string{
		`@click.stop="window.ComiGoShelf?.rescanAllStores?.()"`,
		`:aria-label="i18next.t('rescan_all_stores')"`,
		`flex justify-center items-center w-8 h-8 mx-1 my-0 rounded hover:ring`,
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("ShelfHeaderTitle missing %q in %s", want, rendered)
		}
	}
}

func TestGetShelfSortByUsesCookieForNormalHTTP(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?sort_by=filesize_reverse", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	req.AddCookie(&http.Cookie{Name: "ShelfSortBy", Value: "filename"})
	c := e.NewContext(req, httptest.NewRecorder())

	if got := getShelfSortBy(c); got != "filename" {
		t.Fatalf("getShelfSortBy() = %q, want cookie sort", got)
	}
}

func TestGetShelfSortByOnlyPrefersQueryForWailsWebView(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?sort_by=filesize_reverse", nil)
	req.RemoteAddr = "192.0.2.1:1234"
	req.AddCookie(&http.Cookie{Name: "ShelfSortBy", Value: "filename"})
	c := e.NewContext(req, httptest.NewRecorder())

	want := "filename"
	if assets.IsWailsWebViewRequest(req) {
		want = "filesize_reverse"
	}
	if got := getShelfSortBy(c); got != want {
		t.Fatalf("getShelfSortBy() = %q, want %q", got, want)
	}
}
