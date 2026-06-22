package scan

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/yumenaka/comigo/model"
)

// 验证本地 EPUB 元数据会写回书籍信息，并按阅读顺序排序页面。
func TestApplyEpubInfoFromLocalFileUpdatesMetadataAndSortsPages(t *testing.T) {
	epubPath := filepath.Join(t.TempDir(), "book.epub")
	writeMinimalEpub(t, epubPath)

	book := &model.Book{
		BookInfo: model.BookInfo{
			Type: model.TypeEpub,
		},
		PageInfos: model.PageInfos{
			{Name: "images/002.jpg", PageNum: 2},
			{Name: "images/001.jpg", PageNum: 1},
		},
	}

	applyEpubInfoFromLocalFile(epubPath, book)

	if book.Author != "EPUB Author" {
		t.Fatalf("Author = %q, want EPUB Author", book.Author)
	}
	if book.Press != "EPUB Publisher" {
		t.Fatalf("Press = %q, want EPUB Publisher", book.Press)
	}
	if got := book.PageInfos[0].Name; got != "images/001.jpg" {
		t.Fatalf("first page after EPUB spine sort = %q", got)
	}
}

func writeMinimalEpub(t *testing.T, epubPath string) {
	t.Helper()
	file, err := os.Create(epubPath)
	if err != nil {
		t.Fatalf("create epub: %v", err)
	}
	defer file.Close()

	zw := zip.NewWriter(file)
	defer zw.Close()

	files := map[string]string{
		"META-INF/container.xml": `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`,
		"OEBPS/content.opf": `<?xml version="1.0" encoding="UTF-8"?>
<package version="3.0" unique-identifier="bookid" xmlns="http://www.idpf.org/2007/opf">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Test EPUB</dc:title>
    <dc:creator>EPUB Author</dc:creator>
    <dc:publisher>EPUB Publisher</dc:publisher>
    <dc:language>zh-CN</dc:language>
  </metadata>
  <manifest>
    <item id="page1" href="page1.xhtml" media-type="application/xhtml+xml"/>
    <item id="page2" href="page2.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine>
    <itemref idref="page1"/>
    <itemref idref="page2"/>
  </spine>
</package>`,
		"OEBPS/page1.xhtml": `<html><body><img src="images/001.jpg"/></body></html>`,
		"OEBPS/page2.xhtml": `<html><body><img src="images/002.jpg"/></body></html>`,
	}

	for name, content := range files {
		writer, err := zw.Create(name)
		if err != nil {
			t.Fatalf("create epub entry %s: %v", name, err)
		}
		if _, err := writer.Write([]byte(content)); err != nil {
			t.Fatalf("write epub entry %s: %v", name, err)
		}
	}
}
