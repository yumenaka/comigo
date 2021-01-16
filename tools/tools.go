package tools

import (
	"fmt"
	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/yumenaka/comi/locale"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

)

//打印阅读链接
func PrintAllReaderURL(Port int,OpenBrowserFlag bool,EnableFrpcServer bool,PrintAllIP bool,ServerHost string,ServerAddr string,FrpRemotePort int,DisableLAN bool) {
	localURL := "http://127.0.0.1:" + strconv.Itoa(Port)
	fmt.Println(locale.GetString("local_reading") + localURL)
	//PrintQRCode(localURL)
	//打开浏览器
	if OpenBrowserFlag {
		OpenBrowser("http://127.0.0.1:" + strconv.Itoa(Port))
		if EnableFrpcServer {
			OpenBrowser("http://" + ServerAddr + ":" + strconv.Itoa(FrpRemotePort))
		}
	}
	if !DisableLAN {
		printURLAndQRCode(Port,EnableFrpcServer,PrintAllIP,ServerHost,ServerAddr,FrpRemotePort)
	}
}

func printURLAndQRCode(port int,EnableFrpcServer bool,PrintAllIP bool,ServerHost string, ServerAddr string, FrpRemotePort int) {
	//启用Frp的时候
	if EnableFrpcServer {
		readURL := "http://" + ServerAddr + ":" + strconv.Itoa(FrpRemotePort)
		fmt.Println(locale.GetString("frp_reading_url_is")  + readURL)
		PrintQRCode(readURL)
	}
	if ServerHost != "" {
		readURL := "http://" + ServerHost + ":" + strconv.Itoa(port)
		fmt.Println(locale.GetString("reading_url_maybe")  + readURL)
		PrintQRCode(readURL)
		return
	}
	//打印所有可用网卡IP
	if PrintAllIP {
		IPList, err := GetIPList()
		if err != nil {
			fmt.Printf(locale.GetString("get_ip_error")+" %v", err)
		}
		for _, IP := range IPList {
			readURL := "http://" + IP + ":" + strconv.Itoa(port)
			fmt.Println(locale.GetString("reading_url_maybe") + readURL)
			PrintQRCode(readURL)
		}
	} else {
		//只打印出口IP
		OutIP := GetOutboundIP().String()
		readURL := "http://" + OutIP + ":" + strconv.Itoa(port)
		fmt.Println(locale.GetString("reading_url_maybe") + readURL)
		PrintQRCode(readURL)
	}

}

func PrintQRCode(text string) {
	obj := qrcodeTerminal.New()
	obj.Get(text).Print()
}

//检测端口是否可用
func CheckPort(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Fprintf(os.Stderr, locale.GetString("cannot_listen") +"%q: %s", port, err)
		return false
	}
	err = ln.Close()
	if err != nil {
		fmt.Println( locale.GetString("check_pork_error")+ strconv.Itoa(port))
	}
	//fmt.Printf("TCP Port %q is available", port)
	return true
}

//获取本机IP列表
func GetIPList() (IPList []string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf(locale.GetString("get_ip_error")+"%v", err)
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

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

// 获取mac地址列表,暂时用不着
func GetMacAddrList() (macAddrList []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf(locale.GetString("check_mac_error")+": %v", err)
		return macAddrList
	}
	//for _, netInterface := range netInterfaces {
	//	macAddr := netInterface.HardwareAddr.String()
	//	if len(macAddr) == 0 {
	//		continue
	//	}
	//	macAddrList = append(macAddrList, macAddr)
	//}
	for _, netInterface := range netInterfaces {
		flags := netInterface.Flags.String()
		if strings.Contains(flags, "up") && strings.Contains(flags, "broadcast") {
			macAddrList = append(macAddrList, netInterface.HardwareAddr.String())
		}
	}
	return macAddrList

}

//判断所给路径文件或文件夹是否存在
func ChickFileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func OpenBrowser(uri string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("CMD", "/C", "start", uri)
		if err := cmd.Start(); err != nil {
			fmt.Println(locale.GetString("open_browser_error"))
			fmt.Println(err.Error())
		}
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", uri)
		if err := cmd.Start(); err != nil {
			fmt.Println(locale.GetString("open_browser_error"))
			fmt.Println(err.Error())
		}
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("xdg-open", uri)
	}
}




