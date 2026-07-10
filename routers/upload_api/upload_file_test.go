package upload_api

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
)

// TestParseUploadMultipartFormRejectsOversizedBody 验证请求在落完整临时文件前被限制。
func TestParseUploadMultipartFormRejectsOversizedBody(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("files", "book.cbz")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("uploaded")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	e := echo.New()
	if _, err := parseUploadMultipartForm(e.NewContext(req, httptest.NewRecorder()), int64(body.Len()-1)); err == nil {
		t.Fatal("oversized multipart body should fail")
	}
}

// TestSaveUploadedFileUsesUniqueName 验证并发安全的创建方式不会覆盖同名文件。
func TestSaveUploadedFileUsesUniqueName(t *testing.T) {
	storeDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(storeDir, "book.cbz"), []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("files", "book.cbz")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("uploaded")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}
	filename, err := saveUploadedFile(storeDir, "book.cbz", req.MultipartForm.File["files"][0])
	if err != nil {
		t.Fatal(err)
	}
	if filename != "book_1.cbz" {
		t.Fatalf("filename = %q, want book_1.cbz", filename)
	}
	data, err := os.ReadFile(filepath.Join(storeDir, filename))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "uploaded" {
		t.Fatalf("uploaded data = %q", data)
	}
}
