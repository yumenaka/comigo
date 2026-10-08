package scheduler

import (
	"testing"
	"time"
)

// 正在扫描时停止调度器，应能等到任务退出，不因运行状态锁而超时。
func TestStopDuringScan(t *testing.T) {
	scanner := NewLibraryScanner()
	started := make(chan struct{})
	finish := make(chan struct{})
	if err := scanner.Start(1, func() error {
		close(started)
		<-finish
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	defer scanner.Stop()
	if err := scanner.scheduler.Jobs()[0].RunNow(); err != nil {
		close(finish)
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		close(finish)
		t.Fatal("扫描任务未启动")
	}
	if !scanner.IsRunning() {
		close(finish)
		t.Fatal("扫描状态未更新")
	}
	time.AfterFunc(20*time.Millisecond, func() { close(finish) })
	if err := scanner.Stop(); err != nil {
		t.Fatal(err)
	}
	if scanner.IsRunning() || scanner.GetInterval() != 0 {
		t.Fatal("停止后仍保留扫描状态")
	}
}
