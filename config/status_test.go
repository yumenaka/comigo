package config

import (
	"os"
	"path/filepath"
	"testing"
)

// 文件位置、运行配置类型及删除状态以实际运行配置为准。
func TestConfigFileStatus(t *testing.T) {
	old := CopyCfg()
	oldProfile := runtimeConfigProfile
	t.Cleanup(func() { *GetCfg() = old; runtimeConfigProfile = oldProfile })
	home := t.TempDir()
	t.Setenv("HOME", home)
	work := t.TempDir()
	t.Chdir(work)
	GetCfg().TemporaryReaderMode = false
	GetCfg().ConfigFile = ""
	homeConfig := filepath.Join(home, ".config", "comigo", "config.toml")
	if err := os.MkdirAll(filepath.Dir(homeConfig), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(homeConfig, []byte("Port = 1234"), 0600); err != nil {
		t.Fatal(err)
	}
	var status Status
	if err := status.SetConfigStatus(); err != nil {
		t.Fatal(err)
	}
	if status.Current.Path != "" || status.In != "None" || status.Path.HomeDirectory != homeConfig {
		t.Fatalf("内存配置不能误认成磁盘文件: %+v", status)
	}
	for _, profile := range []string{"cli", "desktop", "tray"} {
		runtimeConfigProfile = profile
		for _, tc := range []struct{ path, location string }{
			{homeConfig, HomeDirectory},
			{filepath.Join(work, "custom.toml"), WorkingDirectory},
			{filepath.Join(t.TempDir(), "reader.toml"), "Custom"},
		} {
			if err := os.WriteFile(tc.path, []byte("Port = 1234"), 0600); err != nil {
				t.Fatal(err)
			}
			GetCfg().ConfigFile = tc.path
			if err := status.SetConfigStatus(); err != nil {
				t.Fatal(err)
			}
			if status.Current.Path != tc.path || status.Current.Location != tc.location || status.Current.Type != profile || status.Current.Format != "toml" || !status.Current.Exists {
				t.Fatalf("文件信息不符: %+v", status)
			}
			if err := os.Remove(tc.path); err != nil {
				t.Fatal(err)
			}
			if err := status.SetConfigStatus(); err != nil {
				t.Fatal(err)
			}
			if status.Current.Exists || status.Current.Path != tc.path {
				t.Fatalf("删除状态不符: %+v", status)
			}
		}
	}
	GetCfg().TemporaryReaderMode = true
	if err := status.SetConfigStatus(); err != nil {
		t.Fatal(err)
	}
	if status.Current.Path != "" || status.Current.Exists {
		t.Fatalf("临时阅读不应报告持久配置: %+v", status)
	}
}
