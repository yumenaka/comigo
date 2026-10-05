package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/yumenaka/comigo/assets/locale"
)

var processCommandsOnce sync.Once

// addProcessCommands 只注册一级子命令；run 与旧的无参数入口共用初始化流程。
func addProcessCommands() {
	RootCmd.AddCommand(&cobra.Command{
		Use: "run [paths...]", Short: locale.GetString("cli_run"),
		Args: cobra.ArbitraryArgs, Run: RootCmd.Run,
	})
	RootCmd.AddCommand(&cobra.Command{
		Use: "start [paths...]", Short: locale.GetString("cli_start"),
		Args: cobra.ArbitraryArgs, RunE: startBackground,
	})
	for _, name := range []string{"stop", "reload"} {
		RootCmd.AddCommand(&cobra.Command{
			Use: name, Short: locale.GetString("cli_" + name), Args: cobra.NoArgs,
			RunE: func(command *cobra.Command, _ []string) error {
				return controlProcess(command.Name())
			},
		})
	}
	RootCmd.AddCommand(&cobra.Command{
		Use: "upgrade", Short: locale.GetString("cli_upgrade"), Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error { return runSelfUpgrade() },
	})
	RootCmd.AddCommand(&cobra.Command{
		Use: "version", Short: locale.GetString("cli_version"), Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			// 与 --version 使用同一模板，避免版本输出逐渐分叉。
			tmpl, err := template.New("version").Parse(RootCmd.VersionTemplate())
			if err != nil {
				return err
			}
			return tmpl.Execute(command.OutOrStdout(), RootCmd)
		},
	})
}

// RunProcessCommand 提前执行管理命令，防止打印版本、升级或停止之后又进入 TUI。
// run 返回未处理，由原来的前台入口继续运行；Cobra 负责解析 flag，避免手工猜测参数位置。
func RunProcessCommand(args []string, out io.Writer) (bool, error) {
	processCommandsOnce.Do(addProcessCommands)
	InitFlags()
	command, remaining, err := RootCmd.Find(args)
	if err != nil {
		return false, nil
	}
	command.InitDefaultHelpFlag()
	command.InitDefaultVersionFlag()
	if err := command.ParseFlags(remaining); err != nil {
		return true, err
	}
	handled := false
	switch command.Name() {
	case "start", "stop", "reload", "upgrade", "version":
		handled = true
	}
	// 让 Cobra 区分 flag 与 -- 后的文件名，避免扫描 os.Args 将书库误作 --version。
	for _, name := range []string{"help", "version", "upgrade"} {
		value, _ := command.Flags().GetBool(name)
		handled = handled || value
	}
	if !handled {
		return false, nil
	}
	RootCmd.SetArgs(args)
	RootCmd.SetOut(out)
	return true, RootCmd.Execute()
}

// backgroundArgs 从已经解析的 flag 重建子进程参数，保留带空格的路径和 -- 后的文件名。
// 不直接替换 os.Args 中的 "start"，因为它也可能只是 --config 的值或书库名。
func backgroundArgs(command *cobra.Command, paths []string) []string {
	args := []string{"run"}
	command.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Name != "no-tui" && flag.Name != "open-browser" {
			args = append(args, "--"+flag.Name+"="+flag.Value.String())
		}
	})
	args = append(args, "--no-tui", "--open-browser=false", "--")
	return append(args, paths...)
}

// startBackground 脱离终端启动当前程序，并等待该子进程完成初始化。
// 标准输入为空，输出写入配置目录的日志文件，退出启动终端不会关闭后台服务。
func startBackground(command *cobra.Command, paths []string) error {
	statePath, err := processStatePath()
	if err != nil {
		return err
	}
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	logPath := statePath + ".log"
	log, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer log.Close()
	child := exec.Command(executable, backgroundArgs(command, paths)...)
	child.Stdout, child.Stderr = log, log
	child.SysProcAttr = sysProcAttrForBackground()
	if err := child.Start(); err != nil {
		return err
	}
	// Wait 同时回收失败的子进程；成功时父命令退出，后台进程继续由系统管理。
	exited := make(chan error, 1)
	go func() { exited <- child.Wait() }()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.NewTimer(2 * time.Minute)
	defer timeout.Stop()
	for {
		select {
		case err := <-exited:
			return fmt.Errorf(locale.GetString("cli_start_failed"), logPath, err)
		case <-timeout.C:
			_ = child.Process.Kill()
			<-exited
			return fmt.Errorf(locale.GetString("cli_start_failed"), logPath, "timeout")
		case <-ticker.C:
			state, err := readProcessState(statePath)
			// 不能把此前残留的状态文件或另一个实例误当作本次启动成功。
			if err == nil && state.PID == child.Process.Pid && processRequest(state, "ready", "") == nil {
				fmt.Fprintf(command.OutOrStdout(), locale.GetString("cli_started")+"\n", state.PID, logPath)
				return nil
			}
		}
	}
}
