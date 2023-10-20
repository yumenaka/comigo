package routers

import (
	"github.com/yumenaka/comi/logger"
	"strings"

	"github.com/sanity-io/litter"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/types"
	"github.com/yumenaka/comi/util"
)

func showQRCode() {
	//cmd打印链接二维码.如果只有一本书，就直接打开那本书.
	etcStr := ""
	//只有一本书的时候，URL需要附加的参数
	if types.GetBooksNumber() == 1 {
		bookList, err := types.GetAllBookInfoList("name")
		if err != nil {
			logger.Info(err)
		}
		if len(bookList.BookInfos) == 1 {
			etcStr = "/#/scroll/" + bookList.BookInfos[0].BookID
		}
		if config.Config.DefaultMode != "" {
			etcStr = "/#/" + strings.ToLower(config.Config.DefaultMode) + "/" + bookList.BookInfos[0].BookID
		}
	}
	enableTls := config.Config.CertFile != "" && config.Config.KeyFile != ""
	OutIP := config.Config.Host
	if config.Config.Host == "DefaultHost" {
		OutIP = util.GetOutboundIP().String()
	}
	util.PrintAllReaderURL(
		config.Config.Port,
		config.Config.OpenBrowser,
		config.Config.PrintAllPossibleQRCode,
		OutIP,
		config.Config.DisableLAN,
		enableTls,
		etcStr)
	//打印配置，调试用
	if config.Config.Debug {
		litter.Dump(config.Config)
	}
	logger.Info(locale.GetString("ctrl_c_hint"))
}
