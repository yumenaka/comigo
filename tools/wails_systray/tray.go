//go:build wails && !js && !bindings

package wails_systray

import (
	"context"
	"sync"
	"sync/atomic"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Tray 封装 Wails v2 临时托盘逻辑，迁到 Wails v3 时可整包替换。
type Tray struct {
	mu       sync.RWMutex
	ctx      context.Context
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

// SetContext 保存 Wails runtime context，供托盘点击恢复窗口或退出。
func (t *Tray) SetContext(ctx context.Context) {
	if t == nil {
		return
	}
	t.mu.Lock()
	t.ctx = ctx
	t.mu.Unlock()
}

// Stop 停止托盘；由 Wails OnShutdown 调用。
func (t *Tray) Stop() {
	if t == nil || !t.enabled {
		return
	}
	t.stopOnce.Do(func() {
		t.SetContext(nil)
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
func (t *Tray) HandleBeforeClose(ctx context.Context) bool {
	if t == nil || !t.enabled || t.quitting.Load() {
		return false
	}
	t.hideWindow(ctx)
	return true
}

// context 读取当前 Wails context。
func (t *Tray) context() context.Context {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.ctx
}

// showWindow 恢复已隐藏或最小化的 Wails 窗口。
func (t *Tray) showWindow() {
	ctx := t.context()
	if ctx == nil {
		return
	}
	setPlatformWindowVisible(true)
	wailsruntime.Show(ctx)
	wailsruntime.WindowShow(ctx)
	wailsruntime.WindowUnminimise(ctx)
}

// hideWindow 把窗口收进托盘，并同步隐藏 macOS Dock 图标。
func (t *Tray) hideWindow(ctx context.Context) {
	if ctx == nil {
		return
	}
	wailsruntime.WindowHide(ctx)
	setPlatformWindowVisible(false)
}

// quit 通过 Wails 退出应用，尚未启动完成时退回到 systray 退出。
func (t *Tray) quit() {
	ctx := t.context()
	if ctx == nil {
		quitPlatformFallback()
		return
	}
	t.quitting.Store(true)
	setPlatformWindowVisible(true)
	wailsruntime.Quit(ctx)
}
