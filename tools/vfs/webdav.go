package vfs

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/studio-b12/gowebdav"
	"github.com/yumenaka/comigo/tools/logger"
)

// WebDAVFS WebDAV 文件系统实现
type WebDAVFS struct {
	client   *gowebdav.Client
	baseURL  string
	basePath string // WebDAV 服务器上的基础路径
	options  Options
	cache    *FileCache
	mu       sync.RWMutex
}

// NewWebDAVFS 创建 WebDAV 文件系统实例
// urlStr 格式: webdav://user:pass@host:port/path 或 http://user:pass@host:port/path
func NewWebDAVFS(urlStr string, opts ...Options) (*WebDAVFS, error) {
	// 解析 URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("无法解析 WebDAV URL: %w", err)
	}

	// 确定协议
	scheme := parsedURL.Scheme
	if scheme == "webdav" || scheme == "dav" {
		scheme = "http"
	} else if scheme == "davs" {
		scheme = "https"
	}

	// 构建服务器 URL（不包含路径）
	serverURL := fmt.Sprintf("%s://%s", scheme, parsedURL.Host)

	// 获取基础路径
	basePath := parsedURL.Path
	if basePath == "" {
		basePath = "/"
	}
	// 确保路径以 / 开头
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	// 移除尾部的 /（除非是根路径）
	if basePath != "/" {
		basePath = strings.TrimSuffix(basePath, "/")
	}

	// 获取认证信息
	username := ""
	password := ""
	if parsedURL.User != nil {
		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password()
	}

	// 创建 WebDAV 客户端
	client := gowebdav.NewClient(serverURL, username, password)

	// 设置超时
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}
	client.SetTimeout(time.Duration(options.Timeout) * time.Second)

	// 测试连接
	_, err = client.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("无法连接 WebDAV 服务器: %w", err)
	}

	wfs := &WebDAVFS{
		client:   client,
		baseURL:  urlStr,
		basePath: basePath,
		options:  options,
	}

	// 初始化缓存
	if options.CacheEnabled && options.CacheDir != "" {
		wfs.cache = NewFileCache(options.CacheDir, options.Debug)
	}

	if options.Debug {
		logger.Infof("WebDAV 文件系统已连接: %s, 基础路径: %s", serverURL, basePath)
	}

	return wfs, nil
}

// resolvePath 将相对路径解析为 WebDAV 服务器上的完整路径
func (w *WebDAVFS) resolvePath(p string) string {
	// 如果是绝对路径（以 / 开头），检查是否已经包含 basePath
	if strings.HasPrefix(p, "/") {
		if strings.HasPrefix(p, w.basePath) {
			return p
		}
		// 绝对路径但不包含 basePath，需要拼接
		return path.Join(w.basePath, p)
	}
	// 相对路径，直接拼接
	return path.Join(w.basePath, p)
}

// Open 打开文件用于读取
// 如果启用了缓存，会先检查缓存，未命中则从服务器读取并缓存
func (w *WebDAVFS) Open(p string) (File, error) {
	fullPath := w.resolvePath(p)

	// 检查缓存
	if w.cache != nil {
		if data, ok := w.cache.Get(fullPath); ok {
			return newWebDAVFile(data, fullPath, w), nil
		}
	}

	// 从服务器读取
	reader, err := w.client.ReadStream(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 WebDAV 文件 %s: %w", fullPath, err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("无法读取 WebDAV 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存（如果启用）
	if w.cache != nil {
		w.cache.Set(fullPath, data)
	}

	return newWebDAVFile(data, fullPath, w), nil
}

// Stat 获取文件信息
func (w *WebDAVFS) Stat(p string) (FileInfo, error) {
	fullPath := w.resolvePath(p)
	info, err := w.client.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法获取 WebDAV 文件信息 %s: %w", fullPath, err)
	}
	return &webDAVFileInfo{info: info}, nil
}

// ReadDir 读取目录
func (w *WebDAVFS) ReadDir(p string) ([]DirEntry, error) {
	fullPath := w.resolvePath(p)
	files, err := w.client.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取 WebDAV 目录 %s: %w", fullPath, err)
	}

	entries := make([]DirEntry, len(files))
	for i, f := range files {
		entries[i] = &webDAVDirEntry{info: &webDAVFileInfo{info: f}}
	}
	return entries, nil
}

// ReadFile 读取文件内容
func (w *WebDAVFS) ReadFile(p string) ([]byte, error) {
	fullPath := w.resolvePath(p)

	// 检查缓存
	if w.cache != nil {
		if data, ok := w.cache.Get(fullPath); ok {
			return data, nil
		}
	}

	data, err := w.client.Read(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取 WebDAV 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存
	if w.cache != nil {
		w.cache.Set(fullPath, data)
	}

	return data, nil
}

// Type 返回后端类型
func (w *WebDAVFS) Type() BackendType {
	return WebDAV
}

// BaseURL 返回基础 URL
func (w *WebDAVFS) BaseURL() string {
	return w.baseURL
}

// Close 关闭连接
func (w *WebDAVFS) Close() error {
	// gowebdav 客户端不需要显式关闭
	// 清理缓存
	if w.cache != nil {
		w.cache.Clear()
	}
	return nil
}

// IsRemote 返回是否为远程文件系统
func (w *WebDAVFS) IsRemote() bool {
	return true
}

// JoinPath 连接路径（使用 URL 风格的路径分隔符）
func (w *WebDAVFS) JoinPath(elem ...string) string {
	return path.Join(elem...)
}

// RelPath 计算相对路径
func (w *WebDAVFS) RelPath(basePath, targetPath string) (string, error) {
	// 确保路径格式一致
	basePath = path.Clean(basePath)
	targetPath = path.Clean(targetPath)

	// 如果 targetPath 不以 basePath 开头，无法计算相对路径
	if !strings.HasPrefix(targetPath, basePath) {
		return "", fmt.Errorf("目标路径 %s 不在基础路径 %s 下", targetPath, basePath)
	}

	rel := strings.TrimPrefix(targetPath, basePath)
	rel = strings.TrimPrefix(rel, "/")
	if rel == "" {
		rel = "."
	}
	return rel, nil
}

// Exists 检查路径是否存在
func (w *WebDAVFS) Exists(p string) (bool, error) {
	fullPath := w.resolvePath(p)
	_, err := w.client.Stat(fullPath)
	if err != nil {
		// 检查是否是 404 错误
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsDir 检查是否为目录
func (w *WebDAVFS) IsDir(p string) (bool, error) {
	fullPath := w.resolvePath(p)
	info, err := w.client.Stat(fullPath)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// OpenReaderAtSeeker 打开文件并返回支持 Seek 的 Reader
// WebDAV 协议不支持随机访问，所以需要将整个文件读取到内存中
// 如果启用了缓存，会先检查缓存，未命中则从服务器读取并缓存
// 返回的 bytes.Reader 实现了 io.Reader, io.ReaderAt, io.Seeker 接口
func (w *WebDAVFS) OpenReaderAtSeeker(p string) (ReaderAtSeeker, error) {
	fullPath := w.resolvePath(p)

	// 检查缓存
	if w.cache != nil {
		if data, ok := w.cache.Get(fullPath); ok {
			return bytes.NewReader(data), nil
		}
	}

	// 从服务器读取整个文件到内存
	reader, err := w.client.ReadStream(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 WebDAV 文件 %s: %w", fullPath, err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("无法读取 WebDAV 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存（如果启用）
	if w.cache != nil {
		w.cache.Set(fullPath, data)
	}

	// bytes.Reader 实现了 io.Reader, io.ReaderAt, io.Seeker
	return bytes.NewReader(data), nil
}

// GetBasePath 返回 WebDAV 基础路径
func (w *WebDAVFS) GetBasePath() string {
	return w.basePath
}

// webDAVFileInfo 实现 FileInfo 接口
type webDAVFileInfo struct {
	info fs.FileInfo
}

func (fi *webDAVFileInfo) Name() string       { return fi.info.Name() }
func (fi *webDAVFileInfo) Size() int64        { return fi.info.Size() }
func (fi *webDAVFileInfo) Mode() fs.FileMode  { return fi.info.Mode() }
func (fi *webDAVFileInfo) ModTime() time.Time { return fi.info.ModTime() }
func (fi *webDAVFileInfo) IsDir() bool        { return fi.info.IsDir() }
func (fi *webDAVFileInfo) Sys() any           { return fi.info.Sys() }

// webDAVDirEntry 实现 DirEntry 接口
type webDAVDirEntry struct {
	info *webDAVFileInfo
}

func (de *webDAVDirEntry) Name() string               { return de.info.Name() }
func (de *webDAVDirEntry) IsDir() bool                { return de.info.IsDir() }
func (de *webDAVDirEntry) Type() fs.FileMode          { return de.info.Mode().Type() }
func (de *webDAVDirEntry) Info() (fs.FileInfo, error) { return de.info, nil }

// webDAVFile 实现 File 接口
type webDAVFile struct {
	data   []byte
	reader *bytes.Reader
	path   string
	wfs    *WebDAVFS
}

func newWebDAVFile(data []byte, path string, wfs *WebDAVFS) *webDAVFile {
	return &webDAVFile{
		data:   data,
		reader: bytes.NewReader(data),
		path:   path,
		wfs:    wfs,
	}
}

func (f *webDAVFile) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

func (f *webDAVFile) Close() error {
	// 释放引用，让 GC 回收内存
	f.data = nil
	f.reader = nil
	return nil
}

func (f *webDAVFile) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

func (f *webDAVFile) Stat() (FileInfo, error) {
	return f.wfs.Stat(f.path)
}

func (f *webDAVFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.reader.ReadAt(p, off)
}

// 确保 WebDAVFS 实现了 FileSystem 接口
var _ FileSystem = (*WebDAVFS)(nil)
