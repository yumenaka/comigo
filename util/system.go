package util

import (
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/logger"
)

// CheckPort 检测端口是否可用
func CheckPort(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		logger.Info(os.Stderr, locale.GetString("cannot_listen"), port, err)
		return false
	}
	err = ln.Close()
	if err != nil {
		logger.Infof(locale.GetString("check_pork_error")+"%d", port)
		return false
	}
	//logger.Infof("TCP Port %q is available", port)
	return true
}

// GetFreePort 获取一个空闲可用的端口号
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func(l *net.TCPListener) {
		err := l.Close()
		if err != nil {
			log.Println(err)
		}
	}(l)
	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetIPList 获取本机IP列表
func GetIPList() (IPList []string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if i.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := i.Addrs()
		if err != nil {
			logger.Infof(locale.GetString("get_ip_error")+"%v", err)
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			IPList = append(IPList, ip.String())
		}
	}
	return IPList, err
}

// GetOutboundIP 获取本机的首选出站IP
// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		time.Sleep(3 * time.Second)
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func OpenBrowser(uri string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("CMD", "/C", "start", uri)
		if err := cmd.Start(); err != nil {
			logger.Infof(locale.GetString("open_browser_error")+"%s", err.Error())
		}
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", uri)
		if err := cmd.Start(); err != nil {
			logger.Infof(locale.GetString("open_browser_error")+"%s", err.Error())
		}
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("xdg-open", uri)
	}
}

// SystemStatus Documentation: https://pkg.go.dev/github.com/shirou/gopsutil
// 获取服务器当前状况
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
		logger.Infof("%s", err)
	} else {
		sys.CPUNumLogical = CPUNumLogical
	}
	CPUNumPhysical, err := cpu.Counts(false)
	if err != nil {
		logger.Infof("%s", err)
	} else {
		sys.CPUNumPhysical = CPUNumPhysical
	}
	CPUUsedPercent, err := cpu.Percent(0, false)
	if err != nil {
		logger.Infof("%s", err)
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
		logger.Infof("%s", err)
	} else {
		sys.MemoryTotal = v.Total
		sys.MemoryFree = v.Free
		sys.MemoryUsedPercent = v.UsedPercent
	}
	//// almost every return value is a struct
	//logger.Infof("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	//// convert to JSON. String() is also implemented
	//logger.Infof(v)

	//hostname, err := os.Hostname()
	//if err == nil {
	//	logger.Infof(hostname)
	//}
	return sys
}

//// 获取mac地址列表,暂时用不着
//func GetMacAddrList() (macAddrList []string) {
//	netInterfaces, err := net.Interfaces()
//	if err != nil {
//		logger.Infof(locale.GetString("check_mac_error")+": %v", err)
//		return macAddrList
//	}
//	//for _, netInterface := range netInterfaces {
//	//	macAddr := netInterface.HardwareAddr.String()
//	//	if len(macAddr) == 0 {
//	//		continue
//	//	}
//	//	macAddrList = append(macAddrList, macAddr)
//	//}
//	for _, netInterface := range netInterfaces {
//		flags := netInterface.Flags.String()
//		if strings.Contains(flags, "up") && strings.Contains(flags, "broadcast") {
//			macAddrList = append(macAddrList, netInterface.HardwareAddr.String())
//		}
//	}
//	return macAddrList
//}
