package epub

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"strings"
	"testing"
)

func TestGenerateEscapesMetadataAndSanitizesImageNames(t *testing.T) {
	imageData, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAFgwJ/lD9q2wAAAABJRU5ErkJggg==")
	if err != nil {
		t.Fatal(err)
	}

	images := []ImageFile{{Name: "001 & cover.png", Data: imageData}}
	bookData := CreateBookData("book&1", "A & B EPUB", "Tom & Jerry", images)

	generator, err := NewGenerator()
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := generator.Generate(&buf, bookData, images); err != nil {
		t.Fatal(err)
	}

	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}

	content := readZipFile(t, reader, "OEBPS/content.opf")
	if strings.Contains(string(content), "A & B") || !strings.Contains(string(content), "A &amp; B") {
		t.Fatalf("content.opf title was not XML-escaped:\n%s", content)
	}
	var root struct {
		XMLName xml.Name
	}
	if err := xml.Unmarshal(content, &root); err != nil {
		t.Fatalf("content.opf should be valid XML: %v", err)
	}

	readZipFile(t, reader, "OEBPS/images/image_0001.png")
	page := string(readZipFile(t, reader, "OEBPS/text/page_1.xhtml"))
	if !strings.Contains(page, `src="../images/image_0001.png"`) {
		t.Fatalf("page did not use sanitized image name:\n%s", page)
	}
}

// readZipFile 读取生成包内文件，测试缺文件时直接失败。
func readZipFile(t *testing.T, reader *zip.Reader, name string) []byte {
	t.Helper()
	for _, file := range reader.File {
		if file.Name != name {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		defer rc.Close()
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(rc); err != nil {
			t.Fatal(err)
		}
		return buf.Bytes()
	}
	t.Fatalf("zip file not found: %s", name)
	return nil
}
