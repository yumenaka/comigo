package vfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// LocalFS 本地文件系统实现
type LocalFS struct {
	basePath string
	options  Options
}

// NewLocalFS 创建本地文件系统实例
func NewLocalFS(basePath string, opts ...Options) (*LocalFS, error) {
	// 转换为绝对路径
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("无法获取绝对路径: %w", err)
	}

	// 清理路径
	absPath = filepath.Clean(absPath)

	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	return &LocalFS{
		basePath: absPath,
		options:  options,
	}, nil
}

// resolvePath 将相对路径解析为完整路径
func (l *LocalFS) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(l.basePath, path)
}

// Open 打开文件
func (l *LocalFS) Open(path string) (File, error) {
	fullPath := l.resolvePath(path)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	return &localFile{file: file}, nil
}

// Stat 获取文件信息
func (l *LocalFS) Stat(path string) (FileInfo, error) {
	fullPath := l.resolvePath(path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// ReadDir 读取目录
func (l *LocalFS) ReadDir(path string) ([]DirEntry, error) {
	fullPath := l.resolvePath(path)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	result := make([]DirEntry, len(entries))
	for i, entry := range entries {
		result[i] = entry
	}
	return result, nil
}

// ReadFile 读取文件内容
func (l *LocalFS) ReadFile(path string) ([]byte, error) {
	fullPath := l.resolvePath(path)
	return os.ReadFile(fullPath)
}

// Type 返回后端类型
func (l *LocalFS) Type() BackendType {
	return LocalDisk
}

// BaseURL 返回基础路径
func (l *LocalFS) BaseURL() string {
	return l.basePath
}

// Close 关闭文件系统（本地文件系统无需关闭）
func (l *LocalFS) Close() error {
	return nil
}

// IsRemote 返回是否为远程文件系统
func (l *LocalFS) IsRemote() bool {
	return false
}

// JoinPath 连接路径
func (l *LocalFS) JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// RelPath 计算相对路径
func (l *LocalFS) RelPath(basePath, targetPath string) (string, error) {
	return filepath.Rel(basePath, targetPath)
}

// Exists 检查路径是否存在
func (l *LocalFS) Exists(path string) (bool, error) {
	fullPath := l.resolvePath(path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsDir 检查是否为目录
func (l *LocalFS) IsDir(path string) (bool, error) {
	fullPath := l.resolvePath(path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// OpenReaderAtSeeker 打开文件并返回支持 Seek 的 Reader
func (l *LocalFS) OpenReaderAtSeeker(path string) (ReaderAtSeeker, error) {
	fullPath := l.resolvePath(path)
	return os.Open(fullPath)
}

// localFile 本地文件实现
type localFile struct {
	file *os.File
}

func (f *localFile) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *localFile) Close() error {
	return f.file.Close()
}

func (f *localFile) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *localFile) Stat() (FileInfo, error) {
	return f.file.Stat()
}

func (f *localFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.file.ReadAt(p, off)
}

// 确保 LocalFS 实现了 FileSystem 接口
var _ FileSystem = (*LocalFS)(nil)

// 确保 os.FileInfo 实现了 FileInfo 接口
var _ FileInfo = (fs.FileInfo)(nil)

// 确保 os.DirEntry 实现了 DirEntry 接口
var _ DirEntry = (fs.DirEntry)(nil)
