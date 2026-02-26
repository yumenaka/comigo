package service

import (
	"testing"

	"github.com/yumenaka/comigo/config"
)

func TestBuildConfigChangeAction(t *testing.T) {
	oldCfg := &config.Config{
		Port:                      1234,
		StoreUrls:                 []string{"/a"},
		EnableTailscale:           false,
		AutoRescanIntervalMinutes: 0,
	}
	newCfg := &config.Config{
		Port:                      5678,
		StoreUrls:                 []string{"/a", "/b"},
		EnableTailscale:           true,
		AutoRescanIntervalMinutes: 10,
	}

	action := BuildConfigChangeAction(oldCfg, newCfg)
	if !action.ReScanStores {
		t.Fatalf("expected ReScanStores=true")
	}
	if !action.ReStartWebServer {
		t.Fatalf("expected ReStartWebServer=true")
	}
	if !action.StartTailscale {
		t.Fatalf("expected StartTailscale=true")
	}
	if !action.UpdateAutoRescan {
		t.Fatalf("expected UpdateAutoRescan=true")
	}
}
