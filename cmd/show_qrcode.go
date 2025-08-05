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
	// 如果只有一本书，URL 需要附加的参数
	etcStr := ""
	if model.MainStores.GetBooksNumber() == 1 {
		bookList := model.MainStores.ListBooks()
		if len(bookList) == 1 {
			etcStr = fmt.Sprintf("/#/%s/%s", config.GetDefaultMode(), bookList[0].BookID)
		}
	}

	enableTLS := config.GetCertFile() != "" && config.GetKeyFile() != ""
	outIP := config.GetHost()
	if config.GetHost() == "" {
		outIP = util.GetOutboundIP().String()
	}

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

	fmt.Println(locale.GetString("ctrl_c_hint"))
}
