//go:build wails && !js && !bindings

package wails_systray

import (
	"context"
	"testing"
)

// TestHandleBeforeCloseWithoutTray 确认无托盘环境仍使用 Wails 默认关闭行为。
func TestHandleBeforeCloseWithoutTray(t *testing.T) {
	if (&Tray{}).HandleBeforeClose(context.Background()) {
		t.Fatal("disabled tray should not intercept window close")
	}
}
