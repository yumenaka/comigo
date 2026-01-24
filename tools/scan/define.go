package scan

import (
	"path/filepath"
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

func InitConfig(c ConfigInterface) {
	cfg = c
}

// IsSupportTemplate 判断压缩包内的文件是否是支持的模板文件
func IsSupportTemplate(checkPath string) bool {
	// 如果是以 . 开头的隐藏文件，跳过
	if strings.HasPrefix(filepath.Base(checkPath), ".") {
		return false
	}
	for _, ex := range cfg.GetSupportTemplateFile() {
		suffix := strings.ToLower(filepath.Ext(checkPath)) // strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportMedia 判断文件是否需要展示
func IsSupportMedia(checkPath string) bool {
	// 如果是以 . 开头的隐藏文件，跳过
	if strings.HasPrefix(filepath.Base(checkPath), ".") {
		return false
	}
	for _, ex := range cfg.GetSupportMediaType() {
		suffix := strings.ToLower(filepath.Ext(checkPath)) // strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportFile 判断压缩包文件是否是支持的文件类型
func IsSupportFile(checkPath string) bool {
	for _, ex := range cfg.GetSupportFileType() {
		suffix := strings.ToLower(filepath.Ext(checkPath)) // strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSkipDir  检查路径是否应该跳过（排除文件，文件夹列表）。
func IsSkipDir(path string) bool {
	// 说明：
	// - 调用方有时传入“目录名”(如 "node_modules")，有时传入“完整路径”(如 "/a/b/node_modules")
	// - ExcludePath 的配置语义是“需要排除的文件或文件夹的名字”
	// 因此这里统一按 basename 匹配，同时兼容完整路径以 separator+name 结尾的情况，避免纯 HasSuffix 带来的误判。
	clean := filepath.Clean(path)
	base := filepath.Base(clean)
	for _, name := range cfg.GetExcludePath() {
		if name == "" {
			continue
		}
		// 1) 传入目录名时，直接按名称匹配
		if base == name {
			return true
		}
		// 2) 传入完整路径时，匹配以 "/name" 结尾（跨平台使用 filepath.Separator）
		if strings.HasSuffix(clean, string(filepath.Separator)+name) {
			return true
		}
		// 3) 兼容用户直接写完整路径作为排除规则的情况
		if clean == filepath.Clean(name) {
			return true
		}
	}
	return false
}
