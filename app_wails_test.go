//go:build wails && !js

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yumenaka/comigo/model"
)

func TestTrashableBookPathAllowsOnlyLocalStoreBook(t *testing.T) {
	storeDir := t.TempDir()
	bookPath := filepath.Join(storeDir, "book.cbz")
	if err := os.WriteFile(bookPath, []byte("book"), 0o644); err != nil {
		t.Fatal(err)
	}
	book := &model.Book{BookInfo: model.BookInfo{
		BookID:   "book",
		BookPath: bookPath,
		Type:     model.TypeCbz,
	}}

	got, isDir, err := trashableBookPath(book, []string{storeDir})
	if err != nil {
		t.Fatal(err)
	}
	if got != bookPath || isDir {
		t.Fatalf("trashableBookPath = %q, %v; want %q, false", got, isDir, bookPath)
	}

	book.IsRemote = true
	if _, _, err := trashableBookPath(book, []string{storeDir}); err == nil {
		t.Fatal("remote book should not be trashable")
	}
}
