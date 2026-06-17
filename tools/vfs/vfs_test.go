package vfs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yumenaka/comigo/tools"
)

// TestLocalFS 测试本地文件系统实现
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

	// 创建 LocalFS 实例
	fs, err := NewLocalFS(tempDir)
	if err != nil {
		t.Fatalf("创建 LocalFS 失败: %v", err)
	}
	defer fs.Close()

	// 测试 IsRemote()
	if fs.IsRemote() {
		t.Error("IsRemote() = true, 期望 false")
	}

	// 测试 ReadFile()
	data, err := fs.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("ReadFile() 失败: %v", err)
	}
	if string(data) != string(testContent) {
		t.Errorf("ReadFile() = %q, 期望 %q", string(data), string(testContent))
	}

	// 测试 Stat()
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

	// 测试 ReadDir()
	entries, err := fs.ReadDir(".")
	if err != nil {
		t.Fatalf("ReadDir() 失败: %v", err)
	}
	if len(entries) != 2 { // test.txt 和 subdir
		t.Errorf("ReadDir() 返回 %d 项, 期望 2 项", len(entries))
	}

	// 测试 Exists()
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

	// 测试 IsDir()
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

	// 测试 Open()
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

// TestParseStoreURL 测试 URL 解析
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

// TestFileCache 测试文件缓存
func TestFileCache(t *testing.T) {
	// 创建临时缓存目录
	tempDir, err := os.MkdirTemp("", "vfs_cache_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cache := NewFileCache(tempDir, false)

	// 测试 Set 和 Get
	testPath := "/remote/path/to/file.txt"
	testData := []byte("cached content")

	cache.Set(testPath, testData)

	data, ok := cache.Get(testPath)
	if !ok {
		t.Error("Get() 返回 false, 期望 true")
	}
	if string(data) != string(testData) {
		t.Errorf("Get() = %q, 期望 %q", string(data), string(testData))
	}

	// 测试 Size
	if cache.Size() != 1 {
		t.Errorf("Size() = %d, 期望 1", cache.Size())
	}

	// 测试 Delete
	cache.Delete(testPath)
	_, ok = cache.Get(testPath)
	if ok {
		t.Error("Delete 后 Get() 返回 true, 期望 false")
	}

	// 测试 Clear
	cache.Set("/path1", []byte("data1"))
	cache.Set("/path2", []byte("data2"))
	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("Clear 后 Size() = %d, 期望 0", cache.Size())
	}
}
