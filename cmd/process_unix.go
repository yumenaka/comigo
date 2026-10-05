//go:build !windows

package cmd

import "syscall"

// sysProcAttrForBackground 创建独立会话，使后台进程脱离启动终端。
func sysProcAttrForBackground() *syscall.SysProcAttr {
	return sysProcAttrForUpgradeRestart()
}
