package cmd

import (
	"github.com/spf13/viper"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/handlers"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
	"github.com/yumenaka/comigo/util/scan"
)

// 用于由客户端发送消息的队列，扮演通道的角色。后面定义了一个 goroutine 来从这个通道读取新消息，然后将它们发送给其它连接到服务器的客户端。
var rescanBroadcast = make(chan string) // broadcast channel
func init() {
	// Start listening for incoming chat messages
	go waitRescanMessages()
	handlers.LocalRescanBroadcast = &rescanBroadcast
}

// 一个简单循环，从“broadcast”中连续读取数据，然后通过各自的 WebSocket 连接将消息传播到客户端。
func waitRescanMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-rescanBroadcast //广播频道
		// Send it out to every client that is currently connected
		switch msg {
		case "upload":
			logger.Infof("扫描上传文件夹：%s", msg)
			ReScanUploadPath()
			//保存扫描结果到数据库
			if config.GetEnableDatabase() {
				err := scan.SaveResultsToDatabase(viper.ConfigFileUsed(), config.GetClearDatabaseWhenExit())
				if err != nil {
					return
				}
			}
		case "AnotherPath":
			logger.Infof("收到重新扫描消息：%s", msg)
			ReScanPath(msg, false)
		default:
			continue
		}
	}
}

// ReScanUploadPath 重新扫描上传目录,因为需要设置下载路径，gin 初始化后才能执行
func ReScanUploadPath() {
	//没启用上传，则不扫描
	if !config.GetEnableUpload() {
		return
	}
	ReScanPath(config.GetUploadPath(), true)
}

func ReScanPath(path string, reScanFile bool) {
	//扫描上传目录的文件
	option := scan.NewOption(
		reScanFile,
		config.GetCfg(),
	)
	books, err := scan.InitStore(path, option)
	if err != nil {
		logger.Infof(locale.GetString("scan_error")+"path:%s  %s", path, err)
		return
	}
	scan.AddBooksToStore(books, path, config.GetMinImageNum())
	model.ResetBookGroupData()
}
