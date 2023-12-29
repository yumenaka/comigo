package cmd

import (
	"github.com/spf13/viper"
	"github.com/yumenaka/comi/arch/scan"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/routers/handlers"
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
			logger.Info("扫描上传文件夹：", msg)
			ReScanUploadPath()
			//保存扫描结果到数据库
			if config.Config.EnableDatabase {
				err := scan.SaveResultsToDatabase(viper.ConfigFileUsed(), config.Config.ClearDatabaseWhenExit)
				if err != nil {
					return
				}
			}
		case "SomePath":
			logger.Info("收到重新扫描消息：", msg)
			ReScanPath(msg, false)
		default:
			continue
		}
	}
}

// ReScanUploadPath 重新扫描上传目录,因为需要设置下载路径，gin 初始化后才能执行
func ReScanUploadPath() {
	//没启用上传，则不扫描
	if !config.Config.EnableUpload {
		return
	}
	uploadPath := "upload"
	if config.Config.UploadPath != "" {
		uploadPath = config.Config.UploadPath
	}
	ReScanPath(uploadPath, false)
}

func ReScanPath(path string, reScanFile bool) {
	//扫描上传目录的文件
	option := scan.NewScanOption(
		reScanFile,
		config.Config.StoresPath,
		config.Config.MaxScanDepth,
		config.Config.MinImageNum,
		config.Config.TimeoutLimitForScan,
		config.Config.ExcludePath,
		config.Config.SupportMediaType,
		config.Config.SupportFileType,
		config.Config.ZipFileTextEncoding,
		config.Config.EnableDatabase,
		config.Config.ClearDatabaseWhenExit,
		config.Config.Debug,
	)
	addList, err := scan.ScanAndGetBookList(path, option)
	if err != nil {
		logger.Info(locale.GetString("scan_error"), path, err)
	} else {
		scan.AddBooksToStore(addList, path, config.Config.MinImageNum)
	}
}
