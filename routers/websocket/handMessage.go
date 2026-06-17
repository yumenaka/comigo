package websocket

import (
	"encoding/json"
	"math"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

var WsDebug *bool

// 处理乒乓信息
func handPingMessage(msg Message) Message {
	msg.Detail = "pong!"
	return msg
}

// 处理心跳信息
func handHeartbeatMessage(msg Message) Message {
	msg.Detail = " 服务器收到心跳消息。" + time.Now().Format("2006.01.02 15:04:05")
	return msg
}

// handSyncPageMessageToFlipMode 处理翻页阅读消息
func handSyncPageMessageToFlipMode(msg Message, clientID string) (Message, bool) {
	msg.Detail = "同步页数。"
	// Message 定义一个对象来管理消息，反引号包含的文本是 Go 在对象和 JSON 之间进行序列化和反序列化时需要的元数据。
	type syncData struct {
		BookID             string `json:"book_id"`
		NowPageNum         int    `json:"now_page_num"`
		NeedDoublePageMode bool   `json:"need_double_page_mode"` // 需要切换为双页模式
	}
	var data syncData
	if err := json.Unmarshal([]byte(msg.DataString), &data); err != nil {
		logger.Infof("handSyncPageMessageToFlipMode error: %v", err)
		return msg, false
	}
	// log_syncpage_message_to_flipmode 翻页阅读同步页数消息 data: %v, clientID: %v
	if isWsDebug() {
		logger.Infof(locale.GetString("log_syncpage_message_to_flipmode"), data, clientID)
	}
	// 验证收到的数据 SyncPage消息发送到FlipMode
	if data.BookID == "" || data.NowPageNum < 0 || data.NowPageNum > math.MaxInt {
		logger.Infof("handSyncPage_ToFlipode data error: %v", data)
		return msg, false
	}
	return msg, true
}

// handSyncPageMessageToScrollMode 处理翻页信息(下拉阅读模式)
func handSyncPageMessageToScrollMode(msg Message, clientID string) (Message, bool) {
	msg.Detail = "同步页数。"
	type syncData struct {
		BookID            string  `json:"book_id"`
		NowPageNum        int     `json:"now_page_num"`
		NowPageNumPercent float64 `json:"now_page_num_percent"`
		StartLoadPageNum  int     `json:"start_load_page_num"`
		EndLoadPageNum    int     `json:"end_load_page_num"`
	}
	var data syncData
	if err := json.Unmarshal([]byte(msg.DataString), &data); err != nil {
		logger.Infof("handSyncPage_ToFlipode error: %v", err)
		return msg, false
	}
	if isWsDebug() {
		logger.Infof(locale.GetString("log_syncpage_message_to_scrollmode"), data, clientID)
	}
	if data.BookID == "" || data.NowPageNum < 0 || data.NowPageNum > math.MaxInt || data.NowPageNumPercent > 1 {
		logger.Infof("handSyncPage_ToFlipode data error: %v", data)
		return msg, false
	}
	return msg, true
}

func isWsDebug() bool {
	return WsDebug != nil && *WsDebug
}
