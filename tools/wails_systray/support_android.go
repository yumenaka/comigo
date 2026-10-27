//go:build wails && android && !js

package wails_systray

// traySupported 移动端不创建系统托盘。
func traySupported() bool { return false }
