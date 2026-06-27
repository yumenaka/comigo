//go:build wails && !js && !bindings && !linux && !android

package wails_systray

// traySupported 非 Linux 平台默认有可恢复入口。
func traySupported() bool {
	return true
}
