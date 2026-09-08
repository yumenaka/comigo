// Package traffic 统计本进程提供 HTTP 服务时实际读写的正文，不统计协议开销。
package traffic

import (
	"sync"
	"time"
)

type Snapshot struct {
	StartedAt             time.Time `json:"startedAt"`
	ReceivedBytes         uint64    `json:"receivedBytes"`
	SentBytes             uint64    `json:"sentBytes"`
	ReceiveBytesPerSecond float64   `json:"receiveBytesPerSecond"`
	SendBytesPerSecond    float64   `json:"sendBytesPerSecond"`
}
type bucket struct {
	second         int64
	received, sent uint64
}

// Counter 用固定大小的时间桶提供滚动速度，避免为每个连接启动采样协程。
type Counter struct {
	mu             sync.Mutex
	started        time.Time
	received, sent uint64
	buckets        [5]bucket
}

func New() *Counter { return &Counter{started: time.Now()} }

var Default = New()

func (c *Counter) Add(received, sent int) { c.addAt(received, sent, time.Now()) }
func (c *Counter) addAt(received, sent int, now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	second := int64(now.Sub(c.started) / time.Second)
	if second < 0 {
		return
	}
	b := &c.buckets[second%5]
	if b.second != second {
		*b = bucket{second: second}
	}
	if received > 0 {
		c.received += uint64(received)
		b.received += uint64(received)
	}
	if sent > 0 {
		c.sent += uint64(sent)
		b.sent += uint64(sent)
	}
}
func (c *Counter) Snapshot() Snapshot { return c.snapshotAt(time.Now()) }
func (c *Counter) snapshotAt(now time.Time) Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	elapsed := now.Sub(c.started).Seconds()
	second := int64(elapsed)
	var received, sent uint64
	for _, b := range c.buckets {
		if b.second <= second && b.second > second-5 {
			received += b.received
			sent += b.sent
		}
	}
	window := min(5.0, max(1.0, elapsed))
	return Snapshot{c.started, c.received, c.sent, float64(received) / window, float64(sent) / window}
}
