package vfs

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// FTPFS FTP 文件系统实现
type FTPFS struct {
	conn     *ftp.ServerConn
	baseURL  string
	basePath string
	useTLS   bool // 是否使用 FTPS
	options  Options
	cache    *FileCache
	mu       sync.RWMutex
}

// NewFTPFS 创建 FTP 文件系统实例
// urlStr 格式: ftp://user:pass@host:port/path 或 ftps://user:pass@host:port/path
func NewFTPFS(urlStr string, opts ...Options) (*FTPFS, error) {
	// 解析 URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("无法解析 FTP URL: %w", err)
	}

	// 判断是否使用 TLS
	useTLS := strings.ToLower(parsedURL.Scheme) == "ftps"

	// 获取主机和端口
	host := parsedURL.Hostname()
	portStr := parsedURL.Port()
	defaultPort := 21
	if useTLS {
		defaultPort = 990
	}
	port := defaultPort
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("无效的端口号: %s", portStr)
		}
		port = p
	}

	// 获取认证信息
	username := "anonymous"
	password := "anonymous"
	if parsedURL.User != nil {
		if u := parsedURL.User.Username(); u != "" {
			username = u
		}
		if p, ok := parsedURL.User.Password(); ok {
			password = p
		}
	}

	// 获取基础路径
	basePath := parsedURL.Path
	if basePath == "" {
		basePath = "/"
	}
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	if basePath != "/" {
		basePath = strings.TrimSuffix(basePath, "/")
	}

	// 获取配置选项
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	// 连接地址
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	timeout := time.Duration(options.Timeout) * time.Second

	if options.Debug {
		logger.Infof("正在连接 FTP 服务器 %s (TLS: %v, 超时: %v)", addr, useTLS, timeout)
	}

	// 建立 FTP 连接
	var conn *ftp.ServerConn
	dialOpts := []ftp.DialOption{
		ftp.DialWithTimeout(timeout),
	}

	if useTLS {
		// FTPS: 使用显式 TLS
		dialOpts = append(dialOpts, ftp.DialWithExplicitTLS(&tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证（与其他后端保持一致）
		}))
	}

	conn, err = ftp.Dial(addr, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("无法连接 FTP 服务器 %s: %w", addr, err)
	}

	// 登录
	if err := conn.Login(username, password); err != nil {
		conn.Quit()
		return nil, fmt.Errorf("FTP 登录失败 (用户: %s): %w", username, err)
	}

	// 测试连接：检查基础路径是否存在
	entries, err := conn.List(basePath)
	if err != nil {
		conn.Quit()
		return nil, fmt.Errorf("无法访问 FTP 路径 %s: %w", basePath, err)
	}
	_ = entries // 验证路径可访问即可

	ftpfs := &FTPFS{
		conn:     conn,
		baseURL:  urlStr,
		basePath: basePath,
		useTLS:   useTLS,
		options:  options,
	}

	// 初始化缓存
	if options.CacheEnabled {
		ftpfs.cache = NewFileCache(options.CacheDir, options.Debug)
	}

	if options.Debug {
		logger.Infof(locale.GetString("log_ftp_filesystem_connected"), addr, basePath)
	}

	return ftpfs, nil
}

// resolvePath 将相对路径解析为 FTP 服务器上的完整路径
func (f *FTPFS) resolvePath(p string) string {
	// 如果是绝对路径（以 / 开头），检查是否已经包含 basePath
	if strings.HasPrefix(p, "/") {
		if strings.HasPrefix(p, f.basePath) {
			return p
		}
		return path.Join(f.basePath, p)
	}
	// 相对路径，直接拼接
	return path.Join(f.basePath, p)
}

// Open 打开文件用于读取
func (f *FTPFS) Open(p string) (File, error) {
	fullPath := f.resolvePath(p)

	// 检查缓存
	if f.cache != nil {
		if data, ok := f.cache.Get(fullPath); ok {
			return newFTPFile(data, fullPath, f), nil
		}
	}

	f.mu.Lock()
	resp, err := f.conn.Retr(fullPath)
	f.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("无法打开 FTP 文件 %s: %w", fullPath, err)
	}
	defer resp.Close()

	data, err := io.ReadAll(resp)
	if err != nil {
		return nil, fmt.Errorf("无法读取 FTP 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存
	if f.cache != nil {
		f.cache.Set(fullPath, data)
	}

	return newFTPFile(data, fullPath, f), nil
}

// Stat 获取文件信息
func (f *FTPFS) Stat(p string) (FileInfo, error) {
	fullPath := f.resolvePath(p)

	f.mu.Lock()
	entry, err := f.conn.GetEntry(fullPath)
	f.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("无法获取 FTP 文件信息 %s: %w", fullPath, err)
	}

	return &ftpFileInfo{entry: entry, filePath: fullPath}, nil
}

// ReadDir 读取目录
func (f *FTPFS) ReadDir(p string) ([]DirEntry, error) {
	fullPath := f.resolvePath(p)

	f.mu.Lock()
	entries, err := f.conn.List(fullPath)
	f.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("无法读取 FTP 目录 %s: %w", fullPath, err)
	}

	var result []DirEntry
	for _, e := range entries {
		// 跳过 "." 和 ".." 目录项
		if e.Name == "." || e.Name == ".." {
			continue
		}
		result = append(result, &ftpDirEntry{entry: e})
	}
	return result, nil
}

// ReadFile 读取文件内容
func (f *FTPFS) ReadFile(p string) ([]byte, error) {
	fullPath := f.resolvePath(p)

	// 检查缓存
	if f.cache != nil {
		if data, ok := f.cache.Get(fullPath); ok {
			return data, nil
		}
	}

	f.mu.Lock()
	resp, err := f.conn.Retr(fullPath)
	f.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("无法打开 FTP 文件 %s: %w", fullPath, err)
	}
	defer resp.Close()

	data, err := io.ReadAll(resp)
	if err != nil {
		return nil, fmt.Errorf("无法读取 FTP 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存
	if f.cache != nil {
		f.cache.Set(fullPath, data)
	}

	return data, nil
}

// Type 返回后端类型
func (f *FTPFS) Type() BackendType {
	return FTP
}

// BaseURL 返回基础 URL
func (f *FTPFS) BaseURL() string {
	return f.baseURL
}

// Close 关闭连接
func (f *FTPFS) Close() error {
	// 清理缓存
	if f.cache != nil {
		f.cache.Clear()
	}
	if f.conn != nil {
		return f.conn.Quit()
	}
	return nil
}

// IsRemote 返回是否为远程文件系统
func (f *FTPFS) IsRemote() bool {
	return true
}

// JoinPath 连接路径（使用 Unix 风格的路径分隔符）
func (f *FTPFS) JoinPath(elem ...string) string {
	return path.Join(elem...)
}

// RelPath 计算相对路径
func (f *FTPFS) RelPath(basePath, targetPath string) (string, error) {
	basePath = path.Clean(basePath)
	targetPath = path.Clean(targetPath)

	// 如果 basePath 是 "." 或 ""，表示根目录，直接返回 targetPath
	if basePath == "." || basePath == "" {
		rel := strings.TrimPrefix(targetPath, "./")
		rel = strings.TrimPrefix(rel, "/")
		if rel == "" {
			rel = "."
		}
		return rel, nil
	}

	// 确保 basePath 以 "/" 结尾以便正确匹配
	basePathWithSlash := basePath
	if !strings.HasSuffix(basePathWithSlash, "/") {
		basePathWithSlash = basePathWithSlash + "/"
	}

	if strings.HasPrefix(targetPath, basePathWithSlash) {
		rel := strings.TrimPrefix(targetPath, basePathWithSlash)
		if rel == "" {
			rel = "."
		}
		return rel, nil
	}

	if strings.HasPrefix(targetPath, basePath) {
		rel := strings.TrimPrefix(targetPath, basePath)
		rel = strings.TrimPrefix(rel, "/")
		if rel == "" {
			rel = "."
		}
		return rel, nil
	}

	return "", fmt.Errorf("目标路径 %s 不在基础路径 %s 下", targetPath, basePath)
}

// Exists 检查路径是否存在
func (f *FTPFS) Exists(p string) (bool, error) {
	fullPath := f.resolvePath(p)

	f.mu.Lock()
	_, err := f.conn.GetEntry(fullPath)
	f.mu.Unlock()
	if err != nil {
		// FTP 文件不存在通常返回 550 错误
		errStr := err.Error()
		if strings.Contains(errStr, "550") ||
			strings.Contains(errStr, "does not exist") ||
			strings.Contains(errStr, "not found") ||
			strings.Contains(errStr, "No such file") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsDir 检查是否为目录
func (f *FTPFS) IsDir(p string) (bool, error) {
	fullPath := f.resolvePath(p)

	f.mu.Lock()
	entry, err := f.conn.GetEntry(fullPath)
	f.mu.Unlock()
	if err != nil {
		return false, err
	}
	return entry.Type == ftp.EntryTypeFolder, nil
}

// OpenReaderAtSeeker 打开文件并返回支持 Seek 的 Reader
// FTP 不支持随机访问（ReadAt），统一下载到内存并缓存以避免重复下载
func (f *FTPFS) OpenReaderAtSeeker(p string) (ReaderAtSeeker, error) {
	fullPath := f.resolvePath(p)

	// 优先使用缓存
	if f.cache != nil {
		if data, ok := f.cache.Get(fullPath); ok {
			return bytes.NewReader(data), nil
		}
	}

	// 下载文件到内存
	f.mu.Lock()
	resp, err := f.conn.Retr(fullPath)
	f.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("无法打开 FTP 文件 %s: %w", fullPath, err)
	}
	defer resp.Close()

	data, err := io.ReadAll(resp)
	if err != nil {
		return nil, fmt.Errorf("无法读取 FTP 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存，避免重复下载
	if f.cache != nil {
		f.cache.Set(fullPath, data)
	}

	return bytes.NewReader(data), nil
}

// GetBasePath 返回 FTP 基础路径
func (f *FTPFS) GetBasePath() string {
	return f.basePath
}

// --- 辅助类型 ---

// ftpFileInfo 实现 FileInfo 接口
type ftpFileInfo struct {
	entry    *ftp.Entry
	filePath string
}

func (fi *ftpFileInfo) Name() string {
	if fi.entry.Name != "" {
		return fi.entry.Name
	}
	return path.Base(fi.filePath)
}
func (fi *ftpFileInfo) Size() int64 { return int64(fi.entry.Size) }
func (fi *ftpFileInfo) Mode() fs.FileMode {
	if fi.entry.Type == ftp.EntryTypeFolder {
		return fs.ModeDir | 0o755
	}
	return 0o644
}
func (fi *ftpFileInfo) ModTime() time.Time { return fi.entry.Time }
func (fi *ftpFileInfo) IsDir() bool        { return fi.entry.Type == ftp.EntryTypeFolder }
func (fi *ftpFileInfo) Sys() any           { return fi.entry }

// ftpDirEntry 实现 DirEntry 接口
type ftpDirEntry struct {
	entry *ftp.Entry
}

func (de *ftpDirEntry) Name() string { return de.entry.Name }
func (de *ftpDirEntry) IsDir() bool  { return de.entry.Type == ftp.EntryTypeFolder }
func (de *ftpDirEntry) Type() fs.FileMode {
	if de.entry.Type == ftp.EntryTypeFolder {
		return fs.ModeDir
	}
	return 0
}
func (de *ftpDirEntry) Info() (fs.FileInfo, error) {
	return &ftpFileInfo{entry: de.entry}, nil
}

// ftpFile 实现 File 接口（基于内存的 bytes.Reader 包装）
type ftpFile struct {
	data   []byte
	reader *bytes.Reader
	path   string
	ffs    *FTPFS
}

func newFTPFile(data []byte, path string, ffs *FTPFS) *ftpFile {
	return &ftpFile{
		data:   data,
		reader: bytes.NewReader(data),
		path:   path,
		ffs:    ffs,
	}
}

func (fi *ftpFile) Read(p []byte) (n int, err error) {
	return fi.reader.Read(p)
}

func (fi *ftpFile) Close() error {
	fi.data = nil
	fi.reader = nil
	return nil
}

func (fi *ftpFile) Seek(offset int64, whence int) (int64, error) {
	return fi.reader.Seek(offset, whence)
}

func (fi *ftpFile) Stat() (FileInfo, error) {
	return fi.ffs.Stat(fi.path)
}

func (fi *ftpFile) ReadAt(p []byte, off int64) (n int, err error) {
	return fi.reader.ReadAt(p, off)
}

// 确保 FTPFS 实现了 FileSystem 接口
var _ FileSystem = (*FTPFS)(nil)
