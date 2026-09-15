package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/yumenaka/comigo/config"
)

// 验证默认使用当前目录（包括非终端启动），以及已有配置和参数的优先级。
func TestDefaultScanPath(t *testing.T) {
	cfg := config.GetCfg()
	saved, savedArgs, savedStdin := *cfg, Args, os.Stdin
	// 使用管道模拟非终端启动，默认目录选择不应依赖终端。
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = reader
	t.Cleanup(func() {
		*cfg, Args, os.Stdin = saved, savedArgs, savedStdin
		_ = reader.Close()
		_ = writer.Close()
	})
	for _, tc := range []struct {
		name       string
		dirs       []string
		file       string
		args       []string
		configured bool
		existing   bool
		disabled   bool
	}{
		{name: "disabled", disabled: true},
		{name: "ignore home directories", dirs: []string{"Pictures", "Documents", "Downloads"}},
		{name: "ignore documents", dirs: []string{"Documents", "Downloads"}},
		{name: "ignore downloads", dirs: []string{"Downloads"}},
		{name: "no directories"},
		{name: "skip regular file", dirs: []string{"Documents"}, file: "Pictures"},
		{name: "configuration exists", dirs: []string{"Pictures"}, configured: true},
		{name: "arguments exist", dirs: []string{"Pictures"}, args: []string{"missing"}},
		{name: "library exists", dirs: []string{"Pictures"}, existing: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home, cwd := t.TempDir(), t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("USERPROFILE", home)
			t.Chdir(cwd)
			*cfg, Args = saved, tc.args
			cfg.StoreUrls, cfg.ConfigFile = nil, ""
			cfg.NoDefaultLibrary = tc.disabled
			if tc.configured {
				cfg.ConfigFile = filepath.Join(home, "config.toml")
			}
			if tc.existing {
				cfg.StoreUrls = []string{cwd}
			}
			for _, dir := range tc.dirs {
				if err := os.Mkdir(filepath.Join(home, dir), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			if tc.file != "" {
				if err := os.WriteFile(filepath.Join(home, tc.file), nil, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			SetCwdAsScanPathIfNeed()
			want := cwd
			// 统一解析 macOS 临时目录可能包含的符号链接。
			want, err := filepath.EvalSymlinks(want)
			if err != nil {
				t.Fatal(err)
			}
			got := cfg.StoreUrls
			if len(got) == 1 {
				got = []string{got[0]}
				got[0], err = filepath.EvalSymlinks(got[0])
				if err != nil {
					t.Fatal(err)
				}
			}
			expected := []string{want}
			if tc.disabled || tc.configured || len(tc.args) > 0 {
				expected = nil
			}
			if !reflect.DeepEqual(got, expected) {
				t.Fatalf("书库 = %v，期望 %v", got, expected)
			}
			entries, err := os.ReadDir(home)
			if err != nil {
				t.Fatal(err)
			}
			count := len(tc.dirs)
			if tc.file != "" {
				count++
			}
			if len(entries) != count {
				t.Fatalf("不应创建用户目录: %v", entries)
			}
		})
	}
}
