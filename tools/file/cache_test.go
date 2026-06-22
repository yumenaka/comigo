package file

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// 验证封面缓存文件名包含缩放高度，避免不同尺寸互相覆盖。
func TestCoverCacheUsesResizeHeightInFilename(t *testing.T) {
	cacheDir := t.TempDir()
	bookID := "book123"

	if err := SaveCoverToLocal(cacheDir, bookID, 128, []byte("small")); err != nil {
		t.Fatalf("SaveCoverToLocal small error: %v", err)
	}
	if err := SaveCoverToLocal(cacheDir, bookID, 352, []byte("large")); err != nil {
		t.Fatalf("SaveCoverToLocal large error: %v", err)
	}

	small, err := GetCoverFromLocal(cacheDir, bookID, 128)
	if err != nil {
		t.Fatalf("GetCoverFromLocal small error: %v", err)
	}
	large, err := GetCoverFromLocal(cacheDir, bookID, 352)
	if err != nil {
		t.Fatalf("GetCoverFromLocal large error: %v", err)
	}
	if string(small) != "small" {
		t.Fatalf("小尺寸封面缓存内容不正确: got %q", string(small))
	}
	if string(large) != "large" {
		t.Fatalf("大尺寸封面缓存内容不正确: got %q", string(large))
	}
	if !CoverFileCacheExists(cacheDir, bookID, 128) || !CoverFileCacheExists(cacheDir, bookID, 352) {
		t.Fatalf("不同尺寸的封面缓存都应存在")
	}
}

// 验证删除封面缓存会清理同一封面的所有尺寸变体。
func TestDeleteCoverCacheRemovesAllSizeVariants(t *testing.T) {
	cacheDir := t.TempDir()
	bookID := "book123"

	files := []string{
		bookID + ".jpg",
		bookID + "_h128.jpg",
		bookID + "_h352.jpg",
		bookID + "_original.jpg",
		"book1234_h128.jpg",
	}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(cacheDir, name), []byte(name), 0o644); err != nil {
			t.Fatalf("write test cache %s: %v", name, err)
		}
	}

	if err := DeleteCoverCache(cacheDir, bookID); err != nil {
		t.Fatalf("DeleteCoverCache error: %v", err)
	}

	removed := []string{
		bookID + ".jpg",
		bookID + "_h128.jpg",
		bookID + "_h352.jpg",
		bookID + "_original.jpg",
	}
	for _, name := range removed {
		if _, err := os.Stat(filepath.Join(cacheDir, name)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("封面缓存 %s 应已删除，stat err=%v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(cacheDir, "book1234_h128.jpg")); err != nil {
		t.Fatalf("其他书籍的封面缓存不应被删除: %v", err)
	}
}
