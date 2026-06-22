//go:build wails && !js

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers"
)

// 验证 Wails 删除源文件时只允许处理本地书库内的普通书籍文件。
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

	got, isDir, err := routers.TrashableBookPathForWails(book, []string{storeDir})
	if err != nil {
		t.Fatal(err)
	}
	if got != bookPath || isDir {
		t.Fatalf("trashableBookPath = %q, %v; want %q, false", got, isDir, bookPath)
	}

	book.IsRemote = true
	if _, _, err := routers.TrashableBookPathForWails(book, []string{storeDir}); err == nil {
		t.Fatal("remote book should not be trashable")
	}

	dirBook := &model.Book{BookInfo: model.BookInfo{
		BookID:   "dir",
		BookPath: storeDir,
		Type:     model.TypeDir,
	}}
	if _, _, err := routers.TrashableBookPathForWails(dirBook, []string{storeDir}); err == nil {
		t.Fatal("directory book should not be trashable")
	}
}
