//go:build !windows

package tools

// RegisterComigoAsDefaultArchiveHandler 是 Windows 专用功能的空实现。
// 在非 Windows 平台上调用该函数不会执行任何操作，仅用于保持 API 一致性。
func RegisterComigoAsDefaultArchiveHandler(exts []string) error {
	return nil
}

// UnregisterComigoAsDefaultArchiveHandler 是 Windows 专用功能的空实现。
// 在非 Windows 平台上调用该函数不会执行任何操作，仅用于保持 API 一致性。
func UnregisterComigoAsDefaultArchiveHandler(exts []string) error {
	return nil
}

// AddComigoToFolderContextMenu 是 Windows 专用功能的空实现。
// 在非 Windows 平台上调用该函数不会执行任何操作，仅用于保持 API 一致性。
func AddComigoToFolderContextMenu() error {
	return nil
}

// RemoveComigoFromFolderContextMenu 是 Windows 专用功能的空实现。
// 在非 Windows 平台上调用该函数不会执行任何操作，仅用于保持 API 一致性。
func RemoveComigoFromFolderContextMenu() error {
	return nil
}

// HasComigoFolderContextMenu 是 Windows 专用功能的空实现。
// 在非 Windows 平台上始终返回 false。
func HasComigoFolderContextMenu() bool {
	return false
}
