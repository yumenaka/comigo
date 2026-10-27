//go:build !wails && !js && !linux && !android

package system_tray

// traySupported 非 Linux 平台默认有可恢复入口。
func traySupported() bool {
	return true
}
