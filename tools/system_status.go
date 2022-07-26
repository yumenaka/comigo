package tools

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Documentation: https://pkg.go.dev/github.com/shirou/gopsutil

// SystemStatus 服务器当前状况
type SystemStatus struct {
	//CPU相关
	CPUNumLogical  int     `json:"cpu_num_logical_total"`
	CPUNumPhysical int     `json:"cpu_num_physical"`
	CPUUsedPercent float64 `json:"cpu_used_percent"`
	//内存相关
	MemoryTotal       uint64  `json:"memory_total"`
	MemoryFree        uint64  `json:"memory_free"`
	MemoryUsedPercent float64 `json:"memory_used_percent"`
	//设备描述
	Description string `json:"description"`
}

func GetSystemStatus() SystemStatus {
	sys := SystemStatus{
		Description:       runtime.GOOS + " " + runtime.GOARCH,
		CPUNumLogical:     -1,
		CPUNumPhysical:    -1,
		CPUUsedPercent:    0,
		MemoryTotal:       0.0,
		MemoryFree:        0.0,
		MemoryUsedPercent: 0,
	}
	//获取物理和逻辑核数,以及CPU、内存整体使用率
	CPUNumLogical, err := cpu.Counts(true)
	if err != nil {
		fmt.Println(err)
	} else {
		sys.CPUNumLogical = CPUNumLogical
	}
	CPUNumPhysical, err := cpu.Counts(false)
	if err != nil {
		fmt.Println(err)
	} else {
		sys.CPUNumPhysical = CPUNumPhysical
	}
	CPUUsedPercent, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println(err)
	} else {
		//p := 0.0
		//if len(CPUUsedPercent) > 1 {
		//	for _, value := range CPUUsedPercent {
		//		p += value
		//	}
		//	p = p / float64(len(CPUUsedPercent))
		//} else if len(CPUUsedPercent) == 1 {
		//	p = CPUUsedPercent[0]
		//}
		sys.CPUUsedPercent = CPUUsedPercent[0]
	}
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
	} else {
		sys.MemoryTotal = v.Total
		sys.MemoryFree = v.Free
		sys.MemoryUsedPercent = v.UsedPercent
	}
	//// almost every return value is a struct
	//fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	//// convert to JSON. String() is also implemented
	//fmt.Println(v)

	//hostname, err := os.Hostname()
	//if err == nil {
	//	fmt.Println(hostname)
	//}

	return sys
}
