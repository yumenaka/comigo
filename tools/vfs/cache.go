package vfs

import (
	"sync"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// FileCache 缓存远程文件内容，减少重复下载。
type FileCache struct {
	debug  bool
	memory sync.Map // key: path, value: []byte
}

// NewFileCache 创建内存文件缓存。
func NewFileCache(debug bool) *FileCache {
	return &FileCache{debug: debug}
}

// Get 读取缓存，未命中时由调用方下载。
func (c *FileCache) Get(path string) ([]byte, bool) {
	data, ok := c.memory.Load(path)
	if !ok {
		return nil, false
	}
	if c.debug {
		logger.Infof(locale.GetString("log_cache_hit_memory"), path)
	}
	return data.([]byte), true
}

// Set 保存文件内容。
func (c *FileCache) Set(path string, data []byte) {
	c.memory.Store(path, data)
}

// Clear 释放所有缓存内容。
func (c *FileCache) Clear() {
	c.memory.Clear()
}
