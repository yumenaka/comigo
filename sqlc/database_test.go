package sqlc

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
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
	if _, err := db.ExecContext(context.Background(), ddl); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	if err := migrateDatabase(context.Background(), db); err != nil {
		t.Fatalf("migrate schema: %v", err)
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
	if len(got.PageInfos) != 2 || got.PageInfos[0].Name != "001.jpg" {
		t.Fatalf("page infos were not restored in default order: %#v", got.PageInfos)
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
		if book.BookID == "old-group" {
			t.Fatalf("old group was not removed before regenerating")
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
