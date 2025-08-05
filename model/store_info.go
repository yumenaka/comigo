package model

import "time"

// StoreInfo 书库基本信息
type StoreInfo struct {
	URL           string
	Name          string
	Description   string
	FileBackendID int64 // 文件后端ID, -1 表示未设置
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
