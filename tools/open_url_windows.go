//go:build windows

package tools

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// openURL 在 Windows 上使用 ShellExecute 打开 URL，避免通过 CMD 启动导致黑框闪烁。
func openURL(uri string) error {
	verbPtr, err := windows.UTF16PtrFromString("open")
	if err != nil {
		return fmt.Errorf("open url verb utf16: %w", err)
	}
	uriPtr, err := windows.UTF16PtrFromString(uri)
	if err != nil {
		return fmt.Errorf("open url utf16: %w", err)
	}

	// https://learn.microsoft.com/windows/win32/api/shellapi/nf-shellapi-shellexecutew
	// 返回值 > 32 表示成功
	r, _, callErr := windows.NewLazySystemDLL("shell32.dll").
		NewProc("ShellExecuteW").
		Call(
			0,
			uintptr(unsafe.Pointer(verbPtr)),
			uintptr(unsafe.Pointer(uriPtr)),
			0,
			0,
			uintptr(windows.SW_SHOWNORMAL),
		)
	if r <= 32 {
		// callErr 可能为 Errno(0)，这里也一并返回便于定位
		return fmt.Errorf("ShellExecuteW failed: ret=%d err=%v", r, callErr)
	}
	return nil
}
