package util

import (
	"fmt"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
	"strconv"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
)

// PrintAllReaderURL 打印阅读链接
func PrintAllReaderURL(Port int, OpenBrowserFlag bool, PrintAllPossibleQRCode bool, ServerHost string, DisableLAN bool, enableTls bool, etcStr string) {
	protocol := "http://"
	if enableTls {
		protocol = "https://"
	}
	localURL := protocol + "127.0.0.1:" + strconv.Itoa(Port) + etcStr
	fmt.Println(locale.GetString("local_reading") + localURL + etcStr)
	//打开浏览器
	if OpenBrowserFlag {
		OpenBrowser(protocol + "127.0.0.1:" + strconv.Itoa(Port) + etcStr)
	}
	if !DisableLAN {
		printURLAndQRCode(Port, PrintAllPossibleQRCode, ServerHost, protocol, etcStr)
	}
}

func printURLAndQRCode(port int, PrintAllPossibleQRCode bool, ServerHost string, protocol string, etcStr string) {

	if ServerHost != "" {
		readURL := protocol + ServerHost + ":" + strconv.Itoa(port) + etcStr
		fmt.Println(locale.GetString("reading_url_maybe") + readURL)
		PrintQRCode(readURL)
		return
	}
	//打印所有可用网卡IP
	if PrintAllPossibleQRCode {
		IPList, err := GetIPList()
		if err != nil {
			logger.Infof(locale.GetString("get_ip_error")+" %v", err)
		}
		for _, IP := range IPList {
			readURL := protocol + IP + ":" + strconv.Itoa(port) + etcStr
			fmt.Println(locale.GetString("reading_url_maybe") + readURL)
			PrintQRCode(readURL)
		}
	} else {
		//只打印本机的首选出站IP
		OutIP := GetOutboundIP().String()
		readURL := protocol + OutIP + ":" + strconv.Itoa(port) + etcStr
		fmt.Println(locale.GetString("reading_url_maybe") + readURL)
		PrintQRCode(readURL)
	}
}

func PrintQRCode(text string) {
	obj := qrcodeTerminal.New()
	obj.Get(text).Print()
}
