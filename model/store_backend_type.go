package model

// BackendType 文件存储类型
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
