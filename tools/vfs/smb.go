package vfs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudsoda/go-smb2"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// SMBFS SMB 文件系统实现
type SMBFS struct {
	session   *smb2.Session
	share     *smb2.Share
	baseURL   string
	basePath  string  // 共享内的基础路径
	shareName string  // SMB 共享名称
	server    string  // 服务器地址（host:port）
	options   Options
	cache     *FileCache
	mu        sync.RWMutex
}

// NewSMBFS 创建 SMB 文件系统实例
// urlStr 格式: smb://workgroup;user:password@server/share/folder/books
func NewSMBFS(urlStr string, opts ...Options) (*SMBFS, error) {
	// 解析 URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("无法解析 SMB URL: %w", err)
	}

	// 获取主机和端口
	host := parsedURL.Hostname()
	portStr := parsedURL.Port()
	port := 445 // 默认 SMB 端口
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("无效的端口号: %s", portStr)
		}
		port = p
	}
	server := net.JoinHostPort(host, strconv.Itoa(port))

	// 获取认证信息
	username := ""
	password := ""
	workgroup := ""
	if parsedURL.User != nil {
		usernameRaw := parsedURL.User.Username()
		// SMB格式可能是 workgroup;user 或 user
		if strings.Contains(usernameRaw, ";") {
			parts := strings.SplitN(usernameRaw, ";", 2)
			if len(parts) > 1 {
				workgroup = parts[0]
				username = parts[1]
			} else {
				username = usernameRaw
			}
		} else {
			username = usernameRaw
		}
		password, _ = parsedURL.User.Password()
		// 如果密码为空，确保是空字符串（用于访客访问）
		if password == "" {
			password = ""
		}
	}

	if username == "" {
		return nil, fmt.Errorf("SMB URL 必须包含用户名")
	}

	// 解析路径部分：/share/folder/books
	path := parsedURL.Path
	path = strings.TrimPrefix(path, "/")
	pathParts := strings.SplitN(path, "/", 2)
	if len(pathParts) < 1 {
		return nil, fmt.Errorf("SMB URL 必须包含共享名称")
	}

	shareName := pathParts[0]
	basePath := ""
	if len(pathParts) > 1 {
		basePath = pathParts[1]
	}
	if basePath == "" {
		basePath = "."
	}
	// 标准化路径：使用反斜杠（Windows 风格）
	basePath = filepath.Clean(basePath)
	// 将反斜杠转换为正斜杠（go-smb2 使用正斜杠）
	basePath = strings.ReplaceAll(basePath, "\\", "/")

	// 获取配置选项
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	// 创建 SMB Dialer
	// 对于访客访问（空密码），确保使用空字符串而不是 nil
	ntlmInitiator := &smb2.NTLMInitiator{
		User:     username,
		Password: password,
		Domain:   workgroup, // workgroup 作为 domain
	}
	// 如果密码为空，确保是空字符串（访客访问）
	if password == "" {
		ntlmInitiator.Password = ""
	}

	dialer := &smb2.Dialer{
		Initiator: ntlmInitiator,
	}

	// 建立 SMB 会话
	// go-smb2 的 Dial 方法接受 context 和地址
	// SMB 连接可能需要更长时间，特别是第一次连接，使用至少 60 秒超时
	timeout := options.Timeout
	if timeout < 60 {
		timeout = 60 // 至少 60 秒超时
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if options.Debug {
		logger.Infof("正在连接 SMB 服务器 %s (超时: %d秒, 用户: %s, 共享: %s)", server, timeout, username, shareName)
	}

	session, err := dialer.Dial(ctx, server)
	if err != nil {
		// 提供更详细的错误信息
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("连接 SMB 服务器 %s 超时（%d秒）。请检查网络连接和服务器是否可达: %w", server, timeout, err)
		}
		return nil, fmt.Errorf("无法连接 SMB 服务器 %s: %w", server, err)
	}

	// 挂载共享
	if options.Debug {
		logger.Infof("正在挂载 SMB 共享: %s", shareName)
	}
	share, err := session.Mount(shareName)
	if err != nil {
		session.Logoff()
		return nil, fmt.Errorf("无法挂载 SMB 共享 %s: %w", shareName, err)
	}

	// 测试连接：检查基础路径是否存在
	fullPath := basePath
	if fullPath == "." {
		fullPath = ""
	}
	_, err = share.Stat(fullPath)
	if err != nil {
		share.Umount()
		session.Logoff()
		return nil, fmt.Errorf("无法访问 SMB 路径 %s: %w", fullPath, err)
	}

	smbfs := &SMBFS{
		session:   session,
		share:     share,
		baseURL:   urlStr,
		basePath:  basePath,
		shareName: shareName,
		server:    server,
		options:   options,
	}

	// 初始化缓存
	if options.CacheEnabled {
		smbfs.cache = NewFileCache(options.CacheDir, options.Debug)
	}

	if options.Debug {
		logger.Infof(locale.GetString("log_smb_filesystem_connected"), server, basePath)
	}

	return smbfs, nil
}

// resolvePath 将相对路径解析为共享内的完整路径
func (s *SMBFS) resolvePath(p string) string {
	// 标准化路径：统一使用正斜杠（SMB 使用正斜杠）
	p = strings.ReplaceAll(p, "\\", "/")
	// 清理路径中的多余斜杠和 . 和 ..
	p = filepath.Clean(p)
	// 再次转换为正斜杠（filepath.Clean 可能返回反斜杠）
	p = strings.ReplaceAll(p, "\\", "/")
	// 移除开头的 /（如果有）
	p = strings.TrimPrefix(p, "/")

	// 如果 basePath 是 "." 或空，直接返回路径（表示根目录）
	if s.basePath == "." || s.basePath == "" {
		if p == "" {
			return "."
		}
		return p
	}

	// 标准化 basePath（确保使用正斜杠）
	basePath := strings.ReplaceAll(s.basePath, "\\", "/")
	basePath = strings.TrimPrefix(basePath, "/")
	basePath = strings.TrimSuffix(basePath, "/")

	// 如果路径已经是完整路径（以 basePath 开头），直接返回
	// 需要确保匹配时考虑路径边界（避免部分匹配）
	if strings.HasPrefix(p, basePath+"/") || p == basePath {
		return p
	}

	// 拼接 basePath 和相对路径
	if p == "" || p == "." {
		return basePath
	}
	
	// 使用正斜杠拼接路径
	if basePath == "." || basePath == "" {
		return p
	}
	return basePath + "/" + p
}

// Open 打开文件用于读取
func (s *SMBFS) Open(p string) (File, error) {
	fullPath := s.resolvePath(p)

	// 检查缓存
	if s.cache != nil {
		if data, ok := s.cache.Get(fullPath); ok {
			return newSMBFile(data, fullPath, s), nil
		}
	}

	// 从共享读取
	file, err := s.share.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 SMB 文件 %s: %w", fullPath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("无法读取 SMB 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存（如果启用）
	if s.cache != nil {
		s.cache.Set(fullPath, data)
	}

	return newSMBFile(data, fullPath, s), nil
}

// Stat 获取文件信息
func (s *SMBFS) Stat(p string) (FileInfo, error) {
	fullPath := s.resolvePath(p)
	info, err := s.share.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法获取 SMB 文件信息 %s: %w", fullPath, err)
	}
	return &smbFileInfo{info: info}, nil
}

// ReadDir 读取目录
func (s *SMBFS) ReadDir(p string) ([]DirEntry, error) {
	fullPath := s.resolvePath(p)
	files, err := s.share.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取 SMB 目录 %s: %w", fullPath, err)
	}

	entries := make([]DirEntry, len(files))
	for i, f := range files {
		entries[i] = &smbDirEntry{info: &smbFileInfo{info: f}}
	}
	return entries, nil
}

// ReadFile 读取文件内容
func (s *SMBFS) ReadFile(p string) ([]byte, error) {
	fullPath := s.resolvePath(p)

	// 检查缓存
	if s.cache != nil {
		if data, ok := s.cache.Get(fullPath); ok {
			return data, nil
		}
	}

	file, err := s.share.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 SMB 文件 %s: %w", fullPath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("无法读取 SMB 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存
	if s.cache != nil {
		s.cache.Set(fullPath, data)
	}

	return data, nil
}

// Type 返回后端类型
func (s *SMBFS) Type() BackendType {
	return SMB
}

// BaseURL 返回基础 URL
func (s *SMBFS) BaseURL() string {
	return s.baseURL
}

// Close 关闭连接
func (s *SMBFS) Close() error {
	var errs []error
	if s.share != nil {
		if err := s.share.Umount(); err != nil {
			errs = append(errs, err)
		}
	}
	if s.session != nil {
		if err := s.session.Logoff(); err != nil {
			errs = append(errs, err)
		}
	}
	// 清理缓存
	if s.cache != nil {
		s.cache.Clear()
	}
	if len(errs) > 0 {
		return fmt.Errorf("关闭 SMB 连接时出错: %v", errs)
	}
	return nil
}

// IsRemote 返回是否为远程文件系统
func (s *SMBFS) IsRemote() bool {
	return true
}

// JoinPath 连接路径（使用正斜杠，go-smb2 使用正斜杠）
func (s *SMBFS) JoinPath(elem ...string) string {
	// 使用 filepath.Join 然后转换为正斜杠
	joined := filepath.Join(elem...)
	return strings.ReplaceAll(joined, "\\", "/")
}

// RelPath 计算相对路径
func (s *SMBFS) RelPath(basePath, targetPath string) (string, error) {
	// 标准化路径
	basePath = filepath.Clean(basePath)
	targetPath = filepath.Clean(targetPath)
	basePath = strings.ReplaceAll(basePath, "\\", "/")
	targetPath = strings.ReplaceAll(targetPath, "\\", "/")

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
func (s *SMBFS) Exists(p string) (bool, error) {
	fullPath := s.resolvePath(p)
	_, err := s.share.Stat(fullPath)
	if err != nil {
		// 检查是否是文件不存在的错误
		if strings.Contains(err.Error(), "does not exist") ||
			strings.Contains(err.Error(), "not found") ||
			strings.Contains(err.Error(), "no such file") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsDir 检查是否为目录
func (s *SMBFS) IsDir(p string) (bool, error) {
	fullPath := s.resolvePath(p)
	info, err := s.share.Stat(fullPath)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// OpenReaderAtSeeker 打开文件并返回支持 Seek 的 Reader
// SMB 文件句柄支持 Seek，但不一定支持 ReadAt
// 所有文件都会下载到内存以确保兼容性，并缓存以避免重复下载
func (s *SMBFS) OpenReaderAtSeeker(p string) (ReaderAtSeeker, error) {
	fullPath := s.resolvePath(p)

	// 检查完整文件缓存（优先使用缓存）
	if s.cache != nil {
		if data, ok := s.cache.Get(fullPath); ok {
			return bytes.NewReader(data), nil
		}
	}

	// 获取文件大小
	info, err := s.share.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法获取文件信息 %s: %w", fullPath, err)
	}
	fileSize := info.Size()

	// 对于小文件（< 1MB），下载到内存并缓存
	smallFileThreshold := int64(1024 * 1024) // 1MB
	if fileSize < smallFileThreshold {
		file, err := s.share.Open(fullPath)
		if err != nil {
			return nil, fmt.Errorf("无法打开 SMB 文件 %s: %w", fullPath, err)
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("无法读取 SMB 文件 %s: %w", fullPath, err)
		}

		// 保存到缓存（如果启用）
		if s.cache != nil {
			s.cache.Set(fullPath, data)
		}

		return bytes.NewReader(data), nil
	}

	// 对于大文件，尝试使用 SMB 文件句柄（如果支持 Seek）
	// 注意：go-smb2 的 File 可能不支持 ReadAt，所以我们需要下载到内存
	file, err := s.share.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 SMB 文件 %s: %w", fullPath, err)
	}

	// 读取整个文件到内存（SMB 文件句柄可能不支持 ReadAt）
	data, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		return nil, fmt.Errorf("无法读取 SMB 文件 %s: %w", fullPath, err)
	}

	// 保存到缓存（如果启用），避免重复下载
	if s.cache != nil {
		s.cache.Set(fullPath, data)
	}

	return bytes.NewReader(data), nil
}

// GetBasePath 返回 SMB 基础路径
func (s *SMBFS) GetBasePath() string {
	return s.basePath
}

// smbFileInfo 实现 FileInfo 接口
type smbFileInfo struct {
	info fs.FileInfo
}

func (fi *smbFileInfo) Name() string       { return fi.info.Name() }
func (fi *smbFileInfo) Size() int64        { return fi.info.Size() }
func (fi *smbFileInfo) Mode() fs.FileMode  { return fi.info.Mode() }
func (fi *smbFileInfo) ModTime() time.Time { return fi.info.ModTime() }
func (fi *smbFileInfo) IsDir() bool        { return fi.info.IsDir() }
func (fi *smbFileInfo) Sys() any           { return fi.info.Sys() }

// smbDirEntry 实现 DirEntry 接口
type smbDirEntry struct {
	info *smbFileInfo
}

func (de *smbDirEntry) Name() string               { return de.info.Name() }
func (de *smbDirEntry) IsDir() bool                { return de.info.IsDir() }
func (de *smbDirEntry) Type() fs.FileMode          { return de.info.Mode().Type() }
func (de *smbDirEntry) Info() (fs.FileInfo, error) { return de.info, nil }

// smbFile 实现 File 接口
type smbFile struct {
	data   []byte
	reader *bytes.Reader
	path   string
	sfs    *SMBFS
}

func newSMBFile(data []byte, path string, sfs *SMBFS) *smbFile {
	return &smbFile{
		data:   data,
		reader: bytes.NewReader(data),
		path:   path,
		sfs:    sfs,
	}
}

func (f *smbFile) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

func (f *smbFile) Close() error {
	// 释放引用，让 GC 回收内存
	f.data = nil
	f.reader = nil
	return nil
}

func (f *smbFile) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

func (f *smbFile) Stat() (FileInfo, error) {
	return f.sfs.Stat(f.path)
}

func (f *smbFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.reader.ReadAt(p, off)
}

// 确保 SMBFS 实现了 FileSystem 接口
var _ FileSystem = (*SMBFS)(nil)
