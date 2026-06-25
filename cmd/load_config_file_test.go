package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// 验证只有全部位置参数都是本地普通文件时，才进入临时阅读模式。
func TestTemporaryReaderMode(t *testing.T) {
	dir := t.TempDir()
	fileA := filepath.Join(dir, "a.cbz")
	fileB := filepath.Join(dir, "b.pdf")
	if err := os.WriteFile(fileA, []byte("a"), 0o644); err != nil {
		t.Fatalf("写入测试文件失败: %v", err)
	}
	if err := os.WriteFile(fileB, []byte("b"), 0o644); err != nil {
		t.Fatalf("写入测试文件失败: %v", err)
	}

	tests := []struct {
		name       string
		args       []string
		configFile string
		force      bool
		want       bool
	}{
		{name: "multiple files", args: []string{fileA, fileB}, want: true},
		{name: "directory mixed", args: []string{fileA, dir}, want: false},
		{name: "no args", args: nil, want: false},
		{name: "explicit config", args: []string{fileA}, configFile: "config.toml", want: false},
		{name: "missing file", args: []string{filepath.Join(dir, "missing.cbz")}, want: false},
		{name: "forced empty args", force: true, want: true},
		{name: "forced directory and config", args: []string{dir}, configFile: "config.toml", force: true, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := temporaryReaderMode(tt.args, tt.configFile, tt.force); got != tt.want {
				t.Fatalf("temporaryReaderMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
