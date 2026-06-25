//go:build wails

package config

var runtimeConfigProfile = "cli"

// UseTrayConfigProfile 标记当前进程使用系统托盘启动链路。
func UseTrayConfigProfile() {
	runtimeConfigProfile = "tray"
}

// UseDesktopConfigProfile 标记当前进程正在 Wails v2 桌面环境内运行。
func UseDesktopConfigProfile() {
	runtimeConfigProfile = "desktop"
}

// configProfile 返回 Wails 构建中的运行时配置类型；未进入桌面环境时仍按 CLI 处理。
func configProfile() string {
	return runtimeConfigProfile
}
