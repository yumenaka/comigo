//go:build wails

package assets

import (
	"net"
	"net/http"
)

// IsWailsBuild 告诉模板和静态脚本当前是否为 Wails 桌面壳构建。
func IsWailsBuild() bool {
	return true
}

func isWailsBuild() bool {
	return IsWailsBuild()
}

// IsWailsWebViewRequest 判断请求是否来自 Wails WebView 的 assetserver 转发。
func IsWailsWebViewRequest(r *http.Request) bool {
	if r == nil {
		return false
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	return err == nil && host == "192.0.2.1"
}
