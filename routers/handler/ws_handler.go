package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSMsg struct {
	UserUUID     string  `json:"user_id"`
	ServerStatus string  `json:"server_status"`
	NowBookUUID  string  `json:"now_book_id"`
	ReadPercent  float64 `json:"read_percent"`
	Message      string  `json:"msg"`
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//webSocket请求ping 返回pong
func WsHandler(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("err = %s\n", err)
		return
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("read fail = %v\n", err)
			break
		}
		msg := &WSMsg{}
		json.Unmarshal(message, msg)
		fmt.Println(msg)
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
		fmt.Println(msg)
		err = ws.WriteJSON(msg)
		if err != nil {
			fmt.Printf("write fail = %v\n", err)
			break
		}
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
