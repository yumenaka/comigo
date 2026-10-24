package config

import (
	"path/filepath"
	"sync"
	"testing"
)

// 配置写入与状态轮询并发时，旧快照及其切片仍保持不变。
func TestConfigSnapshotsConcurrentUpdates(t *testing.T) {
	old := CopyCfg()
	t.Cleanup(func() { ReplaceConfig(old) })
	next := CopyCfg()
	next.ExcludePath = []string{"keep"}
	ReplaceConfig(next)
	snapshot := GetCfg()
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if err := UpdateConfigByJson(`{"Debug":true,"ExcludePath":["new"]}`); err != nil {
					t.Error(err)
				}
				copy := CopyCfg()
				copy.ExcludePath[0] = "private"
				_ = GetCfg().Debug
			}
		}()
	}
	wg.Wait()
	if snapshot.Debug != next.Debug || snapshot.ExcludePath[0] != "keep" || GetCfg().ExcludePath[0] != "new" {
		t.Fatal("配置快照被其他更新修改")
	}
}

// 保存失败不能发布候选配置，避免界面成功而重启后丢失设置。
func TestChangeConfigDoesNotPublishOnSaveFailure(t *testing.T) {
	old := CopyCfg()
	t.Cleanup(func() { ReplaceConfig(old) })
	next := CopyCfg()
	next.TemporaryReaderMode = false
	next.ConfigFile = t.TempDir()
	ReplaceConfig(next)
	snapshot := GetCfg()
	if _, err := ChangeConfig(true, func(c *Config) error { c.Debug = !c.Debug; return nil }); err == nil {
		t.Fatal("期望目录目标写入失败")
	}
	if GetCfg() != snapshot {
		t.Fatal("保存失败发布了配置")
	}
	_, err := ChangeConfig(true, func(c *Config) error {
		c.ConfigFile = filepath.Join(t.TempDir(), "config.toml")
		c.Debug = !c.Debug
		return nil
	})
	if err != nil || GetCfg().Debug == snapshot.Debug {
		t.Fatalf("合法保存失败: %v", err)
	}
}
