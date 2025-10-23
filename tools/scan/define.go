package scan

import (
	"path"
	"strings"
)

type ConfigInterface interface {
	GetStoreUrls() []string
	GetMaxScanDepth() int
	GetMinImageNum() int
	GetTimeoutLimitForScan() int
	GetExcludePath() []string
	GetSupportMediaType() []string
	GetSupportFileType() []string
	GetSupportTemplateFile() []string
	GetZipFileTextEncoding() string
	GetEnableDatabase() bool
	GetClearDatabaseWhenExit() bool
	GetDebug() bool
}

var cfg ConfigInterface
var minImageNum int

func init() {
	minImageNum = 1
}
func InitConfig(c ConfigInterface) {
	cfg = c
	minImageNum = cfg.GetMinImageNum()
}

// IsSupportTemplate 判断压缩包内的文件是否是支持的模板文件
func IsSupportTemplate(checkPath string) bool {
	for _, ex := range cfg.GetSupportTemplateFile() {
		suffix := strings.ToLower(path.Ext(checkPath)) // strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportMedia 判断文件是否需要展示
func IsSupportMedia(checkPath string) bool {
	for _, ex := range cfg.GetSupportMediaType() {
		suffix := strings.ToLower(path.Ext(checkPath)) // strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportFile 判断压缩包内的文件是否需要展示
func IsSupportFile(checkPath string) bool {
	for _, ex := range cfg.GetSupportFileType() {
		suffix := strings.ToLower(path.Ext(checkPath)) // strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportArchiver 是否是支持的压缩文件
func IsSupportArchiver(checkPath string) bool {
	for _, ex := range cfg.GetSupportFileType() {
		suffix := path.Ext(checkPath)
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSkipDir  检查路径是否应该跳过（排除文件，文件夹列表）。
func IsSkipDir(path string) bool {
	for _, substr := range cfg.GetExcludePath() {
		if strings.HasSuffix(path, substr) {
			return true
		}
	}
	return false
}
