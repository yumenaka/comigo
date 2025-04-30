package websocket

import (
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yumenaka/comigo/util/logger"
)

var WsDebug *bool

// 默认的处理方式：原样返回
func handDefaultMessage(client *websocket.Conn, msg Message, clientID string) {
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		// 如果写入 Websocket 时出现错误，关闭连接，并将其从“clients” 映射中删除。
		err := client.Close()
		if err != nil {
			return
		}
		delete(clients, client)
	}
}

// 处理乒乓信息
func handPingMessage(client *websocket.Conn, msg Message, clientID string) {
	msg.Detail = "pong!"
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		// 如果写入 Websocket 时出现错误，关闭连接，并将其从“clients” 映射中删除。
		err := client.Close()
		if err != nil {
			return
		}
		delete(clients, client)
	}
}

// 处理心跳信息
func handHeartbeatMessage(client *websocket.Conn, msg Message, clientID string) {
	msg.Detail = " 服务器收到心跳消息。" + time.Now().Format("2006.01.02 15:04:05")
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		// 如果写入 Websocket 时出现错误，关闭连接，并将其从“clients” 映射中删除。
		err := client.Close()
		if err != nil {
			return
		}
		delete(clients, client)
	}
}

// handSyncPageMessageToFlipMode 处理翻页消息(翻页模式)
func handSyncPageMessageToFlipMode(client *websocket.Conn, msg Message, clientID string) {
	msg.Detail = "同步页数。"
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("handSyncPageMessageToFlipMode error: %v", err)
		// 如果写入 Websocket 时出现错误，关闭连接，并将其从“clients” 映射中删除。
		err := client.Close()
		if err != nil {
			return
		}
		delete(clients, client)
	}
	// Message 定义一个对象来管理消息，反引号包含的文本是 Go 在对象和 JSON 之间进行序列化和反序列化时需要的元数据。
	type syncData struct {
		BookID             string `json:"book_id"`
		NowPageNum         int    `json:"now_page_num"`
		NeedDoublePageMode bool   `json:"need_double_page_mode"` // 需要切换为双页模式
	}
	var data syncData
	if err := json.Unmarshal([]byte(msg.DataString), &data); err != nil {
		log.Printf("handSyncPageMessageToFlipMode error: %v", err)
		return
	}
	if *WsDebug {
		logger.Infof(" SyncPage message toFlipMode: %s %s", data, clientID)
	}
	// 验证收到的数据
	if data.BookID == "" || data.NowPageNum < 0 || data.NowPageNum > math.MaxInt {
		log.Printf("handSyncPage_ToFlipode data error: %v", data)
		return
	}
	// 向所有其他在线客户端发送翻页信息
	for c := range clients {
		// if id == clientID {
		//	log.Printf("跳过一个客户端 clientID: %v", id)
		//	continue
		// }
		err := c.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			err := c.Close()
			if err != nil {
				return
			}
			delete(clients, c)
		}
	}
}

// handSyncPageMessageToScrollMode 处理翻页信息(下拉阅读模式)
func handSyncPageMessageToScrollMode(client *websocket.Conn, msg Message, clientID string) {
	msg.Detail = "同步页数。"
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("handSyncPage_ToFlipode error: %v", err)
		err := client.Close()
		if err != nil {
			return
		}
		delete(clients, client)
	}
	type syncData struct {
		BookID            string  `json:"book_id"`
		NowPageNum        int     `json:"now_page_num"`
		NowPageNumPercent float64 `json:"now_page_num_percent"`
		StartLoadPageNum  int     `json:"start_load_page_num"`
		EndLoadPageNum    int     `json:"end_load_page_num"`
	}
	var data syncData
	if err := json.Unmarshal([]byte(msg.DataString), &data); err != nil {
		log.Printf("handSyncPage_ToFlipode error: %v", err)
		return
	}
	if *WsDebug {
		logger.Infof(" SyncPage message to ScrollMode:%s %s", data, clientID)
	}
	if data.BookID == "" || data.NowPageNum < 0 || data.NowPageNum > math.MaxInt || data.NowPageNumPercent > 1 {
		log.Printf("handSyncPage_ToFlipode data error: %v", data)
		return
	}
	// 向所有其他在线客户端发送翻页信息
	for c := range clients {
		// if id == clientID {
		//	log.Printf("跳过一个客户端 clientID: %v", id)
		//	continue
		// }
		err := c.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			err := c.Close()
			if err != nil {
				return
			}
			delete(clients, c)
		}
	}
}
