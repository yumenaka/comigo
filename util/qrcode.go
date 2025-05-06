package util

import (
	"fmt"
	"strconv"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/util/logger"
)

// PrintAllReaderURL 打印阅读链接
func PrintAllReaderURL(Port int, OpenBrowserFlag bool, PrintAllPossibleQRCode bool, ServerHost string, DisableLAN bool, enableTLS bool, etcStr string) {
	protocol := "http://"
	if enableTLS {
		protocol = "https://"
	}
	localURL := protocol + "127.0.0.1:" + strconv.Itoa(Port) + etcStr
	fmt.Println(locale.GetString("local_reading") + localURL + etcStr)
	// 打开浏览器
	if OpenBrowserFlag {
		go OpenBrowser(protocol + "127.0.0.1:" + strconv.Itoa(Port) + etcStr)
	}
	if !DisableLAN {
		printURLAndQRCode(Port, PrintAllPossibleQRCode, ServerHost, protocol, etcStr)
	}
}

func printURLAndQRCode(port int, PrintAllPossibleQRCode bool, ServerHost string, protocol string, etcStr string) {
	// 打印指定的服务器地址
	if ServerHost != "" {
		readURL := protocol + ServerHost + ":" + strconv.Itoa(port) + etcStr
		fmt.Println(locale.GetString("reading_url_maybe") + readURL)
		PrintQRCode(readURL)
		return
	}
	// 打印所有可用网卡IP
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
		// 只打印本机的首选出站IP
		OutIP := GetOutboundIP().String()
		readURL := protocol + OutIP + ":" + strconv.Itoa(port) + etcStr
		fmt.Println(locale.GetString("reading_url_maybe") + readURL)
		PrintQRCode(readURL)
	}
}

func PrintQRCode(text string) {
	// or https://github.com/mdp/qrterminal
	obj := qrcodeTerminal.New()
	obj.Get(text).Print()
}
