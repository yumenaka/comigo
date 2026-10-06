package sqlc

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/yumenaka/comigo/model"
)

func newTestStoreDatabase(t *testing.T) (*sql.DB, *StoreDatabase) {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite memory database: %v", err)
	}
	if err := configureSQLitePragmas(context.Background(), db); err != nil {
		t.Fatalf("configure sqlite pragmas: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), ddl); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	store := NewDBStore(db)

	oldStore := DbStore
	oldModelStore := model.IStore
	DbStore = store
	model.IStore = store
	t.Cleanup(func() {
		DbStore = oldStore
		model.IStore = oldModelStore
		_ = db.Close()
	})
	return db, store
}

// 验证 SQLite 初始化会启用增量自动清理。
func TestConfigureSQLitePragmasEnablesIncrementalAutoVacuum(t *testing.T) {
	db, _ := newTestStoreDatabase(t)

	var autoVacuum int
	if err := db.QueryRowContext(context.Background(), "PRAGMA auto_vacuum").Scan(&autoVacuum); err != nil {
		t.Fatalf("query auto_vacuum pragma: %v", err)
	}
	if autoVacuum != 2 {
		t.Fatalf("auto_vacuum = %d, want 2 (INCREMENTAL)", autoVacuum)
	}
}

// 验证数据库类型必须显式支持。
func TestOpenDatabaseRejectsUnsupportedType(t *testing.T) {
	if err := OpenDatabase(DBOptions{Type: "unsupported"}); err == nil {
		t.Fatal("OpenDatabase should reject unsupported database type")
	}
}

// 验证书籍写入数据库再读出时保留 JSON 元数据字段。
func TestStoreBookRoundTripKeepsJSONMetadataFields(t *testing.T) {
	_, store := newTestStoreDatabase(t)
	now := time.Date(2026, 4, 27, 10, 30, 0, 0, time.UTC)
	book := &model.Book{
		BookInfo: model.BookInfo{
			BookID:           "book-one",
			Title:            "Book One",
			Author:           "Author Name",
			Type:             model.TypeZip,
			BookPath:         "/library/book-one.zip",
			StoreUrl:         "/library",
			ParentFolder:     "library",
			FileSize:         1234,
			Modified:         now,
			PageCount:        2,
			CreatedByVersion: "v1.2.3",
		},
		PageInfos: model.PageInfos{
			{Name: "002.jpg", Path: "/library/002.jpg", PageNum: 2},
			{Name: "001.jpg", Path: "/library/001.jpg", PageNum: 1},
		},
		BookMarks: model.BookMarks{
			{
				Type:        model.UserMark,
				BookID:      "book-one",
				BookStoreID: "store-id",
				PageIndex:   1,
				Description: "note",
				CreatedAt:   now.Add(-time.Hour),
				UpdatedAt:   now,
			},
		},
	}

	if err := store.StoreBook(book); err != nil {
		t.Fatalf("store book: %v", err)
	}

	got, err := store.GetBook("book-one")
	if err != nil {
		t.Fatalf("get book: %v", err)
	}
	if got.Author != book.Author {
		t.Fatalf("author was not restored: got %q want %q", got.Author, book.Author)
	}
	if got.CreatedByVersion != book.CreatedByVersion {
		t.Fatalf("created version was not restored: got %q want %q", got.CreatedByVersion, book.CreatedByVersion)
	}
	if len(got.PageInfos) != 2 || got.PageInfos[0].Name != "002.jpg" {
		t.Fatalf("page infos were not restored in stored order: %#v", got.PageInfos)
	}
	if len(got.BookMarks) != 1 {
		t.Fatalf("bookmarks were not restored: %#v", got.BookMarks)
	}
	if got.BookMarks[0].BookStoreID != "store-id" {
		t.Fatalf("bookmark store id was not restored: got %q", got.BookMarks[0].BookStoreID)
	}
	if got.BookMarks[0].CreatedAt.IsZero() || got.BookMarks[0].UpdatedAt.IsZero() {
		t.Fatalf("bookmark times were not restored: %#v", got.BookMarks[0])
	}
}

// 验证保存空页面列表会清掉数据库中的旧页面记录。
func TestStoreBookWithEmptyPageInfosClearsOldRows(t *testing.T) {
	db, store := newTestStoreDatabase(t)
	book := &model.Book{
		BookInfo: model.BookInfo{
			BookID:   "book-clear-pages",
			Title:    "Book Clear Pages",
			Type:     model.TypeZip,
			BookPath: "/library/book-clear-pages.zip",
			StoreUrl: "/library",
		},
		PageInfos: model.PageInfos{{Name: "001.jpg", PageNum: 1}},
	}
	if err := store.StoreBook(book); err != nil {
		t.Fatalf("store book with pages: %v", err)
	}
	book.PageInfos = nil
	if err := store.StoreBook(book); err != nil {
		t.Fatalf("store book without pages: %v", err)
	}
	count, err := New(db).CountPageInfosByBookID(context.Background(), "book-clear-pages")
	if err != nil {
		t.Fatalf("count page infos: %v", err)
	}
	if count != 0 {
		t.Fatalf("old page infos were not cleared: got %d", count)
	}
}

// 验证生成书籍组会覆盖所有书库，并忽略旧的自动分组。
func TestGenerateBookGroupProcessesAllStoresAndIgnoresOldGroups(t *testing.T) {
	_, store := newTestStoreDatabase(t)
	root := t.TempDir()
	storeA := filepath.Join(root, "store-a")
	storeB := filepath.Join(root, "store-b")
	seriesA := filepath.Join(storeA, "series")
	seriesB := filepath.Join(storeB, "series")
	for _, dir := range []string{seriesA, seriesB} {
		if err := mkdirAllForTest(dir); err != nil {
			t.Fatalf("create test directory %s: %v", dir, err)
		}
	}

	books := []*model.Book{
		testBook("a1", filepath.Join(seriesA, "a1.zip"), storeA, 1),
		testBook("a2", filepath.Join(seriesA, "a2.zip"), storeA, 1),
		testBook("b1", filepath.Join(seriesB, "b1.zip"), storeB, 1),
		testBook("b2", filepath.Join(seriesB, "b2.zip"), storeB, 1),
		{
			BookInfo: model.BookInfo{
				BookID:        "old-group",
				Title:         "old",
				Type:          model.TypeBooksGroup,
				BookPath:      seriesA,
				StoreUrl:      storeA,
				Depth:         0,
				ChildBooksID:  []string{"a1", "a2"},
				ChildBooksNum: 2,
			},
		},
	}
	for _, book := range books {
		if err := store.StoreBook(book); err != nil {
			t.Fatalf("store %s: %v", book.BookID, err)
		}
	}
	if err := store.GenerateBookGroup(); err != nil {
		t.Fatalf("generate book groups: %v", err)
	}

	allBooks, err := store.ListBooks()
	if err != nil {
		t.Fatalf("list books: %v", err)
	}
	groupsByStore := map[string]int{}
	for _, book := range allBooks {
		if book.Type != model.TypeBooksGroup {
			continue
		}
		groupsByStore[book.StoreUrl]++
		if book.Title == "old" {
			t.Fatalf("old group topology was not refreshed")
		}
		if book.ChildBooksNum != 2 {
			t.Fatalf("generated group included stale children: %#v", book)
		}
	}
	if groupsByStore[storeA] != 1 || groupsByStore[storeB] != 1 {
		t.Fatalf("groups were not generated for all stores: %#v", groupsByStore)
	}
}

func testBook(id string, path string, storeURL string, depth int) *model.Book {
	return &model.Book{
		BookInfo: model.BookInfo{
			BookID:       id,
			Title:        filepath.Base(path),
			Type:         model.TypeZip,
			BookPath:     path,
			StoreUrl:     storeURL,
			Depth:        depth,
			ParentFolder: filepath.Base(filepath.Dir(path)),
			Modified:     time.Now(),
		},
	}
}

func mkdirAllForTest(path string) error {
	return os.MkdirAll(path, 0o755)
}

// 验证完整字段、EPUB 顺序、失败回滚和并发书签，防止数据库与内存功能分叉。
func TestSQLiteCompleteMetadataAndAtomicWrites(t *testing.T) {
	db, store := newTestStoreDatabase(t)
	now := time.Date(2026, 5, 1, 2, 3, 4, 567890123, time.UTC)
	page := model.PageInfo{Name: "cover.png", Path: "internal/cover.png", Size: 42, ModTime: now, Url: "/cover", PageNum: 7, Blurhash: "hash", Height: 30, Width: 20, ImgType: "png", InsertHtml: "<p>chapter</p>"}
	book := &model.Book{BookInfo: model.BookInfo{
		BookID: "full", Title: "完整元数据", Author: "作者", Type: model.TypeEpub, BookPath: "internal/book.epub", ParentFolder: "internal", StoreUrl: "https://example.test/library",
		IsRemote: true, RemoteURL: "https://example.test/library", RemoteBookID: "original", RemoteStoreKey: "store", RemoteShelfKey: "shelf", RemoteShelfName: "书架",
		FileSize: 42, Modified: now, PageCount: 2, Cover: page, ISBN: "isbn", Press: "出版社", PublishedAt: "2026", ChildBooksNum: 2, ChildBooksID: []string{"one", "two"}, Depth: 3,
		ExtractPath: "cache/book", ExtractNum: 2, NonUTF8Zip: true, ZipTextEncoding: "ShiftJIS", InitComplete: true, BookComplete: true, CreatedByVersion: "v1.0.0",
	}, PageInfos: model.PageInfos{page, {Name: "chapter.png", PageNum: 1, ModTime: now}}, BookMarks: model.BookMarks{{Type: model.AutoMark, BookID: "full", BookStoreID: "store", PageIndex: 1, Description: "进度", CreatedAt: now, UpdatedAt: now}}}
	check := func(want *model.Book) {
		t.Helper()
		got, err := store.GetBook(want.BookID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("round trip mismatch:\ngot  %#v\nwant %#v", got, want)
		}
		list, err := store.ListBooks()
		if err != nil || len(list) != 1 || !reflect.DeepEqual(list[0], want) {
			t.Fatalf("list mismatch: %#v, %v", list, err)
		}
	}
	if err := store.StoreBook(book); err != nil {
		t.Fatal(err)
	}
	check(book)
	book.Title = "更新后的标题"
	book.Modified = now.Add(-time.Hour)
	if err := store.StoreBook(book); err != nil {
		t.Fatal(err)
	}
	check(book)
	// 在最后一步注入写入失败，确认书籍、页面和旧书签都没有被部分覆盖。
	if _, err := db.Exec("CREATE TRIGGER reject_bookmark BEFORE INSERT ON bookmarks WHEN NEW.description = 'reject' BEGIN SELECT RAISE(ABORT, 'test failure'); END"); err != nil {
		t.Fatal(err)
	}
	broken := book.CloneForView()
	broken.Title = "不应保存"
	broken.PageInfos = nil
	broken.BookMarks[0].Description = "reject"
	if err := store.StoreBook(broken); err == nil {
		t.Fatal("expected failed transaction")
	}
	check(book)
	var pageID int64
	if err := db.QueryRow("SELECT min(id) FROM page_infos").Scan(&pageID); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 1; i <= 16; i++ {
		wg.Go(func() {
			if err := store.StoreBookMark(model.NewBookMark(model.UserMark, book.BookID, "store", i, "并发书签")); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	var updatedPageID int64
	if err := db.QueryRow("SELECT min(id) FROM page_infos").Scan(&updatedPageID); err != nil || updatedPageID != pageID {
		t.Fatalf("bookmark update rewrote page data: %d, %d, %v", pageID, updatedPageID, err)
	}
	marks, err := store.GetBookMarks(book.BookID)
	if err != nil || len(*marks) != 17 {
		t.Fatalf("lost concurrent bookmarks: %v, %v", marks, err)
	}
	if err := store.DeleteBook(book.BookID); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"books", "page_infos", "bookmarks"} {
		var count int
		if err := db.QueryRow("SELECT count(*) FROM " + table).Scan(&count); err != nil || count != 0 {
			t.Fatalf("orphaned %s: %d, %v", table, count, err)
		}
	}
}

// 重建书组的计算、写入任一阶段失败，都不能丢失原拓扑；成功重建保持 ID。
func TestGenerateBookGroupRollback(t *testing.T) {
	db, store := newTestStoreDatabase(t)
	root := t.TempDir()
	dir := filepath.Join(root, "series")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	book := testBook("child", filepath.Join(dir, "book.zip"), root, 1)
	if err := store.StoreBook(book); err != nil {
		t.Fatal(err)
	}
	if err := store.GenerateBookGroup(); err != nil {
		t.Fatal(err)
	}
	var group *model.Book
	books, err := store.ListBooks()
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range books {
		if b.Type == model.TypeBooksGroup {
			group = b
		}
	}
	if group == nil {
		t.Fatal("missing group")
	}
	check := func() {
		t.Helper()
		got, err := store.GetBook(group.BookID)
		if err != nil || !reflect.DeepEqual(got, group) {
			t.Fatalf("group changed after failure: %#v, %v", got, err)
		}
	}
	if _, err := db.Exec("CREATE TRIGGER reject_group BEFORE INSERT ON books WHEN NEW.type = '" + string(model.TypeBooksGroup) + "' BEGIN SELECT RAISE(ABORT, 'test failure'); END"); err != nil {
		t.Fatal(err)
	}
	if err := store.GenerateBookGroup(); err == nil {
		t.Fatal("expected write failure")
	}
	check()
	if _, err := db.Exec("DROP TRIGGER reject_group"); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(dir); err != nil {
		t.Fatal(err)
	}
	if err := store.GenerateBookGroup(); err == nil {
		t.Fatal("expected filesystem failure")
	}
	check()
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := store.GenerateBookGroup(); err != nil {
		t.Fatal(err)
	}
	got, err := store.GetBook(group.BookID)
	if err != nil || !reflect.DeepEqual(got.ChildBooksID, group.ChildBooksID) {
		t.Fatalf("group ID or children changed: %#v, %v", got, err)
	}
}
