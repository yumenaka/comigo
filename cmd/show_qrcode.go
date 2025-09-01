package cmd

import (
	"fmt"

	"github.com/sanity-io/litter"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util"
)

func ShowQRCode() {
	// 如果只有一本书，二维码展示的 URL 需要附加参数，让读者可以直接去读这本书
	etcStr := ""
	if model.MainStoreGroup.GetBooksNumber() == 1 {
		bookList := model.MainStoreGroup.ListBooks()
		if len(bookList) == 1 {
			etcStr = fmt.Sprintf("/#/%s/%s", config.GetDefaultMode(), bookList[0].BookID)
		}
	}
	// 判断是否启用 TLS
	// 如果配置文件中有证书和密钥文件，则启用 TLS
	enableTLS := config.GetCertFile() != "" && config.GetKeyFile() != ""
	outIP := config.GetHost()
	if config.GetHost() == "" {
		outIP = util.GetOutboundIP().String()
	}
	// 打印二维码
	util.PrintAllReaderURL(
		config.GetPort(),
		config.GetOpenBrowser(),
		config.GetPrintAllPossibleQRCode(),
		outIP,
		config.GetDisableLAN(),
		enableTLS,
		etcStr,
	)

	// 打印配置，调试用
	if config.GetDebug() {
		litter.Dump(config.GetCfg())
	}
	// ”如何快捷键退出“的文字提示
	fmt.Println(locale.GetString("ctrl_c_hint"))
}
