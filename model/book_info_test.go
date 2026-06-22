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

// 验证短书籍 ID 冲突时会在同一快照内继续扩展长度。
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

// 验证相同路径和类型的书籍不会重复生成新 ID。
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

// 验证按最近阅读排序时，没有阅读记录的书籍会回退到修改时间。
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

// 验证书籍列表的反向排序选项符合预期顺序。
func TestSortBooksReverseOrders(t *testing.T) {
	byTitle := BookInfos{
		{Title: "001.zip", Type: TypeZip},
		{Title: "003.zip", Type: TypeZip},
		{Title: "002.zip", Type: TypeZip},
	}
	byTitle.SortBooks("filename_reverse")
	if byTitle[0].Title != "003.zip" || byTitle[1].Title != "002.zip" || byTitle[2].Title != "001.zip" {
		t.Fatalf("filename_reverse order = %#v", byTitle)
	}

	bySize := BookInfos{
		{Title: "large", Type: TypeZip, FileSize: 30},
		{Title: "small", Type: TypeZip, FileSize: 10},
		{Title: "middle", Type: TypeZip, FileSize: 20},
	}
	bySize.SortBooks("filesize_reverse")
	if bySize[0].Title != "small" || bySize[1].Title != "middle" || bySize[2].Title != "large" {
		t.Fatalf("filesize_reverse order = %#v", bySize)
	}

	byAuthor := BookInfos{
		{Title: "a", Author: "Ann", Type: TypeZip},
		{Title: "b", Author: "Cat", Type: TypeZip},
		{Title: "c", Author: "Bob", Type: TypeZip},
	}
	byAuthor.SortBooks("author_reverse")
	if byAuthor[0].Author != "Cat" || byAuthor[1].Author != "Bob" || byAuthor[2].Author != "Ann" {
		t.Fatalf("author_reverse order = %#v", byAuthor)
	}
}
