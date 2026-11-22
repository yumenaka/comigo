package cmd

import (
	"os"

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

// TestExecuteReturnsStartupError 启动配置失败必须向调用方返回，不能直接结束进程。
func TestExecuteReturnsStartupError(t *testing.T) {
	root := processTestCommand(t)
	file := filepath.Join(t.TempDir(), "invalid.toml")
	if err := os.WriteFile(file, []byte("Port = ["), 0600); err != nil {
		t.Fatal(err)
	}
	root.SetArgs([]string{"--config", file})
	if err := Execute(); err == nil {
		t.Fatal("无效配置没有向调用方返回错误")
	}
}
