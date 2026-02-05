package vfs

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
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
		// 如果调用方传入了 options，且希望启用缓存/更新超时，则尝试“升级”已有实例配置
		// 典型场景：某处曾用 CacheEnabled=false 创建了 WebDAVFS，后续读取图片希望 CacheEnabled=true 来复用下载结果
		if len(opts) > 0 {
			desired := opts[0]
			// WebDAVFS 和 SFTPFS 支持 FileCache 和配置升级
			if wfs, ok := fs.(*WebDAVFS); ok {
				// 升级超时（以最新传入为准，避免旧实例超时过短/过长）
				if desired.Timeout > 0 && desired.Timeout != wfs.options.Timeout {
					wfs.options.Timeout = desired.Timeout
					wfs.client.SetTimeout(time.Duration(desired.Timeout) * time.Second)
				}

				// 升级缓存：如果期望启用缓存且当前未初始化缓存，则初始化（内存缓存始终可用）
				if desired.CacheEnabled && wfs.cache == nil {
					wfs.options.CacheEnabled = true
					// 仅在调用方显式传入 CacheDir 时才启用磁盘缓存；否则只用内存缓存，避免落盘
					if desired.CacheDir != "" {
						wfs.options.CacheDir = desired.CacheDir
					}
					// Debug 取并集：任一处开启 debug 都可以输出缓存命中日志
					wfs.options.Debug = wfs.options.Debug || desired.Debug
					wfs.cache = NewFileCache(wfs.options.CacheDir, wfs.options.Debug)
				}
			} else if sfs, ok := fs.(*SFTPFS); ok {
				// SFTPFS 配置升级：更新超时和缓存设置
				// 注意：SSH 连接的超时在创建时设置，无法动态更新
				// 但可以更新缓存设置
				if desired.CacheEnabled && sfs.cache == nil {
					sfs.options.CacheEnabled = true
					if desired.CacheDir != "" {
						sfs.options.CacheDir = desired.CacheDir
					}
					sfs.options.Debug = sfs.options.Debug || desired.Debug
					sfs.cache = NewFileCache(sfs.options.CacheDir, sfs.options.Debug)
				}
			}
		}
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
		return NewSFTPFS(storeURL, options)

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
	case "sftp":
		// SFTP scheme 保持不变
		scheme = "sftp"
	}

	// 重建 URL（不包含认证信息，用于作为 key）
	// 对于 SFTP，包含端口号以确保唯一性
	if scheme == "sftp" && u.Port() != "" {
		result := fmt.Sprintf("%s://%s:%s%s", scheme, u.Hostname(), u.Port(), u.Path)
		return result
	}
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
