//go:build !wails && !js

package main

import (
	"testing"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// 独立托盘恢复后仍通过配置事务切换语言，并遵守只读配置。
func TestTrayLanguage(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { config.ReplaceConfig(old); locale.InitLanguageFromConfig(old.Language) })
	next := config.CopyCfg()
	next.TemporaryReaderMode, next.ReadOnlyMode = true, false
	config.ReplaceConfig(next)
	if err := setLanguage("ja-JP"); err != nil || config.GetCfg().Language != "ja-JP" {
		t.Fatalf("语言未更新: %v", err)
	}
	next = config.CopyCfg()
	next.ReadOnlyMode = true
	config.ReplaceConfig(next)
	if err := setLanguage("en-US"); err == nil || config.GetCfg().Language != "ja-JP" {
		t.Fatal("只读配置不应允许切换语言")
	}
}
