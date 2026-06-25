//go:build !wails

package config

import (
	"os"
	"testing"
)

// 验证启动壳配置类型由入口显式标记，不再从可执行文件名猜测。
func TestTrayConfigProfileUsesExplicitRuntimeMarker(t *testing.T) {
	oldProfile := runtimeConfigProfile
	oldArgs := os.Args
	t.Cleanup(func() {
		runtimeConfigProfile = oldProfile
		os.Args = oldArgs
	})

	runtimeConfigProfile = "cli"
	os.Args = []string{"comigo-tray"}
	if got := PlatformConfigFilename(); got != cliConfigFilename {
		t.Fatalf("不应通过文件名切换到 tray 配置: got %q want %q", got, cliConfigFilename)
	}

	UseTrayConfigProfile()
	if got := PlatformConfigFilename(); got != trayConfigFilename {
		t.Fatalf("tray 入口标记后配置文件名不正确: got %q want %q", got, trayConfigFilename)
	}
}
