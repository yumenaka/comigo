package routers

import (
	"fmt"
	"strings"

	"github.com/sanity-io/litter"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
)

func showQRCode() {
	// 如果只有一本书，URL 需要附加的参数
	etcStr := ""
	if entity.GetBooksNumber() == 1 {
		bookList, err := entity.GetAllBookInfoList("name")
		if err != nil {
			logger.Infof("Error getting book list: %s", err)
			return
		}
		if len(bookList.BookInfos) == 1 {
			mode := "scroll"
			if config.Config.DefaultMode != "" {
				mode = strings.ToLower(config.Config.DefaultMode)
			}
			etcStr = fmt.Sprintf("/#/%s/%s", mode, bookList.BookInfos[0].BookID)
		}
	}

	enableTLS := config.Config.CertFile != "" && config.Config.KeyFile != ""
	outIP := config.Config.Host
	if config.Config.Host == "DefaultHost" {
		outIP = util.GetOutboundIP().String()
	}

	util.PrintAllReaderURL(
		config.Config.Port,
		config.Config.OpenBrowser,
		config.Config.PrintAllPossibleQRCode,
		outIP,
		config.Config.DisableLAN,
		enableTLS,
		etcStr,
	)

	// 打印配置，调试用
	if config.Config.Debug {
		litter.Dump(config.Config)
	}

	fmt.Println(locale.GetString("ctrl_c_hint"))
}
