package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

//默认的处理方式：原样返回
func handDefaultMessage(client *websocket.Conn, msg Message) {
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		//如果写入 Websocket 时出现错误，关闭连接，并将其从“clients” 映射中删除。
		client.Close()
		delete(clients, client)
	}
}

//处理乒乓信息
func handPingMessage(client *websocket.Conn, msg Message) {
	msg.Msg = "pong!"
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		//如果写入 Websocket 时出现错误，关闭连接，并将其从“clients” 映射中删除。
		client.Close()
		delete(clients, client)
	}
}

//处理心跳信息
func handHeartbeatMessage(client *websocket.Conn, msg Message) {
	msg.Msg = " 服务器收到心跳消息。" + time.Now().Format("2006.01.02 15:04:05")
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		//如果写入 Websocket 时出现错误，关闭连接，并将其从“clients” 映射中删除。
		client.Close()
		delete(clients, client)
	}
}

//处理翻页消息
func handSyncPageMessage(client *websocket.Conn, msg Message) {
	msg.Msg = "同步页数。"
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("handSyncPageMessage error: %v", err)
		//如果写入 Websocket 时出现错误，关闭连接，并将其从“clients” 映射中删除。
		client.Close()
		delete(clients, client)
	}
	// Message 定义一个对象来管理消息，反引号包含的文本是 Go 在对象和 JSON 之间进行序列化和反序列化时需要的元数据。
	type syncData struct {
		BookID            string  `json:"book_id"`
		NowPageNum        int     `json:"now_page_num"`
		NowPageNumPercent float64 `json:"now_page_num_percent"`
		ReadPercent       float64 `json:"read_percent"`
	}
	var data syncData
	if err := json.Unmarshal([]byte(msg.data), &data); err != nil {
		log.Printf("handSyncPageMessage syncData error: %v", err)
		return
	}
	fmt.Println(data.BookID)

}
