//go:build windows

package locale

import (
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// getWindowsLocale 在 Windows 上通过 WinAPI 获取当前 UI 语言（例如 zh-CN / en-US），避免启动 powershell。
func getWindowsLocale() (lang string, loc string, ok bool) {
	// https://learn.microsoft.com/windows/win32/api/winnls/nf-winnls-getuserdefaultlocalename
	const localeNameMaxLength = 85 // LOCALE_NAME_MAX_LENGTH

	buf := make([]uint16, localeNameMaxLength)
	r, _, _ := windows.NewLazySystemDLL("kernel32.dll").
		NewProc("GetUserDefaultLocaleName").
		Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if r == 0 {
		return "", "", false
	}

	name := strings.TrimSpace(windows.UTF16ToString(buf))
	if name == "" {
		return "", "", false
	}

	// 统一处理：zh-CN / en-US / zh-Hans-CN -> zh_CN / en_US / zh_Hans_CN
	name = strings.Split(name, ".")[0]
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "-", "_")
	parts := strings.Split(name, "_")
	if len(parts) == 0 || parts[0] == "" {
		return "", "", false
	}

	lang = parts[0]
	loc = lang
	// Windows 上可能出现 zh_CN、zh_Hans_CN 等形式，这里统一取最后一段作为地区信息
	if len(parts) > 1 && parts[len(parts)-1] != "" {
		loc = parts[len(parts)-1]
	}
	return lang, loc, true
}
