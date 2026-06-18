//go:build wails && !js

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
)

// DeleteBookFile 由 Wails WebView 绑定调用，确认后把书籍源文件移到系统垃圾桶。
func (a *App) DeleteBookFile(bookID string) (bool, error) {
	if a.ctx == nil || bookID == "" {
		return false, errors.New(locale.GetString("wails_delete_file_not_allowed"))
	}
	book, err := model.IStore.GetBook(bookID)
	if err != nil {
		return false, err
	}
	trashPath, isDir, err := trashableBookPath(book, config.GetCfg().StoreUrls)
	if err != nil {
		return false, err
	}
	ok, err := a.confirmDeleteBookFile(book, trashPath)
	if err != nil || !ok {
		return false, err
	}
	if err := movePathToSystemTrash(trashPath, isDir); err != nil {
		return false, err
	}
	if err := model.IStore.DeleteBook(book.BookID); err != nil {
		return false, err
	}
	if err := model.IStore.GenerateBookGroup(); err != nil {
		return false, err
	}
	// 删除后立刻保存，避免重启后从元数据里复活。
	cmd.SaveMetadata()
	return true, nil
}

// confirmDeleteBookFile 用系统弹窗二次确认，取消时不会触碰源文件。
func (a *App) confirmDeleteBookFile(book *model.Book, trashPath string) (bool, error) {
	deleteButton := locale.GetString("wails_delete_file_confirm_button")
	cancelButton := locale.GetString("cancel")
	result, err := wailsruntime.MessageDialog(a.ctx, wailsruntime.MessageDialogOptions{
		Type:          wailsruntime.QuestionDialog,
		Title:         locale.GetString("wails_delete_file_confirm_title"),
		Message:       fmt.Sprintf(locale.GetString("wails_delete_file_confirm_message"), trashPath),
		Buttons:       []string{deleteButton, cancelButton},
		DefaultButton: cancelButton,
		CancelButton:  cancelButton,
	})
	return result == deleteButton, err
}

// trashableBookPath 只允许删除当前本地书库内的真实书籍，远程书和书籍组不触碰磁盘。
func trashableBookPath(book *model.Book, storeUrls []string) (string, bool, error) {
	if book == nil || book.BookPath == "" {
		return "", false, errors.New(locale.GetString("wails_delete_file_not_allowed"))
	}
	if book.Type == model.TypeBooksGroup || book.IsRemote || book.RemoteURL != "" || book.RemoteBookID != "" || book.RemoteStoreKey != "" {
		return "", false, errors.New(locale.GetString("wails_delete_file_not_allowed"))
	}
	trashPath, err := filepath.Abs(book.BookPath)
	if err != nil {
		return "", false, err
	}
	info, err := os.Stat(trashPath)
	if err != nil {
		return "", false, err
	}
	if !pathInsideLocalStore(trashPath, storeUrls) {
		return "", false, errors.New(locale.GetString("wails_delete_file_not_allowed"))
	}
	return trashPath, info.IsDir(), nil
}

// pathInsideLocalStore 判断目标是否位于某个本地书库目录内，避免误删外部路径。
func pathInsideLocalStore(target string, storeUrls []string) bool {
	target = filepath.Clean(target)
	for _, storeURL := range storeUrls {
		info := tools.DetectStoreURL(storeURL)
		if info.Type != tools.StoreBackendLocalDisk {
			continue
		}
		storePath := info.LocalPath
		if storePath == "" {
			storePath = info.URL
		}
		storeAbs, err := filepath.Abs(storePath)
		if err != nil {
			continue
		}
		if tools.IsSubPath(filepath.Clean(storeAbs), target) {
			return true
		}
	}
	return false
}

// movePathToSystemTrash 使用系统自带能力进垃圾桶；找不到能力时返回错误，不退化为永久删除。
func movePathToSystemTrash(target string, isDir bool) error {
	switch runtime.GOOS {
	case "darwin":
		return runTrashCommand("osascript",
			"-e", "on run argv",
			"-e", `tell application "Finder" to delete (POSIX file (item 1 of argv) as alias)`,
			"-e", "end run",
			target,
		)
	case "windows":
		script := `Add-Type -AssemblyName Microsoft.VisualBasic
$p = $args[0]
if (` + boolPowerShell(isDir) + `) {
  [Microsoft.VisualBasic.FileIO.FileSystem]::DeleteDirectory($p, [Microsoft.VisualBasic.FileIO.UIOption]::OnlyErrorDialogs, [Microsoft.VisualBasic.FileIO.RecycleOption]::SendToRecycleBin)
} else {
  [Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile($p, [Microsoft.VisualBasic.FileIO.UIOption]::OnlyErrorDialogs, [Microsoft.VisualBasic.FileIO.RecycleOption]::SendToRecycleBin)
}`
		return runTrashCommand("powershell", "-NoProfile", "-NonInteractive", "-Command", script, target)
	case "linux":
		if path, err := exec.LookPath("gio"); err == nil {
			return runTrashCommand(path, "trash", target)
		}
		if path, err := exec.LookPath("trash-put"); err == nil {
			return runTrashCommand(path, target)
		}
	}
	return errors.New(locale.GetString("wails_delete_file_unsupported"))
}

func boolPowerShell(v bool) string {
	if v {
		return "$true"
	}
	return "$false"
}

// runTrashCommand 执行系统垃圾桶命令，并把命令输出补到错误里方便定位失败原因。
func runTrashCommand(name string, args ...string) error {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err == nil {
		return nil
	}
	message := strings.TrimSpace(string(out))
	if message == "" {
		return err
	}
	return fmt.Errorf("%w: %s", err, message)
}
