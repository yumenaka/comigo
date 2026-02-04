package vfs

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/yumenaka/comigo/tools/logger"
)

// FileCache 文件缓存管理器
// 用于缓存远程文件系统读取的文件，减少网络请求
type FileCache struct {
	cacheDir string
	debug    bool
	memory   sync.Map // 内存缓存，key: path, value: *cacheEntry
}

// cacheEntry 缓存条目
type cacheEntry struct {
	data      []byte
	timestamp time.Time
}

// NewFileCache 创建文件缓存管理器
func NewFileCache(cacheDir string, debug bool) *FileCache {
	// 确保缓存目录存在
	if cacheDir != "" {
		if err := os.MkdirAll(cacheDir, 0o755); err != nil {
			logger.Infof("创建缓存目录失败: %v", err)
		}
	}

	return &FileCache{
		cacheDir: cacheDir,
		debug:    debug,
	}
}

// hashPath 计算路径的哈希值，用于生成缓存文件名
func hashPath(path string) string {
	hash := md5.Sum([]byte(path))
	return hex.EncodeToString(hash[:])
}

// getCacheFilePath 获取缓存文件的完整路径
func (c *FileCache) getCacheFilePath(remotePath string) string {
	if c.cacheDir == "" {
		return ""
	}
	hash := hashPath(remotePath)
	return filepath.Join(c.cacheDir, hash)
}

// Get 从缓存获取文件内容
// 优先从内存缓存获取，然后从磁盘缓存获取
func (c *FileCache) Get(path string) ([]byte, bool) {
	// 先检查内存缓存
	if entry, ok := c.memory.Load(path); ok {
		ce := entry.(*cacheEntry)
		if c.debug {
			logger.Infof("从内存缓存命中: %s", path)
		}
		return ce.data, true
	}

	// 检查磁盘缓存
	if c.cacheDir != "" {
		cachePath := c.getCacheFilePath(path)
		if data, err := os.ReadFile(cachePath); err == nil {
			// 加载到内存缓存
			c.memory.Store(path, &cacheEntry{
				data:      data,
				timestamp: time.Now(),
			})
			if c.debug {
				logger.Infof("从磁盘缓存命中: %s", path)
			}
			return data, true
		}
	}

	return nil, false
}

// Set 设置缓存
func (c *FileCache) Set(path string, data []byte) {
	// 保存到内存缓存
	c.memory.Store(path, &cacheEntry{
		data:      data,
		timestamp: time.Now(),
	})

	// 保存到磁盘缓存
	if c.cacheDir != "" {
		cachePath := c.getCacheFilePath(path)
		if err := os.WriteFile(cachePath, data, 0o644); err != nil {
			if c.debug {
				logger.Infof("写入磁盘缓存失败: %v", err)
			}
		} else if c.debug {
			logger.Infof("已缓存到磁盘: %s -> %s", path, cachePath)
		}
	}
}

// Delete 删除缓存
func (c *FileCache) Delete(path string) {
	// 从内存缓存删除
	c.memory.Delete(path)

	// 从磁盘缓存删除
	if c.cacheDir != "" {
		cachePath := c.getCacheFilePath(path)
		_ = os.Remove(cachePath)
	}
}

// Clear 清空所有缓存
func (c *FileCache) Clear() {
	// 清空内存缓存
	c.memory.Range(func(key, value interface{}) bool {
		c.memory.Delete(key)
		return true
	})

	// 清空磁盘缓存目录
	if c.cacheDir != "" {
		entries, err := os.ReadDir(c.cacheDir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				_ = os.Remove(filepath.Join(c.cacheDir, entry.Name()))
			}
		}
	}
}

// ClearOlderThan 清除超过指定时间的缓存
func (c *FileCache) ClearOlderThan(duration time.Duration) {
	cutoff := time.Now().Add(-duration)

	// 清理内存缓存
	c.memory.Range(func(key, value interface{}) bool {
		entry := value.(*cacheEntry)
		if entry.timestamp.Before(cutoff) {
			c.memory.Delete(key)
		}
		return true
	})

	// 清理磁盘缓存
	if c.cacheDir != "" {
		entries, err := os.ReadDir(c.cacheDir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				_ = os.Remove(filepath.Join(c.cacheDir, entry.Name()))
			}
		}
	}
}

// Size 返回缓存条目数量
func (c *FileCache) Size() int {
	count := 0
	c.memory.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
