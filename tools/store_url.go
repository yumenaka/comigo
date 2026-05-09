package tools

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

// StoreBackendType 是跨 config/store/vfs 共用的书库后端类型。
type StoreBackendType int

const (
	StoreBackendLocalDisk StoreBackendType = iota
	StoreBackendSMB
	StoreBackendSFTP
	StoreBackendWebDAV
	StoreBackendS3
	StoreBackendFTP
)

// StoreURLInfo 保存解析后的 store URL 信息，避免各包重复解析同一套规则。
type StoreURLInfo struct {
	URL            string
	Scheme         string
	Type           StoreBackendType
	LocalPath      string
	ServerHost     string
	ServerPort     int
	NeedAuth       bool
	AuthUsername   string
	AuthPassword   string
	RemotePath     string
	SMBShareName   string
	SMBPath        string
	S3BucketName   string
	S3ObjectPrefix string
}

// ParseStoreURL 解析本地路径和远程 store URL，并对已知远程协议做严格校验。
func ParseStoreURL(storeURL string) (StoreURLInfo, error) {
	if storeURL == "" {
		return StoreURLInfo{}, fmt.Errorf("URL cannot be empty")
	}

	if strings.HasPrefix(strings.ToLower(storeURL), "file://") {
		localPath := parseFileStoreURLPath(storeURL)
		return StoreURLInfo{
			URL:       localPath,
			Scheme:    "file",
			Type:      StoreBackendLocalDisk,
			LocalPath: localPath,
		}, nil
	}

	if IsLocalStorePath(storeURL) {
		return StoreURLInfo{
			URL:       storeURL,
			Type:      StoreBackendLocalDisk,
			LocalPath: storeURL,
		}, nil
	}

	parsedURL, err := url.Parse(storeURL)
	if err != nil {
		return StoreURLInfo{}, fmt.Errorf("invalid URL: %w", err)
	}

	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme == "" {
		return StoreURLInfo{
			URL:       storeURL,
			Type:      StoreBackendLocalDisk,
			LocalPath: storeURL,
		}, nil
	}

	backendType, ok := remoteStoreBackendTypeFromScheme(scheme)
	if !ok {
		return StoreURLInfo{}, fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}
	switch backendType {
	case StoreBackendSMB:
		return parseSMBStoreURL(parsedURL)
	case StoreBackendSFTP:
		return parseSFTPStoreURL(parsedURL)
	case StoreBackendWebDAV:
		return parseWebDAVStoreURL(parsedURL)
	case StoreBackendFTP:
		return parseFTPStoreURL(parsedURL)
	case StoreBackendS3:
		return parseS3StoreURL(parsedURL)
	default:
		return StoreURLInfo{}, fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}
}

// DetectStoreURL 用宽松规则识别 store URL；未知或无效 URL 按本地路径处理。
func DetectStoreURL(storeURL string) StoreURLInfo {
	info, err := ParseStoreURL(storeURL)
	if err == nil {
		return info
	}
	if parsedURL, parseErr := url.Parse(storeURL); parseErr == nil {
		if backendType, ok := remoteStoreBackendTypeFromScheme(parsedURL.Scheme); ok {
			return StoreURLInfo{
				URL:  storeURL,
				Type: backendType,
			}
		}
	}
	return StoreURLInfo{
		URL:       storeURL,
		Type:      StoreBackendLocalDisk,
		LocalPath: storeURL,
	}
}

// IsLocalStorePath 判断字符串是否已经是操作系统本地路径。
func IsLocalStorePath(storeURL string) bool {
	return filepath.IsAbs(storeURL) ||
		(len(storeURL) > 2 && storeURL[1] == ':' && (storeURL[2] == '\\' || storeURL[2] == '/')) ||
		strings.HasPrefix(storeURL, `\\`)
}

// IsRemoteStoreURL 判断 storeURL 是否为受支持的远程书库 URL。
func IsRemoteStoreURL(storeURL string) bool {
	if IsLocalStorePath(storeURL) || strings.HasPrefix(strings.ToLower(storeURL), "file://") {
		return false
	}
	parsedURL, err := url.Parse(storeURL)
	if err != nil {
		return false
	}
	_, ok := remoteStoreBackendTypeFromScheme(parsedURL.Scheme)
	return ok
}

// StoreURLHost 返回远程 URL 的主机名，解析失败或本地路径时返回空字符串。
func StoreURLHost(storeURL string) string {
	parsedURL, err := url.Parse(storeURL)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}

// NormalizeStoreURLKey 生成 VFS 缓存 key，去掉认证信息并规范 WebDAV 等价协议。
func NormalizeStoreURLKey(storeURL string) string {
	if strings.HasPrefix(strings.ToLower(storeURL), "file://") {
		return parseFileStoreURLPath(storeURL)
	}
	if IsLocalStorePath(storeURL) {
		return storeURL
	}

	parsedURL, err := url.Parse(storeURL)
	if err != nil || parsedURL.Scheme == "" {
		return storeURL
	}

	scheme := strings.ToLower(parsedURL.Scheme)
	switch scheme {
	case "webdav", "dav":
		scheme = "http"
	case "davs":
		scheme = "https"
	}

	host := parsedURL.Host
	if parsedURL.Port() != "" {
		switch scheme {
		case "sftp", "smb", "ftp", "ftps":
			host = fmt.Sprintf("%s:%s", parsedURL.Hostname(), parsedURL.Port())
		}
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, parsedURL.Path)
}

// NormalizeStoreURLForCompare 统一书库路径重叠检查的比较值。
func NormalizeStoreURLForCompare(storeURL string) (normalized string, remote bool, err error) {
	if IsRemoteStoreURL(storeURL) {
		return storeURL, true, nil
	}

	normalized, err = NormalizeAbsPath(storeURL)
	if err != nil {
		return "", false, err
	}
	return normalized, false, nil
}

// IsSubPath 判断 child 是否位于 parent 目录内部。
func IsSubPath(parent, child string) bool {
	rel, err := filepath.Rel(filepath.Clean(parent), filepath.Clean(child))
	if err != nil {
		return false
	}
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

func parseFileStoreURLPath(fileURL string) string {
	path := fileURL
	if strings.HasPrefix(strings.ToLower(fileURL), "file://") {
		path = fileURL[len("file://"):]
	}
	if len(path) > 0 && path[0] == '/' {
		if len(path) > 1 && path[1] == '/' {
			return path
		}
		if len(path) > 2 && path[1] != '/' && path[2] == ':' {
			return path[1:]
		}
	}
	return path
}

func remoteStoreBackendTypeFromScheme(scheme string) (StoreBackendType, bool) {
	switch strings.ToLower(scheme) {
	case "webdav", "dav", "davs", "http", "https":
		return StoreBackendWebDAV, true
	case "smb":
		return StoreBackendSMB, true
	case "sftp":
		return StoreBackendSFTP, true
	case "s3":
		return StoreBackendS3, true
	case "ftp", "ftps":
		return StoreBackendFTP, true
	default:
		return StoreBackendLocalDisk, false
	}
}

func parseSMBStoreURL(parsedURL *url.URL) (StoreURLInfo, error) {
	if parsedURL.Host == "" {
		return StoreURLInfo{}, fmt.Errorf("SMB URL requires a host")
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		return StoreURLInfo{}, fmt.Errorf("SMB URL requires a share name")
	}

	info := StoreURLInfo{
		URL:          parsedURL.String(),
		Scheme:       strings.ToLower(parsedURL.Scheme),
		Type:         StoreBackendSMB,
		ServerHost:   parsedURL.Hostname(),
		ServerPort:   parseStoreURLPort(parsedURL.Port(), 445),
		SMBShareName: pathParts[0],
	}
	if len(pathParts) > 1 {
		info.SMBPath = strings.Join(pathParts[1:], "/")
	}
	applyStoreURLAuth(&info, parsedURL)
	if strings.Contains(info.AuthUsername, ";") {
		parts := strings.SplitN(info.AuthUsername, ";", 2)
		info.AuthUsername = parts[1]
	}
	return info, nil
}

func parseSFTPStoreURL(parsedURL *url.URL) (StoreURLInfo, error) {
	if parsedURL.Host == "" {
		return StoreURLInfo{}, fmt.Errorf("SFTP URL requires a host")
	}

	info := StoreURLInfo{
		URL:          parsedURL.String(),
		Scheme:       strings.ToLower(parsedURL.Scheme),
		Type:         StoreBackendSFTP,
		ServerHost:   parsedURL.Hostname(),
		ServerPort:   parseStoreURLPort(parsedURL.Port(), 22),
		RemotePath:   strings.Trim(parsedURL.Path, "/"),
		NeedAuth:     parsedURL.User != nil,
		AuthUsername: "",
	}
	applyStoreURLAuth(&info, parsedURL)
	return info, nil
}

func parseWebDAVStoreURL(parsedURL *url.URL) (StoreURLInfo, error) {
	if parsedURL.Host == "" {
		return StoreURLInfo{}, fmt.Errorf("WebDAV URL requires a host")
	}

	defaultPort := 80
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme == "https" || scheme == "davs" {
		defaultPort = 443
	}

	info := StoreURLInfo{
		URL:        parsedURL.String(),
		Scheme:     scheme,
		Type:       StoreBackendWebDAV,
		ServerHost: parsedURL.Hostname(),
		ServerPort: parseStoreURLPort(parsedURL.Port(), defaultPort),
		RemotePath: strings.Trim(parsedURL.Path, "/"),
	}
	applyStoreURLAuth(&info, parsedURL)
	return info, nil
}

func parseFTPStoreURL(parsedURL *url.URL) (StoreURLInfo, error) {
	if parsedURL.Host == "" {
		return StoreURLInfo{}, fmt.Errorf("FTP URL requires a host")
	}

	defaultPort := 21
	if strings.ToLower(parsedURL.Scheme) == "ftps" {
		defaultPort = 990
	}

	info := StoreURLInfo{
		URL:        parsedURL.String(),
		Scheme:     strings.ToLower(parsedURL.Scheme),
		Type:       StoreBackendFTP,
		ServerHost: parsedURL.Hostname(),
		ServerPort: parseStoreURLPort(parsedURL.Port(), defaultPort),
		RemotePath: strings.Trim(parsedURL.Path, "/"),
	}
	applyStoreURLAuth(&info, parsedURL)
	return info, nil
}

func parseS3StoreURL(parsedURL *url.URL) (StoreURLInfo, error) {
	if parsedURL.Host == "" {
		return StoreURLInfo{}, fmt.Errorf("S3 URL requires an endpoint")
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		return StoreURLInfo{}, fmt.Errorf("S3 URL requires a bucket name")
	}

	info := StoreURLInfo{
		URL:          parsedURL.String(),
		Scheme:       strings.ToLower(parsedURL.Scheme),
		Type:         StoreBackendS3,
		ServerHost:   parsedURL.Hostname(),
		ServerPort:   parseStoreURLPort(parsedURL.Port(), 443),
		S3BucketName: pathParts[0],
	}
	if len(pathParts) > 1 {
		info.S3ObjectPrefix = strings.Join(pathParts[1:], "/")
	}
	applyStoreURLAuth(&info, parsedURL)
	return info, nil
}

func applyStoreURLAuth(info *StoreURLInfo, parsedURL *url.URL) {
	if parsedURL.User == nil {
		return
	}
	info.NeedAuth = true
	info.AuthUsername = parsedURL.User.Username()
	info.AuthPassword, _ = parsedURL.User.Password()
}

func parseStoreURLPort(port string, defaultPort int) int {
	if port == "" {
		return defaultPort
	}
	if parsedPort, err := strconv.Atoi(port); err == nil {
		return parsedPort
	}
	return defaultPort
}
