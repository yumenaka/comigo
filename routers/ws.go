package routers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yumenaka/comi/common"
	"net/http"
)

type WSMsg struct {
	UserUUID     string  `json:"user_uuid"`
	ServerStatus string  `json:"server_status"`
	NowBookUUID  string  `json:"now_book_uuid"`
	ReadPercent  float64 `json:"read_percent"`
	Message      string  `json:"msg"`
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//webSocket请求ping 返回pong
func wsHandler(c *gin.Context) {
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
		if msg.Message == "ChangeBook" && msg.NowBookUUID != common.ReadingBook.FileID {
			if changeReadingBook(msg.NowBookUUID) {
				fmt.Println("正在切换书籍：", common.ReadingBook.FilePath)
				if err != nil {
					fmt.Println("无法初始化书籍。", err, common.ReadingBook)
				} else {
					msg.Message = "Extracting"
				}
			} else {
				msg.Message = "BookNotFund"
				fmt.Println("没找到这本书：", msg)
			}
		}
		//返回漫画压缩包解压状态
		if msg.Message == "CheckExtract" && msg.NowBookUUID == common.ReadingBook.FileID {
			if common.ReadingBook.ExtractComplete {
				msg.Message = "ExtractComplete"
			} else {
				msg.Message = "Extracting"
			}
		}
		//
		if msg.Message == "MasterDevicesSync" && msg.NowBookUUID == common.ReadingBook.FileID {
			common.ReadingBook.ReadPercent = msg.ReadPercent
		}
		//
		if msg.Message == "SlaveDevicesSync" && msg.NowBookUUID == common.ReadingBook.FileID {
			common.ReadingBook.ReadPercent = msg.ReadPercent
		}
		fmt.Println(msg)
		err = ws.WriteJSON(msg)
		if err != nil {
			fmt.Printf("write fail = %v\n", err)
			break
		}
	}
}

func changeReadingBook(u string) bool {
	for i := 0; i < len(common.BookList); i++ {
		if common.BookList[i].FileID == u {
			common.ReadingBook = common.BookList[i]
			//初始化书籍
			err := common.InitReadingBook()
			if err != nil {
				fmt.Println("无法初始化书籍。", err, common.ReadingBook)
				return false
			}
			return true
		}
	}
	return false
}
