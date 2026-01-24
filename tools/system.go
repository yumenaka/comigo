package tools

import (
	"context"
	"log"
	"net"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
	httpChecker "wait4x.dev/v3/checker/http"
	"wait4x.dev/v3/waiter"
)

// TrackTIme 计算耗时
// 使用时只需要写一行：defer TrackTIme(time.Now())
func TrackTIme(pre time.Time) time.Duration {
	elapsed := time.Since(pre)
	logger.Infof(locale.GetString("log_time_elapsed")+"\n", elapsed)
	return elapsed
}

// CheckPort 检测端口是否可用
func CheckPort(port uint16) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		logger.Info(locale.GetString("cannot_listen"), port, err)
		return false
	}
	err = ln.Close()
	if err != nil {
		logger.Infof(locale.GetString("check_port_error"), port)
		return false
	}
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

func OpenBrowser(isHTTPS bool, host string, port int) {
	if isHTTPS {
		OpenBrowserByURL("https://" + host + ":" + strconv.Itoa(port))
	} else {
		OpenBrowserByURL("http://" + host + ":" + strconv.Itoa(port))
	}
}

// OpenBrowserByURL 打开浏览器，为了防止阻塞，需要使用go关键字调用
func OpenBrowserByURL(uri string) {
	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create an HTTP checker with multiple validations
	checker := httpChecker.New(
		uri,
		httpChecker.WithTimeout(3*time.Second),
		httpChecker.WithExpectStatusCode(200),
		// httpChecker.WithExpectBodyJSON("status"), // Check that 'status' field exists in JSON
		// httpChecker.WithExpectBodyRegex(`"healthy":\s*true`), // Regex to check response
		// httpChecker.WithExpectHeader("Content-Type=application/json"),
		// httpChecker.WithRequestHeaders(headers),
		// httpChecker.WithRequestBody(requestBody),
		httpChecker.WithInsecureSkipTLSVerify(true), // Skip TLS verification
	)

	// Wait for the API to be available and responding correctly
	logger.Info(locale.GetString("log_waiting_for_api_health"))

	err := waiter.WaitContext(
		ctx,
		checker,
		waiter.WithTimeout(2*time.Minute),
		waiter.WithInterval(500*time.Millisecond),
		waiter.WithBackoffPolicy(waiter.BackoffPolicyExponential),
	)
	if err != nil {
		log.Fatalf("API health check failed: %v", err)
	}
	logger.Info(locale.GetString("log_api_healthy_ready"))

	// 打开浏览器（Windows 使用 ShellExecute，避免闪黑框）
	if err := openURL(uri); err != nil {
		logger.Infof(locale.GetString("open_browser_error")+"%s", err.Error())
	}
}

// SystemStatus Documentation: https://pkg.go.dev/github.com/shirou/gopsutil
// 获取服务器当前状况
type SystemStatus struct {
	// CPU相关
	CPUNumLogical  int     `json:"cpu_num_logical_total"`
	CPUNumPhysical int     `json:"cpu_num_physical"`
	CPUUsedPercent float64 `json:"cpu_used_percent"`
	// 内存相关
	MemoryTotal       uint64  `json:"memory_total"`
	MemoryFree        uint64  `json:"memory_free"`
	MemoryUsedPercent float64 `json:"memory_used_percent"`
	// 设备描述
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
	// 获取物理和逻辑核数,以及CPU、内存整体使用率
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
		// p := 0.0
		// if len(CPUUsedPercent) > 1 {
		//	for _, value := range CPUUsedPercent {
		//		p += value
		//	}
		//	p = p / float64(len(CPUUsedPercent))
		// } else if len(CPUUsedPercent) == 1 {
		//	p = CPUUsedPercent[0]
		// }
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
	// // almost every return value is a struct
	// logger.Infof("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	// // convert to JSON. String() is also implemented
	// logger.Infof(v)

	// hostname, err := os.Hostname()
	// if err == nil {
	//	logger.Infof(hostname)
	// }
	return sys
}

// // 获取mac地址列表,暂时用不着
// func GetMacAddrList() (macAddrList []string) {
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
// }

// ServerStatus 服务器当前状况
type ServerStatus struct {
	ServerName            string       // 服务器描述
	ServerHost            string       // ServerHost 服务器主机或 IP 地址。
	ServerPort            uint16       // ServerPort 服务运行的端口号。
	TailscaleAuthURL      string       // Tailscale身份验证URL（如果适用）
	TailscaleUrl          string       // Tailscale阅读地址（如果有）
	NumberOfBooks         int          // 当前拥有的书籍总数
	SupportUploadFile     bool         // 是否支持上传文件
	ClientIP              string       // 客户端IP
	OSInfo                SystemStatus // 系统信息
	ReScanServiceEnable   bool         // 是否启用自动扫描服务
	ReScanServiceInterval int          // 自动扫描服务间隔（分钟）
}

type ConfigInterface interface {
	GetHost() string
	GetPort() int
	GetEnableUpload() bool
}

type ServerInfoParams struct {
	Cfg                   ConfigInterface
	Version               string
	AllBooksNumber        int
	ClientIP              string
	ReScanServiceEnable   bool
	ReScanServiceInterval int
}

func GetServerInfo(params ServerInfoParams) *ServerStatus {
	serverName := "Comigo " + params.Version
	configHost := params.Cfg.GetHost()
	port := params.Cfg.GetPort()
	enableUpload := params.Cfg.GetEnableUpload()
	// 本机首选出站IP
	host := ""
	if configHost == "" {
		host = GetOutboundIP().String()
	} else {
		host = configHost
	}
	serverStatus := ServerStatus{
		ServerName:            serverName,
		ServerHost:            host,
		ServerPort:            uint16(port),
		SupportUploadFile:     enableUpload,
		NumberOfBooks:         params.AllBooksNumber,
		ClientIP:              params.ClientIP,
		ReScanServiceEnable:   params.ReScanServiceEnable,
		ReScanServiceInterval: params.ReScanServiceInterval,
		OSInfo:                GetSystemStatus(),
	}
	return &serverStatus
}
