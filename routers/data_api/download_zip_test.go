package data_api

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/comigo_remote"
	"golang.org/x/net/webdav"
)

type zipDownloadStore struct {
	model.StoreInterface
	book *model.Book
}

func (s *zipDownloadStore) GetBook(string) (*model.Book, error) { return s.book, nil }

// 本地与 WebDAV 图片目录均应打包实际页面，覆盖没有 Path 时的文件名回退。
func TestDownloadZipLocalAndWebDAV(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "books")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "001.jpg"), []byte("image fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(&webdav.Handler{FileSystem: webdav.Dir(root), LockSystem: webdav.NewMemLS()})
	defer server.Close()
	old := model.IStore
	defer func() { model.IStore = old }()
	for _, remote := range []bool{false, true} {
		t.Run(map[bool]string{false: "local", true: "webdav"}[remote], func(t *testing.T) {
			book := &model.Book{BookInfo: model.BookInfo{BookID: "fixture", Type: model.TypeDir, BookPath: dir, IsRemote: remote}, PageInfos: model.PageInfos{{Name: "001.jpg"}}}
			if remote {
				book.BookPath = "/books"
				book.RemoteURL = strings.Replace(server.URL, "http://", "webdav://", 1) + "/books"
			}
			model.IStore = &zipDownloadStore{book: book}
			recorder := httptest.NewRecorder()
			c := echo.New().NewContext(httptest.NewRequest("GET", "/api/download-zip?id=fixture", nil), recorder)
			if err := DownloadZip(c); err != nil {
				t.Fatal(err)
			}
			if recorder.Code != 200 {
				t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
			}
			archive, err := zip.NewReader(bytes.NewReader(recorder.Body.Bytes()), int64(recorder.Body.Len()))
			if err != nil {
				t.Fatal(err)
			}
			if len(archive.File) != 1 || archive.File[0].Name != "001.jpg" {
				t.Fatalf("entries=%v", archive.File)
			}
			file, err := archive.File[0].Open()
			if err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(file)
			file.Close()
			if err != nil || string(data) != "image fixture" {
				t.Fatalf("data=%q err=%v", data, err)
			}
		})
	}
}

// 远端 Comigo 下载保留 BasePath、认证和附件头，使用远端 ID，不读取空的本地路径。
func TestDownloadZipProxiesRemoteComigo(t *testing.T) {
	requested := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/reader/api/login" {
			if err := r.ParseForm(); err != nil {
				t.Error(err)
			}
			if r.Form.Get("username") != "fixture" || r.Form.Get("password") != "secret" {
				t.Error("missing credentials")
			}
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "fixture", Path: "/"})
			w.Write([]byte(`{"success":true}`))
			return
		}
		requested = true
		if r.URL.Path != "/reader/api/download-zip" || r.URL.Query().Get("id") != "original" || r.URL.Query().Get("remote_store") != "" {
			t.Errorf("unexpected target: %s", r.URL)
		}
		if _, err := r.Cookie("session"); err != nil {
			t.Error("missing session")
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", `attachment; filename="remote.zip"`)
		w.Write([]byte("remote zip bytes"))
	}))
	defer server.Close()
	saved, old := config.CopyCfg(), model.IStore
	defer func() { config.ReplaceConfig(saved); model.IStore = old }()
	storeURL := strings.Replace(server.URL, "http://", "http://fixture:secret@", 1) + "/reader"
	config.GetCfg().StoreUrls = []string{storeURL}
	key, err := comigo_remote.StoreKey(storeURL)
	if err != nil {
		t.Fatal(err)
	}
	model.IStore = &zipDownloadStore{book: &model.Book{BookInfo: model.BookInfo{BookID: "local", RemoteBookID: "original", RemoteStoreKey: key, IsRemote: true, Type: model.TypeDir}}}
	recorder := httptest.NewRecorder()
	c := echo.New().NewContext(httptest.NewRequest("GET", "/api/download-zip?id=local&remote_store="+url.QueryEscape(key), nil), recorder)
	if err := DownloadZip(c); err != nil {
		t.Fatal(err)
	}
	if !requested || recorder.Code != 200 || recorder.Body.String() != "remote zip bytes" || recorder.Header().Get("Content-Disposition") != `attachment; filename="remote.zip"` {
		t.Fatalf("proxy response: %d %v %s", recorder.Code, recorder.Header(), recorder.Body.String())
	}
}
