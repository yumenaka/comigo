package vfs

import (
	"bytes"
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

	"github.com/pkg/sftp"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
	"golang.org/x/crypto/ssh"
)

// SFTPFS SFTP 文件系统实现
type SFTPFS struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
	baseURL    string
	basePath   string
	options    Options
	cache      *FileCache
	mu         sync.RWMutex
}

// NewSFTPFS 创建 SFTP 文件系统实例
// urlStr 格式: sftp://user:pass@host:port/path
func NewSFTPFS(urlStr string, opts ...Options) (*SFTPFS, error) {
	// 解析 URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("无法解析 SFTP URL: %w", err)
	}

	// 获取主机和端口
	host := parsedURL.Hostname()
	portStr := parsedURL.Port()
	port := 22 // 默认 SFTP 端口
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("无效的端口号: %s", portStr)
		}
		port = p
	}

	// 获取认证信息
	username := ""
	password := ""
	if parsedURL.User != nil {
		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password()
	}

	if username == "" {
		return nil, fmt.Errorf("SFTP URL 必须包含用户名")
	}

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

	// 获取配置选项
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	// 建立 SSH 连接
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：生产环境应验证主机密钥
		Timeout:         time.Duration(options.Timeout) * time.Second,
	}

	// 连接地址
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	// 建立 SSH 连接
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("无法连接 SFTP 服务器 %s: %w", addr, err)
	}

	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, fmt.Errorf("无法创建 SFTP 客户端: %w", err)
	}

	// 测试连接：检查基础路径是否存在
	_, err = sftpClient.Stat(basePath)
	if err != nil {
		sftpClient.Close()
		sshClient.Close()
		return nil, fmt.Errorf("无法访问 SFTP 路径 %s: %w", basePath, err)
	}

	sfs := &SFTPFS{
		sshClient:  sshClient,
		sftpClient: sftpClient,
		baseURL:    urlStr,
		basePath:   basePath,
		options:    options,
	}

	// 初始化缓存
	if options.CacheEnabled {
		sfs.cache = NewFileCache(options.CacheDir, options.Debug)
	}

	if options.Debug {
		logger.Infof(locale.GetString("log_sftp_filesystem_connected"), addr, basePath)
	}

	return sfs, nil
}

// resolvePath 将相对路径解析为 SFTP 服务器上的完整路径
func (s *SFTPFS) resolvePath(p string) string {
	// 如果是绝对路径（以 / 开头），检查是否已经包含 basePath
	if strings.HasPrefix(p, "/") {
		if strings.HasPrefix(p, s.basePath) {
			return p
		}
		// 绝对路径但不包含 basePath，需要拼接
		return path.Join(s.basePath, p)
	}
	// 相对路径，直接拼接
	return path.Join(s.basePath, p)
}

// Open 打开文件用于读取
// 如果启用了缓存，会先检查缓存，未命中则从服务器读取并缓存
func (s *SFTPFS) Open(p string) (File, error) {
	fullPath := s.resolvePath(p)

	// 检查缓存
	if s.cache != nil {
		if data, ok := s.cache.Get(fullPath); ok {
			return newSFTPFile(data, fullPath, s), nil
		}
	}

	// 从服务器读取
	file, err := s.sftpClient.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 SFTP 文件 %s: %w", fullPath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("无法读取 SFTP 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存（如果启用）
	if s.cache != nil {
		s.cache.Set(fullPath, data)
	}

	return newSFTPFile(data, fullPath, s), nil
}

// Stat 获取文件信息
func (s *SFTPFS) Stat(p string) (FileInfo, error) {
	fullPath := s.resolvePath(p)
	info, err := s.sftpClient.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法获取 SFTP 文件信息 %s: %w", fullPath, err)
	}
	return &sftpFileInfo{info: info}, nil
}

// ReadDir 读取目录
func (s *SFTPFS) ReadDir(p string) ([]DirEntry, error) {
	fullPath := s.resolvePath(p)
	files, err := s.sftpClient.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取 SFTP 目录 %s: %w", fullPath, err)
	}

	entries := make([]DirEntry, len(files))
	for i, f := range files {
		entries[i] = &sftpDirEntry{info: &sftpFileInfo{info: f}}
	}
	return entries, nil
}

// ReadFile 读取文件内容
func (s *SFTPFS) ReadFile(p string) ([]byte, error) {
	fullPath := s.resolvePath(p)

	// 检查缓存
	if s.cache != nil {
		if data, ok := s.cache.Get(fullPath); ok {
			return data, nil
		}
	}

	file, err := s.sftpClient.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 SFTP 文件 %s: %w", fullPath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("无法读取 SFTP 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存
	if s.cache != nil {
		s.cache.Set(fullPath, data)
	}

	return data, nil
}

// Type 返回后端类型
func (s *SFTPFS) Type() BackendType {
	return SFTP
}

// BaseURL 返回基础 URL
func (s *SFTPFS) BaseURL() string {
	return s.baseURL
}

// Close 关闭连接
func (s *SFTPFS) Close() error {
	var errs []error
	if s.sftpClient != nil {
		if err := s.sftpClient.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if s.sshClient != nil {
		if err := s.sshClient.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	// 清理缓存
	if s.cache != nil {
		s.cache.Clear()
	}
	if len(errs) > 0 {
		return fmt.Errorf("关闭 SFTP 连接时出错: %v", errs)
	}
	return nil
}

// IsRemote 返回是否为远程文件系统
func (s *SFTPFS) IsRemote() bool {
	return true
}

// JoinPath 连接路径（使用 Unix 风格的路径分隔符）
func (s *SFTPFS) JoinPath(elem ...string) string {
	return path.Join(elem...)
}

// RelPath 计算相对路径
func (s *SFTPFS) RelPath(basePath, targetPath string) (string, error) {
	// 确保路径格式一致
	basePath = path.Clean(basePath)
	targetPath = path.Clean(targetPath)

	// 如果 basePath 是 "." 或 ""，表示根目录，直接返回 targetPath
	if basePath == "." || basePath == "" {
		// 移除开头的 "./" 或 "/"
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

	// 如果 targetPath 以 basePath 或 basePathWithSlash 开头
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

	// 如果 targetPath 不以 basePath 开头，无法计算相对路径
	return "", fmt.Errorf("目标路径 %s 不在基础路径 %s 下", targetPath, basePath)
}

// Exists 检查路径是否存在
func (s *SFTPFS) Exists(p string) (bool, error) {
	fullPath := s.resolvePath(p)
	_, err := s.sftpClient.Stat(fullPath)
	if err != nil {
		// 检查是否是文件不存在的错误
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsDir 检查是否为目录
func (s *SFTPFS) IsDir(p string) (bool, error) {
	fullPath := s.resolvePath(p)
	info, err := s.sftpClient.Stat(fullPath)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// OpenReaderAtSeeker 打开文件并返回支持 Seek 的 Reader
// SFTP 原生支持随机访问，可以直接使用 SFTP 文件句柄
func (s *SFTPFS) OpenReaderAtSeeker(p string) (ReaderAtSeeker, error) {
	fullPath := s.resolvePath(p)

	// 检查完整文件缓存（用于小文件）
	if s.cache != nil {
		if data, ok := s.cache.Get(fullPath); ok {
			return bytes.NewReader(data), nil
		}
	}

	// 获取文件大小
	info, err := s.sftpClient.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法获取文件信息 %s: %w", fullPath, err)
	}
	fileSize := info.Size()

	// 对于小文件（< 1MB），可以考虑缓存以提高性能
	smallFileThreshold := int64(1024 * 1024) // 1MB
	if fileSize < smallFileThreshold {
		// 读取整个文件并缓存
		file, err := s.sftpClient.Open(fullPath)
		if err != nil {
			return nil, fmt.Errorf("无法打开 SFTP 文件 %s: %w", fullPath, err)
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("无法读取 SFTP 文件 %s: %w", fullPath, err)
		}

		// 保存到缓存（如果启用）
		if s.cache != nil {
			s.cache.Set(fullPath, data)
		}

		return bytes.NewReader(data), nil
	}

	// 对于大文件，直接使用 SFTP 文件句柄（支持随机访问）
	file, err := s.sftpClient.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 SFTP 文件 %s: %w", fullPath, err)
	}

	// SFTP 文件句柄实现了 io.ReaderAt 和 io.Seeker
	// 包装为 sftpReaderAtSeeker 以确保实现 ReaderAtSeeker 接口
	return &sftpReaderAtSeeker{file: file}, nil
}

// GetBasePath 返回 SFTP 基础路径
func (s *SFTPFS) GetBasePath() string {
	return s.basePath
}

// sftpReaderAtSeeker 包装 SFTP 文件句柄以实现 ReaderAtSeeker 接口
type sftpReaderAtSeeker struct {
	file *sftp.File
	pos  int64
}

func (r *sftpReaderAtSeeker) Read(p []byte) (n int, err error) {
	n, err = r.file.Read(p)
	r.pos += int64(n)
	return n, err
}

func (r *sftpReaderAtSeeker) ReadAt(p []byte, off int64) (n int, err error) {
	return r.file.ReadAt(p, off)
}

func (r *sftpReaderAtSeeker) Seek(offset int64, whence int) (int64, error) {
	return r.file.Seek(offset, whence)
}

// Close 关闭 SFTP 文件句柄（实现 io.Closer，虽然不是 ReaderAtSeeker 接口的一部分，但调用方会检查）
func (r *sftpReaderAtSeeker) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

// sftpFileInfo 实现 FileInfo 接口
type sftpFileInfo struct {
	info fs.FileInfo
}

func (fi *sftpFileInfo) Name() string       { return fi.info.Name() }
func (fi *sftpFileInfo) Size() int64        { return fi.info.Size() }
func (fi *sftpFileInfo) Mode() fs.FileMode  { return fi.info.Mode() }
func (fi *sftpFileInfo) ModTime() time.Time { return fi.info.ModTime() }
func (fi *sftpFileInfo) IsDir() bool        { return fi.info.IsDir() }
func (fi *sftpFileInfo) Sys() any           { return fi.info.Sys() }

// sftpDirEntry 实现 DirEntry 接口
type sftpDirEntry struct {
	info *sftpFileInfo
}

func (de *sftpDirEntry) Name() string               { return de.info.Name() }
func (de *sftpDirEntry) IsDir() bool                { return de.info.IsDir() }
func (de *sftpDirEntry) Type() fs.FileMode          { return de.info.Mode().Type() }
func (de *sftpDirEntry) Info() (fs.FileInfo, error) { return de.info, nil }

// sftpFile 实现 File 接口
type sftpFile struct {
	data   []byte
	reader *bytes.Reader
	path   string
	sfs    *SFTPFS
}

func newSFTPFile(data []byte, path string, sfs *SFTPFS) *sftpFile {
	return &sftpFile{
		data:   data,
		reader: bytes.NewReader(data),
		path:   path,
		sfs:    sfs,
	}
}

func (f *sftpFile) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

func (f *sftpFile) Close() error {
	// 释放引用，让 GC 回收内存
	f.data = nil
	f.reader = nil
	return nil
}

func (f *sftpFile) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

func (f *sftpFile) Stat() (FileInfo, error) {
	return f.sfs.Stat(f.path)
}

func (f *sftpFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.reader.ReadAt(p, off)
}

// 确保 SFTPFS 实现了 FileSystem 接口
var _ FileSystem = (*SFTPFS)(nil)
