//go:build wails && !js

package wails_systray

import "testing"

// 无托盘环境必须保留正常关闭，避免窗口隐藏后无入口恢复。
func TestCloseWithoutTray(t *testing.T) {
	if (&Tray{}).HandleBeforeClose() {
		t.Fatal("无托盘不应拦截关闭")
	}
}
