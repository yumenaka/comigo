package file

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

// TestUnArchiveAutoRejectsPathTraversal 确认不可信归档不能写出目标目录。
func TestUnArchiveAutoRejectsPathTraversal(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "traversal.zip")
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(file)
	entry, err := zw.Create("../escaped.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte("escaped")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	extractDir := filepath.Join(root, "output")
	if err := UnArchiveAuto(archivePath, extractDir, ""); err == nil {
		t.Fatal("包含父目录跳转的归档应解压失败")
	}
	if _, err := os.Stat(filepath.Join(root, "escaped.txt")); !os.IsNotExist(err) {
		t.Fatalf("归档内容写出了目标目录: %v", err)
	}
}
