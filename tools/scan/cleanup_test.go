package scan

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/sqlc"
	"github.com/yumenaka/comigo/store"
)

// cleanupFixture 使用临时书库和配置目录，禁止测试接触真实缓存。
func cleanupFixture(t *testing.T) (*store.StoreInRam, string, string) {
	t.Helper()
	oldCfg, oldStore := config.CopyCfg(), model.IStore
	root := t.TempDir()
	library := filepath.Join(root, "library")
	if err := os.MkdirAll(library, 0755); err != nil {
		t.Fatal(err)
	}
	cfg := oldCfg
	cfg.ConfigFile = filepath.Join(root, "config", "config.toml")
	cfg.CacheDir = filepath.Join(root, "cache")
	cfg.StoreUrls = []string{library}
	cfg.EnableDatabase = false
	config.ReplaceConfig(cfg)
	ram := &store.StoreInRam{}
	model.IStore = ram
	t.Cleanup(func() { config.ReplaceConfig(oldCfg); model.IStore = oldStore })
	return ram, library, root
}

// writeCleanupFile 创建可精确统计的测试文件。
func writeCleanupFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("fixture"), 0644); err != nil {
		t.Fatal(err)
	}
}

// TestCleanupPreviewAndExecute 验证升级遗留、源文件删除、确认过期，以及有效数据的保留。
func TestCleanupPreviewAndExecute(t *testing.T) {
	ram, library, root := cleanupFixture(t)
	active := &model.Book{BookInfo: model.BookInfo{BookID: "Active1", BookPath: filepath.Join(library, "active.cbz"), StoreUrl: library, Type: model.TypeCbz, CreatedByVersion: config.GetVersion()}}
	gone := &model.Book{BookInfo: model.BookInfo{BookID: "Deleted", BookPath: filepath.Join(library, "gone.cbz"), StoreUrl: library, Type: model.TypeCbz, CreatedByVersion: config.GetVersion()}}
	writeCleanupFile(t, active.BookPath)
	for _, b := range []*model.Book{active, gone} {
		if err := ram.StoreBook(b); err != nil {
			t.Fatal(err)
		}
	}
	// 待升级索引仍是书签恢复来源，不能因未加载到内存而清掉。
	pending := &model.Book{BookInfo: model.BookInfo{BookID: "Pending", BookPath: filepath.Join(library, "pending.cbz"), StoreUrl: library, Type: model.TypeCbz, CreatedByVersion: "v1.3.0"}}
	writeCleanupFile(t, pending.BookPath)
	if err := store.SaveMetaJson(pending); err != nil {
		t.Fatal(err)
	}
	ram.PendingBooks.Store([2]string{library, pending.BookPath}, pending)
	// 已有新索引的旧 ID 才能连同待恢复条目一起清理。
	legacy := *active
	legacy.BookID = "Legacy1"
	legacy.CreatedByVersion = "v1.3.0"
	if err := store.SaveMetaJson(&legacy); err != nil {
		t.Fatal(err)
	}
	ram.PendingBooks.Store([2]string{library, legacy.BookPath}, &legacy)
	metadata := filepath.Join(root, "config", "metadata", active.GetStoreID())
	keep := []string{active.BookPath, pending.BookPath, filepath.Join(metadata, "Active1.json"), filepath.Join(metadata, "Pending.json"), filepath.Join(metadata, "Pending_h300.jpg"), filepath.Join(root, "cache", "Active1", "page.data"), filepath.Join(root, "cache", "Pending", "page.data"), filepath.Join(root, "cache", "notes.txt")}
	remove := []string{filepath.Join(metadata, "Legacy1.json"), filepath.Join(metadata, "Deleted.json"), filepath.Join(metadata, "OldId12_original.jpg"), filepath.Join(root, "cache", "OldId12", "page.data"), filepath.Join(root, "cache", "Deleted", "page.data"), filepath.Join(root, "cache", "remote_books", "OldId12", "remote.pdf")}
	for _, p := range append(append([]string{}, keep[4:]...), remove[2:]...) {
		writeCleanupFile(t, p)
	}
	preview, err := CleanupInvalidData("")
	if err != nil {
		t.Fatal(err)
	}
	if preview.Files != len(remove) || preview.Books != 1 {
		t.Fatalf("preview: %+v", preview)
	}
	for _, p := range remove {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("preview mutated %s: %v", p, err)
		}
	}
	// 预览后的新增缓存必须先重新展示，不能混入已确认的清单。
	extra := filepath.Join(root, "cache", "OldId12", "new.data")
	writeCleanupFile(t, extra)
	if _, err := CleanupInvalidData(preview.Token); !errors.Is(err, ErrCleanupChanged) {
		t.Fatalf("stale preview: %v", err)
	}
	preview, err = CleanupInvalidData("")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CleanupInvalidData(preview.Token); err != nil {
		t.Fatal(err)
	}
	for _, p := range append(remove, extra) {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Fatalf("not removed %s: %v", p, err)
		}
	}
	for _, p := range keep {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("valid file lost %s: %v", p, err)
		}
	}
	if _, err := ram.GetBook(gone.BookID); err == nil {
		t.Fatal("invalid record remains")
	}
	if _, ok := ram.PendingBooks.Load([2]string{library, legacy.BookPath}); ok {
		t.Fatal("migrated pending index remains")
	}
	if _, ok := ram.PendingBooks.Load([2]string{library, pending.BookPath}); !ok {
		t.Fatal("pending index lost")
	}
	next, err := CleanupInvalidData("")
	if err != nil || next.Files != 0 || next.Books != 0 {
		t.Fatalf("cleanup not idempotent: %+v %v", next, err)
	}
}

// TestCleanupPreservesOfflineAndSymlink 验证离线书库与缓存内软链接不导致误删源文件。
func TestCleanupPreservesOfflineAndSymlink(t *testing.T) {
	ram, library, root := cleanupFixture(t)
	offline := filepath.Join(root, "offline")
	cfg := config.CopyCfg()
	cfg.StoreUrls = append(cfg.StoreUrls, offline)
	config.ReplaceConfig(cfg)
	book := &model.Book{BookInfo: model.BookInfo{BookID: "Offline", BookPath: filepath.Join(offline, "book.cbz"), StoreUrl: offline, Type: model.TypeCbz, CreatedByVersion: config.GetVersion()}}
	if err := ram.StoreBook(book); err != nil {
		t.Fatal(err)
	}
	cache := filepath.Join(root, "cache", "Offline", "page.data")
	writeCleanupFile(t, cache)
	source := filepath.Join(library, "original.jpg")
	writeCleanupFile(t, source)
	if err := os.Symlink(library, filepath.Join(root, "cache", "Symlink")); err != nil {
		t.Fatal(err)
	}
	p, err := CleanupInvalidData("")
	if err != nil || p.Files != 0 || p.Books != 0 || p.Skipped == 0 {
		t.Fatalf("unsafe preview: %+v %v", p, err)
	}
	if _, err := CleanupInvalidData(p.Token); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{source, cache} {
		if _, err := os.Stat(path); err != nil {
			t.Fatal(err)
		}
	}
	// 即便缓存目录配置为书库目录，也不能把名似 ID 的源目录当缓存。
	cfg.CacheDir = library
	config.ReplaceConfig(cfg)
	nested := filepath.Join(library, "OldId12", "source.jpg")
	writeCleanupFile(t, nested)
	p, err = CleanupInvalidData("")
	if err != nil || p.Files != 0 {
		t.Fatalf("overlap: %+v %v", p, err)
	}
}

// TestCleanupPreservesRestoredSource 源文件在预览后恢复时，旧确认不得删除其索引。
func TestCleanupPreservesRestoredSource(t *testing.T) {
	ram, library, _ := cleanupFixture(t)
	b := &model.Book{BookInfo: model.BookInfo{BookID: "Restore", BookPath: filepath.Join(library, "restored.cbz"), StoreUrl: library, Type: model.TypeCbz, CreatedByVersion: config.GetVersion()}}
	if err := ram.StoreBook(b); err != nil {
		t.Fatal(err)
	}
	p, err := CleanupInvalidData("")
	if err != nil {
		t.Fatal(err)
	}
	writeCleanupFile(t, b.BookPath)
	if _, err := CleanupInvalidData(p.Token); !errors.Is(err, ErrCleanupChanged) {
		t.Fatalf("restored source accepted stale plan: %v", err)
	}
	if _, err := ram.GetBook(b.BookID); err != nil {
		t.Fatal(err)
	}
}

// TestCleanupCacheAtConfigDirectory 保护配置目录中的 metadata 和用户插件，包含相对配置路径。
func TestCleanupCacheAtConfigDirectory(t *testing.T) {
	ram, library, root := cleanupFixture(t)
	t.Chdir(root)
	cfg := config.CopyCfg()
	cfg.ConfigFile = filepath.Join("config", "config.toml")
	config.ReplaceConfig(cfg)
	b := &model.Book{BookInfo: model.BookInfo{BookID: "Active1", BookPath: filepath.Join(library, "book.cbz"), StoreUrl: library, Type: model.TypeCbz, CreatedByVersion: config.GetVersion()}}
	writeCleanupFile(t, b.BookPath)
	if err := ram.StoreBook(b); err != nil {
		t.Fatal(err)
	}
	plugin := filepath.Join(root, "config", "plugins", "user.js")
	writeCleanupFile(t, plugin)
	for _, cache := range []string{filepath.Join(root, "config"), filepath.Join(root, "config", "metadata")} {
		cfg.CacheDir = cache
		config.ReplaceConfig(cfg)
		p, err := CleanupInvalidData("")
		if err != nil || p.Files != 0 {
			t.Fatalf("config mistaken for cache: %+v %v", p, err)
		}
	}
	// 清理后必须同步移除待恢复占位，不能因相对路径匹配失败而残留。
	legacy := b.CloneForView()
	legacy.BookID, legacy.CreatedByVersion = "Legacy1", "v1.3.0"
	if err := store.SaveMetaJson(legacy); err != nil {
		t.Fatal(err)
	}
	key := [2]string{library, legacy.BookPath}
	ram.PendingBooks.Store(key, legacy)
	p, err := CleanupInvalidData("")
	if err != nil || p.Files != 1 {
		t.Fatalf("legacy preview: %+v %v", p, err)
	}
	if _, err := CleanupInvalidData(p.Token); err != nil {
		t.Fatal(err)
	}
	if _, ok := ram.PendingBooks.Load(key); ok {
		t.Fatal("relative metadata path left stale pending entry")
	}
	if _, err := os.Stat(plugin); err != nil {
		t.Fatal(err)
	}
}

// TestCleanupRemoteIdentity 不同远程书籍不互相替代，移除书库后可清理无本地路径的索引和缓存。
func TestCleanupRemoteIdentity(t *testing.T) {
	ram, _, root := cleanupFixture(t)
	cfg := config.CopyCfg()
	cfg.StoreUrls = []string{"https://example.test"}
	config.ReplaceConfig(cfg)
	old := &model.Book{BookInfo: model.BookInfo{BookID: "OldId12", StoreUrl: cfg.StoreUrls[0], Type: model.TypeCbz, IsRemote: true, RemoteBookID: "A", RemoteStoreKey: "remote", CreatedByVersion: "v1.3.0"}}
	current := old.CloneForView()
	current.BookID, current.RemoteBookID, current.CreatedByVersion = "NewId12", "B", config.GetVersion()
	if err := ram.StoreBook(current); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveMetaJson(old); err != nil {
		t.Fatal(err)
	}
	key := [2]string{"", old.BookID}
	ram.PendingBooks.Store(key, old)
	cover := filepath.Join(root, "config", "metadata", old.GetStoreID(), old.BookID+"_h300.jpg")
	cache := filepath.Join(root, "cache", old.BookID, "page.data")
	writeCleanupFile(t, cover)
	writeCleanupFile(t, cache)
	p, err := CleanupInvalidData("")
	if err != nil || p.Files != 0 || p.Books != 0 {
		t.Fatalf("unrelated remote book treated as replacement: %+v %v", p, err)
	}
	cfg.StoreUrls = nil
	config.ReplaceConfig(cfg)
	p, err = CleanupInvalidData("")
	if err != nil || p.Files != 4 || p.Books != 1 {
		t.Fatalf("removed remote library not cleaned: %+v %v", p, err)
	}
	if _, err := CleanupInvalidData(p.Token); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{cover, cache, filepath.Join(root, "config", "metadata", old.GetStoreID(), old.BookID+".json")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("remote artifact remains: %s: %v", path, err)
		}
	}
	if _, ok := ram.PendingBooks.Load(key); ok {
		t.Fatal("removed remote library left pending entry")
	}
}

// TestCleanupSQLiteRecords 覆盖数据库模式下书籍、页面、书签与扫描失败记录的实际清理。
func TestCleanupSQLiteRecords(t *testing.T) {
	_, library, root := cleanupFixture(t)
	cfg := config.CopyCfg()
	cfg.EnableDatabase = true
	config.ReplaceConfig(cfg)
	if err := sqlc.OpenDatabase(sqlc.DBOptions{Type: "sqlite", ConfigDir: filepath.Join(root, "config")}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(sqlc.CloseDatabase)
	model.IStore = sqlc.DbStore
	b := &model.Book{BookInfo: model.BookInfo{BookID: "Deleted", BookPath: filepath.Join(library, "gone.cbz"), StoreUrl: library, Type: model.TypeCbz, CreatedByVersion: config.GetVersion()}, PageInfos: model.PageInfos{{Name: "page.jpg", PageNum: 1}}, BookMarks: model.BookMarks{{BookID: "Deleted", Type: model.AutoMark, PageIndex: 1}}}
	if err := model.IStore.StoreBook(b); err != nil {
		t.Fatal(err)
	}
	failures := scanFailureCache{"missing": {StoreURL: library, FilePath: b.BookPath, CreatedByVersion: config.GetVersion()}}
	if err := saveScanFailureCache(failures); err != nil {
		t.Fatal(err)
	}
	p, err := CleanupInvalidData("")
	if err != nil || p.Books != 1 || p.Pages != 1 || p.Bookmarks != 1 || p.Failures != 1 {
		t.Fatalf("database preview: %+v %v", p, err)
	}
	if _, err := CleanupInvalidData(p.Token); err != nil {
		t.Fatal(err)
	}
	// 清理入口必须实际调用回收，而不仅仅提供数据库方法。
	wal, err := os.Stat(filepath.Join(root, "config", "comigo.sqlite-wal"))
	if err != nil || wal.Size() != 0 {
		t.Fatalf("cleanup did not truncate WAL: %v %v", wal, err)
	}
	if _, err := model.IStore.GetBook(b.BookID); err == nil {
		t.Fatal("book remains")
	}
	rows, err := sqlc.DbStore.ListOrphanRows()
	if err != nil || len(rows["page_infos"])+len(rows["bookmarks"]) != 0 {
		t.Fatalf("associated rows remain: %+v %v", rows, err)
	}
	failures, err = loadScanFailureCache()
	if err != nil || len(failures) != 0 {
		t.Fatalf("failures remain: %+v %v", failures, err)
	}
}
