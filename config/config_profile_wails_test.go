//go:build wails

package config

import "testing"

// 验证 Wails 构建本身不等于桌面运行环境。
func TestDesktopConfigProfileUsesExplicitRuntimeMarker(t *testing.T) {
	oldProfile := runtimeConfigProfile
	t.Cleanup(func() {
		runtimeConfigProfile = oldProfile
	})

	runtimeConfigProfile = "cli"
	if got := PlatformConfigFilename(); got != cliConfigFilename {
		t.Fatalf("未进入 Wails 桌面环境时不应使用 desktop 配置: got %q want %q", got, cliConfigFilename)
	}

	UseTrayConfigProfile()
	if got := PlatformConfigFilename(); got != trayConfigFilename {
		t.Fatalf("tray 入口标记后配置文件名不正确: got %q want %q", got, trayConfigFilename)
	}

	UseDesktopConfigProfile()
	if got := PlatformConfigFilename(); got != desktopConfigFilename {
		t.Fatalf("Wails 桌面环境标记后配置文件名不正确: got %q want %q", got, desktopConfigFilename)
	}
}
