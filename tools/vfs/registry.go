package vfs

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

// registry 全局文件系统注册表
var registry = &fsRegistry{
	instances: make(map[string]FileSystem),
}

// fsRegistry 文件系统注册表
type fsRegistry struct {
	mu        sync.RWMutex
	instances map[string]FileSystem // key: 标准化后的 URL/路径
}

// Register 注册文件系统实例
func Register(key string, fs FileSystem) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.instances[key] = fs
}

// Unregister 注销文件系统实例
func Unregister(key string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if fs, ok := registry.instances[key]; ok {
		_ = fs.Close()
		delete(registry.instances, key)
	}
}

// Get 获取已注册的文件系统实例
func Get(key string) (FileSystem, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	fs, ok := registry.instances[key]
	return fs, ok
}

// GetOrCreate 获取或创建文件系统实例
// 如果已存在则返回现有实例，否则创建新实例并注册
func GetOrCreate(storeURL string, opts ...Options) (FileSystem, error) {
	// 标准化 URL 作为 key
	key := normalizeURL(storeURL)

	// 检查是否已存在
	if fs, ok := Get(key); ok {
		return fs, nil
	}

	// 创建新实例
	fs, err := New(storeURL, opts...)
	if err != nil {
		return nil, err
	}

	// 注册实例
	Register(key, fs)
	return fs, nil
}

// New 根据 URL 创建对应的文件系统实例
func New(storeURL string, opts ...Options) (FileSystem, error) {
	// 解析后端类型
	backendType, localPath := parseStoreURL(storeURL)

	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	switch backendType {
	case LocalDisk:
		return NewLocalFS(localPath, options)

	case WebDAV:
		return NewWebDAVFS(storeURL, options)

	case SMB:
		// TODO: 实现 SMB 支持
		return nil, fmt.Errorf("SMB 支持尚未实现")

	case SFTP:
		// TODO: 实现 SFTP 支持
		return nil, fmt.Errorf("SFTP 支持尚未实现")

	case S3:
		// TODO: 实现 S3 支持
		return nil, fmt.Errorf("S3 支持尚未实现")

	case FTP:
		// TODO: 实现 FTP 支持
		return nil, fmt.Errorf("FTP 支持尚未实现")

	default:
		return nil, fmt.Errorf("不支持的后端类型: %v", backendType)
	}
}

// parseStoreURL 解析存储 URL，返回后端类型和路径
func parseStoreURL(storeURL string) (BackendType, string) {
	// 检查是否为 file:// 协议
	if strings.HasPrefix(storeURL, "file://") {
		path := strings.TrimPrefix(storeURL, "file://")
		// 处理 Windows 路径格式 file:///C:/path
		if len(path) > 2 && path[0] == '/' && path[2] == ':' {
			path = path[1:]
		}
		return LocalDisk, path
	}

	// 检查是否为本地绝对路径
	if strings.HasPrefix(storeURL, "/") ||
		(len(storeURL) > 2 && storeURL[1] == ':' && (storeURL[2] == '\\' || storeURL[2] == '/')) ||
		strings.HasPrefix(storeURL, "\\\\") {
		return LocalDisk, storeURL
	}

	// 尝试解析为 URL
	u, err := url.Parse(storeURL)
	if err != nil {
		// 解析失败，当作本地路径
		return LocalDisk, storeURL
	}

	// 如果没有 scheme，当作本地路径
	if u.Scheme == "" {
		return LocalDisk, storeURL
	}

	// 根据 scheme 判断后端类型
	switch strings.ToLower(u.Scheme) {
	case "webdav", "dav", "davs", "http", "https":
		return WebDAV, storeURL
	case "smb":
		return SMB, storeURL
	case "sftp":
		return SFTP, storeURL
	case "s3":
		return S3, storeURL
	case "ftp", "ftps":
		return FTP, storeURL
	default:
		// 未知协议，当作本地路径
		return LocalDisk, storeURL
	}
}

// normalizeURL 标准化 URL
func normalizeURL(storeURL string) string {
	// 尝试解析为 URL
	u, err := url.Parse(storeURL)
	if err != nil {
		// 可能是本地路径，直接返回
		return storeURL
	}

	// 如果没有 scheme，可能是本地路径
	if u.Scheme == "" {
		return storeURL
	}

	// 标准化 scheme
	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case "webdav", "dav":
		scheme = "http"
	case "davs":
		scheme = "https"
	}

	// 重建 URL（不包含认证信息，用于作为 key）
	result := fmt.Sprintf("%s://%s%s", scheme, u.Host, u.Path)
	return result
}

// CloseAll 关闭所有注册的文件系统
func CloseAll() {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	for key, fs := range registry.instances {
		_ = fs.Close()
		delete(registry.instances, key)
	}
}

// List 列出所有注册的文件系统
func List() []string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	keys := make([]string, 0, len(registry.instances))
	for key := range registry.instances {
		keys = append(keys, key)
	}
	return keys
}

// IsRemoteURL 判断 URL 是否为远程存储
func IsRemoteURL(storeURL string) bool {
	backendType, _ := parseStoreURL(storeURL)
	return backendType != LocalDisk
}

// GetBackendType 获取 URL 对应的后端类型
func GetBackendType(storeURL string) BackendType {
	backendType, _ := parseStoreURL(storeURL)
	return backendType
}
