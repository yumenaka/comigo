package scan

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
)

// testCfgScan 仅用于扫描相关的单元测试
type testCfgScan struct {
	excludePath       []string
	supportMediaType  []string
	supportFileType   []string
	supportTemplate   []string
	maxScanDepth      int
	minImageNum       int
	timeoutLimitForSc int
}

func (c *testCfgScan) GetStoreUrls() []string           { return nil }
func (c *testCfgScan) GetMaxScanDepth() int             { return c.maxScanDepth }
func (c *testCfgScan) GetMinImageNum() int              { return c.minImageNum }
func (c *testCfgScan) GetTimeoutLimitForScan() int      { return c.timeoutLimitForSc }
func (c *testCfgScan) GetExcludePath() []string         { return c.excludePath }
func (c *testCfgScan) GetSupportMediaType() []string    { return c.supportMediaType }
func (c *testCfgScan) GetSupportFileType() []string     { return c.supportFileType }
func (c *testCfgScan) GetSupportTemplateFile() []string { return c.supportTemplate }
func (c *testCfgScan) GetZipFileTextEncoding() string   { return "utf-8" }
func (c *testCfgScan) GetEnableDatabase() bool          { return false }
func (c *testCfgScan) GetClearDatabaseWhenExit() bool   { return false }
func (c *testCfgScan) GetDebug() bool                   { return false }

// TestHandleDirectory_ShouldCollectSupportedFiles
// 回归用例：目录递归扫描时，应该能收集到支持的文件（例如 .zip），否则 InitStore 的“处理文件”阶段会漏扫。
func TestHandleDirectory_ShouldCollectSupportedFiles(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "test")
	if err := os.MkdirAll(filepath.Join(root, "TestDir"), 0o755); err != nil {
		t.Fatalf("mkdir TestDir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "TestDir 2"), 0o755); err != nil {
		t.Fatalf("mkdir TestDir 2: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "TestDir3"), 0o755); err != nil {
		t.Fatalf("mkdir TestDir3: %v", err)
	}
	// 放一个假的 zip 文件（不需要有效内容；这里只验证 HandleDirectory 是否能发现它）
	zipPath := filepath.Join(root, "TestDir", "a.zip")
	if err := os.WriteFile(zipPath, []byte("not-a-real-zip"), 0o644); err != nil {
		t.Fatalf("write a.zip: %v", err)
	}

	InitConfig(&testCfgScan{
		excludePath:       []string{},
		supportMediaType:  []string{".jpg", ".png", ".webp"},
		supportFileType:   []string{".zip", ".cbz", ".rar", ".cbr", ".tar", ".epub", ".pdf"},
		supportTemplate:   []string{".html"},
		maxScanDepth:      -1,
		minImageNum:       1,
		timeoutLimitForSc: 0,
	})

	_, _, foundFiles, err := HandleDirectory(root, 0)
	if err != nil {
		t.Fatalf("HandleDirectory error: %v", err)
	}
	found := false
	for _, f := range foundFiles {
		if f.Path == zipPath {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected to find %s in foundFiles, but it was missing", zipPath)
	}
}

func TestInitStoreRescansChangedDirectoryBook(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "library")
	bookDir := filepath.Join(root, "book")
	if err := os.MkdirAll(bookDir, 0o755); err != nil {
		t.Fatalf("mkdir bookDir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bookDir, "001.jpg"), []byte("one"), 0o644); err != nil {
		t.Fatalf("write first image: %v", err)
	}
	initialModTime := time.Unix(1700000000, 0)
	if err := os.Chtimes(bookDir, initialModTime, initialModTime); err != nil {
		t.Fatalf("chtimes initial bookDir: %v", err)
	}

	oldStore := model.IStore
	oldCfg := config.CopyCfg()
	model.IStore = &store.StoreInRam{}
	*config.GetCfg() = oldCfg
	config.GetCfg().ConfigFile = filepath.Join(tmp, "config", "config.toml")
	config.GetCfg().MinImageNum = 1
	t.Cleanup(func() {
		model.IStore = oldStore
		*config.GetCfg() = oldCfg
	})

	scanCfg := &testCfgScan{
		excludePath:       []string{},
		supportMediaType:  []string{".jpg"},
		supportFileType:   []string{".zip", ".cbz", ".rar", ".cbr", ".tar", ".epub", ".pdf"},
		supportTemplate:   []string{".html"},
		maxScanDepth:      -1,
		minImageNum:       1,
		timeoutLimitForSc: 0,
	}
	if err := InitStore(root, scanCfg); err != nil {
		t.Fatalf("initial InitStore: %v", err)
	}
	book := mustFindBookByPath(t, bookDir)
	oldBookID := book.BookID
	book.BookMarks = append(book.BookMarks, model.BookMark{
		Type:        model.UserMark,
		BookID:      "old-book-id",
		BookStoreID: "old-store-id",
		PageIndex:   1,
		Description: "keep",
		CreatedAt:   initialModTime,
		UpdatedAt:   initialModTime,
	})
	if err := model.IStore.StoreBook(book); err != nil {
		t.Fatalf("store bookmark state: %v", err)
	}

	if err := os.WriteFile(filepath.Join(bookDir, "002.jpg"), []byte("two"), 0o644); err != nil {
		t.Fatalf("write second image: %v", err)
	}
	changedModTime := initialModTime.Add(time.Hour)
	if err := os.Chtimes(bookDir, changedModTime, changedModTime); err != nil {
		t.Fatalf("chtimes changed bookDir: %v", err)
	}
	if err := InitStore(root, scanCfg); err != nil {
		t.Fatalf("rescan InitStore: %v", err)
	}

	refreshed := mustFindBookByPath(t, bookDir)
	if refreshed.PageCount != 2 {
		t.Fatalf("PageCount after rescan = %d, want 2", refreshed.PageCount)
	}
	if !refreshed.Modified.Equal(changedModTime) {
		t.Fatalf("Modified after rescan = %v, want %v", refreshed.Modified, changedModTime)
	}
	if len(refreshed.BookMarks) != 1 || refreshed.BookMarks[0].Description != "keep" {
		t.Fatalf("BookMarks after rescan = %+v, want preserved bookmark", refreshed.BookMarks)
	}
	if refreshed.BookMarks[0].BookID != refreshed.BookID {
		t.Fatalf("migrated bookmark BookID = %q, want %q", refreshed.BookMarks[0].BookID, refreshed.BookID)
	}
	if refreshed.BookMarks[0].BookStoreID != refreshed.GetStoreID() {
		t.Fatalf("migrated bookmark BookStoreID = %q, want %q", refreshed.BookMarks[0].BookStoreID, refreshed.GetStoreID())
	}
	if oldBookID == "" || refreshed.BookID == "" {
		t.Fatalf("book IDs should not be empty: old=%q new=%q", oldBookID, refreshed.BookID)
	}
}

func TestDeleteStaleComigoRemoteBooksRemovesGeneratedGroups(t *testing.T) {
	oldStore := model.IStore
	model.IStore = &store.StoreInRam{}
	t.Cleanup(func() {
		model.IStore = oldStore
	})

	const remoteStore = "http://example.com"
	currentBook := &model.Book{BookInfo: model.BookInfo{
		BookID:         "current",
		StoreUrl:       remoteStore,
		Type:           model.TypeZip,
		IsRemote:       true,
		RemoteBookID:   "remote-current",
		RemoteStoreKey: "remote-key",
	}}
	missingBook := &model.Book{BookInfo: model.BookInfo{
		BookID:         "missing",
		StoreUrl:       remoteStore,
		Type:           model.TypeZip,
		IsRemote:       true,
		RemoteBookID:   "remote-missing",
		RemoteStoreKey: "remote-key",
	}}
	generatedGroup := &model.Book{BookInfo: model.BookInfo{
		BookID:         "generated-group",
		StoreUrl:       remoteStore,
		Type:           model.TypeBooksGroup,
		IsRemote:       true,
		RemoteStoreKey: "remote-key",
	}}
	originalRemoteGroup := &model.Book{BookInfo: model.BookInfo{
		BookID:         "original-remote-group",
		StoreUrl:       remoteStore,
		Type:           model.TypeBooksGroup,
		IsRemote:       true,
		RemoteBookID:   "remote-group",
		RemoteStoreKey: "remote-key",
	}}
	localGroup := &model.Book{BookInfo: model.BookInfo{
		BookID:   "local-group",
		StoreUrl: "/local/books",
		Type:     model.TypeBooksGroup,
	}}
	for _, book := range []*model.Book{currentBook, missingBook, generatedGroup, originalRemoteGroup, localGroup} {
		if err := model.IStore.StoreBook(book); err != nil {
			t.Fatalf("StoreBook(%s): %v", book.BookID, err)
		}
	}

	deleteStaleComigoRemoteBooks(remoteStore, map[string]bool{"current": true, "original-remote-group": true})

	if _, err := model.IStore.GetBook("current"); err != nil {
		t.Fatalf("当前远程书被误删: %v", err)
	}
	if _, err := model.IStore.GetBook("local-group"); err != nil {
		t.Fatalf("本地书组被误删: %v", err)
	}
	if _, err := model.IStore.GetBook("original-remote-group"); err != nil {
		t.Fatalf("远端原始书组被误删: %v", err)
	}
	if _, err := model.IStore.GetBook("missing"); err == nil {
		t.Fatal("缺失的远程书没有被删除")
	}
	if _, err := model.IStore.GetBook("generated-group"); err == nil {
		t.Fatal("本地生成的远程书组没有被删除")
	}
}

func mustFindBookByPath(t *testing.T, bookPath string) *model.Book {
	t.Helper()
	absBookPath, err := filepath.Abs(bookPath)
	if err != nil {
		t.Fatalf("abs bookPath: %v", err)
	}
	books, err := model.IStore.ListBooks()
	if err != nil {
		t.Fatalf("ListBooks: %v", err)
	}
	for _, book := range books {
		if book.BookPath == absBookPath && book.Type == model.TypeDir {
			return book
		}
	}
	t.Fatalf("cannot find TypeDir book at %s in %+v", absBookPath, books)
	return nil
}
