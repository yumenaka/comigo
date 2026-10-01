package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
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

// TestLoadMetadataMigratesOldJSONHistory 用 test 下真实压缩包验证升级重扫、ID 变化和重启后的历史恢复。
func TestLoadMetadataMigratesOldJSONHistory(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "test", "Public Domain Image Archive.zip"))
	if os.IsNotExist(err) {
		t.Skip("test 书籍未提供")
	}
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	bookPath := filepath.Join(root, "book.zip")
	if err := os.WriteFile(bookPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	oldCfg, oldStore, oldRamStore := config.CopyCfg(), model.IStore, store.RamStore
	t.Cleanup(func() { *config.GetCfg(), model.IStore, store.RamStore = oldCfg, oldStore, oldRamStore })
	config.GetCfg().ConfigFile = filepath.Join(t.TempDir(), "config.toml")
	config.GetCfg().StoreUrls = []string{root}
	config.GetCfg().MinImageNum = 1
	ramStore := &store.StoreInRam{}
	model.IStore, store.RamStore = ramStore, ramStore
	ScanStore()
	// 从实际扫描结果中定位测试书籍。
	findBook := func() *model.Book {
		books, _ := model.IStore.ListBooks()
		for _, book := range books {
			if book.BookPath == bookPath {
				return book
			}
		}
		t.Fatal("扫描结果缺少测试书籍")
		return nil
	}
	book := findBook()
	book.BookID = "old-id"
	book.CreatedByVersion = "v1.2.99"
	book.BookComplete = true
	book.BookMarks = model.BookMarks{
		{Type: model.AutoMark, BookID: book.BookID, PageIndex: 3, UpdatedAt: time.Unix(1700000000, 0).UTC()},
		{Type: model.UserMark, BookID: book.BookID, PageIndex: 2, Description: "保留书签", CreatedAt: time.Unix(1690000000, 0).UTC()},
	}
	if err := store.SaveMetaJson(book); err != nil {
		t.Fatal(err)
	}
	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	// 只留下旧版 JSON，模拟升级前磁盘状态。
	entries, err := filepath.Glob(filepath.Join(configDir, "metadata", book.GetStoreID(), "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if filepath.Base(entry) != book.BookID+".json" {
			if err := os.Remove(entry); err != nil {
				t.Fatal(err)
			}
		}
	}
	ramStore = &store.StoreInRam{}
	model.IStore, store.RamStore = ramStore, ramStore
	// 只调用加载入口，验证未请求启动扫描时也自动迁移。
	LoadMetadata()
	refreshed := findBook()
	wantMarks := append(model.BookMarks(nil), book.BookMarks...)
	for i := range wantMarks {
		wantMarks[i].BookID, wantMarks[i].BookStoreID = refreshed.BookID, refreshed.GetStoreID()
	}
	if refreshed.BookID == book.BookID || refreshed.CreatedByVersion != config.GetVersion() || !refreshed.BookComplete || !reflect.DeepEqual(refreshed.BookMarks, wantMarks) {
		t.Fatalf("升级迁移失败: %+v", refreshed)
	}
	ramStore = &store.StoreInRam{}
	model.IStore, store.RamStore = ramStore, ramStore
	LoadMetadata()
	reloaded := findBook()
	if !reflect.DeepEqual(reloaded.BookMarks, wantMarks) || !reloaded.BookComplete || reloaded.GetLastReadPage() != 3 {
		t.Fatalf("重启恢复阅读历史失败: %+v", reloaded)
	}
}
