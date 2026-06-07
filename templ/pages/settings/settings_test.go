package settings

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
)

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
