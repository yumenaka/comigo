//go:build wails && android && !js && !bindings

package wails_systray

// traySupported 在 Android 上固定关闭桌面托盘，避免编译桌面-only 依赖。
func traySupported() bool {
	return false
}

func startPlatform(*Tray) func() {
	return func() {}
}

func setPlatformWindowVisible(bool) {
}

func quitPlatformFallback() {
}
