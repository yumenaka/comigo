package scroll

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
)

// TestParseScrollBookmarkPage 验证 page 始终表示精确书页，limit 才启用分页加载。
func TestParseScrollBookmarkPage(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantMode  string
		wantChunk int
	}{
		{name: "无限卷轴精确页", url: "/scroll/book?page=50", wantMode: scrollLoadModeInfinite, wantChunk: -1},
		{name: "分页块第一页", url: "/scroll/book?page=32&limit=32", wantMode: scrollLoadModePaged, wantChunk: 1},
		{name: "分页块第二页", url: "/scroll/book?page=33&limit=32", wantMode: scrollLoadModePaged, wantChunk: 2},
		{name: "分页精确页", url: "/scroll/book?page=50&limit=32", wantMode: scrollLoadModePaged, wantChunk: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			c := e.NewContext(httptest.NewRequest(http.MethodGet, tt.url, nil), httptest.NewRecorder())
			mode := parseScrollLoadMode(c)
			if mode != tt.wantMode {
				t.Fatalf("load mode = %q, want %q", mode, tt.wantMode)
			}
			if got := parseScrollPageIndex(c, mode, 32); got != tt.wantChunk {
				t.Fatalf("chunk = %d, want %d", got, tt.wantChunk)
			}
		})
	}
}

// TestGetScrollPaginationURL 验证分页链接仍使用统一的精确 page 参数。
func TestGetScrollPaginationURL(t *testing.T) {
	book := &model.Book{BookInfo: model.BookInfo{BookID: "book", PageCount: 100}}
	if got := getScrollPaginationURL(book, 2, 32); got != "/scroll/book?page=33&limit=32" {
		t.Fatalf("pagination URL = %q", got)
	}

	oldBasePath := config.GetCfg().BasePath
	t.Cleanup(func() { config.GetCfg().BasePath = oldBasePath })
	config.GetCfg().BasePath = "/proxy"
	book.RemoteStoreKey = "remote store"
	if got := getScrollPaginationURL(book, 2, 32); got != "/proxy/scroll/book?page=33&limit=32&remote_store=remote+store" {
		t.Fatalf("remote pagination URL = %q", got)
	}
}
