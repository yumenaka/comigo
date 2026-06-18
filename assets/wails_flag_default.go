//go:build !wails

package assets

import "net/http"

// IsWailsBuild 告诉模板和静态脚本当前是否为 Wails 桌面壳构建。
func IsWailsBuild() bool {
	return false
}

func isWailsBuild() bool {
	return IsWailsBuild()
}

// IsWailsWebViewRequest 非 Wails 构建不存在 WebView 请求。
func IsWailsWebViewRequest(_ *http.Request) bool {
	return false
}
