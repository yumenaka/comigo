package vfs

import (
	"fmt"
	"sync"
	"time"

	"github.com/yumenaka/comigo/tools"
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
	key := tools.NormalizeStoreURLKey(storeURL)

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

				enableCache(&wfs.options, &wfs.cache, desired)
			} else if sfs, ok := fs.(*SFTPFS); ok {
				// SFTPFS 配置升级：更新超时和缓存设置
				// 注意：SSH 连接的超时在创建时设置，无法动态更新
				// 但可以更新缓存设置
				enableCache(&sfs.options, &sfs.cache, desired)
			} else if smbfs, ok := fs.(*SMBFS); ok {
				// SMBFS 配置升级：更新缓存设置
				// 注意：SMB 连接的超时在创建时设置，无法动态更新
				// 但可以更新缓存设置
				enableCache(&smbfs.options, &smbfs.cache, desired)
			} else if ftpfs, ok := fs.(*FTPFS); ok {
				// FTPFS 配置升级：更新缓存设置
				enableCache(&ftpfs.options, &ftpfs.cache, desired)
			} else if s3fs, ok := fs.(*S3FS); ok {
				// S3FS 配置升级：更新缓存设置
				enableCache(&s3fs.options, &s3fs.cache, desired)
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

// enableCache 按调用方传入的最新配置补开缓存；CacheDir 为空时只启用内存缓存。
func enableCache(options *Options, cache **FileCache, desired Options) {
	if !desired.CacheEnabled || *cache != nil {
		return
	}
	options.CacheEnabled = true
	if desired.CacheDir != "" {
		options.CacheDir = desired.CacheDir
	}
	options.Debug = options.Debug || desired.Debug
	*cache = NewFileCache(options.CacheDir, options.Debug)
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
	case tools.StoreBackendLocalDisk:
		return NewLocalFS(localPath, options)

	case tools.StoreBackendWebDAV:
		return NewWebDAVFS(storeURL, options)

	case tools.StoreBackendSMB:
		return NewSMBFS(storeURL, options)

	case tools.StoreBackendSFTP:
		return NewSFTPFS(storeURL, options)

	case tools.StoreBackendS3:
		return NewS3FS(storeURL, options)

	case tools.StoreBackendFTP:
		return NewFTPFS(storeURL, options)

	default:
		return nil, fmt.Errorf("不支持的后端类型: %v", backendType)
	}
}

// parseStoreURL 解析存储 URL，返回后端类型和路径
func parseStoreURL(storeURL string) (tools.StoreBackendType, string) {
	info := tools.DetectStoreURL(storeURL)
	if info.Type == tools.StoreBackendLocalDisk {
		if info.LocalPath != "" {
			return tools.StoreBackendLocalDisk, info.LocalPath
		}
		return tools.StoreBackendLocalDisk, info.URL
	}
	return info.Type, info.URL
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
