package settings

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
)

// 提供设置页书库计数测试所需的最小内存书库。
type storeBookCountsTestStore struct {
	books         map[string]*model.Book
	deleteCalls   int
	generateCalls int
}

func (s *storeBookCountsTestStore) StoreBook(b *model.Book) error {
	s.books[b.BookID] = b
	return nil
}

func (s *storeBookCountsTestStore) GetBook(id string) (*model.Book, error) {
	if b, ok := s.books[id]; ok {
		return b, nil
	}
	return nil, errors.New("book not found")
}

func (s *storeBookCountsTestStore) DeleteBook(id string) error {
	s.deleteCalls++
	delete(s.books, id)
	return nil
}

func (s *storeBookCountsTestStore) ListBooks() ([]*model.Book, error) {
	books := make([]*model.Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	return books, nil
}

func (s *storeBookCountsTestStore) GenerateBookGroup() error {
	s.generateCalls++
	return nil
}

func (s *storeBookCountsTestStore) StoreBookMark(mark *model.BookMark) error { return nil }

func (s *storeBookCountsTestStore) GetBookMarks(bookID string) (*model.BookMarks, error) {
	marks := model.BookMarks{}
	return &marks, nil
}

func (s *storeBookCountsTestStore) DeleteBookMark(bookID string, markType model.MarkType, pageIndex int) error {
	return nil
}

// 验证书库计数会清理已经不存在的本地书籍。
func TestGetStoreBookCountsCleansMissingBooks(t *testing.T) {
	oldCfg := config.CopyCfg()
	t.Cleanup(func() {
		*config.GetCfg() = oldCfg
	})
	oldStore := model.IStore
	t.Cleanup(func() {
		model.IStore = oldStore
	})

	storeDir := t.TempDir()
	existingPath := filepath.Join(storeDir, "exists.zip")
	if err := os.WriteFile(existingPath, []byte("zip"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	missingPath := filepath.Join(storeDir, "missing.zip")
	config.GetCfg().StoreUrls = []string{storeDir}

	testStore := &storeBookCountsTestStore{books: map[string]*model.Book{
		"existing": {
			BookInfo: model.BookInfo{
				BookID:   "existing",
				BookPath: existingPath,
				StoreUrl: storeDir,
				Type:     model.TypeZip,
			},
		},
		"missing": {
			BookInfo: model.BookInfo{
				BookID:   "missing",
				BookPath: missingPath,
				StoreUrl: storeDir,
				Type:     model.TypeZip,
			},
		},
	}}
	model.IStore = testStore

	counts := GetStoreBookCounts()
	if got := counts[storeBookCountKey(storeDir)]; got != 1 {
		t.Fatalf("store count = %d, want 1 after missing book cleanup", got)
	}
	if testStore.deleteCalls != 1 {
		t.Fatalf("DeleteBook calls = %d, want 1", testStore.deleteCalls)
	}
}

// 验证重新扫描能正确计算被移除的书籍数量。
func TestRescanBookDeltaReportsRemovedBooks(t *testing.T) {
	newBooksCount, removedBooksCount := rescanBookDelta(5, 3)
	if newBooksCount != 0 || removedBooksCount != 2 {
		t.Fatalf("delta = new %d removed %d, want new 0 removed 2", newBooksCount, removedBooksCount)
	}
}

// 验证重新扫描结果提示使用自然的中文文案。
func TestRescanBookDeltaMessageUsesNaturalChinese(t *testing.T) {
	locale.SetLanguage("zh-CN")

	tests := []struct {
		name    string
		added   int
		removed int
		want    string
	}{
		{name: "no change", want: "数量没变化"},
		{name: "added", added: 1, want: "多了 1 本书"},
		{name: "removed", removed: 1, want: "少了 1 本书"},
		{name: "both", added: 2, removed: 1, want: "新加 2 本书，少了 1 本书"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rescanBookDeltaMessage(tt.added, tt.removed); got != tt.want {
				t.Fatalf("message = %q, want %q", got, tt.want)
			}
		})
	}
}

// 验证单个书库计数只统计目标书库内的真实书籍。
func TestGetStoreRealBookCountOnlyCountsTargetStore(t *testing.T) {
	oldStore := model.IStore
	t.Cleanup(func() {
		model.IStore = oldStore
	})

	storeDir := t.TempDir()
	otherStoreDir := t.TempDir()
	model.IStore = &storeBookCountsTestStore{books: map[string]*model.Book{
		"target": {
			BookInfo: model.BookInfo{
				BookID:   "target",
				StoreUrl: storeDir,
				Type:     model.TypeZip,
			},
		},
		"group": {
			BookInfo: model.BookInfo{
				BookID:   "group",
				StoreUrl: storeDir,
				Type:     model.TypeBooksGroup,
			},
		},
		"other": {
			BookInfo: model.BookInfo{
				BookID:   "other",
				StoreUrl: otherStoreDir,
				Type:     model.TypeZip,
			},
		},
	}}

	if got := getStoreRealBookCount(storeDir); got != 1 {
		t.Fatalf("target store count = %d, want 1", got)
	}
}

// 验证 Wails 环境可显示系统目录选择入口。
func TestStoreConfigRendersWailsFolderPicker(t *testing.T) {
	var html bytes.Buffer
	if err := StoreConfig("StoreUrls", nil, "StoreUrls_Description", nil, false, false).Render(context.Background(), &html); err != nil {
		t.Fatalf("render StoreConfig: %v", err)
	}
	text := html.String()
	for _, want := range []string{
		`x-show="$store.global.wailsBook"`,
		`selectStoreFolder`,
		`/api/wails/select-directory`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("StoreConfig missing %q in %s", want, text)
		}
	}
}

// 验证 Wails 构建会隐藏会误报连接关闭的实时日志面板。
func TestMainAreaHidesServerLogInWailsBuild(t *testing.T) {
	if got, want := showServerLogInSettings(), !assets.IsWailsBuild(); got != want {
		t.Fatalf("showServerLogInSettings() = %v, want %v", got, want)
	}
}

// 验证远端 Comigo 版本更旧时会显示兼容性提示。
func TestRemoteComigoVersionWarningWhenRemoteOlder(t *testing.T) {
	server := newServerInfoTestServer(t, `{"Version":"v0.1.0","ServerName":"Comigo v0.1.0"}`)

	warning := remoteComigoVersionWarning(server.URL)
	if warning == "" {
		t.Fatal("expected warning for older remote Comigo version")
	}
	if !strings.Contains(warning, "v0.1.0") || !strings.Contains(warning, config.GetVersion()) {
		t.Fatalf("warning = %q, want both remote and local version", warning)
	}
}

// 验证远端 Comigo 版本更新时不显示兼容性提示。
func TestRemoteComigoVersionWarningAllowsNewerRemote(t *testing.T) {
	server := newServerInfoTestServer(t, `{"Version":"v99.0.0","ServerName":"Comigo v99.0.0"}`)

	if warning := remoteComigoVersionWarning(server.URL); warning != "" {
		t.Fatalf("warning = %q, want empty for newer remote Comigo", warning)
	}
}

// 创建只响应 /api/server-info 的测试服务器。
func newServerInfoTestServer(t *testing.T, body string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/server-info" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server
}
