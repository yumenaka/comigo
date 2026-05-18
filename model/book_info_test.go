package model

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/tools"
)

type bookInfoTestStore struct {
	books     []*Book
	bookMarks map[string]BookMarks
	listCalls int
}

func (s *bookInfoTestStore) StoreBook(b *Book) error {
	s.books = append(s.books, b)
	return nil
}

func (s *bookInfoTestStore) GetBook(id string) (*Book, error) {
	for _, book := range s.books {
		if book.BookID == id {
			return book, nil
		}
	}
	return nil, errors.New("book not found")
}

func (s *bookInfoTestStore) DeleteBook(id string) error { return nil }

func (s *bookInfoTestStore) ListBooks() ([]*Book, error) {
	s.listCalls++
	return s.books, nil
}

func (s *bookInfoTestStore) GenerateBookGroup() error { return nil }

func (s *bookInfoTestStore) StoreBookMark(mark *BookMark) error { return nil }

func (s *bookInfoTestStore) GetBookMarks(bookID string) (*BookMarks, error) {
	if marks, ok := s.bookMarks[bookID]; ok {
		return &marks, nil
	}
	bookMarks := BookMarks{}
	return &bookMarks, nil
}

func (s *bookInfoTestStore) DeleteBookMark(bookID string, markType MarkType, pageIndex int) error {
	return nil
}

func withBookInfoTestStore(t *testing.T, store *bookInfoTestStore) {
	t.Helper()
	originalStore := IStore
	IStore = store
	t.Cleanup(func() {
		IStore = originalStore
	})
}

func TestInitBookIDExpandsShortIDConflictFromSingleSnapshot(t *testing.T) {
	store := &bookInfoTestStore{}
	withBookInfoTestStore(t, store)

	const (
		bookPath = "fixtures/series/book.cbz"
		storeURL = "fixtures/library"
	)
	modified := time.Unix(1700000000, 0)
	fileSize := int64(12345)

	baseline, err := NewBook(bookPath, modified, fileSize, storeURL, 1, TypeCbz)
	if err != nil {
		t.Fatalf("baseline NewBook failed: %v", err)
	}
	idSource := baseline.BookPath + strconv.Itoa(int(baseline.FileSize)) + string(baseline.Type) + baseline.ParentFolder + baseline.StoreUrl
	fullID := base62.EncodeToString([]byte(tools.Md5string(tools.Md5string(idSource))))
	const minLength = 7
	if len(fullID) <= minLength {
		t.Fatalf("test setup expected fullID longer than %d, got %q", minLength, fullID)
	}

	// 只制造 BookID 前 7 位冲突，BookPath 故意不同，避免触发“同路径同类型已存在”的分支。
	store.books = []*Book{{
		BookInfo: BookInfo{
			BookID:   fullID[:minLength],
			BookPath: "fixtures/other/book.cbz",
			Type:     TypeCbz,
		},
	}}
	store.listCalls = 0

	got, err := NewBook(bookPath, modified, fileSize, storeURL, 1, TypeCbz)
	if err != nil {
		t.Fatalf("NewBook with short ID conflict failed: %v", err)
	}
	wantID := fullID[:minLength+1]
	if got.BookID != wantID {
		t.Fatalf("BookID = %q, want %q", got.BookID, wantID)
	}
	if store.listCalls != 1 {
		t.Fatalf("ListBooks calls = %d, want 1", store.listCalls)
	}
}

func TestInitBookIDRejectsExistingSamePathAndType(t *testing.T) {
	store := &bookInfoTestStore{}
	withBookInfoTestStore(t, store)

	const bookPath = "fixtures/existing/book.cbz"
	absBookPath, err := filepath.Abs(bookPath)
	if err != nil {
		t.Fatalf("filepath.Abs failed: %v", err)
	}
	store.books = []*Book{{
		BookInfo: BookInfo{
			BookID:   "existing-book",
			BookPath: absBookPath,
			Type:     TypeCbz,
		},
	}}

	_, err = NewBook(bookPath, time.Unix(1700000000, 0), 1, "fixtures/library", 1, TypeCbz)
	if err == nil {
		t.Fatal("NewBook expected duplicate error, got nil")
	}
	if !strings.Contains(err.Error(), "existing-book") {
		t.Fatalf("duplicate error = %q, want existing BookID", err.Error())
	}
	if store.listCalls != 1 {
		t.Fatalf("ListBooks calls = %d, want 1", store.listCalls)
	}
}

func TestSortBooksByLastReadFallsBackToModifiedTime(t *testing.T) {
	store := &bookInfoTestStore{
		bookMarks: map[string]BookMarks{
			"read-old": {{
				Type:      AutoMark,
				BookID:    "read-old",
				UpdatedAt: time.Unix(1700000200, 0),
			}},
			"read-new": {{
				Type:      AutoMark,
				BookID:    "read-new",
				UpdatedAt: time.Unix(1700000300, 0),
			}},
			"user-only": {{
				Type:      UserMark,
				BookID:    "user-only",
				UpdatedAt: time.Unix(1700000400, 0),
			}},
		},
	}
	withBookInfoTestStore(t, store)

	books := BookInfos{
		{BookID: "unread-old", Title: "unread-old", Modified: time.Unix(1700000100, 0)},
		{BookID: "read-old", Title: "read-old", Modified: time.Unix(1700000500, 0)},
		{BookID: "user-only", Title: "user-only", Modified: time.Unix(1700000700, 0)},
		{BookID: "read-new", Title: "read-new", Modified: time.Unix(1700000000, 0)},
		{BookID: "unread-new", Title: "unread-new", Modified: time.Unix(1700000600, 0)},
	}

	books.SortBooks("last_read")

	got := make([]string, 0, len(books))
	for _, book := range books {
		got = append(got, book.BookID)
	}
	want := []string{"read-new", "read-old", "user-only", "unread-new", "unread-old"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("last_read order = %v, want %v", got, want)
	}
}
