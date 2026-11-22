package file

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/yumenaka/comigo/config"
)

// TestPDFErrorsReturnFailure 缺失和损坏的 PDF 都不能作为成功的空图片或零页书籍返回。
func TestPDFErrorsReturnFailure(t *testing.T) {
	dir := t.TempDir()
	broken := filepath.Join(dir, "broken.pdf")
	if err := os.WriteFile(broken, []byte("not a PDF"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{filepath.Join(dir, "missing.pdf"), broken} {
		if data, err := GetImageFromPDF(name, 1, false); err == nil || len(data) != 0 {
			t.Fatalf("错误 PDF 返回成功: %q %d %v", name, len(data), err)
		}
		if _, err := CountPagesOfPDF(name); err == nil {
			t.Fatalf("错误 PDF 页数返回成功: %q", name)
		}
	}
}

// TestEPUBRelativePaths 覆盖正常相对路径及超出根目录的路径，避免下标越界。
func TestEPUBRelativePaths(t *testing.T) {
	for _, tc := range []struct{ src, base, want string }{
		{"../image.jpg", "text/chapter.xhtml", "image.jpg"},
		{"../../image.jpg", "chapter.xhtml", "../../image.jpg"},
		{"/../../image.jpg", "text/chapter.xhtml", "/image.jpg"},
		{"/image.jpg", "text/chapter.xhtml", "/image.jpg"},
		{"https://example.com/image.jpg", "chapter.xhtml", "https://example.com/image.jpg"},
	} {
		if got := absUrl(tc.src, tc.base); got != tc.want {
			t.Errorf("%s / %s = %s; want %s", tc.src, tc.base, got, tc.want)
		}
	}
}

// TestPictureContentTypeMatchesOutput 缩放失败保留 PNG，成功变换为 JPEG，BlurHash 文本继续保留文本类型。
func TestPictureContentTypeMatchesOutput(t *testing.T) {
	dir := t.TempDir()
	var original bytes.Buffer
	if err := png.Encode(&original, image.NewRGBA(image.Rect(0, 0, 4, 4))); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "image.png"), original.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		option GetPictureDataOption
		want   string
	}{
		{"original", GetPictureDataOption{}, "image/png"},
		{"max width unchanged", GetPictureDataOption{ResizeMaxWidth: 10}, "image/png"},
		{"max height unchanged", GetPictureDataOption{ResizeMaxHeight: 10}, "image/png"},
		{"resize", GetPictureDataOption{ResizeWidth: 2}, "image/jpeg"},
		{"gray", GetPictureDataOption{Gray: true}, "image/jpeg"},
		{"blurhash image", GetPictureDataOption{BlurHashImage: 1}, "image/jpeg"},
		{"blurhash text", GetPictureDataOption{BlurHash: 1, BlurHashImage: 1}, "text/plain"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opt := tc.option
			opt.BookIsDir = true
			opt.BookPath = dir
			opt.PictureName = "image.png"
			data, contentType, err := GetPictureData(opt)
			if err != nil || contentType != tc.want {
				t.Fatalf("MIME=%s want=%s err=%v", contentType, tc.want, err)
			}
			if tc.want != "text/plain" && http.DetectContentType(data) != contentType {
				t.Fatalf("编码与 MIME 不符: %s / %s", http.DetectContentType(data), contentType)
			}
		})
	}
	// 无法解码时继续返回原始字节和原 MIME，不冒充 JPEG。
	if err := os.WriteFile(filepath.Join(dir, "image.png"), []byte("broken"), 0600); err != nil {
		t.Fatal(err)
	}
	data, mime, err := GetPictureData(GetPictureDataOption{BookIsDir: true, BookPath: dir, PictureName: "image.png", ResizeWidth: 2, Gray: true, BlurHashImage: 1})
	if err != nil || mime != "image/png" || string(data) != "broken" {
		t.Fatalf("失败回退错误: %s %s %v", data, mime, err)
	}
}

// TestRemotePDFUsesScanCache 读取复用扫描缓存，缺失时向调用方传递后端错误。
func TestRemotePDFUsesScanCache(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { config.ReplaceConfig(old) })
	cfg := old
	cfg.CacheDir = t.TempDir()
	config.ReplaceConfig(cfg)
	cached := filepath.Join(cfg.CacheDir, "remote_books", "book", "sample.pdf")
	if err := os.MkdirAll(filepath.Dir(cached), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cached, []byte("cached"), 0600); err != nil {
		t.Fatal(err)
	}
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		http.Error(w, "offline", http.StatusServiceUnavailable)
	}))
	defer server.Close()
	remote := "webdav" + strings.TrimPrefix(server.URL, "http")
	got, err := getRemotePdfCachePath("book", "/folder/sample.pdf", remote)
	if err != nil || got != cached || requests.Load() != 0 {
		t.Fatalf("未命中扫描缓存: %q %v", got, err)
	}
	if got, err := getRemotePdfCachePath("book", "/folder/missing.pdf", remote); err == nil || got != "" || requests.Load() == 0 {
		t.Fatalf("未返回远端读取失败: %q %v", got, err)
	}
}
