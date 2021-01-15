package common

import (
	"fmt"
	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/yumenaka/comi/locale"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

//打印阅读链接
func PrintAllReaderURL() {
	localURL := "http://127.0.0.1:" + strconv.Itoa(Config.Port)
	fmt.Println(locale.GetString("local_reading") + localURL)
	//PrintQRCode(localURL)
	//打开浏览器
	if Config.OpenBrowser {
		OpenBrowser("http://127.0.0.1:" + strconv.Itoa(Config.Port))
		if Config.EnableFrpcServer {
			OpenBrowser("http://" + Config.FrpConfig.ServerAddr + ":" + strconv.Itoa(Config.FrpConfig.RemotePort))
		}
	}
	if !Config.DisableLAN {
		printURLAndQRCode(Config.Port)
	}
}

func printURLAndQRCode(port int) {
	//启用Frp的时候
	if Config.EnableFrpcServer {
		readURL := "http://" + Config.FrpConfig.ServerAddr + ":" + strconv.Itoa(Config.FrpConfig.RemotePort)
		fmt.Println(locale.GetString("frp_reading_url_is")  + readURL)
		PrintQRCode(readURL)
	}
	if Config.ServerHost != "" {
		readURL := "http://" + Config.ServerHost + ":" + strconv.Itoa(port)
		fmt.Println(locale.GetString("reading_url_maybe")  + readURL)
		PrintQRCode(readURL)
		return
	}
	//打印所有可用网卡IP
	if Config.PrintAllIP {
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

//中断处理：程序被中断的时候，清理临时文件
func SetupCloseHander() {
	c := make(chan os.Signal, 2)
	//SIGHUP（挂起）, SIGINT（中断）或 SIGTERM（终止）默认会使得程序退出。
	//1、SIGHUP 信号在用户终端连接(正常或非正常)结束时发出。
	//2、syscall.SIGINT 和 os.Interrupt 是同义词,按下 CTRL+C 时发出。
	//3、SIGTERM（终止）:kill终止进程,允许程序处理问题后退出。
	//4.syscall.SIGHUP,终端控制进程结束(终端连接断开)
	//5、syscall.SIGQUIT，CTRL+\ 退出
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		<-c
		fmt.Println("\r"+locale.GetString("start_clear_file"))
		deleteTempFiles()
		os.Exit(0)
	}()
}

func InitReadingBook() (err error) {
	//准备解压，设置图片文件夹
	if ReadingBook.IsFolder {
		PictureDir = ReadingBook.FilePath
		ReadingBook.ExtractComplete = true
		ReadingBook.ExtractNum = ReadingBook.PageNum
	} else {
		err = SetTempDir()
		if err != nil {
			fmt.Println(locale.GetString("temp_folder_error"), err)
			return err
		}
		PictureDir = TempDir
		err = ExtractArchive(&ReadingBook)
		if err != nil {
			fmt.Println(locale.GetString("file_not_found"))
			return err
		}
		ReadingBook.SetArchiveBookName(ReadingBook.FilePath) //设置书名
	}
	//服务器分析图片分辨率
	if Config.CheckImageInServer {
		ReadingBook.ScanAllImageGo() //扫描所有图片，取得分辨率信息，使用了协程
	}
	return err
}

//设置临时文件夹，退出时会被清理
func SetTempDir() (err error) {
	if TempDir != "" {
		return err
	}
	TempDir, err = ioutil.TempDir("", "comic_cache_A8cG")
	if err != nil {
		println(locale.GetString("temp_folder_create_error"))
	} else {
		fmt.Println(locale.GetString("temp_folder_path") + TempDir)
	}
	return err
}

func deleteTempFiles() {
	fmt.Println(locale.GetString("clear_temp_file_start"))
	if strings.Contains(TempDir, "comic_cache_A8cG") { //判断文件夹前缀，避免删错文件
		err := os.RemoveAll(TempDir)
		if err != nil {
			fmt.Println(locale.GetString("clear_temp_file_error") + TempDir)
		} else {
			fmt.Println(locale.GetString("clear_temp_file_completed") + TempDir)
		}
	}
	deleteOldTempFiles()
}

//根据权限，清理老文件可能失败
func deleteOldTempFiles() {
	tempDirUpperFolder := TempDir
	post := strings.LastIndex(TempDir, "/") //Unix风格的路径分隔符
	if post == -1 {
		post = strings.LastIndex(TempDir, "\\") //windows风格的分隔符
	}
	if post != -1 {
		tempDirUpperFolder = string([]rune(TempDir)[:post]) //为了防止中文字符被错误截断，先转换成rune，再转回来
		fmt.Println(locale.GetString("temp_folder_path"), tempDirUpperFolder)
	}
	files, err := ioutil.ReadDir(tempDirUpperFolder)
	if err != nil {
		fmt.Println(err)
	}
	for _, fi := range files {
		if fi.IsDir() {
			oldTempDir := tempDirUpperFolder + "/" + fi.Name()
			if strings.Contains(oldTempDir, "comic_cache_A8cG") { //判断文件夹前缀，避免删错文件
				err := os.RemoveAll(oldTempDir)
				if err != nil {
					fmt.Println(locale.GetString("clear_temp_file_error") + oldTempDir)
				} else {
					fmt.Println(locale.GetString("clear_temp_file_completed") + oldTempDir)
				}
			}
		}
	}
}
