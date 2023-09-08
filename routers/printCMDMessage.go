package routers

import (
	"fmt"
	"github.com/sanity-io/litter"
	"github.com/yumenaka/comi/book"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/tools"
	"strings"
)

// 6、printQRCodeInCMD
func printQRCodeInCMD() {
	//cmd打印链接二维码.如果只有一本书，就直接打开那本书.
	etcStr := ""
	//只有一本书的时候，URL需要附加的参数
	if book.GetBooksNumber() == 1 {
		bookList, err := book.GetAllBookInfoList("name")
		if err != nil {
			fmt.Println(err)
		}
		if len(bookList.BookInfos) == 1 {
			etcStr = "/#/scroll/" + bookList.BookInfos[0].BookID
		}
		if common.Config.DefaultMode != "" {
			etcStr = "/#/" + strings.ToLower(common.Config.DefaultMode) + "/" + bookList.BookInfos[0].BookID
		}
	}
	enableTls := common.Config.CertFile != "" && common.Config.KeyFile != ""
	OutIP := common.Config.Host
	if common.Config.Host == "DefaultHost" {
		OutIP = tools.GetOutboundIP().String()
	}
	tools.PrintAllReaderURL(
		common.Config.Port,
		common.Config.OpenBrowser,
		common.Config.PrintAllPossibleQRCode,
		OutIP,
		common.Config.DisableLAN,
		enableTls,
		etcStr)
	//打印配置，调试用
	if common.Config.Debug {
		litter.Dump(common.Config)
	}
	fmt.Println(locale.GetString("ctrl_c_hint"))
}
