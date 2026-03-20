//go:build !windows

package cmd

import "syscall"

func sysProcAttrForUpgradeRestart() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}
