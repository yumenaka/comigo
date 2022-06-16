package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//收发消息的时候，用来区分消息类型
type MessageType int

const (
	TestOnlineMessage MessageType = iota // 开始生成枚举值, 默认为0
	SyncPageNum                          //翻页用
	ErrorHint                            //书籍不存在、文件已移除的提示
	Rifle
	Blower
)

//json格式的消息体，用来传递比较复杂的信息
type WSMsg struct {
	MessageType  int     `json:"message_type"`
	UserUUID     string  `json:"user_id"`
	ServerStatus string  `json:"server_status"`
	NowBookUUID  string  `json:"now_book_id"`
	ReadPercent  float64 `json:"read_percent"`
	Message      string  `json:"msg"`
}

//Upgrader 用于升级 http 请求，把 http 请求升级为长连接的 WebSocket。
//Gorilla的工作是转换原始HTTP连接进入一个有状态的websocket连接。
//使用全局升级程序变量通过来帮助我们将任何传入的HTTP连接转换为websocket协议
var upGrader = websocket.Upgrader{
	//ReadBufferSize:  4096,//读缓存区大小 单位是 bytes，依需求設定（设为 0，则不限制大小）
	//WriteBufferSize: 1024,// 写缓存区大小 同上
	// use default options
	// 检测请求来源 //检查是否跨域
	CheckOrigin: func(r *http.Request) bool {
		////验证方法，这是只支持Get
		//if r.Method != "GET" {
		//	fmt.Println("method is not GET")
		//	return false
		//}
		//验证路径
		if r.URL.Path != "/api/ws" {
			fmt.Println("path error")
			return false
		}
		return true
	},
}

//webSocket请求ping 返回pong
//  ws://127.0.0.1:1234/api/ws
func WebSocketHandler(c *gin.Context) {
	//Upgrade 函数将 http get请求升级到 WebSocket 协议。
	//   responseHeader包含在对客户端升级请求的响应中。
	//// 使用responseHeader指定cookie（Set-Cookie）和应用程序协商的子协议（Sec-WebSocket-Protocol）。
	//// 如果升级失败，则升级将使用HTTP错误响应回复客户端
	//// 返回一个 Conn 指针(wsConn)，拿到他后，可使用 Conn 读写数据与客户端通信。
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("Error during connection up gradation:", err) //连接升级出错
		return
	}
	//比 defer wsConn.Close() 更好的写法
	defer func() {
		closeSocketErr := wsConn.Close()
		if closeSocketErr != nil {
			fmt.Println(err)
		}
	}()

	// The event loop
	for {
		////读取ws中的数据
		messageType, message, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)
		fmt.Printf("Message Type: %d, Message: %s\n", messageType, string(message))

		//测试用的乒乓逻辑
		if string(message) == "ping!" || string(message) == "ping" || string(message) == "乒" {
			message = []byte("pang!")
		}
		//写入ws数据
		err = wsConn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}

		////读取ws中的数据,转换为json

		//msg := &WSMsg{}
		//json.Unmarshal(message, msg)
		//fmt.Println(msg)

		//客户端请求更换书籍
		//if msg.Message == "ChangeBook" && msg.NowBookUUID != common.ReadingBook.GetBookID() {
		//	if changeReadingBook(msg.NowBookUUID) {
		//		fmt.Println("正在切换书籍：", common.ReadingBook.GetFilePath())
		//		if err != nil {
		//			fmt.Println("无法初始化书籍。", err, common.ReadingBook)
		//		} else {
		//			msg.Message = "Extracting"
		//		}
		//	} else {
		//		msg.Message = "BookNotFund"
		//		fmt.Println("没找到这本书：", msg)
		//	}
		//}
		////返回漫画压缩包解压状态
		//if msg.Message == "CheckExtract" && msg.NowBookUUID == common.ReadingBook.GetBookID() {
		//	if common.ReadingBook.InitComplete {
		//		msg.Message = "InitComplete"
		//	} else {
		//		msg.Message = "Extracting"
		//	}
		//}
		//
		//if msg.Message == "MasterDevicesSync" && msg.NowBookUUID == common.ReadingBook.GetBookID() {
		//	common.ReadingBook.ReadPercent = msg.ReadPercent
		//}
		////
		//if msg.Message == "SlaveDevicesSync" && msg.NowBookUUID == common.ReadingBook.GetBookID() {
		//	common.ReadingBook.ReadPercent = msg.ReadPercent
		//}
		//fmt.Println(msg)
		//err = wsConn.WriteJSON(msg)
		//if err != nil {
		//	fmt.Printf("write fail = %v\n", err)
		//	break
		//}
	}
}

func changeReadingBook(u string) bool {
	//b, err := common.GetBookByID(u, false)
	//if err != nil {
	//	return false
	//}
	//common.ReadingBook = b
	return false
}
