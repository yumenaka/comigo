package cmd

import (
	"fmt"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/routers/handler"
)

//用于由客户端发送消息的队列，扮演通道的角色。后面定义了一个 goroutine 来从这个通道读取新消息，然后将它们发送给其它连接到服务器的客户端。
var rescanBroadcast = make(chan string) // broadcast channel
func init() {
	// Start listening for incoming chat messages
	go waitRescanMessages()
	handler.LocalRescanBroadcast = &rescanBroadcast
}

//一个简单循环，从“broadcast”中连续读取数据，然后通过各自的 WebSocket 连接将消息传播到客户端。
func waitRescanMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-rescanBroadcast //广播频道
		// Send it out to every client that is currently connected
		switch msg {
		case "ComigoUpload":
			fmt.Println("收到重新扫描消息：", msg)
			ReScanUploadPath(msg)
		default:
			continue
		}
	}
}

//ReScanUploadPath 6、扫描上传目录的文件
func ReScanUploadPath(p string) {
	addList, err := common.ScanAndGetBookList(p, databaseBookList)
	if err != nil {
		fmt.Println(locale.GetString("scan_error"), p)
	} else {
		common.AddBooksToStore(addList, p)
	}
	//4，保存扫描结果到数据库
	SaveResultsToDatabase()
}
