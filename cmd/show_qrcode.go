package cmd

import (
	"context"
	"fmt"

	"github.com/sanity-io/litter"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
)

func ShowQRCode() {
	// 如果只有一本书，二维码展示的 URL 需要附加参数，让读者可以直接去读这本书
	etcStr := ""
	if model.IStore.GetAllBooksNumber() == 1 {
		bookList := model.IStore.ListBooks()
		if len(bookList) == 1 {
			etcStr = fmt.Sprintf("/#/%s/%s", config.GetCfg().DefaultMode, bookList[0].BookID)
		}
	}
	// 判断是否启用 TLS
	// 如果配置文件中有证书和密钥文件，则启用 TLS
	enableTLS := config.GetCfg().CertFile != "" && config.GetCfg().KeyFile != ""
	outIP := config.GetCfg().Host
	if config.GetCfg().Host == "" {
		outIP = tools.GetOutboundIP().String()
	}
	// 打印二维码
	tools.PrintAllReaderURL(
		config.GetCfg().Port,
		config.GetCfg().OpenBrowser,
		config.GetCfg().PrintAllPossibleQRCode,
		outIP,
		config.GetCfg().DisableLAN,
		enableTLS,
		etcStr,
	)

	// 打印配置，调试用
	if config.GetCfg().Debug {
		litter.Dump(config.GetCfg())
	}
	// ”如何快捷键退出“的文字提示
	fmt.Println(locale.GetString("ctrl_c_hint"))
	// 打印 Tailscale 访问地址的二维码
	ShowQRCodeTailscale(context.Background())
}
