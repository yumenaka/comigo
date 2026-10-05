//go:build windows

package cmd

import "syscall"

// sysProcAttrForBackground 脱离父控制台；仅 HideWindow 仍会随终端关闭而被终止。
func sysProcAttrForBackground() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: 0x00000008 | 0x00000200} // DETACHED_PROCESS | CREATE_NEW_PROCESS_GROUP
}
