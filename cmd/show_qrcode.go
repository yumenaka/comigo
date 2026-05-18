package cmd

import (
	"context"
	"fmt"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

type readerLinkContext struct {
	etcStr    string
	outIP     string
	customTLS bool
	autoTLS   bool
}

// buildReaderLinkContext 汇总命令行阅读链接需要的参数，避免“打印二维码”和“打开浏览器”重复拼接 URL。
func buildReaderLinkContext() readerLinkContext {
	// 如果只有一本书，二维码展示的 URL 需要附加参数，让读者可以直接去读这本书
	etcStr := config.PrefixPath("/")
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}
	if len(allBooks) == 1 {
		etcStr = config.PrefixPath(fmt.Sprintf("/#/%s/%s", "scroll", allBooks[0].BookID))
	}

	// 判断是否启用 TLS
	// 如果配置文件中有证书和密钥文件，则启用 TLS
	outIP := config.GetCfg().Host
	if config.GetCfg().Host == "" {
		outIP = tools.GetOutboundIP().String()
	}
	return readerLinkContext{
		etcStr:    etcStr,
		outIP:     outIP,
		customTLS: config.GetCfg().CertFile != "" && config.GetCfg().KeyFile != "",
		autoTLS:   config.GetCfg().AutoTLSCertificate,
	}
}

func ShowQRCode() {
	link := buildReaderLinkContext()
	// 打印二维码
	tools.PrintAllReaderURL(
		config.GetCfg().Port,
		config.GetCfg().PrintAllPossibleQRCode,
		link.outIP,
		config.GetCfg().DisableLAN,
		link.customTLS,
		link.autoTLS,
		link.etcStr,
	)
	// ”如何快捷键退出“的文字提示
	logger.Info(locale.GetString("ctrl_c_hint"))
	// 打印 Tailscale 访问地址的二维码
	ShowQRCodeTailscale(context.Background())
}

// OpenReaderBrowserIfNeeded 只负责按配置打开本机阅读页，避免 ShowQRCode 隐含浏览器副作用。
func OpenReaderBrowserIfNeeded() {
	if !config.GetCfg().OpenBrowser {
		return
	}
	link := buildReaderLinkContext()
	tools.OpenLocalReaderURL(config.GetCfg().Port, link.customTLS, link.autoTLS, link.etcStr)
}
