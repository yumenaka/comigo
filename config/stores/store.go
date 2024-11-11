package stores

// Store 书库设置
type Store struct {
	// 书库的类型,计划支持local，ftp、sftp、webdav、smb（2或3）
	Type string
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
