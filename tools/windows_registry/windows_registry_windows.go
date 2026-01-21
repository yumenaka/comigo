//go:build windows

package windows_registry

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/yumenaka/comigo/assets/locale"
	"golang.org/x/sys/windows/registry"
)

const (
	comigoArchiveProgID    = "Comigo.Archive"
	classesRootCurrentUser = `Software\\Classes`
)

// 默认作为候选打开方式支持的压缩扩展名
var defaultArchiveExts = []string{".zip", ".rar"}

// RegisterComigoAsDefaultArchiveHandler 将当前可执行程序注册为指定扩展名的“候选打开方式”（OpenWithProgids），并提供 Comigo.Archive ProgID。
// exts 为空时，会使用默认扩展名 [.zip, .rar]。
// 该操作仅在当前用户 (HKCU) 范围内生效，不会尝试修改 UserChoice。
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

	// 2. 为 Applications\\<exe> 注册 open 命令，便于系统在“选择其他应用”中展示
	if err := registerComigoApplication(command); err != nil {
		return err
	}

	// 3. 将特定扩展名关联到 Comigo.Archive，并写入 OpenWithProgids
	if len(exts) == 0 {
		exts = defaultArchiveExts
	}
	for _, ext := range exts {
		if err := associateExtensionWithProgID(ext, comigoArchiveProgID); err != nil {
			return err
		}
		if err := registerOpenWithProgID(ext, comigoArchiveProgID); err != nil {
			return err
		}
	}

	return nil
}

// UnregisterComigoAsDefaultArchiveHandler 取消指定扩展名对 Comigo 的候选打开方式及 ProgID 关联。
// 仅在当前用户 (HKCU) 范围内生效。
func UnregisterComigoAsDefaultArchiveHandler(exts []string) error {
	if len(exts) == 0 {
		exts = defaultArchiveExts
	}

	for _, ext := range exts {
		if err := dissociateExtensionFromProgID(ext, comigoArchiveProgID); err != nil {
			return err
		}
		if err := unregisterOpenWithProgID(ext, comigoArchiveProgID); err != nil {
			return err
		}
	}

	// 尝试删除 Comigo.Archive ProgID 相关键（如果仍被其他扩展使用则可能保留）
	_ = deleteComigoArchiveProgID()
	// 尝试删除 Applications\\<exe> 相关键
	_ = unregisterComigoApplication()

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
		// 其他错误视为"未知"，这里保守返回 false，避免误判为存在
		return false
	}
	_ = key.Close()
	return true
}

// AddComigoHereToFolderBackgroundMenu 在文件夹空白处右键菜单中添加"ComiGo Here"菜单项。
// 该操作仅在当前用户 (HKCU) 范围内生效，多次调用会覆盖原有设置。
// 兼容 Win11 新右键菜单和经典右键菜单。
func AddComigoHereToFolderBackgroundMenu() error {
	exePath, err := currentExecutablePath()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	command := fmt.Sprintf("\"%s\" \"%%V\"", exePath)
	menuText := "ComiGo Here"
	iconValue := fmt.Sprintf("%s,0", exePath)

	// HKCU\Software\Classes\Directory\Background\shell\ComiGo
	keyPath := filepath.Join(classesRootCurrentUser, `Directory\\Background\\shell\\ComiGo`)
	if err := setStringValue(registry.CURRENT_USER, keyPath, "", menuText); err != nil {
		return err
	}

	// 设置图标
	if err := setStringValue(registry.CURRENT_USER, keyPath, "Icon", iconValue); err != nil {
		return err
	}

	// HKCU\Software\Classes\Directory\Background\shell\ComiGo\command
	commandKeyPath := filepath.Join(keyPath, "command")
	if err := setStringValue(registry.CURRENT_USER, commandKeyPath, "", command); err != nil {
		return err
	}

	return nil
}

// RemoveComigoHereFromFolderBackgroundMenu 移除文件夹空白处右键菜单中的"ComiGo Here"菜单项。
// 仅在当前用户 (HKCU) 范围内生效。
func RemoveComigoHereFromFolderBackgroundMenu() error {
	// 先删除 command 子键
	commandKeyPath := filepath.Join(classesRootCurrentUser, `Directory\\Background\\shell\\ComiGo\\command`)
	if err := deleteRegistryKey(registry.CURRENT_USER, commandKeyPath); err != nil {
		return err
	}

	// 再删除 ComiGo 主键
	menuKeyPath := filepath.Join(classesRootCurrentUser, `Directory\\Background\\shell\\ComiGo`)
	if err := deleteRegistryKey(registry.CURRENT_USER, menuKeyPath); err != nil {
		return err
	}

	return nil
}

// HasComigoHereFolderBackgroundMenu 检查当前用户下是否已存在 Comigo 文件夹空白处右键菜单。
func HasComigoHereFolderBackgroundMenu() bool {
	menuKeyPath := filepath.Join(classesRootCurrentUser, `Directory\\Background\\shell\\ComiGo`)
	key, err := registry.OpenKey(registry.CURRENT_USER, menuKeyPath, registry.QUERY_VALUE)
	if err != nil {
		if isNotFoundError(err) {
			return false
		}
		// 其他错误视为"未知"，这里保守返回 false，避免误判为存在
		return false
	}
	_ = key.Close()
	return true
}

// HasComigoArchiveAssociation 检查是否已为指定扩展名（或默认扩展名）注册了 Comigo 作为候选打开方式。
// 如果 exts 为空，则使用 defaultArchiveExts。只要任意一个扩展名存在 OpenWithProgids 中的 Comigo.Archive，即视为已注册。
func HasComigoArchiveAssociation(exts []string) bool {
	if len(exts) == 0 {
		exts = defaultArchiveExts
	}
	for _, ext := range exts {
		if hasOpenWithProgID(ext, comigoArchiveProgID) {
			return true
		}
	}
	return false
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

// registerOpenWithProgID 在 HKCU\\Software\\Classes\\.ext\\OpenWithProgids 下登记 Comigo 的 ProgID 作为候选打开方式。
func registerOpenWithProgID(ext, progID string) error {
	if ext == "" {
		return nil
	}
	if ext[0] != '.' {
		ext = "." + ext
	}
	keyPath := filepath.Join(classesRootCurrentUser, ext, "OpenWithProgids")
	// 值内容通常为空字符串即可
	return setStringValue(registry.CURRENT_USER, keyPath, progID, "")
}

// unregisterOpenWithProgID 从 HKCU\\Software\\Classes\\.ext\\OpenWithProgids 中移除 Comigo 的 ProgID。
func unregisterOpenWithProgID(ext, progID string) error {
	if ext == "" {
		return nil
	}
	if ext[0] != '.' {
		ext = "." + ext
	}
	keyPath := filepath.Join(classesRootCurrentUser, ext, "OpenWithProgids")
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return fmt.Errorf("open registry key %s: %w", keyPath, err)
	}
	defer key.Close()

	if err := key.DeleteValue(progID); err != nil && !isNotFoundError(err) {
		return fmt.Errorf("delete OpenWithProgids value %s in %s: %w", progID, keyPath, err)
	}
	return nil
}

// hasOpenWithProgID 判断 HKCU\\Software\\Classes\\.ext\\OpenWithProgids 下是否存在指定 ProgID。
func hasOpenWithProgID(ext, progID string) bool {
	if ext == "" {
		return false
	}
	if ext[0] != '.' {
		ext = "." + ext
	}
	keyPath := filepath.Join(classesRootCurrentUser, ext, "OpenWithProgids")
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()
	_, _, err = key.GetStringValue(progID)
	return err == nil
}

// registerComigoApplication 在 HKCU\\Software\\Classes\\Applications\\<exe>\\shell\\open\\command 下登记 Comigo 的应用信息。
func registerComigoApplication(command string) error {
	exePath, err := currentExecutablePath()
	if err != nil {
		return err
	}
	exeName := filepath.Base(exePath)
	appCommandPath := filepath.Join(classesRootCurrentUser, `Applications\\`+exeName+`\\shell\\open\\command`)
	return setStringValue(registry.CURRENT_USER, appCommandPath, "", command)
}

// unregisterComigoApplication 尝试清理 Applications\\<exe> 相关注册表键。
func unregisterComigoApplication() error {
	exePath, err := currentExecutablePath()
	if err != nil {
		return err
	}
	exeName := filepath.Base(exePath)

	paths := []string{
		filepath.Join(classesRootCurrentUser, `Applications\\`+exeName+`\\shell\\open\\command`),
		filepath.Join(classesRootCurrentUser, `Applications\\`+exeName+`\\shell\\open`),
		filepath.Join(classesRootCurrentUser, `Applications\\`+exeName+`\\shell`),
		filepath.Join(classesRootCurrentUser, `Applications\\`+exeName),
	}
	for _, p := range paths {
		if err := deleteRegistryKey(registry.CURRENT_USER, p); err != nil {
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

// isNotFoundError 判断错误是否为"路径/键不存在"。
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

// CreateDesktopShortcut 在桌面创建 Comigo 的快捷方式。
// 使用 PowerShell 脚本创建快捷方式，如果快捷方式已存在则会被覆盖。
func CreateDesktopShortcut() error {
	exePath, err := currentExecutablePath()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	// 获取桌面路径
	desktopPath, err := getDesktopPath()
	if err != nil {
		return fmt.Errorf("get desktop path: %w", err)
	}

	// 快捷方式文件名
	shortcutName := "Comigo.lnk"
	shortcutPath := filepath.Join(desktopPath, shortcutName)

	// 使用 PowerShell 创建快捷方式
	// $WshShell = New-Object -ComObject WScript.Shell
	// $Shortcut = $WshShell.CreateShortcut("路径")
	// $Shortcut.TargetPath = "目标路径"
	// $Shortcut.Save()
	psScript := fmt.Sprintf(`$WshShell = New-Object -ComObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("%s"); $Shortcut.TargetPath = "%s"; $Shortcut.WorkingDirectory = "%s"; $Shortcut.Save()`,
		shortcutPath,
		exePath,
		filepath.Dir(exePath))

	// 使用隐藏窗口方式执行，避免 powershell 闪黑框
	cmd := exec.Command("powershell", "-NoProfile", "-WindowStyle", "Hidden", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("create shortcut: %w", err)
	}

	return nil
}

// getDesktopPath 获取当前用户的桌面路径。
func getDesktopPath() (string, error) {
	// 尝试从环境变量获取
	userProfile := os.Getenv("USERPROFILE")
	if userProfile != "" {
		desktopPath := filepath.Join(userProfile, "Desktop")
		if _, err := os.Stat(desktopPath); err == nil {
			return desktopPath, nil
		}
	}

	// 尝试从注册表获取
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("open registry key: %w", err)
	}
	defer key.Close()

	desktopPath, _, err := key.GetStringValue("Desktop")
	if err != nil {
		return "", fmt.Errorf("get desktop path from registry: %w", err)
	}

	return desktopPath, nil
}
