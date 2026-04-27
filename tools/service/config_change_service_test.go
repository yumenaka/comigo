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

	action := BuildConfigChangeAction(*oldCfg, newCfg)
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

func TestBuildConfigChangeActionRestartsWebServerWhenBasePathChanges(t *testing.T) {
	oldCfg := &config.Config{
		BasePath: "",
		Port:     1234,
	}
	newCfg := &config.Config{
		BasePath: "/proxy",
		Port:     1234,
	}

	action := BuildConfigChangeAction(*oldCfg, newCfg)
	if !action.ReStartWebServer {
		t.Fatalf("expected ReStartWebServer=true")
	}
}

func TestBuildConfigChangeActionIgnoresEquivalentBasePath(t *testing.T) {
	oldCfg := &config.Config{
		BasePath: "/proxy/",
		Port:     1234,
	}
	newCfg := &config.Config{
		BasePath: "/proxy",
		Port:     1234,
	}

	action := BuildConfigChangeAction(*oldCfg, newCfg)
	if action.ReStartWebServer {
		t.Fatalf("expected ReStartWebServer=false for equivalent BasePath")
	}
}

func TestBuildConfigChangeActionStartTailscale(t *testing.T) {
	oldCfg := &config.Config{
		EnableTailscale:   false,
		TailscaleHostname: "comigo",
		TailscalePort:     443,
	}
	newCfg := &config.Config{
		EnableTailscale:   true,
		TailscaleHostname: "comigo",
		TailscalePort:     443,
	}

	action := BuildConfigChangeAction(*oldCfg, newCfg)
	if !action.StartTailscale {
		t.Fatalf("expected StartTailscale=true")
	}
	if action.StopTailscale {
		t.Fatalf("expected StopTailscale=false")
	}
	if action.ReStartTailscale {
		t.Fatalf("expected ReStartTailscale=false")
	}
}

func TestBuildConfigChangeActionStopTailscale(t *testing.T) {
	oldCfg := &config.Config{
		EnableTailscale:   true,
		TailscaleHostname: "comigo",
		TailscalePort:     443,
	}
	newCfg := &config.Config{
		EnableTailscale:   false,
		TailscaleHostname: "comigo",
		TailscalePort:     443,
	}

	action := BuildConfigChangeAction(*oldCfg, newCfg)
	if !action.StopTailscale {
		t.Fatalf("expected StopTailscale=true")
	}
	if action.StartTailscale {
		t.Fatalf("expected StartTailscale=false")
	}
	if action.ReStartTailscale {
		t.Fatalf("expected ReStartTailscale=false")
	}
}

func TestBuildConfigChangeActionRestartTailscaleWhenConfigChanges(t *testing.T) {
	oldCfg := &config.Config{
		EnableTailscale:   true,
		TailscaleAuthKey:  "tskey-old",
		TailscaleHostname: "comigo",
		TailscalePort:     443,
		FunnelTunnel:      false,
	}

	testCases := []struct {
		name   string
		newCfg *config.Config
	}{
		{
			name: "auth key changed",
			newCfg: &config.Config{
				EnableTailscale:   true,
				TailscaleAuthKey:  "tskey-new",
				TailscaleHostname: "comigo",
				TailscalePort:     443,
				FunnelTunnel:      false,
			},
		},
		{
			name: "hostname changed",
			newCfg: &config.Config{
				EnableTailscale:   true,
				TailscaleAuthKey:  "tskey-old",
				TailscaleHostname: "reader",
				TailscalePort:     443,
				FunnelTunnel:      false,
			},
		},
		{
			name: "port changed",
			newCfg: &config.Config{
				EnableTailscale:   true,
				TailscaleAuthKey:  "tskey-old",
				TailscaleHostname: "comigo",
				TailscalePort:     8443,
				FunnelTunnel:      false,
			},
		},
		{
			name: "funnel changed",
			newCfg: &config.Config{
				EnableTailscale:   true,
				TailscaleAuthKey:  "tskey-old",
				TailscaleHostname: "comigo",
				TailscalePort:     443,
				FunnelTunnel:      true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			action := BuildConfigChangeAction(*oldCfg, tc.newCfg)
			if !action.ReStartTailscale {
				t.Fatalf("expected ReStartTailscale=true")
			}
			if action.StartTailscale {
				t.Fatalf("expected StartTailscale=false")
			}
			if action.StopTailscale {
				t.Fatalf("expected StopTailscale=false")
			}
		})
	}
}
