package cmd

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yumenaka/comigo/config"
)

// TestTailscaleQRPollCancellation 已取消的轮询立即结束，多次并发启动/停止不共享计数状态。
func TestTailscaleQRPollCancellation(t *testing.T) {
	old := config.CopyCfg()
	defer config.ReplaceConfig(old)
	cfg := old
	cfg.EnableTailscale = true
	config.ReplaceConfig(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	go func() { pollTailscaleQRCode(ctx); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("取消没有终止轮询")
	}
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() { defer wg.Done(); ShowQRCodeTailscale(ctx); stopTailscaleQRCode() }()
	}
	wg.Wait()
	stopTailscaleQRCode()
}
