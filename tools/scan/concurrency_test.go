package scan

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
)

type concurrentScanConfig struct {
	*testCfgScan
	active atomic.Int32
	max    atomic.Int32
}

func (c *concurrentScanConfig) GetMaxScanDepth() int {
	active := c.active.Add(1)
	defer c.active.Add(-1)
	for current := c.max.Load(); active > current && !c.max.CompareAndSwap(current, active); current = c.max.Load() {
	}
	time.Sleep(10 * time.Millisecond)
	return c.testCfgScan.GetMaxScanDepth()
}

// TestInitStoreSerializesScans 确认多个入口不会并发改写扫描上下文。
func TestInitStoreSerializesScans(t *testing.T) {
	oldStore := model.IStore
	model.IStore = &store.StoreInRam{}
	defer func() { model.IStore = oldStore }()

	cfg := &concurrentScanConfig{testCfgScan: &testCfgScan{maxScanDepth: -1}}
	dirs := []string{t.TempDir(), t.TempDir()}
	start := make(chan struct{})
	errs := make(chan error, len(dirs))
	var wg sync.WaitGroup
	for _, dir := range dirs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs <- InitStore(dir, cfg)
		}()
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if cfg.max.Load() != 1 {
		t.Fatalf("concurrent scans = %d, want 1", cfg.max.Load())
	}
}
