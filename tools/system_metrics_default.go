//go:build !ios

package tools

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/yumenaka/comigo/tools/logger"
)

func populateSystemMetrics(sys *SystemStatus) {
	// 获取物理和逻辑核数,以及 CPU、内存整体使用率。
	cpuNumLogical, err := cpu.Counts(true)
	if err != nil {
		logger.Infof("%s", err)
	} else {
		sys.CPUNumLogical = cpuNumLogical
	}

	cpuNumPhysical, err := cpu.Counts(false)
	if err != nil {
		logger.Infof("%s", err)
	} else {
		sys.CPUNumPhysical = cpuNumPhysical
	}

	cpuUsedPercent, err := cpu.Percent(0, false)
	if err != nil {
		logger.Infof("%s", err)
	} else if len(cpuUsedPercent) > 0 {
		sys.CPUUsedPercent = cpuUsedPercent[0]
	} else {
		logger.Infof("cpu.Percent returned an empty result")
	}

	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		logger.Infof("%s", err)
	} else {
		sys.MemoryTotal = virtualMemory.Total
		sys.MemoryFree = virtualMemory.Free
		sys.MemoryUsedPercent = virtualMemory.UsedPercent
	}
}
