// Package vfs 提供虚拟文件系统抽象层，支持本地文件系统和远程存储（WebDAV、SFTP等）
package vfs

import (
	"io"
	"io/fs"
	"time"
)

// BackendType 文件存储后端类型
type BackendType int

const (
	LocalDisk BackendType = 1 + iota
	SMB
	SFTP
	WebDAV
	S3
	FTP
)

func (f BackendType) String() string {
	switch f {
	case LocalDisk:
		return "Local Disk"
	case SMB:
		return "SMB Share"
	case SFTP:
		return "SFTP Server"
	case WebDAV:
		return "WebDAV Server"
	case S3:
		return "S3 Storage"
	case FTP:
		return "FTP Server"
	default:
		return "Unknown Backend Type"
	}
}

// FileSystem 虚拟文件系统接口
// 提供统一的文件系统操作，支持本地和远程存储
type FileSystem interface {
	// Open 打开文件用于读取
	Open(path string) (File, error)

	// Stat 获取文件或目录的信息
	Stat(path string) (FileInfo, error)

	// ReadDir 读取目录内容
	ReadDir(path string) ([]DirEntry, error)

	// ReadFile 读取整个文件内容
	ReadFile(path string) ([]byte, error)

	// Type 返回文件系统后端类型
	Type() BackendType

	// BaseURL 返回文件系统的基础URL或路径
	BaseURL() string

	// Close 关闭文件系统连接（对于远程文件系统）
	Close() error

	// IsRemote 返回是否为远程文件系统
	IsRemote() bool

	// JoinPath 连接路径，根据文件系统类型使用正确的分隔符
	JoinPath(elem ...string) string

	// RelPath 计算相对路径
	RelPath(basePath, targetPath string) (string, error)

	// Exists 检查路径是否存在
	Exists(path string) (bool, error)

	// IsDir 检查路径是否为目录
	IsDir(path string) (bool, error)

	// OpenReaderAtSeeker 打开文件并返回支持 Seek 的 Reader
	// 这对于需要随机访问的压缩包格式很重要
	// 如果文件系统不支持 Seek，会将数据读取到内存中
	// 返回的 Reader 必须实现 io.Reader, io.ReaderAt, io.Seeker
	OpenReaderAtSeeker(path string) (ReaderAtSeeker, error)
}

// ReaderAtSeeker 接口，组合了 io.Reader, io.ReaderAt, io.Seeker
// 用于支持需要随机访问的压缩包格式
type ReaderAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

// File 文件接口，提供基本的文件读取操作
type File interface {
	io.Reader
	io.Closer
	io.Seeker

	// Stat 获取文件信息
	Stat() (FileInfo, error)

	// ReadAt 在指定偏移量处读取
	ReadAt(p []byte, off int64) (n int, err error)
}

// FileInfo 文件信息接口，兼容 fs.FileInfo
type FileInfo interface {
	fs.FileInfo
}

// DirEntry 目录项接口，兼容 fs.DirEntry
type DirEntry interface {
	fs.DirEntry
}

// fileInfo 是 FileInfo 的基础实现
type fileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
}

func (fi *fileInfo) Name() string       { return fi.name }
func (fi *fileInfo) Size() int64        { return fi.size }
func (fi *fileInfo) Mode() fs.FileMode  { return fi.mode }
func (fi *fileInfo) ModTime() time.Time { return fi.modTime }
func (fi *fileInfo) IsDir() bool        { return fi.isDir }
func (fi *fileInfo) Sys() any           { return nil }

// NewFileInfo 创建新的 FileInfo
func NewFileInfo(name string, size int64, mode fs.FileMode, modTime time.Time, isDir bool) FileInfo {
	return &fileInfo{
		name:    name,
		size:    size,
		mode:    mode,
		modTime: modTime,
		isDir:   isDir,
	}
}

// dirEntry 是 DirEntry 的基础实现
type dirEntry struct {
	info FileInfo
}

func (de *dirEntry) Name() string               { return de.info.Name() }
func (de *dirEntry) IsDir() bool                { return de.info.IsDir() }
func (de *dirEntry) Type() fs.FileMode          { return de.info.Mode().Type() }
func (de *dirEntry) Info() (fs.FileInfo, error) { return de.info, nil }

// NewDirEntry 创建新的 DirEntry
func NewDirEntry(info FileInfo) DirEntry {
	return &dirEntry{info: info}
}

// Options 文件系统配置选项
type Options struct {
	// CacheEnabled 是否启用文件缓存
	CacheEnabled bool

	// CacheDir 缓存目录路径
	CacheDir string

	// Timeout 连接超时（秒）
	Timeout int

	// Debug 是否启用调试日志
	Debug bool
}

// DefaultOptions 返回默认配置选项
func DefaultOptions() Options {
	return Options{
		CacheEnabled: false,
		CacheDir:     "",
		Timeout:      30,
		Debug:        false,
	}
}
