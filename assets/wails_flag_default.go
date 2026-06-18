//go:build !wails

package assets

// IsWailsBuild 告诉模板和静态脚本当前是否为 Wails 桌面壳构建。
func IsWailsBuild() bool {
	return false
}

func isWailsBuild() bool {
	return IsWailsBuild()
}
