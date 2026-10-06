package cmd

import (
	"path/filepath"
	"reflect"
	"testing"
)

// 书库路径和远程 URL 保持为根命令的位置参数。
func TestRootAcceptsLibraryArguments(t *testing.T) {
	for _, args := range [][]string{
		{filepath.Join(t.TempDir(), "Documents")},
		{"books with spaces", "https://example.com/library/"},
		{"desktop", "info"},
	} {
		command, remaining, err := RootCmd.Find(args)
		if err != nil || command != RootCmd || !reflect.DeepEqual(remaining, args) {
			t.Fatalf("书库参数 %q: command=%v remaining=%q err=%v", args, command, remaining, err)
		}
	}
}
