//go:build wails && !darwin && !js

package wails_systray

// setDockVisible 其他平台由窗口管理任务栏图标。
func setDockVisible(bool) {}
