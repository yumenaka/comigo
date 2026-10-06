package model

import "time"

// ScanFailureRecord 记录压缩文件扫描失败时的文件指纹。
// 后续扫描仅在文件变化，或版本跨度足够大时才会再次尝试。
type ScanFailureRecord struct {
	StoreURL         string    `json:"store_url"`
	FilePath         string    `json:"file_path"`
	FileSize         int64     `json:"file_size"`
	ModifiedUnixNano int64     `json:"modified_unix_nano"`
	CreatedByVersion string    `json:"created_by_version"`
	FailedAt         time.Time `json:"failed_at"`
	Error            string    `json:"error"`
	IsRemote         bool      `json:"is_remote"`
}
