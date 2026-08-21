//go:build !wails

package config

import "testing"

// 验证托盘入口使用独立配置文件。
func TestTrayConfigProfileUsesExplicitRuntimeMarker(t *testing.T) {
	oldProfile := runtimeConfigProfile
	t.Cleanup(func() {
		runtimeConfigProfile = oldProfile
	})

	runtimeConfigProfile = "cli"
	UseTrayConfigProfile()
	if got := PlatformConfigFilename(); got != trayConfigFilename {
		t.Fatalf("tray 入口标记后配置文件名不正确: got %q want %q", got, trayConfigFilename)
	}
}
