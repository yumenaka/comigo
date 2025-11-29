//go:build windows

package windows_registry

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/yumenaka/comigo/assets/locale"
	"golang.org/x/sys/windows/registry"
)

const (
	comigoArchiveProgID    = "Comigo.Archive"
	classesRootCurrentUser = `Software\\Classes`
)

// RegisterComigoAsDefaultArchiveHandler 将当前可执行程序注册为指定扩展名的默认打开方式。
// exts 为空时，会使用默认扩展名 [.zip, .rar]。
// 该操作仅在当前用户 (HKCU) 范围内生效，不需要管理员权限。
func RegisterComigoAsDefaultArchiveHandler(exts []string) error {
	exePath, err := currentExecutablePath()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	command := fmt.Sprintf("\"%s\" \"%%1\"", exePath)

	// 1. 设置 Comigo.Archive ProgID（幂等，多次调用会覆盖旧值）
	if err := setArchiveProgID(command); err != nil {
		return err
	}

	// 2. 将特定扩展名关联到 Comigo.Archive
	if len(exts) == 0 {
		exts = []string{".zip", ".rar"}
	}
	for _, ext := range exts {
		if err := associateExtensionWithProgID(ext, comigoArchiveProgID); err != nil {
			return err
		}
	}

	return nil
}

// UnregisterComigoAsDefaultArchiveHandler 取消指定扩展名对 Comigo 的默认关联。
// 仅在当前用户 (HKCU) 范围内生效。
func UnregisterComigoAsDefaultArchiveHandler(exts []string) error {
	if len(exts) == 0 {
		exts = []string{".zip", ".rar"}
	}

	for _, ext := range exts {
		if err := dissociateExtensionFromProgID(ext, comigoArchiveProgID); err != nil {
			return err
		}
	}

	// 尝试删除 Comigo.Archive ProgID 相关键（如果仍被其他扩展使用则可能保留）
	_ = deleteComigoArchiveProgID()

	return nil
}

// AddComigoToFolderContextMenu 在文件夹右键菜单中添加“使用Comigo打开”菜单项。
// 该操作仅在当前用户 (HKCU) 范围内生效，多次调用会覆盖原有设置。
func AddComigoToFolderContextMenu() error {
	exePath, err := currentExecutablePath()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	command := fmt.Sprintf("\"%s\" \"%%1\"", exePath)
	menuText := locale.GetString("context_menu_open_with_comigo")

	// HKCU\Software\Classes\Directory\shell\ComigoOpen
	keyPath := filepath.Join(classesRootCurrentUser, `Directory\\shell\\ComigoOpen`)
	if err := setStringValue(registry.CURRENT_USER, keyPath, "", menuText); err != nil {
		return err
	}

	// HKCU\Software\Classes\Directory\shell\ComigoOpen\command
	commandKeyPath := filepath.Join(keyPath, "command")
	if err := setStringValue(registry.CURRENT_USER, commandKeyPath, "", command); err != nil {
		return err
	}

	return nil
}

// RemoveComigoFromFolderContextMenu 移除文件夹右键菜单中的“使用Comigo打开”菜单项。
// 仅在当前用户 (HKCU) 范围内生效。
func RemoveComigoFromFolderContextMenu() error {
	// 先删除 command 子键
	commandKeyPath := filepath.Join(classesRootCurrentUser, `Directory\\shell\\ComigoOpen\\command`)
	if err := deleteRegistryKey(registry.CURRENT_USER, commandKeyPath); err != nil {
		return err
	}

	// 再删除 ComigoOpen 主键
	menuKeyPath := filepath.Join(classesRootCurrentUser, `Directory\\shell\\ComigoOpen`)
	if err := deleteRegistryKey(registry.CURRENT_USER, menuKeyPath); err != nil {
		return err
	}

	return nil
}

// HasComigoFolderContextMenu 检查当前用户下是否已存在 Comigo 文件夹右键菜单。
func HasComigoFolderContextMenu() bool {
	menuKeyPath := filepath.Join(classesRootCurrentUser, `Directory\\shell\\ComigoOpen`)
	key, err := registry.OpenKey(registry.CURRENT_USER, menuKeyPath, registry.QUERY_VALUE)
	if err != nil {
		if isNotFoundError(err) {
			return false
		}
		// 其他错误视为“未知”，这里保守返回 false，避免误判为存在
		return false
	}
	_ = key.Close()
	return true
}

// currentExecutablePath 返回当前可执行文件的绝对路径，并解析符号链接。
func currentExecutablePath() (string, error) {
	p, err := os.Executable()
	if err != nil {
		return "", err
	}

	p, err = filepath.EvalSymlinks(p)
	if err != nil {
		// 如果解析符号链接失败，退回原始路径
		return p, nil
	}

	return p, nil
}

// setArchiveProgID 在 HKCU\Software\Classes 下创建/更新 Comigo.Archive ProgID。
func setArchiveProgID(command string) error {
	progIDPath := filepath.Join(classesRootCurrentUser, comigoArchiveProgID)

	// 设置描述
	if err := setStringValue(registry.CURRENT_USER, progIDPath, "", "Comigo Archive"); err != nil {
		return err
	}

	// 设置 shell\open\command
	commandKeyPath := filepath.Join(progIDPath, `shell\\open\\command`)
	if err := setStringValue(registry.CURRENT_USER, commandKeyPath, "", command); err != nil {
		return err
	}

	return nil
}

// deleteComigoArchiveProgID 尝试删除 Comigo.Archive 相关的注册表键。
// 如果键不存在或删除失败，不视为致命错误。
func deleteComigoArchiveProgID() error {
	// 依次删除子键，最后删除主键
	paths := []string{
		filepath.Join(classesRootCurrentUser, `Comigo.Archive\\shell\\open\\command`),
		filepath.Join(classesRootCurrentUser, `Comigo.Archive\\shell\\open`),
		filepath.Join(classesRootCurrentUser, `Comigo.Archive\\shell`),
		filepath.Join(classesRootCurrentUser, `Comigo.Archive`),
	}

	for _, p := range paths {
		if err := deleteRegistryKey(registry.CURRENT_USER, p); err != nil {
			// 忽略路径不存在错误，其它错误向上返回
			if !isNotFoundError(err) {
				return err
			}
		}
	}

	return nil
}

// associateExtensionWithProgID 将指定扩展名关联到给定的 ProgID（在 HKCU\Software\Classes 下）。
func associateExtensionWithProgID(ext, progID string) error {
	if ext == "" {
		return nil
	}

	// 确保扩展名以 '.' 开头
	if ext[0] != '.' {
		ext = "." + ext
	}

	extKeyPath := filepath.Join(classesRootCurrentUser, ext)
	return setStringValue(registry.CURRENT_USER, extKeyPath, "", progID)
}

// dissociateExtensionFromProgID 取消扩展名到指定 ProgID 的默认关联（如果当前默认值为该 ProgID）。
func dissociateExtensionFromProgID(ext, progID string) error {
	if ext == "" {
		return nil
	}

	if ext[0] != '.' {
		ext = "." + ext
	}

	extKeyPath := filepath.Join(classesRootCurrentUser, ext)

	key, err := registry.OpenKey(registry.CURRENT_USER, extKeyPath, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return fmt.Errorf("open registry key %s: %w", extKeyPath, err)
	}
	defer key.Close()

	current, _, err := key.GetStringValue("")
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return fmt.Errorf("get registry value for %s: %w", extKeyPath, err)
	}

	if current != progID {
		// 当前默认值不是 Comigo 的 ProgID，避免误删
		return nil
	}

	if err := key.DeleteValue(""); err != nil && !isNotFoundError(err) {
		return fmt.Errorf("delete default value for %s: %w", extKeyPath, err)
	}

	return nil
}

// setStringValue 在给定注册表根键和路径下创建/打开子键，并设置字符串值。
// name 为空字符串时表示设置默认值。
func setStringValue(root registry.Key, path, name, value string) error {
	// registry.OpenKey 不支持多级自动创建，这里使用 CreateKey 逐级创建。
	key, _, err := registry.CreateKey(root, path, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("create/open registry key %s: %w", path, err)
	}
	defer key.Close()

	if err := key.SetStringValue(name, value); err != nil {
		return fmt.Errorf("set registry value %s in %s: %w", name, path, err)
	}

	return nil
}

// deleteRegistryKey 删除指定路径的注册表键。若键不存在则忽略。
func deleteRegistryKey(root registry.Key, path string) error {
	if err := registry.DeleteKey(root, path); err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return fmt.Errorf("delete registry key %s: %w", path, err)
	}
	return nil
}

// isNotFoundError 判断错误是否为“路径/键不存在”。
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	return errno == syscall.ERROR_FILE_NOT_FOUND || errno == syscall.ERROR_PATH_NOT_FOUND
}
