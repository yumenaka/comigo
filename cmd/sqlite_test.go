package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/data_api"
	"github.com/yumenaka/comigo/sqlc"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/scan"
)

// 验证 SQLite 从扫描到重启、阅读和删除的完整流程，不读取或创建 metadata JSON。
func TestSQLiteLibraryWithoutJSON(t *testing.T) {
	oldCfg, oldStore, oldRAM := config.CopyCfg(), model.IStore, store.RamStore
	t.Cleanup(func() {
		sqlc.CloseDatabase()
		*config.GetCfg(), model.IStore, store.RamStore = oldCfg, oldStore, oldRAM
	})
	root, configDir := t.TempDir(), filepath.Join(t.TempDir(), "config #?目录")
	cfg := config.GetCfg()
	cfg.ConfigFile, cfg.CacheDir = filepath.Join(configDir, "config.toml"), t.TempDir()
	cfg.EnableDatabase, cfg.DBType, cfg.StoreUrls, cfg.MinImageNum = true, "sqlite", []string{root}, 1
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		t.Fatal(err)
	}
	// 用普通文件占住 metadata 目录，任何旧 JSON 路径都会失败。
	if err := os.WriteFile(filepath.Join(configDir, "metadata"), []byte("no JSON storage"), 0o600); err != nil {
		t.Fatal(err)
	}
	store.RamStore = &store.StoreInRam{}
	var img bytes.Buffer
	if err := png.Encode(&img, image.NewRGBA(image.Rect(0, 0, 4, 4))); err != nil {
		t.Fatal(err)
	}
	var archive bytes.Buffer
	zw := zip.NewWriter(&archive)
	for _, name := range []string{"001.png", "002.png"} {
		f, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write(img.Bytes()); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	series := filepath.Join(root, "series")
	if err := os.MkdirAll(series, 0o700); err != nil {
		t.Fatal(err)
	}
	bookPath := filepath.Join(series, "book.cbz")
	for path, data := range map[string][]byte{bookPath: archive.Bytes(), filepath.Join(root, "bad.zip"): []byte("broken archive"), filepath.Join(series, "page.png"): img.Bytes()} {
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := LoadMetadata(); err != nil {
		t.Fatal(err)
	}
	if err := scan.InitAllStore(cfg); err != nil {
		t.Fatal(err)
	}
	SaveMetadata()
	books, err := model.IStore.ListBooks()
	if err != nil {
		t.Fatal(err)
	}
	var book *model.Book
	groups := 0
	for _, b := range books {
		if b.BookPath == bookPath {
			book = b
		}
		if b.Type == model.TypeBooksGroup {
			groups++
		}
	}
	if book == nil || groups != 1 {
		t.Fatalf("scan missing book or group: %#v", books)
	}
	failed, err := sqlc.DbStore.LoadScanFailures()
	if err != nil || len(failed) != 1 {
		t.Fatalf("scan failure not persisted: %v, %v", failed, err)
	}
	e := echo.New()
	e.GET("/books/:id", data_api.GetBook)
	e.GET("/file", data_api.GetFile)
	e.GET("/cover", data_api.GetCover)
	e.GET("/raw/:book_id/:file_name", data_api.GetRawFile)
	e.POST("/bookmarks", data_api.StoreBookmark)
	e.GET("/history", data_api.GetReadingHistory)
	e.DELETE("/books/:id/cache", data_api.DeleteBookCache)
	request := func(method, target, body string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, target, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s %s: %d %s", method, target, rec.Code, rec.Body.String())
		}
		return rec
	}
	request("POST", "/bookmarks", fmt.Sprintf(`{"type":"auto","book_id":%q,"page_index":2}`, book.BookID))
	request("POST", "/bookmarks", fmt.Sprintf(`{"type":"user","book_id":%q,"page_index":1,"description":"书签"}`, book.BookID))
	// 关闭数据库并清空内存，模拟仅凭数据库重启。
	sqlc.CloseDatabase()
	store.RamStore = &store.StoreInRam{}
	model.IStore = store.RamStore
	if err := LoadMetadata(); err != nil {
		t.Fatal(err)
	}
	restored, err := model.IStore.GetBook(book.BookID)
	if err != nil || len(restored.BookMarks) != 2 || !restored.Modified.Equal(book.Modified) || !reflect.DeepEqual(restored.Cover, book.Cover) {
		t.Fatalf("restart lost data: %#v, %v", restored, err)
	}
	request("GET", "/books/"+book.BookID, "")
	file := request("GET", "/file?id="+book.BookID+"&filename="+url.QueryEscape(book.PageInfos[0].Name)+"&disable_cache=true", "")
	if !bytes.Equal(file.Body.Bytes(), img.Bytes()) {
		t.Fatal("reading image differs from archive")
	}
	cover := request("GET", "/cover?id="+book.BookID+"&resize_height=4", "")
	if len(cover.Body.Bytes()) == 0 {
		t.Fatal("empty cover")
	}
	raw := request("GET", "/raw/"+book.BookID+"/book.cbz", "")
	if !bytes.Equal(raw.Body.Bytes(), archive.Bytes()) {
		t.Fatal("download differs from source")
	}
	history := request("GET", "/history", "")
	if !bytes.Contains(history.Body.Bytes(), []byte(book.BookID)) {
		t.Fatal("history lost after restart")
	}
	if err := scan.InitAllStore(cfg); err != nil {
		t.Fatal(err)
	}
	after, err := sqlc.DbStore.LoadScanFailures()
	if err != nil || !reflect.DeepEqual(failed, after) {
		t.Fatalf("unchanged broken archive retried: %v", err)
	}
	restored, err = model.IStore.GetBook(book.BookID)
	if err != nil || len(restored.BookMarks) != 2 {
		t.Fatalf("rescan lost bookmarks: %#v, %v", restored, err)
	}
	request("DELETE", "/books/"+book.BookID+"/cache?delete_cover=false&delete_image_cache=false", "")
	if _, err := model.IStore.GetBook(book.BookID); err == nil {
		t.Fatal("metadata deletion did not reach database")
	}
	SaveMetadata()
	ram, _ := store.RamStore.ListBooks()
	if len(ram) != 0 {
		t.Fatal("database data copied into RAM store")
	}
	if err := filepath.WalkDir(configDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".json" {
			return fmt.Errorf("unexpected JSON: %s", path)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

// 数据库打不开时必须返回错误，不能静默切回 JSON 模式。
func TestSQLiteStartupFailureDoesNotFallback(t *testing.T) {
	oldCfg, oldStore := config.CopyCfg(), model.IStore
	t.Cleanup(func() { sqlc.CloseDatabase(); *config.GetCfg(), model.IStore = oldCfg, oldStore })
	dir := t.TempDir()
	config.GetCfg().ConfigFile = filepath.Join(dir, "config.toml")
	config.GetCfg().EnableDatabase = true
	config.GetCfg().DBType = "sqlite"
	if err := os.Mkdir(filepath.Join(dir, "comigo.sqlite"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := LoadMetadata(); err == nil {
		t.Fatal("expected database open error")
	}
	if !config.GetCfg().EnableDatabase || sqlc.DbStore != nil {
		t.Fatal("database failure changed storage mode")
	}
	if _, err := os.Stat(filepath.Join(dir, "metadata")); !os.IsNotExist(err) {
		t.Fatalf("JSON fallback accessed metadata: %v", err)
	}
}

// CLI 重载也必须拒绝热切换数据库，与配置 API 保持一致。
func TestReloadRejectsDatabaseSwitch(t *testing.T) {
	for _, change := range []string{"EnableDatabase = true", "DBType = 'postgres'", "DBDSN = 'test-dsn'"} {
		t.Run(change, func(t *testing.T) {
			root := processTestCommand(t)
			file := filepath.Join(t.TempDir(), "config.toml")
			if err := os.WriteFile(file, []byte("Port = 2345\nEnableDatabase = false\nDBType = 'sqlite'\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			root.SetArgs([]string{"run", "--config", file})
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			before := config.CopyCfg()
			if err := os.WriteFile(file, []byte("Port = 2345\n"+change+"\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := reloadConfig(file); err == nil || !strings.Contains(err.Error(), "database config requires restart") {
				t.Fatalf("unexpected reload result: %v", err)
			}
			if !reflect.DeepEqual(before, config.CopyCfg()) {
				t.Fatal("failed reload changed active config")
			}
		})
	}
}
