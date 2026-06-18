package model

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type bookCleanupTestStore struct {
	books         map[string]*Book
	storeCalls    int
	deleteCalls   int
	generateCalls int
}

func (s *bookCleanupTestStore) StoreBook(b *Book) error {
	s.storeCalls++
	s.books[b.BookID] = b
	return nil
}

func (s *bookCleanupTestStore) GetBook(id string) (*Book, error) {
	if b, ok := s.books[id]; ok {
		return b, nil
	}
	return nil, errors.New("book not found")
}

func (s *bookCleanupTestStore) DeleteBook(id string) error {
	s.deleteCalls++
	delete(s.books, id)
	return nil
}

func (s *bookCleanupTestStore) ListBooks() ([]*Book, error) {
	books := make([]*Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	return books, nil
}

func (s *bookCleanupTestStore) GenerateBookGroup() error {
	s.generateCalls++
	return nil
}

func (s *bookCleanupTestStore) StoreBookMark(mark *BookMark) error { return nil }

func (s *bookCleanupTestStore) GetBookMarks(bookID string) (*BookMarks, error) {
	marks := BookMarks{}
	return &marks, nil
}

func (s *bookCleanupTestStore) DeleteBookMark(bookID string, markType MarkType, pageIndex int) error {
	return nil
}

func withBookCleanupTestStore(t *testing.T, store *bookCleanupTestStore) {
	t.Helper()
	originalStore := IStore
	IStore = store
	t.Cleanup(func() {
		IStore = originalStore
	})
}

func TestClearBookNotExistRemovesMissingTypeDirPages(t *testing.T) {
	dir := t.TempDir()
	existingPath := filepath.Join(dir, "001.jpg")
	if err := os.WriteFile(existingPath, []byte("image"), 0o644); err != nil {
		t.Fatalf("write fixture image: %v", err)
	}
	missingPath := filepath.Join(dir, "002.jpg")
	book := &Book{
		BookInfo: BookInfo{
			BookID:   "dir-book",
			BookPath: dir,
			StoreUrl: dir,
			Type:     TypeDir,
			Modified: time.Unix(1700000000, 0),
		},
		PageInfos: PageInfos{
			{Name: "001.jpg", Path: existingPath, PageNum: 1},
			{Name: "002.jpg", Path: missingPath, PageNum: 2},
		},
	}
	store := &bookCleanupTestStore{books: map[string]*Book{book.BookID: book}}
	withBookCleanupTestStore(t, store)

	ClearBookNotExist()

	got := store.books[book.BookID]
	if got == nil {
		t.Fatal("book should remain when at least one directory page exists")
	}
	if len(got.PageInfos) != 1 || got.PageInfos[0].Name != "001.jpg" {
		t.Fatalf("PageInfos = %#v, want only existing page", got.PageInfos)
	}
	if got.PageCount != 1 {
		t.Fatalf("PageCount = %d, want 1", got.PageCount)
	}
	if got.Cover.Name != "001.jpg" {
		t.Fatalf("Cover.Name = %q, want existing page", got.Cover.Name)
	}
	if store.storeCalls != 1 {
		t.Fatalf("StoreBook calls = %d, want 1", store.storeCalls)
	}
	if store.deleteCalls != 0 {
		t.Fatalf("DeleteBook calls = %d, want 0", store.deleteCalls)
	}
	if store.generateCalls != 0 {
		t.Fatalf("GenerateBookGroup calls = %d, want 0", store.generateCalls)
	}
}

func TestClearBookNotExistDeletesTypeDirWhenAllPagesMissing(t *testing.T) {
	dir := t.TempDir()
	book := &Book{
		BookInfo: BookInfo{
			BookID:   "empty-dir-book",
			BookPath: dir,
			StoreUrl: dir,
			Type:     TypeDir,
			Modified: time.Unix(1700000000, 0),
		},
		PageInfos: PageInfos{
			{Name: "001.jpg", Path: filepath.Join(dir, "001.jpg"), PageNum: 1},
			{Name: "002.jpg", Path: filepath.Join(dir, "002.jpg"), PageNum: 2},
		},
	}
	store := &bookCleanupTestStore{books: map[string]*Book{book.BookID: book}}
	withBookCleanupTestStore(t, store)

	ClearBookNotExist()

	if _, ok := store.books[book.BookID]; ok {
		t.Fatal("book should be deleted when all directory pages are missing")
	}
	if store.deleteCalls != 1 {
		t.Fatalf("DeleteBook calls = %d, want 1", store.deleteCalls)
	}
	if store.storeCalls != 0 {
		t.Fatalf("StoreBook calls = %d, want 0", store.storeCalls)
	}
	if store.generateCalls != 1 {
		t.Fatalf("GenerateBookGroup calls = %d, want 1", store.generateCalls)
	}
}

func TestClearBookNotExistDeletesTypeDirWithNoPages(t *testing.T) {
	dir := t.TempDir()
	book := &Book{
		BookInfo: BookInfo{
			BookID:   "no-page-dir-book",
			BookPath: dir,
			StoreUrl: dir,
			Type:     TypeDir,
			Modified: time.Unix(1700000000, 0),
		},
	}
	store := &bookCleanupTestStore{books: map[string]*Book{book.BookID: book}}
	withBookCleanupTestStore(t, store)

	ClearBookNotExist()

	if _, ok := store.books[book.BookID]; ok {
		t.Fatal("book should be deleted when directory metadata has no pages")
	}
	if store.deleteCalls != 1 {
		t.Fatalf("DeleteBook calls = %d, want 1", store.deleteCalls)
	}
}

func TestCloneForViewSortDoesNotMutateOriginalPages(t *testing.T) {
	book := &Book{
		PageInfos: PageInfos{
			{Name: "002.jpg"},
			{Name: "001.jpg"},
		},
	}

	clone := book.CloneForView()
	clone.SortPages("filename")

	if clone.PageInfos[0].Name != "001.jpg" {
		t.Fatalf("clone first page = %q, want sorted page", clone.PageInfos[0].Name)
	}
	if book.PageInfos[0].Name != "002.jpg" {
		t.Fatalf("original first page = %q, want original order", book.PageInfos[0].Name)
	}
}

func TestSortImagesReverseOrder(t *testing.T) {
	pages := PageInfos{
		{Name: "001.jpg", Size: 1},
		{Name: "003.jpg", Size: 3},
		{Name: "002.jpg", Size: 2},
	}

	pages.SortImages("filename_reverse")
	if pages[0].Name != "003.jpg" || pages[1].Name != "002.jpg" || pages[2].Name != "001.jpg" {
		t.Fatalf("filename_reverse order = %#v", pages)
	}

	pages.SortImages("filesize_reverse")
	if pages[0].Size != 1 || pages[1].Size != 2 || pages[2].Size != 3 {
		t.Fatalf("filesize_reverse order = %#v", pages)
	}
}
