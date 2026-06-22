package scan

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/comigo_remote"
)

// 提供扫描测试所需的最小配置实现。
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

// 验证目录扫描只收集支持的文件。
// 这是防止初始化书库漏扫压缩包的回归用例。
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

// 验证目录书籍会使用扫描阶段收集到的文件列表。
func TestBookFromLocalDirNodeUsesScannedFiles(t *testing.T) {
	tmp := t.TempDir()
	bookDir := filepath.Join(tmp, "book")
	if err := os.MkdirAll(bookDir, 0o755); err != nil {
		t.Fatalf("mkdir bookDir: %v", err)
	}

	oldStore := model.IStore
	model.IStore = &store.StoreInRam{}
	t.Cleanup(func() {
		model.IStore = oldStore
	})

	node := DirNode{
		Name: "book",
		Path: bookDir,
		Files: []model.PageInfo{{
			Name:    "001.jpg",
			Path:    filepath.Join(bookDir, "001.jpg"),
			PageNum: 1,
		}, {
			Name:    "archive.zip",
			Path:    filepath.Join(bookDir, "archive.zip"),
			PageNum: 2,
		}},
	}
	InitConfig(&testCfgScan{
		supportMediaType: []string{".jpg"},
		supportFileType:  []string{".zip"},
	})
	book, err := bookFromLocalDirNode(node, tmp, 0)
	if err != nil {
		t.Fatalf("bookFromLocalDirNode error: %v", err)
	}
	if len(book.PageInfos) != 1 || book.PageInfos[0].Name != "001.jpg" || book.PageInfos[0].Url == "" {
		t.Fatalf("PageInfos = %+v, want only media page with generated URL", book.PageInfos)
	}
	if book.Type != model.TypeDir || book.BookPath != bookDir {
		t.Fatalf("book = %+v, want dir book at %s", book.BookInfo, bookDir)
	}
}

// 验证超过最大扫描深度的目录不会生成书籍。
func TestAppendLocalDirBooksSkipsOverMaxDepthNode(t *testing.T) {
	tmp := t.TempDir()
	nested := filepath.Join(tmp, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}

	oldStore := model.IStore
	model.IStore = &store.StoreInRam{}
	t.Cleanup(func() {
		model.IStore = oldStore
	})

	InitConfig(&testCfgScan{
		supportMediaType: []string{".jpg"},
		maxScanDepth:     0,
	})
	var books []*model.Book
	appendLocalDirBooks(&books, DirNode{
		Name: "a",
		Path: filepath.Join(tmp, "a"),
		SubDirs: []DirNode{{
			Name: "b",
			Path: nested,
			Files: []model.PageInfo{{
				Name: "001.jpg",
				Path: filepath.Join(nested, "001.jpg"),
			}},
		}},
	}, tmp)

	for _, book := range books {
		if book.BookPath == nested {
			t.Fatalf("unexpected nested book beyond max depth: %+v", book.BookInfo)
		}
	}
}

// 验证目录书籍内容变更后重新扫描会刷新页面列表。
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

// 验证远程 Comigo 书库扫描不会重复收录嵌套远程书籍。
func TestInitComigoStoreSkipsNestedRemoteBooks(t *testing.T) {
	// 模拟远端 Comigo 里同时存在本地书、远程书和混合书组，覆盖嵌套远程书过滤。
	remoteBooks := map[string]model.Book{
		"local": {BookInfo: model.BookInfo{
			BookID: "local",
			Title:  "Local Book",
			Type:   model.TypeZip,
			Depth:  0,
		}, PageInfos: model.PageInfos{
			{Name: "001.jpg", Url: "/api/get-file?id=local&filename=001.jpg"},
			{Name: "002.jpg", Url: "/api/get-file?id=local&filename=002.jpg"},
		}},
		"nested-remote": {BookInfo: model.BookInfo{
			BookID:         "nested-remote",
			Title:          "Nested Remote Book",
			Type:           model.TypeZip,
			Depth:          0,
			IsRemote:       true,
			RemoteStoreKey: "remote-key",
		}, PageInfos: model.PageInfos{
			{Name: "001.jpg", Url: "/api/get-file?id=nested-remote&filename=001.jpg"},
			{Name: "002.jpg", Url: "/api/get-file?id=nested-remote&filename=002.jpg"},
		}},
		"group": {BookInfo: model.BookInfo{
			BookID:       "group",
			Title:        "Mixed Group",
			Type:         model.TypeBooksGroup,
			Depth:        0,
			ChildBooksID: []string{"local", "nested-remote"},
		}},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/top-shelf":
			_ = json.NewEncoder(w).Encode([]model.StoreBookInfo{{
				DisplayName: "Remote Shelf",
				BookInfos: model.BookInfos{
					remoteBooks["local"].BookInfo,
					remoteBooks["nested-remote"].BookInfo,
					remoteBooks["group"].BookInfo,
				},
			}})
		case "/api/get-book":
			book, ok := remoteBooks[r.URL.Query().Get("id")]
			if !ok {
				http.NotFound(w, r)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"data": book})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	oldStore := model.IStore
	oldCfg := config.CopyCfg()
	model.IStore = &store.StoreInRam{}
	// 入库流程读取全局最小图片数，测试里固定成 1，避免本机配置影响断言。
	config.GetCfg().MinImageNum = 1
	t.Cleanup(func() {
		model.IStore = oldStore
		*config.GetCfg() = oldCfg
	})

	scanCfg := &testCfgScan{timeoutLimitForSc: 5}
	if err := InitStore(server.URL, scanCfg); err != nil {
		t.Fatalf("InitStore remote Comigo: %v", err)
	}

	localID := localComigoBookID(t, server.URL, "local")
	groupID := localComigoBookID(t, server.URL, "group")
	nestedRemoteID := localComigoBookID(t, server.URL, "nested-remote")
	if _, err := model.IStore.GetBook(localID); err != nil {
		t.Fatalf("本地书没有导入: %v", err)
	}
	group, err := model.IStore.GetBook(groupID)
	if err != nil {
		t.Fatalf("本地书组没有导入: %v", err)
	}
	if _, err := model.IStore.GetBook(nestedRemoteID); err == nil {
		t.Fatal("嵌套远程书被误导入")
	}
	// 混合书组应同步移除被跳过的远程子书，避免前端拿到不存在的 ChildBooksID。
	if len(group.ChildBooksID) != 1 || group.ChildBooksID[0] != localID {
		t.Fatalf("书组子书 ID = %v，期望只保留 %q", group.ChildBooksID, localID)
	}
}

// 验证远程书库删除失效书籍时会同时清理生成的分组。
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

func localComigoBookID(t *testing.T, storeURL string, remoteBookID string) string {
	t.Helper()
	return comigo_remote.LocalBookID(storeURL, remoteBookID)
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
