package cmd

import (
	"fmt"
	"runtime"

	"github.com/sevlyar/go-daemon"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/logger"
)

// DemonFlag 正确地实装，需要理解守护进程的概念
// 需要去 cmd/init_flags.go 设置flag
var (
	DemonFlag      bool
	StopDaemonFlag bool
)

// SetDaemon 设置守护进程, To terminate the daemon use: kill `cat comigo.pid`
// 该函数会在Unix系统上将当前进程转化为守护进程，并在后台运行。
// 如果当前系统不是Unix系统，或者没有指定启动或停止守护进程的参数，则直接返回。
// https://github.com/sevlyar/go-daemon
// https://github.com/sevlyar/go-daemon/blob/v0.1.6/examples/cmd/gd-simple/simple.go
// go run main.go --start
// go run main.go --stop
func SetDaemon() {
	// 如果不是unix系统，直接返回
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		return
	}
	// 如果没有指定启动或停止守护进程的参数，直接返回
	if !DemonFlag && !StopDaemonFlag {
		return
	}
	cntxt := &daemon.Context{
		PidFileName: "/var/run/comigo.pid",
		PidFilePerm: 0o644,
		LogFileName: "comigo.log",
		LogFilePerm: 0o640,
		WorkDir:     "./",
		Umask:       0o27,
		Args:        []string{fmt.Sprintf("[comigo %s daemon]", config.GetVersion())},
	}
	// Reborn 会在指定的上下文中启动当前进程的第二个副本。
	// 该函数在子进程和父进程中分别执行不同的代码段，并对子进程进行守护化（daemonization）。
	// 它看起来类似于 fork-daemonization，但对 goroutine 是安全的。
	// 调用成功时，父进程返回一个 *os.Process 对象，子进程返回 nil；否则返回错误。
	child, err := cntxt.Reborn()
	if err != nil {
		logger.Fatal("Unable to run: ", err)
	}
	// 父运行运行到这里，child不等于nil，然后会返回
	if child != nil {
		return
	} else {
		// 子进程运行到这里，child等于nil
		logger.Info("- - - - - - - - - - - - - - -")
		logger.Info("child daemon started?")
	}
	// 释放PID文件
	defer cntxt.Release()
	// 这里是子进程运行的代码
	logger.Info("- - - - - - - - - - - - - - -")
	logger.Info("daemon started")
}
