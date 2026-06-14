package store

import (
	"fmt"
	"net/url"

	"github.com/yumenaka/comigo/tools"
)

// Backend 文件后端。 可能是 1: 本地文件系统，2: SMB共享，3: SFTP服务器，4: WebDAV服务器，5: S3存储，6: FTP服务器，7: Comigo服务
// 可能用到的字段都放进一个结构体；调用方只填需要字段
// 核心配置是“URL”，其余大多数参数都解析此URL生成。
// 本地书库为文件路径（/home/pi/books C:\Users\用户名\books或  file://some_path/books），
// 其他类型的书库是对应文件服务的 url 形式，
// 如 smb://workgroup;user:password@server/share/folder/books
// 或 sftp://<user>:<password>@<host>/<path>
// 或 webdav://192.168.1.100/books  <scheme>://[<username>[:<password>]@]<host>[:<port>]/<base-path>/<resource-path>（也接受 dav:// 与 davs:// 作为同义写法）
// 或 https://user:password@example.com/comigo-base 另一个 Comigo 服务主页
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
// - WebDAV: webdav://host/path 或 dav://host/path 或 davs://host/path
// - Comigo: http://host/base 或 https://user:password@host/base
// - FTP: ftp://user:password@host:port/dir 或 ftps://user:password@host:port/dir
// - S3: s3://endpoint/bucket/prefix
func (backend *Backend) ParseStoreURL(urlStr string) error {
	info, err := tools.ParseStoreURL(urlStr)
	if err != nil {
		return err
	}
	backend.applyStoreURLInfo(info)
	return nil
}

// ParseFileURL 解析File URL
// 格式: file:///path/to/books 或 file:///C:/path/to/books
func (backend *Backend) ParseFileURL(urlStr string) error {
	return backend.parseKnownStoreURL(urlStr, tools.StoreBackendLocalDisk)
}

// ParseSMBURL 解析SMB URL
// 格式: smb://workgroup;user:password@server/share/folder/books
func (backend *Backend) ParseSMBURL(u *url.URL) error {
	return backend.parseKnownStoreURL(u.String(), tools.StoreBackendSMB)
}

// ParseSFTPURL 解析SFTP URL
// 格式: sftp://user:password@host/path
func (backend *Backend) ParseSFTPURL(u *url.URL) error {
	return backend.parseKnownStoreURL(u.String(), tools.StoreBackendSFTP)
}

// ParseWebDAVURL 解析WebDAV URL
// 格式: webdav://host/path 或 http://host/path 或 https://host/path
func (backend *Backend) ParseWebDAVURL(u *url.URL) error {
	return backend.parseKnownStoreURL(u.String(), tools.StoreBackendWebDAV)
}

// ParseFTPURL 解析FTP URL
// 格式: ftp://user:password@host:port/dir 或 ftps://user:password@host:port/dir
func (backend *Backend) ParseFTPURL(u *url.URL) error {
	return backend.parseKnownStoreURL(u.String(), tools.StoreBackendFTP)
}

// ParseS3URL 解析S3 URL
// 格式: s3://endpoint/bucket/prefix
func (backend *Backend) ParseS3URL(u *url.URL) error {
	return backend.parseKnownStoreURL(u.String(), tools.StoreBackendS3)
}

func (backend *Backend) parseKnownStoreURL(storeURL string, expectedType tools.StoreBackendType) error {
	info, err := tools.ParseStoreURL(storeURL)
	if err != nil {
		return err
	}
	if info.Type != expectedType {
		return fmt.Errorf("unexpected backend type: %v", info.Type)
	}
	backend.applyStoreURLInfo(info)
	return nil
}

func (backend *Backend) applyStoreURLInfo(info tools.StoreURLInfo) {
	backend.URL = info.URL
	backend.Type = backendTypeFromTools(info.Type)
	backend.ServerHost = info.ServerHost
	backend.ServerPort = info.ServerPort
	backend.NeedAuth = info.NeedAuth
	backend.AuthUsername = info.AuthUsername
	backend.AuthPassword = info.AuthPassword
	backend.SMBShareName = info.SMBShareName
	backend.SMBPath = info.SMBPath
}

func backendTypeFromTools(backendType tools.StoreBackendType) BackendType {
	switch backendType {
	case tools.StoreBackendSMB:
		return SMB
	case tools.StoreBackendSFTP:
		return SFTP
	case tools.StoreBackendWebDAV:
		return WebDAV
	case tools.StoreBackendS3:
		return S3
	case tools.StoreBackendFTP:
		return FTP
	case tools.StoreBackendComigo:
		return ComigoRemote
	default:
		return LocalDisk
	}
}

// IsRemoteURL 判断 URL 是否为远程存储（非本地文件系统）
func IsRemoteURL(storeURL string) bool {
	return tools.IsRemoteStoreURL(storeURL)
}

// GetBackendType 获取 URL 对应的后端类型
func GetBackendType(storeURL string) (BackendType, error) {
	info, err := tools.ParseStoreURL(storeURL)
	if err != nil {
		return 0, err
	}
	return backendTypeFromTools(info.Type), nil
}
