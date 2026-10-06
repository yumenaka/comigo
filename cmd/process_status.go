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
	"github.com/yumenaka/comigo/assets/locale"
)

// showProcessStatus 合并管理通道与系统进程，跨配置目录查询且不把残留文件当成运行实例。
func showProcessStatus(out io.Writer) error {
	base, err := processStatePath()
	if err != nil {
		return err
	}
	bases := []string{base}
	if home, err := os.UserHomeDir(); err == nil {
		bases = append(bases, filepath.Join(home, ".config", "comigo", ".comi-process.json"))
	}
	seen := make(map[int]bool)
	for _, base := range bases {
		paths, err := filepath.Glob(strings.TrimSuffix(base, ".json") + "-*.json")
		if err != nil {
			return err
		}
		for _, path := range paths {
			state, err := readProcessState(path)
			if err != nil || seen[state.PID] {
				continue
			}
			var current processState
			if processRequest(state, "status", "", &current) != nil || current.PID != state.PID {
				continue
			}
			fmt.Fprintf(out, locale.GetString("cli_status_running")+"\n", current.PID, current.Port, current.URL)
			seen[current.PID] = true
		}
	}
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf(locale.GetString("cli_status_scan_failed"), err)
	}
	sort.Slice(processes, func(i, j int) bool { return processes[i].Pid < processes[j].Pid })
	for _, proc := range processes {
		pid := int(proc.Pid)
		if seen[pid] || pid == os.Getpid() {
			continue
		}
		name, err := proc.Name()
		if err != nil || (name != "comi" && !strings.EqualFold(name, "comi.exe")) {
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
		fmt.Fprintf(out, locale.GetString("cli_status_process")+"\n", pid, ports)
		seen[pid] = true
	}
	if len(seen) == 0 {
		fmt.Fprintln(out, locale.GetString("cli_not_running"))
	}
	return nil
}
