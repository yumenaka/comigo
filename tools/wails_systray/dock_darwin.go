//go:build wails && darwin && !js

package wails_systray

import "github.com/wailsapp/wails/v3/pkg/services/dock"

// setDockVisible 复用 Wails 的 Dock 能力，不再维护 Objective-C 桥。
func setDockVisible(visible bool) {
	if visible {
		dock.New().ShowAppIcon()
	} else {
		dock.New().HideAppIcon()
	}
}
