package stores

type StoreType int

const (
	Local StoreType = 1 + iota
	FTP
	SMB
	SFTP
	WebDAV
	S3
)

// Store 书库设置
type Store struct {
	// 书库的类型,计划支持local，smb（2或3）ftp、sftp、webdav
	Type StoreType
	// 本地书库配置
	Local LocalOption
	// smb书库配置
	Smb SMBOption
}

type LocalOption struct {
	// 书库路径
	Path string
}
type SMBOption struct {
	// 书库的地址
	Host string
	// 书库的端口
	Port int
	// 书库的用户名
	Username string
	// 书库的密码
	Password string
	// smb的共享名
	ShareName string
	// 二级路径
	Path string
}
