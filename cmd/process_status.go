package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/pflag"
	"github.com/yumenaka/comigo/assets/locale"
)

// localProcess 同时保留启动端口用于筛选、实际端口用于展示。
type localProcess struct {
	state       processState
	path        string
	startupPort int
	ports       string
}

// processInvocation 复用命令和参数定义，避免把配置值或 -- 后的路径误判为查询命令。
func processInvocation(args []string) (configFile string, transient bool) {
	command, _, err := RootCmd.Find(args)
	if err == nil && command != RootCmd && command.Name() != "run" {
		return "", true
	}
	// 只复制参数语法，不绑定全局配置，发现其他进程时不能改写当前命令的参数。
	flags := pflag.NewFlagSet("discovery", pflag.ContinueOnError)
	flags.SetOutput(io.Discard)
	RootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		flags.StringP(flag.Name, flag.Shorthand, "", "")
		flags.Lookup(flag.Name).NoOptDefVal = flag.NoOptDefVal
	})
	flags.BoolP("help", "h", false, "")
	flags.BoolP("version", "v", false, "")
	if flags.Parse(args) != nil {
		return "", false
	}
	for _, name := range []string{"help", "version", "upgrade"} {
		if enabled, _ := strconv.ParseBool(flags.Lookup(name).Value.String()); enabled {
			return "", true
		}
	}
	configFile, _ = flags.GetString("config")
	return configFile, false
}

// showProcessStatus 与 stop 复用实例发现，ps/ls 直接作为此命令的别名。
func showProcessStatus(out io.Writer) error {
	targets, err := localProcesses()
	if err != nil {
		return err
	}
	for _, target := range targets {
		if target.path != "" {
			fmt.Fprintf(out, locale.GetString("cli_status_running")+"\n", target.state.PID, target.state.Port, target.state.URL)
		} else {
			fmt.Fprintf(out, locale.GetString("cli_status_process")+"\n", target.state.PID, target.ports)
		}
	}
	if len(targets) == 0 {
		fmt.Fprintln(out, locale.GetString("cli_not_running"))
	}
	return nil
}

// localProcesses 从本机进程找到配置目录，再用令牌及 PID 验证管理记录。
func localProcesses() ([]localProcess, error) {
	base, err := processStatePath()
	if err != nil {
		return nil, err
	}
	bases := []string{base}
	if home, err := os.UserHomeDir(); err == nil {
		bases = append(bases, filepath.Join(home, ".config", "comigo", ".comi-process.json"))
	}
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf(locale.GetString("cli_status_scan_failed"), err)
	}
	var visible []*process.Process
	for _, proc := range processes {
		if int(proc.Pid) == os.Getpid() {
			continue
		}
		name, err := proc.Name()
		if err != nil || (name != "comi" && !strings.EqualFold(name, "comi.exe")) {
			continue
		}
		// 只解析路径，不打印命令行或环境，避免暴露其他启动参数中的凭据。
		args, _ := proc.CmdlineSlice()
		var file string
		if len(args) > 1 {
			var transient bool
			file, transient = processInvocation(args[1:])
			if transient {
				continue
			}
		}
		visible = append(visible, proc)
		cwd, _ := proc.Cwd()
		executable, _ := proc.Exe()
		dirs := []string{cwd, filepath.Dir(executable)}
		if file != "" {
			if !filepath.IsAbs(file) {
				file = filepath.Join(cwd, file)
			}
			dirs = append(dirs, filepath.Dir(file))
		}
		env, _ := proc.Environ()
		for _, entry := range env {
			if dir, ok := strings.CutPrefix(entry, "COMIGO_CONFIG_DIR="); ok && dir != "" {
				if !filepath.IsAbs(dir) {
					dir = filepath.Join(cwd, dir)
				}
				dirs = append(dirs, dir)
			}
		}
		for _, dir := range dirs {
			if dir != "" {
				bases = append(bases, filepath.Join(dir, ".comi-process.json"))
			}
		}
	}
	var targets []localProcess
	seenBases := make(map[string]bool)
	seen := make(map[int]bool)
	for _, base := range bases {
		base, err = filepath.Abs(base)
		if err != nil {
			return nil, err
		}
		if seenBases[base] {
			continue
		}
		seenBases[base] = true
		paths, err := filepath.Glob(strings.TrimSuffix(base, ".json") + "-*.json")
		if err != nil {
			return nil, err
		}
		for _, path := range paths {
			state, err := readProcessState(path)
			if err != nil || state.PID == os.Getpid() || seen[state.PID] {
				continue
			}
			var current processState
			if processRequest(state, "status", "", &current) != nil || current.PID != state.PID {
				continue
			}
			target := localProcess{state: state, path: path, startupPort: state.Port, ports: strconv.Itoa(current.Port)}
			target.state.Port, target.state.URL = current.Port, current.URL
			targets = append(targets, target)
			seen[current.PID] = true
		}
	}
	for _, proc := range visible {
		pid := int(proc.Pid)
		if seen[pid] {
			continue
		}
		// 未找到管理记录时只报告系统可确认的监听端口，不推断阅读服务已就绪。
		ports := locale.GetString("cli_status_no_listener")
		connections, err := proc.Connections()
		if err != nil {
			ports = locale.GetString("cli_status_ports_unavailable")
		} else {
			unique := make(map[uint32]bool)
			var numbers []int
			for _, connection := range connections {
				if connection.Status == "LISTEN" && !unique[connection.Laddr.Port] {
					unique[connection.Laddr.Port] = true
					numbers = append(numbers, int(connection.Laddr.Port))
				}
			}
			sort.Ints(numbers)
			var values []string
			for _, port := range numbers {
				values = append(values, strconv.Itoa(port))
			}
			if len(values) != 0 {
				ports = strings.Join(values, ", ")
			}
		}
		targets = append(targets, localProcess{state: processState{PID: pid}, ports: ports})
	}
	sort.Slice(targets, func(i, j int) bool { return targets[i].state.PID < targets[j].state.PID })
	return targets, nil
}
