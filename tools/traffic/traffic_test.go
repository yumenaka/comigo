package traffic

import (
	"sync"
	"testing"
	"time"
)

// 长传输在请求结束前可见，五秒无读写后速度归零而累计值保留。
func TestRollingCounter(t *testing.T) {
	c := New()
	start := c.started
	c.addAt(100, 200, start)
	s := c.snapshotAt(start.Add(time.Second))
	if s.ReceivedBytes != 100 || s.SentBytes != 200 || s.SendBytesPerSecond != 200 {
		t.Fatalf("snapshot: %+v", s)
	}
	c.addAt(0, 300, start.Add(3*time.Second))
	s = c.snapshotAt(start.Add(4 * time.Second))
	if s.SendBytesPerSecond != 125 {
		t.Fatalf("rate: %+v", s)
	}
	s = c.snapshotAt(start.Add(9 * time.Second))
	if s.SendBytesPerSecond != 0 || s.ReceiveBytesPerSecond != 0 || s.SentBytes != 500 {
		t.Fatalf("idle: %+v", s)
	}
	if New().Snapshot().SentBytes != 0 {
		t.Fatal("new process counter should start empty")
	}
}

// 并行上传与下载必须无丢计，读快照不得产生数据竞争。
func TestConcurrentCounter(t *testing.T) {
	c := New()
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				c.Add(1, 2)
				c.Snapshot()
			}
		}()
	}
	wg.Wait()
	s := c.Snapshot()
	if s.ReceivedBytes != 50000 || s.SentBytes != 100000 {
		t.Fatalf("totals: %+v", s)
	}
}
