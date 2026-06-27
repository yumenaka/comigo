//go:build wails && !js && !bindings

package wails_systray

import "testing"

// TestHideOnCloseReflectsTraySupport 确认关闭按钮行为只跟托盘可用性绑定。
func TestHideOnCloseReflectsTraySupport(t *testing.T) {
	if (&Tray{enabled: true}).HideOnClose() != true {
		t.Fatal("enabled tray should hide window on close")
	}
	if (&Tray{enabled: false}).HideOnClose() != false {
		t.Fatal("disabled tray should not hide window on close")
	}
}
