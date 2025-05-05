package cmd

import (
	"fmt"

	"github.com/sanity-io/litter"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
)

func ShowQRCode() {
	// 如果只有一本书，URL 需要附加的参数
	etcStr := ""
	if model.GetBooksNumber() == 1 {
		bookList, err := model.GetAllBookInfoList("name")
		if err != nil {
			logger.Infof("Error getting book list: %s", err)
			return
		}
		if len(bookList.BookInfos) == 1 {
			etcStr = fmt.Sprintf("/#/%s/%s", config.GetDefaultMode(), bookList.BookInfos[0].BookID)
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
