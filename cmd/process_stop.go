package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// stopProcesses 先选择目标，再逐个停止；一个实例失败不妨碍其余实例清理。
func stopProcesses(command *cobra.Command, _ []string) error {
	targets, err := localProcesses()
	if err != nil {
		return fmt.Errorf(locale.GetString("cli_stop_error"), err)
	}
	file := config.GetCfg().ConfigFile
	if file != "" {
		file, err = filepath.Abs(file)
		if err != nil {
			return fmt.Errorf(locale.GetString("cli_stop_error"), err)
		}
	}
	selected := targets[:0]
	for _, target := range targets {
		if file != "" && target.state.ConfigFile != file {
			continue
		}
		if command.Flags().Changed("port") && target.startupPort != config.GetCfg().Port {
			continue
		}
		selected = append(selected, target)
	}
	out := command.OutOrStdout()
	if len(selected) == 0 {
		fmt.Fprintln(out, locale.GetString("cli_stop_none"))
		return nil
	}
	if len(selected) > 1 {
		selected, err = selectStopProcesses(command, selected)
		if err != nil {
			return fmt.Errorf(locale.GetString("cli_stop_error"), err)
		}
		if len(selected) == 0 {
			fmt.Fprintln(out, locale.GetString("cli_stop_cancelled"))
			return nil
		}
	}
	var failures []error
	for _, target := range selected {
		if err := stopProcess(target); err != nil {
			failures = append(failures, fmt.Errorf(locale.GetString("cli_stop_failed"), target.state.PID, target.ports, err))
		} else {
			fmt.Fprintf(out, locale.GetString("cli_stopped")+"\n", target.state.PID, target.ports)
		}
	}
	return errors.Join(failures...)
}

// selectStopProcesses 不预选目标；Ctrl+C 或输入结束都取消，不发送停止请求。
func selectStopProcesses(command *cobra.Command, targets []localProcess) ([]localProcess, error) {
	ctx, cancel := signal.NotifyContext(command.Context(), os.Interrupt)
	defer cancel()
	out := command.OutOrStdout()
	fmt.Fprintln(out, locale.GetString("cli_stop_select"))
	for i, target := range targets {
		fmt.Fprintf(out, "%d. "+locale.GetString("cli_stop_candidate")+"\n", i+1, target.state.PID, target.ports)
	}
	fmt.Fprintln(out, "0. "+locale.GetString("cli_stop_all"))
	reader := bufio.NewReader(command.InOrStdin())
	for {
		fmt.Fprint(out, locale.GetString("cli_stop_prompt"))
		// 读取终端时仍监听取消；命令退出后不再使用这次输入。
		result := make(chan struct {
			line string
			err  error
		}, 1)
		go func() {
			line, err := reader.ReadString('\n')
			result <- struct {
				line string
				err  error
			}{line, err}
		}()
		select {
		case <-ctx.Done():
			fmt.Fprintln(out)
			return nil, nil
		case input := <-result:
			if ctx.Err() != nil || (input.err == io.EOF && strings.TrimSpace(input.line) == "") {
				return nil, nil
			}
			if input.err != nil && input.err != io.EOF {
				return nil, input.err
			}
			choice, err := strconv.Atoi(strings.TrimSpace(input.line))
			if err == nil && choice >= 0 && choice <= len(targets) {
				if choice == 0 {
					return targets, nil
				}
				return targets[choice-1 : choice], nil
			}
			fmt.Fprintln(out, locale.GetString("cli_stop_invalid_selection"))
		}
	}
}

// stopProcess 等待原实例删除管理记录后才报告成功，避免误报和 PID 复用误杀。
func stopProcess(target localProcess) error {
	if target.path == "" {
		return errors.New(locale.GetString("cli_stop_unmanaged"))
	}
	if err := processRequest(target.state, "stop", ""); err != nil {
		return err
	}
	for deadline := time.Now().Add(15 * time.Second); time.Now().Before(deadline); time.Sleep(50 * time.Millisecond) {
		current, err := readProcessState(target.path)
		if os.IsNotExist(err) || (err == nil && current.Token != target.state.Token) {
			return nil
		}
	}
	return errors.New(locale.GetString("cli_stop_timeout"))
}
