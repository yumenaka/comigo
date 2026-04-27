package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
)

// RestartCurrentExecutable 启动当前可执行文件的新进程（参数与环境继承），用于托盘升级后由旧进程退出、新进程接管。
func RestartCurrentExecutable() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return err
	}
	c := exec.Command(exe, os.Args[1:]...)
	c.Env = os.Environ()
	c.SysProcAttr = sysProcAttrForUpgradeRestart()
	return c.Start()
}

// PrepareTrayUpgradeRestart 替换二进制成功后：可选释放单实例锁、停止服务、拉起新进程。调用方在返回 nil 后应立即 systray.Quit 与 os.Exit。
func PrepareTrayUpgradeRestart(shutdownServer func(), releaseSingleInstance func()) error {
	if releaseSingleInstance != nil {
		releaseSingleInstance()
	}
	if shutdownServer != nil {
		shutdownServer()
	}
	if err := RestartCurrentExecutable(); err != nil {
		return err
	}
	return nil
}
