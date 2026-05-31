//go:build !ios

package tools

import (
	"strings"
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/yumenaka/comigo/tools/logger"
)

var systemMetricLoggedErrors sync.Map

func populateSystemMetrics(sys *SystemStatus) {
	// 获取物理和逻辑核数,以及 CPU、内存整体使用率。
	cpuNumLogical, err := cpu.Counts(true)
	if err != nil {
		logSystemMetricErrorOnce("cpu_logical_count", err)
	} else {
		sys.CPUNumLogical = cpuNumLogical
	}

	cpuNumPhysical, err := cpu.Counts(false)
	if err != nil {
		logSystemMetricErrorOnce("cpu_physical_count", err)
	} else {
		sys.CPUNumPhysical = cpuNumPhysical
	}

	cpuUsedPercents, err := cpu.Percent(0, true)
	if err != nil {
		logSystemMetricErrorOnce("cpu_percent", err)
	} else if averagePercent, ok := averageCPUPercent(cpuUsedPercents); ok {
		sys.CPUUsedPercent = averagePercent
	} else {
		logger.Infof("cpu.Percent returned an empty result")
	}

	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		logSystemMetricErrorOnce("memory", err)
	} else {
		sys.MemoryTotal = virtualMemory.Total
		sys.MemoryFree = virtualMemory.Free
		sys.MemoryUsedPercent = virtualMemory.UsedPercent
	}
}

// logSystemMetricErrorOnce 对系统指标采集失败做降噪，避免未实现平台反复刷日志。
func logSystemMetricErrorOnce(metric string, err error) {
	if err == nil {
		return
	}
	if strings.Contains(strings.ToLower(err.Error()), "not implemented") {
		return
	}
	key := metric + ":" + err.Error()
	if _, loaded := systemMetricLoggedErrors.LoadOrStore(key, struct{}{}); loaded {
		return
	}
	logger.Infof("%s", err)
}

// averageCPUPercent 将每个逻辑核的 CPU 占用率合并成一个平均值，避免 TUI 在多核机器上按核心显示多行。
func averageCPUPercent(cpuUsedPercents []float64) (float64, bool) {
	if len(cpuUsedPercents) == 0 {
		return 0, false
	}
	total := 0.0
	for _, percent := range cpuUsedPercents {
		total += percent
	}
	return total / float64(len(cpuUsedPercents)), true
}
