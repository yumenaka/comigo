//go:build !js

package system_tray

import (
	"testing"
	"time"

	"fyne.io/systray"
)

// 验证切换托盘库后，菜单点击仍能连续调用原有操作。
func TestMenuClickAction(t *testing.T) {
	item := &systray.MenuItem{ClickedCh: make(chan struct{}, 1)}
	defer close(item.ClickedCh)
	called := make(chan struct{}, 1)
	onMenuClick(item, func() { called <- struct{}{} })
	for range 2 {
		item.ClickedCh <- struct{}{}
		select {
		case <-called:
		case <-time.After(time.Second):
			t.Fatal("菜单操作没有执行")
		}
	}
}
