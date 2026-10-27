package scan

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
)

// 纯文字 EPUB 必须能入库，而空 ZIP 漫画仍受最少图片数限制。
func TestTextEPUBIsStoredWithoutImages(t *testing.T) {
	oldStore, oldScanCfg, oldCfg := model.IStore, cfg, config.CopyCfg()
	t.Cleanup(func() { model.IStore = oldStore; cfg = oldScanCfg; config.ReplaceConfig(oldCfg) })
	model.IStore = &store.StoreInRam{}
	InitConfig(&testCfgScan{supportMediaType: []string{".jpg"}, supportFileType: []string{".epub", ".zip"}})
	_, err := config.ChangeConfig(false, func(c *config.Config) error { c.MinImageNum = 5; return nil })
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	filename := filepath.Join(dir, "text.epub")
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	z := zip.NewWriter(f)
	for name, content := range map[string]string{
		"mimetype":               "application/epub+zip",
		"META-INF/container.xml": `<container><rootfiles><rootfile full-path="book.opf"/></rootfiles></container>`,
		"book.opf":               `<package xmlns="http://www.idpf.org/2007/opf" version="3.0"><metadata xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:title>Text book</dc:title><dc:creator>Author</dc:creator></metadata><manifest><item id="text" href="text.xhtml" media-type="application/xhtml+xml"/></manifest><spine><itemref idref="text"/></spine></package>`,
		"text.xhtml":             `<html xmlns="http://www.w3.org/1999/xhtml"><head><title>Text</title></head><body><h1>纯文字</h1><p>没有图片也可以阅读。</p></body></html>`,
	} {
		w, err := z.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = w.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	book, err := scanFileGetBook(filename, dir, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(book.PageInfos) != 0 {
		t.Fatal("纯文字不应伪造图片页")
	}
	comic := &model.Book{BookInfo: model.BookInfo{BookID: "empty-comic", Type: model.TypeZip}}
	if err := AddBooksToStore([]*model.Book{book, comic}); err != nil {
		t.Fatal(err)
	}
	if _, err := model.IStore.GetBook(book.BookID); err != nil {
		t.Fatalf("纯文字 EPUB 被过滤: %v", err)
	}
	if _, err := model.IStore.GetBook(comic.BookID); err == nil {
		t.Fatal("空 ZIP 漫画不应入库")
	}
}
