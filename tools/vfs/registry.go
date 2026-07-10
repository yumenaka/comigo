package vfs

import (
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/yumenaka/comigo/tools"
)

// registry 全局文件系统注册表
var registry = &fsRegistry{
	instances: make(map[string]FileSystem),
}

// fsRegistry 文件系统注册表
type fsRegistry struct {
	mu        sync.RWMutex
	instances map[string]FileSystem // key: URL 身份哈希与不可变选项
}

// GetOrCreate 获取或创建文件系统实例
// 如果已存在则返回现有实例，否则创建新实例并注册
func GetOrCreate(storeURL string, opts ...Options) (FileSystem, error) {
	key := registryKey(storeURL, opts)
	// ponytail: 创建远程连接时使用全局锁；只有实测并行创建不同书库成为瓶颈时再换 singleflight。
	registry.mu.Lock()
	defer registry.mu.Unlock()

	// 检查是否已存在
	if fs, ok := registry.instances[key]; ok {
		return fs, nil
	}

	// 创建新实例
	fs, err := New(storeURL, opts...)
	if err != nil {
		return nil, err
	}

	// 注册实例
	registry.instances[key] = fs
	return fs, nil
}

// registryKey 让实例配置创建后保持不可变；哈希 URL 可区分凭据且不在 key 中泄露密码。
func registryKey(storeURL string, opts []Options) string {
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}
	identity := sha256.Sum256([]byte(storeURL))
	return fmt.Sprintf("%x|%t|%s|%d|%t|%t", identity, options.CacheEnabled, options.CacheDir, options.Timeout, options.Debug, options.UseRangeRequests)
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
