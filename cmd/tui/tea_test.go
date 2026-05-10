package tui

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/yumenaka/comigo/config"
)

func TestBuildBaseURLUsesConfiguredHost(t *testing.T) {
	restoreConfig(t)
	cfg := config.GetCfg()
	cfg.Host = "example.com"
	cfg.Port = 1234
	cfg.CertFile = ""
	cfg.KeyFile = ""
	cfg.AutoTLSCertificate = false

	if got, want := buildBaseURL(), "http://example.com:1234"; got != want {
		t.Fatalf("buildBaseURL() = %q, want %q", got, want)
	}
}

func TestBuildBaseURLUsesLocalhostWhenLANDisabled(t *testing.T) {
	restoreConfig(t)
	cfg := config.GetCfg()
	cfg.Host = ""
	cfg.Port = 1234
	cfg.DisableLAN = true
	cfg.CertFile = ""
	cfg.KeyFile = ""
	cfg.AutoTLSCertificate = false

	if got, want := buildBaseURL(), "http://127.0.0.1:1234"; got != want {
		t.Fatalf("buildBaseURL() = %q, want %q", got, want)
	}
}

func TestViewLeavesRightmostColumnUnused(t *testing.T) {
	width := 120
	model := &appModel{
		width:          width,
		height:         30,
		logBuffer:      NewLogBuffer(),
		focus:          focusShelf,
		shelfRowToID:   make(map[int]int),
		autoFollowLogs: true,
		status: systemSnapshot{
			CPUPercent: 10,
			RAMPercent: 20,
			StatusText: "running",
		},
	}

	for lineNumber, line := range strings.Split(model.View(), "\n") {
		if got, maxWidth := runewidth.StringWidth(line), width-1; got > maxWidth {
			t.Fatalf("View() line %d width = %d, want <= %d", lineNumber+1, got, maxWidth)
		}
	}
}

func restoreConfig(t *testing.T) {
	t.Helper()
	original := config.CopyCfg()
	t.Cleanup(func() {
		*config.GetCfg() = original
	})
}
