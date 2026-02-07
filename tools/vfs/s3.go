package vfs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// S3FS S3 文件系统实现
type S3FS struct {
	client   *s3.Client
	bucket   string
	baseURL  string
	basePath string // S3 前缀（prefix）
	endpoint string // S3 endpoint 地址
	options  Options
	cache    *FileCache
	mu       sync.RWMutex
}

// NewS3FS 创建 S3 文件系统实例
// urlStr 格式: s3://accessKey:secretKey@endpoint/bucket/prefix
func NewS3FS(urlStr string, opts ...Options) (*S3FS, error) {
	// 解析 URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("无法解析 S3 URL: %w", err)
	}

	// 获取 endpoint（host 部分）
	endpoint := parsedURL.Hostname()
	if endpoint == "" {
		return nil, fmt.Errorf("S3 URL 必须包含 endpoint 地址")
	}
	portStr := parsedURL.Port()

	// 构建完整的 endpoint URL（S3 SDK 需要 HTTP/HTTPS URL）
	endpointURL := "https://" + endpoint
	if portStr != "" {
		endpointURL = "https://" + endpoint + ":" + portStr
	}

	// 获取认证信息
	accessKey := ""
	secretKey := ""
	if parsedURL.User != nil {
		accessKey = parsedURL.User.Username()
		secretKey, _ = parsedURL.User.Password()
	}

	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("S3 URL 必须包含 accessKey 和 secretKey（格式: s3://accessKey:secretKey@endpoint/bucket/prefix）")
	}

	// 解析路径部分：/bucket/prefix
	urlPath := strings.TrimPrefix(parsedURL.Path, "/")
	pathParts := strings.SplitN(urlPath, "/", 2)
	if len(pathParts) < 1 || pathParts[0] == "" {
		return nil, fmt.Errorf("S3 URL 必须包含存储桶名称")
	}

	bucket := pathParts[0]
	basePath := ""
	if len(pathParts) > 1 {
		basePath = pathParts[1]
	}
	// 移除尾部的 /
	basePath = strings.TrimSuffix(basePath, "/")

	// 获取配置选项
	options := DefaultOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	if options.Debug {
		logger.Infof(locale.GetString("log_s3_connecting"), endpointURL, bucket, basePath)
	}

	// 创建 S3 客户端
	client := s3.New(s3.Options{
		Region:       "auto", // 对于自定义 endpoint 使用 auto
		BaseEndpoint: &endpointURL,
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		UsePathStyle: true, // 使用路径风格访问（兼容 MinIO 等非 AWS S3）
	})

	// 测试连接：尝试列出存储桶内容
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(options.Timeout)*time.Second)
	defer cancel()

	prefix := basePath
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	_, err = client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    &bucket,
		Prefix:    &prefix,
		MaxKeys:   aws.Int32(1), // 只需验证可访问性
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, fmt.Errorf("无法访问 S3 存储桶 %s: %w", bucket, err)
	}

	s3fs := &S3FS{
		client:   client,
		bucket:   bucket,
		baseURL:  urlStr,
		basePath: basePath,
		endpoint: endpointURL,
		options:  options,
	}

	// 初始化缓存
	if options.CacheEnabled {
		s3fs.cache = NewFileCache(options.CacheDir, options.Debug)
	}

	if options.Debug {
		logger.Infof(locale.GetString("log_s3_filesystem_connected"), endpointURL, bucket+"/"+basePath)
	}

	return s3fs, nil
}

// resolvePath 将相对路径解析为 S3 对象的完整 key
// S3 的 key 不以 / 开头
func (s *S3FS) resolvePath(p string) string {
	// 清理路径
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	p = strings.TrimPrefix(p, "./")

	if s.basePath == "" || s.basePath == "." {
		if p == "" || p == "." {
			return ""
		}
		return p
	}

	basePath := strings.TrimSuffix(s.basePath, "/")

	// 如果路径已经以 basePath 开头，直接返回
	if strings.HasPrefix(p, basePath+"/") || p == basePath {
		return p
	}

	if p == "" || p == "." {
		return basePath
	}
	return basePath + "/" + p
}

// Open 打开文件用于读取
func (s *S3FS) Open(p string) (File, error) {
	key := s.resolvePath(p)

	// 检查缓存
	if s.cache != nil {
		if data, ok := s.cache.Get(key); ok {
			return newS3File(data, key, s), nil
		}
	}

	// 从 S3 下载
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.options.Timeout)*time.Second)
	defer cancel()

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("无法打开 S3 文件 %s: %w", key, err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("无法读取 S3 文件 %s: %w", key, err)
	}

	// 保存到缓存
	if s.cache != nil {
		s.cache.Set(key, data)
	}

	return newS3File(data, key, s), nil
}

// Stat 获取文件信息
func (s *S3FS) Stat(p string) (FileInfo, error) {
	key := s.resolvePath(p)

	// 先尝试作为文件获取
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.options.Timeout)*time.Second)
	defer cancel()

	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err == nil {
		// 文件存在
		modTime := time.Time{}
		if result.LastModified != nil {
			modTime = *result.LastModified
		}
		size := int64(0)
		if result.ContentLength != nil {
			size = *result.ContentLength
		}
		return NewFileInfo(path.Base(key), size, 0o644, modTime, false), nil
	}

	// 尝试作为"目录"查找（S3 中通过前缀模拟目录）
	prefix := key
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Duration(s.options.Timeout)*time.Second)
	defer cancel2()

	listResult, err := s.client.ListObjectsV2(ctx2, &s3.ListObjectsV2Input{
		Bucket:    &s.bucket,
		Prefix:    &prefix,
		MaxKeys:   aws.Int32(1),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, fmt.Errorf("无法获取 S3 文件信息 %s: %w", key, err)
	}

	if len(listResult.Contents) > 0 || len(listResult.CommonPrefixes) > 0 {
		// 作为目录存在
		name := path.Base(key)
		if name == "" || name == "." {
			name = s.bucket
		}
		return NewFileInfo(name, 0, fs.ModeDir|0o755, time.Time{}, true), nil
	}

	return nil, fmt.Errorf("S3 对象不存在: %s", key)
}

// ReadDir 读取目录（通过 S3 前缀列出对象和公共前缀）
func (s *S3FS) ReadDir(p string) ([]DirEntry, error) {
	key := s.resolvePath(p)

	prefix := key
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.options.Timeout)*time.Second)
	defer cancel()

	var entries []DirEntry
	delimiter := "/"

	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket:    &s.bucket,
		Prefix:    &prefix,
		Delimiter: &delimiter,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("无法读取 S3 目录 %s: %w", key, err)
		}

		// 处理"子目录"（CommonPrefixes）
		for _, cp := range page.CommonPrefixes {
			if cp.Prefix == nil {
				continue
			}
			dirName := strings.TrimSuffix(strings.TrimPrefix(*cp.Prefix, prefix), "/")
			if dirName == "" {
				continue
			}
			info := NewFileInfo(dirName, 0, fs.ModeDir|0o755, time.Time{}, true)
			entries = append(entries, NewDirEntry(info))
		}

		// 处理文件
		for _, obj := range page.Contents {
			if obj.Key == nil {
				continue
			}
			fileName := strings.TrimPrefix(*obj.Key, prefix)
			if fileName == "" {
				continue
			}
			// 跳过子目录中的文件（不应该出现在当前级别）
			if strings.Contains(fileName, "/") {
				continue
			}
			size := int64(0)
			if obj.Size != nil {
				size = *obj.Size
			}
			modTime := time.Time{}
			if obj.LastModified != nil {
				modTime = *obj.LastModified
			}
			info := NewFileInfo(fileName, size, 0o644, modTime, false)
			entries = append(entries, NewDirEntry(info))
		}
	}

	return entries, nil
}

// ReadFile 读取文件内容
func (s *S3FS) ReadFile(p string) ([]byte, error) {
	key := s.resolvePath(p)

	// 检查缓存
	if s.cache != nil {
		if data, ok := s.cache.Get(key); ok {
			return data, nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.options.Timeout)*time.Second)
	defer cancel()

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("无法打开 S3 文件 %s: %w", key, err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("无法读取 S3 文件 %s: %w", key, err)
	}

	// 保存到缓存
	if s.cache != nil {
		s.cache.Set(key, data)
	}

	return data, nil
}

// Type 返回后端类型
func (s *S3FS) Type() BackendType {
	return S3
}

// BaseURL 返回基础 URL
func (s *S3FS) BaseURL() string {
	return s.baseURL
}

// Close 关闭连接（S3 HTTP 客户端无状态，无需特殊处理）
func (s *S3FS) Close() error {
	if s.cache != nil {
		s.cache.Clear()
	}
	return nil
}

// IsRemote 返回是否为远程文件系统
func (s *S3FS) IsRemote() bool {
	return true
}

// JoinPath 连接路径（使用 / 分隔符）
func (s *S3FS) JoinPath(elem ...string) string {
	return path.Join(elem...)
}

// RelPath 计算相对路径
func (s *S3FS) RelPath(basePath, targetPath string) (string, error) {
	basePath = path.Clean(basePath)
	targetPath = path.Clean(targetPath)

	// 如果 basePath 是 "." 或 ""，直接返回 targetPath
	if basePath == "." || basePath == "" {
		rel := strings.TrimPrefix(targetPath, "./")
		rel = strings.TrimPrefix(rel, "/")
		if rel == "" {
			rel = "."
		}
		return rel, nil
	}

	basePathWithSlash := basePath
	if !strings.HasSuffix(basePathWithSlash, "/") {
		basePathWithSlash = basePathWithSlash + "/"
	}

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

	return "", fmt.Errorf("目标路径 %s 不在基础路径 %s 下", targetPath, basePath)
}

// Exists 检查路径是否存在
func (s *S3FS) Exists(p string) (bool, error) {
	key := s.resolvePath(p)

	// 尝试作为文件查找
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.options.Timeout)*time.Second)
	defer cancel()

	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err == nil {
		return true, nil
	}

	// 尝试作为"目录"查找
	prefix := key
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Duration(s.options.Timeout)*time.Second)
	defer cancel2()

	listResult, err := s.client.ListObjectsV2(ctx2, &s3.ListObjectsV2Input{
		Bucket:    &s.bucket,
		Prefix:    &prefix,
		MaxKeys:   aws.Int32(1),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return false, err
	}

	return len(listResult.Contents) > 0 || len(listResult.CommonPrefixes) > 0, nil
}

// IsDir 检查是否为目录（S3 中通过前缀模拟目录）
func (s *S3FS) IsDir(p string) (bool, error) {
	key := s.resolvePath(p)

	// 先检查是否是文件
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.options.Timeout)*time.Second)
	defer cancel()

	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err == nil {
		// 文件存在，不是目录
		return false, nil
	}

	// 检查是否有以此为前缀的对象（即为"目录"）
	prefix := key
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Duration(s.options.Timeout)*time.Second)
	defer cancel2()

	listResult, err := s.client.ListObjectsV2(ctx2, &s3.ListObjectsV2Input{
		Bucket:    &s.bucket,
		Prefix:    &prefix,
		MaxKeys:   aws.Int32(1),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return false, err
	}

	return len(listResult.Contents) > 0 || len(listResult.CommonPrefixes) > 0, nil
}

// OpenReaderAtSeeker 打开文件并返回支持 Seek 的 Reader
// 下载到内存并缓存，避免重复下载 两种策略：
// 小文件（<1MB）: 直接下载到内存+缓存
// 大文件: 也下载到内存+缓存（S3 Range 请求虽可行，但为保持简单先用全量下载+缓存，与 SMB 一致）
func (s *S3FS) OpenReaderAtSeeker(p string) (ReaderAtSeeker, error) {
	key := s.resolvePath(p)

	// 优先使用缓存
	if s.cache != nil {
		if data, ok := s.cache.Get(key); ok {
			return bytes.NewReader(data), nil
		}
	}

	// 下载文件到内存
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.options.Timeout)*time.Second)
	defer cancel()

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("无法打开 S3 文件 %s: %w", key, err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("无法读取 S3 文件 %s: %w", key, err)
	}

	// 保存到缓存，避免重复下载
	if s.cache != nil {
		s.cache.Set(key, data)
	}

	return bytes.NewReader(data), nil
}

// GetBasePath 返回 S3 基础路径（前缀）
func (s *S3FS) GetBasePath() string {
	return s.basePath
}

// --- 辅助类型 ---

// s3File 实现 File 接口（基于内存的 bytes.Reader 包装）
type s3File struct {
	data   []byte
	reader *bytes.Reader
	key    string
	sfs    *S3FS
}

func newS3File(data []byte, key string, sfs *S3FS) *s3File {
	return &s3File{
		data:   data,
		reader: bytes.NewReader(data),
		key:    key,
		sfs:    sfs,
	}
}

func (f *s3File) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

func (f *s3File) Close() error {
	f.data = nil
	f.reader = nil
	return nil
}

func (f *s3File) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

func (f *s3File) Stat() (FileInfo, error) {
	return f.sfs.Stat(f.key)
}

func (f *s3File) ReadAt(p []byte, off int64) (n int, err error) {
	return f.reader.ReadAt(p, off)
}

// 确保 S3FS 实现了 FileSystem 接口
var _ FileSystem = (*S3FS)(nil)
