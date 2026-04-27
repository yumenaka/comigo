package tools

import (
	"errors"
	"strconv"
	"strings"

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
	logger.Info(locale.GetString("local_reading") + localURL)
	// 打开浏览器
	if OpenBrowserFlag {
		go OpenBrowserByURL(localURL)
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

// RenderQRCodeLinesTerminal 将文本渲染为适合终端面板的二维码行内容。
//
// 说明：
//   - 该实现使用 qrcode-terminal-go 生成二维码 bitmap 的终端字符表现；
//   - 再将“单行模块”（█/空）按两行组合为 Unicode 半高块字符（▀▄█），从而保持与原 TUI
//     渲染方式接近的高度比例（避免 QR 面板过高被截断）。
func RenderQRCodeLinesTerminal(text string) ([]string, error) {
	if strings.TrimSpace(text) == "" {
		return []string{locale.GetString("tui_qr_unavailable")}, nil
	}

	obj := qrcodeTerminal.New()
	qr := obj.Get(text)
	if qr == nil {
		return nil, errors.New("qrcodeTerminal: failed to generate qr")
	}

	// qrcode-terminal-go 默认输出包含 ANSI 颜色转义与“两空格模块”；这里将其替换成无 ANSI 的字符。
	front := string(qrcodeTerminal.ConsoleColors.BrightBlack) // module on
	back := string(qrcodeTerminal.ConsoleColors.BrightWhite)  // module off

	raw := strings.TrimRight(string(*qr), "\n")
	raw = strings.ReplaceAll(raw, front, "█")
	raw = strings.ReplaceAll(raw, back, " ")

	rawLines := strings.Split(raw, "\n")
	if len(rawLines) == 0 {
		return []string{locale.GetString("tui_qr_unavailable")}, nil
	}

	// 将 rawLines 解析成模块矩阵（每个字符对应一个模块单元）。
	matrix := make([][]rune, len(rawLines))
	maxCols := 0
	for i, line := range rawLines {
		row := []rune(line)
		matrix[i] = row
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	for i := range matrix {
		if len(matrix[i]) < maxCols {
			pad := make([]rune, maxCols-len(matrix[i]))
			for j := range pad {
				pad[j] = ' '
			}
			matrix[i] = append(matrix[i], pad...)
		}
	}

	// 按两行组合为半高块：与 tea.go 原先 qrBlock 的含义保持一致。
	halfLines := make([]string, 0, (len(matrix)+1)/2)
	for row := 0; row < len(matrix); row += 2 {
		top := matrix[row]
		bottom := make([]rune, maxCols)
		for col := 0; col < maxCols; col++ {
			bottom[col] = ' '
		}
		if row+1 < len(matrix) {
			bottom = matrix[row+1]
		}

		var builder strings.Builder
		for col := 0; col < maxCols; col++ {
			tOn := top[col] == '█'
			bOn := bottom[col] == '█'
			switch {
			case tOn && bOn:
				builder.WriteRune('█')
			case tOn && !bOn:
				builder.WriteRune('▀')
			case !tOn && bOn:
				builder.WriteRune('▄')
			default:
				builder.WriteRune(' ')
			}
		}
		halfLines = append(halfLines, builder.String())
	}

	return halfLines, nil
}
