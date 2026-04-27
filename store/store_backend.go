package store

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Backend 文件后端。 可能是 1: 本地文件系统，2: SMB共享，3: SFTP服务器，4: WebDAV服务器，5: S3存储，6: FTP服务器
// 可能用到的字段都放进一个结构体；调用方只填需要字段
// 核心配置是“URL”，其余大多数参数都解析此URL生成。
// 本地书库为文件路径（/home/pi/books C:\Users\用户名\books或  file://some_path/books），
// 其他类型的书库是对应文件服务的 url 形式，
// 如 smb://workgroup;user:password@server/share/folder/books
// 或 sftp://<user>:<password>@<host>/<path>
// 或 webdav://192.168.1.100/books  <scheme>://[<username>[:<password>]@]<host>[:<port>]/<base-path>/<resource-path> scheme 为 http 或 https（也接受 dav:// 与 davs:// 作为同义写法）
// 或 ftp://<user>:<password>@<host>:<port>/<dir1>  ftps://<user>:<password>@<host>:<port>/<dir1>
// 或 s3://<S3_endpoint>[:<port>]/<bucket_name>/[<S3_prefix>] [region=<S3_region>] [config=<config_file_location> | config_server=<url>] [section=<section_name>]
type Backend struct {
	URL          string      `json:"url"`            // 本地书库为文件路径（/home/pi/books C:\Users\用户名\books或  file://some_path/books），其他类型的书库是对应文件服务的 url 形式，
	Type         BackendType `json:"backend_type"`   // 文件后端类型 1: 本地文件夹，2: SMB共享，3: SFTP服务器，4: WebDAV服务器，5: S3存储，6: FTP服务器
	ServerHost   string      `json:"server_host"`    // 文件服务的服务器地址，smb、ftp、sftp、webdav等类型的书库需要填写。
	ServerPort   int         `json:"server_port"`    // 文件服务的端口号，smb、ftp、sftp、webdav等类型的书库需要填写。
	NeedAuth     bool        `json:"need_auth"`      // 相关文件服务是否需要认证，smb、ftp、sftp、webdav等类型的书库需要填写。
	AuthUsername string      `json:"auth_username"`  // 认证用户名，smb、ftp、sftp、webdav等类型的书库需要填写。
	AuthPassword string      `json:"auth_password"`  // 认证密码，smb、ftp、sftp、webdav等类型的书库需要填写。
	SMBShareName string      `json:"smb_share_name"` // SMB共享名称，smb类型的书库需要填写
	SMBPath      string      `json:"smb_path"`       // SMB共享路径，smb类型的书库需要填写
}

// ParseStoreURL 解析URL字符串并返回FileBackend配置
// 支持的URL格式：
// - 本地文件: file:///path/to/books 或 /path/to/books (Unix路径) 或 C:\path\to\books  D:\path\to\books E:\path\to\books (Windows路径)
// - SMB: smb://workgroup;user:password@server/share/folder/books
// - SFTP: sftp://user:password@host/path
// - WebDAV: webdav://host/path 或 http://host/path 或 https://host/path
// - FTP: ftp://user:password@host:port/dir 或 ftps://user:password@host:port/dir
// - S3: s3://endpoint/bucket/prefix
func (backend *Backend) ParseStoreURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	// 创建一个新的Backend实例
	backend.URL = urlStr
	// 检查是否为本地文件路径（无scheme或file://scheme）
	if strings.HasPrefix(urlStr, "file://") {
		return backend.ParseFileURL(urlStr)
	}

	// 检查是否为绝对路径（本地文件）
	if strings.HasPrefix(urlStr, "/") ||
		(len(urlStr) > 2 && urlStr[1] == ':' && (urlStr[2] == '\\' || urlStr[2] == '/')) ||
		strings.HasPrefix(urlStr, "\\\\") {
		backend.Type = LocalDisk
		backend.URL = urlStr
		return nil
	}

	// 解析带scheme的URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("unable to parse URL: %v", err)
	}

	switch u.Scheme {
	case "smb":
		return backend.ParseSMBURL(u)
	case "sftp":
		return backend.ParseSFTPURL(u)
	case "webdav", "dav", "davs", "http", "https":
		return backend.ParseWebDAVURL(u)
	case "ftp", "ftps":
		return backend.ParseFTPURL(u)
	case "s3":
		return backend.ParseS3URL(u)
	default:
		return fmt.Errorf("unsupported URL scheme: %s", u.Scheme)
	}
}

// ParseFileURL 解析File URL
// 格式: file:///path/to/books 或 file:///C:/path/to/books
func (backend *Backend) ParseFileURL(urlStr string) error {
	// 处理 file:// 格式
	path := strings.TrimPrefix(urlStr, "file://")
	// 对于Unix路径，保留开头的斜杠；对于Windows路径，移除多余的斜杠
	if len(path) > 0 && path[0] == '/' {
		if len(path) > 1 && path[1] == '/' {
			if len(path) > 2 && path[2] == '/' {
				// 移除前三个斜杠，保留路径
				path = path[3:]
			} else {
				// 移除前两个斜杠
				path = path[2:]
			}
		} else if len(path) > 2 && path[1] != '/' && path[2] == ':' {
			// Windows路径格式 like /C:/path，移除前导斜杠
			path = path[1:]
		}
	}
	backend.Type = LocalDisk
	backend.URL = path
	return nil
}

// ParseSMBURL 解析SMB URL
// 格式: smb://workgroup;user:password@server/share/folder/books
func (backend *Backend) ParseSMBURL(u *url.URL) error {
	backend.Type = SMB

	// 解析主机信息
	host := u.Host
	if host == "" {
		return fmt.Errorf("SMB URL missing host address")
	}
	backend.ServerHost = host

	// 默认SMB端口
	backend.ServerPort = 445

	// 解析路径部分
	path := strings.TrimPrefix(u.Path, "/")
	pathParts := strings.SplitN(path, "/", 2)
	if len(pathParts) < 1 {
		return fmt.Errorf("SMB URL missing share name")
	}

	backend.SMBShareName = pathParts[0]
	if len(pathParts) > 1 {
		backend.SMBPath = pathParts[1]
	}

	// 解析认证信息 - SMB格式特殊处理
	if u.User != nil {
		backend.NeedAuth = true
		username := u.User.Username()
		// SMB格式可能是 workgroup;user 或 user
		if strings.Contains(username, ";") {
			parts := strings.SplitN(username, ";", 2)
			if len(parts) > 1 {
				backend.AuthUsername = parts[1] // 取分号后的用户名
			} else {
				backend.AuthUsername = username
			}
		} else {
			backend.AuthUsername = username
		}
		if password, ok := u.User.Password(); ok {
			backend.AuthPassword = password
		}
	}
	return nil
}

// ParseSFTPURL 解析SFTP URL
// 格式: sftp://user:password@host/path
func (backend *Backend) ParseSFTPURL(u *url.URL) error {
	backend.Type = SFTP

	// 解析主机和端口
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("SFTP URL missing host address")
	}
	backend.ServerHost = host

	port := u.Port()
	if port == "" {
		backend.ServerPort = 22 // 默认SFTP端口
	} else {
		if p, err := strconv.Atoi(port); err == nil {
			backend.ServerPort = p
		} else {
			return fmt.Errorf("invalid SFTP port: %s", port)
		}
	}

	// 解析认证信息
	if u.User != nil {
		backend.NeedAuth = true
		backend.AuthUsername = u.User.Username()
		if password, ok := u.User.Password(); ok {
			backend.AuthPassword = password
		}
	}

	return nil
}

// ParseWebDAVURL 解析WebDAV URL
// 格式: webdav://host/path 或 http://host/path 或 https://host/path
func (backend *Backend) ParseWebDAVURL(u *url.URL) error {
	backend.Type = WebDAV

	// 解析主机和端口
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("WebDAV URL missing host address")
	}
	backend.ServerHost = host

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "http", "dav":
			backend.ServerPort = 80
		case "https", "davs":
			backend.ServerPort = 443
		case "webdav":
			backend.ServerPort = 80 // 默认HTTP端口
		}
	} else {
		if p, err := strconv.Atoi(port); err == nil {
			backend.ServerPort = p
		} else {
			return fmt.Errorf("invalid WebDAV port: %s", port)
		}
	}

	// 解析认证信息
	if u.User != nil {
		backend.NeedAuth = true
		backend.AuthUsername = u.User.Username()
		if password, ok := u.User.Password(); ok {
			backend.AuthPassword = password
		}
	}
	return nil
}

// ParseFTPURL 解析FTP URL
// 格式: ftp://user:password@host:port/dir 或 ftps://user:password@host:port/dir
func (backend *Backend) ParseFTPURL(u *url.URL) error {
	backend.Type = FTP

	// 解析主机和端口
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("FTP URL missing host address")
	}
	backend.ServerHost = host

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "ftp":
			backend.ServerPort = 21
		case "ftps":
			backend.ServerPort = 990
		}
	} else {
		if p, err := strconv.Atoi(port); err == nil {
			backend.ServerPort = p
		} else {
			return fmt.Errorf("invalid FTP port: %s", port)
		}
	}

	// 解析认证信息
	if u.User != nil {
		backend.NeedAuth = true
		backend.AuthUsername = u.User.Username()
		if password, ok := u.User.Password(); ok {
			backend.AuthPassword = password
		}
	}

	return nil
}

// ParseS3URL 解析S3 URL
// 格式: s3://endpoint/bucket/prefix
func (backend *Backend) ParseS3URL(u *url.URL) error {
	backend.Type = S3

	// 解析主机（endpoint）
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("S3 URL missing endpoint")
	}
	backend.ServerHost = host

	// 解析端口
	port := u.Port()
	if port == "" {
		backend.ServerPort = 443 // 默认HTTPS端口
	} else {
		if p, err := strconv.Atoi(port); err == nil {
			backend.ServerPort = p
		} else {
			return fmt.Errorf("invalid S3 port: %s", port)
		}
	}

	// 验证路径（bucket和prefix）
	path := strings.TrimPrefix(u.Path, "/")
	pathParts := strings.SplitN(path, "/", 2)
	if len(pathParts) < 1 {
		return fmt.Errorf("S3 URL missing bucket name")
	}

	// S3通常不需要用户名密码认证，而是使用AWS凭证
	// 这里可以根据需要扩展支持AWS凭证配置

	return nil
}

// IsRemoteURL 判断 URL 是否为远程存储（非本地文件系统）
func IsRemoteURL(storeURL string) bool {
	backend := &Backend{}
	if err := backend.ParseStoreURL(storeURL); err != nil {
		return false
	}
	return backend.Type != LocalDisk
}

// GetBackendType 获取 URL 对应的后端类型
func GetBackendType(storeURL string) (BackendType, error) {
	backend := &Backend{}
	if err := backend.ParseStoreURL(storeURL); err != nil {
		return 0, err
	}
	return backend.Type, nil
}
