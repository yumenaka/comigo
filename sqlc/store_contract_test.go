package sqlc

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
)

// 验证内存书库实现满足统一书库接口契约。
func TestStoreInterfaceContractRamStore(t *testing.T) {
	runStoreInterfaceContract(t, func(t *testing.T) model.StoreInterface {
		t.Helper()
		oldCfg := config.CopyCfg()
		oldStore := model.IStore
		tmp := t.TempDir()
		config.GetCfg().ConfigFile = filepath.Join(tmp, "config.toml")
		ramStore := &store.StoreInRam{
			ChildStores:      sync.Map{},
			PendingBookmarks: sync.Map{},
		}
		model.IStore = ramStore
		t.Cleanup(func() {
			*config.GetCfg() = oldCfg
			model.IStore = oldStore
		})
		return ramStore
	})
}

// 验证 SQLite 书库实现满足统一书库接口契约。
func TestStoreInterfaceContractSQLite(t *testing.T) {
	runStoreInterfaceContract(t, func(t *testing.T) model.StoreInterface {
		t.Helper()
		_, store := newTestStoreDatabase(t)
		return store
	})
}

// 验证 PostgreSQL 书库实现满足统一书库接口契约。
func TestStoreInterfaceContractPostgres(t *testing.T) {
	dsn := os.Getenv("COMIGO_TEST_PG_DSN")
	if dsn == "" {
		t.Skip("COMIGO_TEST_PG_DSN is not set")
	}
	runStoreInterfaceContract(t, func(t *testing.T) model.StoreInterface {
		t.Helper()
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			t.Fatalf("open postgres database: %v", err)
		}
		ctx := context.Background()
		if _, err := db.ExecContext(ctx, postgresDDL); err != nil {
			t.Fatalf("create postgres schema: %v", err)
		}
		if err := migratePostgresDatabase(ctx, db); err != nil {
			t.Fatalf("migrate postgres schema: %v", err)
		}
		if _, err := db.ExecContext(ctx, "TRUNCATE bookmarks, page_infos, books, users RESTART IDENTITY CASCADE"); err != nil {
			t.Fatalf("clear postgres test tables: %v", err)
		}
		store := NewPostgresDBStore(db)

		oldStore := DbStore
		oldModelStore := model.IStore
		DbStore = store
		model.IStore = store
		t.Cleanup(func() {
			DbStore = oldStore
			model.IStore = oldModelStore
			_ = db.Close()
		})
		return store
	})
}

func runStoreInterfaceContract(t *testing.T, newStore func(t *testing.T) model.StoreInterface) {
	t.Helper()
	store := newStore(t)
	root := t.TempDir()
	storeURL := filepath.Join(root, "library")
	series := filepath.Join(storeURL, "series")
	if err := os.MkdirAll(series, 0o755); err != nil {
		t.Fatalf("create series dir: %v", err)
	}

	now := time.Date(2026, 6, 2, 10, 30, 0, 0, time.UTC)
	book := &model.Book{
		BookInfo: model.BookInfo{
			BookID:           "contract-book-one",
			Title:            "Contract Book One",
			Author:           "Author Name",
			Type:             model.TypeZip,
			BookPath:         filepath.Join(series, "one.zip"),
			StoreUrl:         storeURL,
			ParentFolder:     "series",
			FileSize:         1234,
			Modified:         now,
			PageCount:        2,
			Depth:            1,
			CreatedByVersion: "v1.2.3",
		},
		PageInfos: model.PageInfos{
			{Name: "001.jpg", Path: filepath.Join(series, "001.jpg"), PageNum: 1},
			{Name: "002.jpg", Path: filepath.Join(series, "002.jpg"), PageNum: 2},
		},
		BookMarks: model.BookMarks{
			{
				Type:        model.UserMark,
				BookID:      "contract-book-one",
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
	got, err := store.GetBook(book.BookID)
	if err != nil {
		t.Fatalf("get book: %v", err)
	}
	if got.Author != book.Author || got.CreatedByVersion != book.CreatedByVersion {
		t.Fatalf("book fields were not restored: got author=%q version=%q", got.Author, got.CreatedByVersion)
	}
	if len(got.PageInfos) != 2 {
		t.Fatalf("page infos were not restored: %#v", got.PageInfos)
	}
	if len(got.BookMarks) != 1 || got.BookMarks[0].BookStoreID != "store-id" {
		t.Fatalf("bookmarks were not restored: %#v", got.BookMarks)
	}

	book.PageInfos = nil
	if err := store.StoreBook(book); err != nil {
		t.Fatalf("store book without pages: %v", err)
	}
	got, err = store.GetBook(book.BookID)
	if err != nil {
		t.Fatalf("get book after clearing pages: %v", err)
	}
	if len(got.PageInfos) != 0 {
		t.Fatalf("old page infos were not cleared: %#v", got.PageInfos)
	}

	if err := store.StoreBookMark(&model.BookMark{Type: model.UserMark, BookID: book.BookID, PageIndex: 3, Description: "first"}); err != nil {
		t.Fatalf("store user bookmark: %v", err)
	}
	if err := store.StoreBookMark(&model.BookMark{Type: model.UserMark, BookID: book.BookID, PageIndex: 3, Description: "updated"}); err != nil {
		t.Fatalf("update user bookmark: %v", err)
	}
	if err := store.StoreBookMark(&model.BookMark{Type: model.UserMark, BookID: book.BookID, PageIndex: 4, Description: "second"}); err != nil {
		t.Fatalf("append user bookmark: %v", err)
	}
	if err := store.StoreBookMark(&model.BookMark{Type: model.AutoMark, BookID: book.BookID, PageIndex: 5}); err != nil {
		t.Fatalf("store auto bookmark: %v", err)
	}
	if err := store.StoreBookMark(&model.BookMark{Type: model.AutoMark, BookID: book.BookID, PageIndex: 6}); err != nil {
		t.Fatalf("update auto bookmark: %v", err)
	}
	marks, err := store.GetBookMarks(book.BookID)
	if err != nil {
		t.Fatalf("get bookmarks: %v", err)
	}
	assertBookmarkShape(t, *marks)
	if err := store.DeleteBookMark(book.BookID, model.UserMark, 3); err != nil {
		t.Fatalf("delete user bookmark: %v", err)
	}
	marks, err = store.GetBookMarks(book.BookID)
	if err != nil {
		t.Fatalf("get bookmarks after delete: %v", err)
	}
	for _, mark := range *marks {
		if mark.Type == model.UserMark && mark.PageIndex == 3 {
			t.Fatalf("user bookmark page 3 was not deleted: %#v", *marks)
		}
	}

	sibling := &model.Book{
		BookInfo: model.BookInfo{
			BookID:       "contract-book-two",
			Title:        "Contract Book Two",
			Type:         model.TypeZip,
			BookPath:     filepath.Join(series, "two.zip"),
			StoreUrl:     storeURL,
			ParentFolder: "series",
			Depth:        1,
			Modified:     now,
		},
	}
	if err := store.StoreBook(sibling); err != nil {
		t.Fatalf("store sibling book: %v", err)
	}
	if err := store.GenerateBookGroup(); err != nil {
		t.Fatalf("generate book group: %v", err)
	}
	allBooks, err := store.ListBooks()
	if err != nil {
		t.Fatalf("list books: %v", err)
	}
	groupCount := 0
	for _, listed := range allBooks {
		if listed.Type == model.TypeBooksGroup && listed.StoreUrl == storeURL {
			groupCount++
			if listed.ChildBooksNum != 2 {
				t.Fatalf("group child count mismatch: %#v", listed)
			}
		}
	}
	if groupCount != 1 {
		t.Fatalf("generated group count = %d, want 1", groupCount)
	}

	if err := store.DeleteBook(book.BookID); err != nil {
		t.Fatalf("delete book: %v", err)
	}
	if _, err := store.GetBook(book.BookID); err == nil {
		t.Fatalf("deleted book was still readable")
	}
}

func assertBookmarkShape(t *testing.T, marks model.BookMarks) {
	t.Helper()
	userByPage := map[int]string{}
	autoCount := 0
	autoPage := -1
	for _, mark := range marks {
		switch mark.Type {
		case model.UserMark:
			userByPage[mark.PageIndex] = mark.Description
		case model.AutoMark:
			autoCount++
			autoPage = mark.PageIndex
		}
	}
	if userByPage[3] != "updated" || userByPage[4] != "second" {
		t.Fatalf("user bookmarks mismatch: %#v", marks)
	}
	if autoCount != 1 || autoPage != 6 {
		t.Fatalf("auto bookmark mismatch: %#v", marks)
	}
}
