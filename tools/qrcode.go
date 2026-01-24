package tools

import (
	"strconv"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// PrintAllReaderURL 打印阅读链接
func PrintAllReaderURL(Port int, OpenBrowserFlag bool, PrintAllPossibleQRCode bool, ServerHost string, DisableLAN bool, customTLS bool, autoTLS bool, etcStr string) {
	protocol := "http://"
	if customTLS || autoTLS {
		protocol = "https://"
	}
	localURL := protocol + "127.0.0.1:" + strconv.Itoa(Port) + etcStr
	logger.Info(locale.GetString("local_reading") + localURL + etcStr)
	// 打开浏览器
	if OpenBrowserFlag {
		go OpenBrowserByURL(protocol + "127.0.0.1:" + strconv.Itoa(Port) + etcStr)
	}
	if !DisableLAN {
		printURLAndQRCode(Port, PrintAllPossibleQRCode, ServerHost, protocol, customTLS, autoTLS, etcStr)
	}
}

func printURLAndQRCode(port int, PrintAllPossibleQRCode bool, ServerHost string, protocol string, customTLS bool, autoTLS bool, etcStr string) {
	// 打印指定的服务器地址
	if ServerHost != "" {
		readURL := protocol + ServerHost + ":" + strconv.Itoa(port) + etcStr
		// 自定义 TLS 时，如果是 443 端口，则不需要加端口号
		if customTLS && port == 443 {
			readURL = protocol + ServerHost + ":" + strconv.Itoa(port) + etcStr
		}
		// 自动 TLS 时，目前只支持443, 不需要加端口号
		if autoTLS {
			readURL = protocol + ServerHost + etcStr
		}
		// 打印指定的服务器地址
		logger.Info(locale.GetString("reading_url_maybe") + readURL)
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
			logger.Info(locale.GetString("reading_url_maybe") + readURL)
			PrintQRCode(readURL)
		}
	} else {
		// 只打印本机的首选出站IP
		OutIP := GetOutboundIP().String()
		readURL := protocol + OutIP + ":" + strconv.Itoa(port) + etcStr
		logger.Info(locale.GetString("reading_url_maybe") + readURL)
		PrintQRCode(readURL)
	}
}

func PrintQRCode(text string) {
	// or https://github.com/mdp/qrterminal
	obj := qrcodeTerminal.New()
	obj.Get(text).Print()
}
