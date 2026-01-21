//go:build !windows

package locale

// getWindowsLocale 在非 Windows 平台上的兜底实现。
func getWindowsLocale() (lang string, loc string, ok bool) {
	return "", "", false
}
