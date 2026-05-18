package sse_hub

import (
	"strings"
	"sync"
)

// Go(Echo) 后端实现 SSE 推送，用于把后端日志和 UI 通知推送给前端。

// --- 简单的广播中心 ---

// Event SSE消息
type Event struct {
	Name string // 可选：自定义事件名（不写则走默认 "message" 事件）
	ID   string // 可选：事件 ID（配合 Last-Event-ID 可断点续传）
	Data string // 必填：消息内容
}

// Hub 广播中心
type Hub struct {
	mu      sync.RWMutex
	clients map[string]chan Event
}

// NewHub 创建一个新的 Hub 实例
func NewHub() *Hub {
	return &Hub{clients: make(map[string]chan Event)}
}

var MessageHub *Hub

func init() {
	// 创建 SSE 广播中心
	MessageHub = NewHub()
}

// Add 注册一个新的客户端
func (h *Hub) Add(id string, ch chan Event) {
	h.mu.Lock()
	h.clients[id] = ch
	h.mu.Unlock()
}

// Remove 注销一个客户端
func (h *Hub) Remove(id string) {
	h.mu.Lock()
	if ch, ok := h.clients[id]; ok {
		close(ch)
		delete(h.clients, id)
	}
	h.mu.Unlock()
}

// CloseAll 关闭所有客户端连接并清空列表（用于优雅关机）
func (h *Hub) CloseAll() {
	h.mu.Lock()
	for id, ch := range h.clients {
		close(ch)
		delete(h.clients, id)
	}
	h.mu.Unlock()
}

// Broadcast 向所有注册的客户端广播事件
func (h *Hub) Broadcast(ev Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, ch := range h.clients {
		select {
		case ch <- ev:
		default:
			// 丢弃过慢的客户端，避免阻塞
		}
	}
}

// --- 辅助函数 ---

// splitLines 按行拆分字符串，兼容各种换行符
func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}
