//go:build wails && !js && !bindings

package wails_systray

import (
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Tray 封装 Wails 桌面壳托盘逻辑。
type Tray struct {
	mu       sync.RWMutex
	app      *application.App
	window   application.Window
	enabled  bool
	end      func()
	quitting atomic.Bool
	stopOnce sync.Once
}

// Start 启动 systray external loop；不支持托盘的环境返回空控制器。
func Start() *Tray {
	t := &Tray{enabled: traySupported()}
	if !t.enabled {
		return t
	}
	t.end = startPlatform(t)
	return t
}

// SetRuntime 保存 Wails 桌面壳对象，供托盘点击恢复窗口或退出。
func (t *Tray) SetRuntime(app *application.App, window application.Window) {
	if t == nil {
		return
	}
	t.mu.Lock()
	t.app = app
	t.window = window
	t.mu.Unlock()
}

// RegisterCloseHook 接管关闭按钮：有托盘时隐藏窗口并阻止 Wails 默认退出。
func (t *Tray) RegisterCloseHook(window application.Window) {
	if t == nil || !t.enabled || window == nil {
		return
	}
	window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if t.HandleBeforeClose() {
			event.Cancel()
		}
	})
}

// Stop 停止托盘；由 Wails OnShutdown 调用。
func (t *Tray) Stop() {
	if t == nil || !t.enabled {
		return
	}
	t.stopOnce.Do(func() {
		t.SetRuntime(nil, nil)
		if t.end != nil {
			t.end()
		}
	})
}

// HideOnClose 仅在托盘可恢复窗口时让关闭按钮隐藏窗口。
func (t *Tray) HideOnClose() bool {
	return t != nil && t.enabled
}

// HandleBeforeClose 接管关闭按钮：有托盘时隐藏窗口并阻止 Wails 默认退出。
func (t *Tray) HandleBeforeClose() bool {
	if t == nil || !t.enabled || t.quitting.Load() {
		return false
	}
	t.hideWindow()
	return true
}

// runtime 读取当前 Wails app/window。
func (t *Tray) runtime() (*application.App, application.Window) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.app, t.window
}

// showWindow 恢复已隐藏或最小化的 Wails 窗口。
func (t *Tray) showWindow() {
	_, window := t.runtime()
	if window == nil {
		return
	}
	setPlatformWindowVisible(true)
	window.Show()
	window.UnMinimise()
}

// hideWindow 把窗口收进托盘，并同步隐藏 macOS Dock 图标。
func (t *Tray) hideWindow() {
	_, window := t.runtime()
	if window == nil {
		return
	}
	window.Hide()
	setPlatformWindowVisible(false)
}

// quit 通过 Wails 退出应用，尚未启动完成时退回到 systray 退出。
func (t *Tray) quit() {
	app, _ := t.runtime()
	if app == nil {
		quitPlatformFallback()
		return
	}
	t.quitting.Store(true)
	setPlatformWindowVisible(true)
	app.Quit()
}
