//go:build linux

package autostart

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kardianos/service"
	"github.com/yumenaka/comigo/config"
)

// 使用临时 HOME 和 systemctl 替身验证真实库生成的单元，不能启停本机服务。
func TestUserUnitLifecycle(t *testing.T) {
	if service.ChosenSystem() == nil || service.ChosenSystem().String() != "linux-systemd" {
		t.Skip("requires the systemd service backend")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", home)
	old := config.CopyCfg()
	t.Cleanup(func() { *config.GetCfg() = old })
	config.GetCfg().ConfigFile = filepath.Join(home, "books \"quoted\" %h ${HOME} {{Path}} \\ name.toml")
	stub := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$HOME/actions\"\n"
	if err := os.WriteFile(filepath.Join(home, "systemctl"), []byte(stub), 0700); err != nil {
		t.Fatal(err)
	}
	cfg, err := serviceConfig()
	if err != nil {
		t.Fatal(err)
	}
	manager, err := service.New(nil, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := manager.Install(); err != nil {
		t.Fatal(err)
	}
	path, err := unitPath()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"# Managed by Comigo\n", "ExecStart=:\"", `"--no-default-library"`, `"--config"`, `books \"quoted\" %%h ${HOME} {{Path}} \\ name.toml`, `"--no-tui"`, "WantedBy=default.target"} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("unit is missing %q: %s", want, data)
		}
	}
	if err := manager.Uninstall(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatalf("unit remains after removal: %v", err)
	}
	actions, err := os.ReadFile(filepath.Join(home, "actions"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "enable --user comigo-config.service\ndaemon-reload --user\ndisable --user comigo-config.service\ndaemon-reload --user\n"; string(actions) != want {
		t.Fatalf("unexpected service actions: %s", actions)
	}
}
