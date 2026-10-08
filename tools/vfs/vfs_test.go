package vfs

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/yumenaka/comigo/tools"
)

// 验证本地文件系统适配器的基础文件操作。
func TestLocalFS(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "vfs_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := []byte("Hello, VFS!")
	if err := os.WriteFile(testFile, testContent, 0o644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建子目录
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("创建子目录失败: %v", err)
	}

	// 创建本地文件系统适配器。
	fs, err := NewLocalFS(tempDir)
	if err != nil {
		t.Fatalf("创建 LocalFS 失败: %v", err)
	}
	defer fs.Close()

	// 本地文件系统不应被识别为远程书库。
	if fs.IsRemote() {
		t.Error("IsRemote() = true, 期望 false")
	}

	// 读取文件应返回原始内容。
	data, err := fs.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("ReadFile() 失败: %v", err)
	}
	if string(data) != string(testContent) {
		t.Errorf("ReadFile() = %q, 期望 %q", string(data), string(testContent))
	}

	// 文件状态应能反映文件名和目录标记。
	info, err := fs.Stat("test.txt")
	if err != nil {
		t.Fatalf("Stat() 失败: %v", err)
	}
	if info.Name() != "test.txt" {
		t.Errorf("Stat().Name() = %q, 期望 %q", info.Name(), "test.txt")
	}
	if info.IsDir() {
		t.Error("Stat().IsDir() = true, 期望 false")
	}

	// 读取目录应返回文件和子目录。
	entries, err := fs.ReadDir(".")
	if err != nil {
		t.Fatalf("ReadDir() 失败: %v", err)
	}
	if len(entries) != 2 { // test.txt 和 subdir
		t.Errorf("ReadDir() 返回 %d 项, 期望 2 项", len(entries))
	}

	// 文件存在性判断应区分已有文件和缺失文件。
	exists, err := fs.Exists("test.txt")
	if err != nil {
		t.Fatalf("Exists() 失败: %v", err)
	}
	if !exists {
		t.Error("Exists('test.txt') = false, 期望 true")
	}

	exists, err = fs.Exists("nonexistent.txt")
	if err != nil {
		t.Fatalf("Exists() 失败: %v", err)
	}
	if exists {
		t.Error("Exists('nonexistent.txt') = true, 期望 false")
	}

	// 目录判断应区分目录和普通文件。
	isDir, err := fs.IsDir("subdir")
	if err != nil {
		t.Fatalf("IsDir() 失败: %v", err)
	}
	if !isDir {
		t.Error("IsDir('subdir') = false, 期望 true")
	}

	isDir, err = fs.IsDir("test.txt")
	if err != nil {
		t.Fatalf("IsDir() 失败: %v", err)
	}
	if isDir {
		t.Error("IsDir('test.txt') = true, 期望 false")
	}

	// 打开文件后应能读出原始内容。
	file, err := fs.Open("test.txt")
	if err != nil {
		t.Fatalf("Open() 失败: %v", err)
	}
	defer file.Close()

	buf := make([]byte, len(testContent))
	n, err := file.Read(buf)
	if err != nil {
		t.Fatalf("File.Read() 失败: %v", err)
	}
	if n != len(testContent) {
		t.Errorf("File.Read() 读取 %d 字节, 期望 %d 字节", n, len(testContent))
	}
}

// TestGetOrCreateReturnsSingleInstance 验证并发请求不会重复创建并泄漏文件系统实例。
func TestGetOrCreateReturnsSingleInstance(t *testing.T) {
	CloseAll()
	t.Cleanup(CloseAll)
	const workers = 20
	storeDir := t.TempDir()
	results := make(chan FileSystem, workers)
	errs := make(chan error, workers)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fs, err := GetOrCreate(storeDir)
			results <- fs
			errs <- err
		}()
	}
	wg.Wait()
	close(results)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	var first FileSystem
	for fs := range results {
		if first == nil {
			first = fs
			continue
		}
		if fs != first {
			t.Fatal("concurrent GetOrCreate returned different instances")
		}
	}
}

// TestRegistryKeySeparatesCredentialsAndOptions 验证不同身份或配置不会共享可变实例。
func TestRegistryKeySeparatesCredentialsAndOptions(t *testing.T) {
	base := registryKey("webdav://alice:secret@example.com/books", nil)
	if base == registryKey("webdav://bob:secret@example.com/books", nil) {
		t.Fatal("different credentials must not share a VFS instance")
	}
	if base == registryKey("webdav://alice:secret@example.com/books", []Options{{CacheEnabled: true, Timeout: 30}}) {
		t.Fatal("different options must not share a mutable VFS instance")
	}
}

// 验证不同书库地址能解析为对应的虚拟文件系统类型。
func TestParseStoreURL(t *testing.T) {
	tests := []struct {
		url          string
		expectedType tools.StoreBackendType
	}{
		// 本地路径
		{"/home/user/books", tools.StoreBackendLocalDisk},
		{"/Users/test/Documents", tools.StoreBackendLocalDisk},
		{"C:\\Users\\test\\Documents", tools.StoreBackendLocalDisk},
		{"D:/Books", tools.StoreBackendLocalDisk},
		{"file:///home/user/books", tools.StoreBackendLocalDisk},

		// Comigo 远程服务
		{"http://localhost/webdav", tools.StoreBackendComigo},
		{"https://example.com/dav", tools.StoreBackendComigo},

		// WebDAV URL
		{"webdav://192.168.1.1/books", tools.StoreBackendWebDAV},
		{"dav://server/path", tools.StoreBackendWebDAV},
		{"davs://secure-server/path", tools.StoreBackendWebDAV},

		// 其他远程协议
		{"smb://server/share", tools.StoreBackendSMB},
		{"sftp://user@server/path", tools.StoreBackendSFTP},
		{"ftp://server/path", tools.StoreBackendFTP},
		{"ftps://server/path", tools.StoreBackendFTP},
		{"s3://bucket/prefix", tools.StoreBackendS3},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			gotType, _ := parseStoreURL(tt.url)
			if gotType != tt.expectedType {
				t.Errorf("parseStoreURL(%q) = %v, 期望 %v", tt.url, gotType, tt.expectedType)
			}
		})
	}
}

// 验证虚拟文件系统缓存的写入、读取和命中行为。
func TestFileCache(t *testing.T) {
	cache := NewFileCache(false)
	cache.Set("file", []byte("content"))
	if data, ok := cache.Get("file"); !ok || string(data) != "content" {
		t.Fatalf("缓存内容错误: %q, %v", data, ok)
	}
	cache.Clear()
	if _, ok := cache.Get("file"); ok {
		t.Fatal("清理后仍命中缓存")
	}
}

// S3 的三种读取入口共用下载结果，缓存关闭后仍能独立下载。
func TestS3ReadsShareCache(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.Method != http.MethodGet || r.URL.Path != "/bucket/books/page.jpg" {
			t.Errorf("错误的请求: %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}
		_, _ = io.WriteString(w, "image-data")
	}))
	defer server.Close()
	fs := &S3FS{
		client: s3.New(s3.Options{Region: "auto", BaseEndpoint: &server.URL, Credentials: aws.AnonymousCredentials{}, UsePathStyle: true}),
		bucket: "bucket", basePath: "books", options: DefaultOptions(), cache: NewFileCache(false),
	}
	data, err := fs.ReadFile("page.jpg")
	if err != nil || string(data) != "image-data" {
		t.Fatalf("ReadFile: %q, %v", data, err)
	}
	file, err := fs.Open("books/page.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	data, err = io.ReadAll(file)
	if err != nil || string(data) != "image-data" {
		t.Fatalf("Open: %q, %v", data, err)
	}
	reader, err := fs.OpenReaderAtSeeker("/books/page.jpg")
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 4)
	if _, err := reader.ReadAt(buf, 6); err != nil || string(buf) != "data" {
		t.Fatalf("ReadAt: %q, %v", buf, err)
	}
	if requests.Load() != 1 {
		t.Fatalf("缓存命中仍重复下载: %d", requests.Load())
	}
	fs.cache = nil
	if _, err := fs.OpenReaderAtSeeker("page.jpg"); err != nil || requests.Load() != 2 {
		t.Fatalf("无缓存下载失败: %d, %v", requests.Load(), err)
	}
}
