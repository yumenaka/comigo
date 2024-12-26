package routers

import (
	"fmt"
	"strings"

	"github.com/sanity-io/litter"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
)

func showQRCode() {
	// 如果只有一本书，URL 需要附加的参数
	etcStr := ""
	if model.GetBooksNumber() == 1 {
		bookList, err := model.GetAllBookInfoList("name")
		if err != nil {
			logger.Infof("Error getting book list: %s", err)
			return
		}
		if len(bookList.BookInfos) == 1 {
			mode := "scroll"
			if config.Cfg.DefaultMode != "" {
				mode = strings.ToLower(config.Cfg.DefaultMode)
			}
			etcStr = fmt.Sprintf("/#/%s/%s", mode, bookList.BookInfos[0].BookID)
		}
	}

	enableTLS := config.Cfg.CertFile != "" && config.Cfg.KeyFile != ""
	outIP := config.Cfg.Host
	if config.Cfg.Host == "DefaultHost" {
		outIP = util.GetOutboundIP().String()
	}

	util.PrintAllReaderURL(
		config.Cfg.Port,
		config.Cfg.OpenBrowser,
		config.Cfg.PrintAllPossibleQRCode,
		outIP,
		config.Cfg.DisableLAN,
		enableTLS,
		etcStr,
	)

	// 打印配置，调试用
	if config.Cfg.Debug {
		litter.Dump(config.Cfg)
	}

	fmt.Println(locale.GetString("ctrl_c_hint"))
}
