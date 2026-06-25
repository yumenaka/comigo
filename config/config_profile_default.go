//go:build !wails

package config

var runtimeConfigProfile = "cli"

// UseTrayConfigProfile 标记当前进程使用系统托盘启动链路。
func UseTrayConfigProfile() {
	runtimeConfigProfile = "tray"
}

// configProfile 返回非 Wails 启动壳的配置类型；CLI 为默认值。
func configProfile() string {
	return runtimeConfigProfile
}
