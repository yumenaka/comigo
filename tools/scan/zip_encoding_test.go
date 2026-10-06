package scan

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/klauspost/compress/zip"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/file"
)

// 非 UTF-8 文件名即使目录遍历没有报错，也必须解码并使用相同编码读取图片。
func TestShiftJISZipNamesAndReading(t *testing.T) {
	old := cfg
	t.Cleanup(func() { cfg = old })
	settings := config.CopyCfg()
	settings.ZipFileTextEncoding = "shiftjis"
	cfg = &settings
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)
	entry, err := writer.Create("xx.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = entry.Write([]byte("image fixture")); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}
	archive := bytes.ReplaceAll(buf.Bytes(), []byte("xx.png"), []byte{0x82, 0xa0, '.', 'p', 'n', 'g'})
	path := filepath.Join(t.TempDir(), "book.zip")
	if err = os.WriteFile(path, archive, 0600); err != nil {
		t.Fatal(err)
	}
	book := &model.Book{}
	if err = handleZipAndEpubFiles(path, book); err != nil {
		t.Fatal(err)
	}
	if len(book.PageInfos) != 1 || book.PageInfos[0].Name != "あ.png" || !book.NonUTF8Zip || book.ZipTextEncoding != "shiftjis" {
		t.Fatalf("book=%+v pages=%+v", book.BookInfo, book.PageInfos)
	}
	data, err := file.GetSingleFile(path, book.PageInfos[0].Name, book.ZipTextEncoding)
	if err != nil || string(data) != "image fixture" {
		t.Fatalf("data=%q err=%v", data, err)
	}
}

// 文件没有变化时，编码改变与旧索引迁移也必须重扫，编码一致才继续复用索引。
func TestZIPEncodingInvalidatesExistingMetadata(t *testing.T) {
	oldCfg, oldStore := cfg, model.IStore
	t.Cleanup(func() { cfg = oldCfg; model.IStore = oldStore })
	settings := config.CopyCfg()
	settings.ZipFileTextEncoding = "shiftjis"
	cfg = &settings
	model.IStore = &store.StoreInRam{}
	dir := t.TempDir()
	path := filepath.Join(dir, "book.zip")
	modified := time.Now()
	for _, encoding := range []string{"", "gbk", "shiftjis"} {
		book := &model.Book{BookInfo: model.BookInfo{BookID: "fixture", Type: model.TypeZip, StoreUrl: dir, BookPath: path, FileSize: 12, Modified: modified, ZipTextEncoding: encoding}}
		if err := model.IStore.StoreBook(book); err != nil {
			t.Fatal(err)
		}
		_, skip := prepareBookPathForScan(dir, path, 12, modified)
		if skip != (encoding == "shiftjis") {
			t.Fatalf("encoding=%q skip=%v", encoding, skip)
		}
	}
}
