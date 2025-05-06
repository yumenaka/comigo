package scan

import (
	"path"
	"strings"

	"github.com/yumenaka/comigo/config/stores"
)

type ConfigInterface interface {
	GetLocalStores() []string
	GetStores() []stores.Store
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

type Option struct {
	Cfg ConfigInterface
}

func NewOption(scanConfig ConfigInterface) Option {
	return Option{
		Cfg: scanConfig,
	}
}

// IsSupportTemplate 判断压缩包内的文件是否是支持的模板文件
func (o *Option) IsSupportTemplate(checkPath string) bool {
	for _, ex := range o.Cfg.GetSupportTemplateFile() {
		suffix := strings.ToLower(path.Ext(checkPath)) // strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportMedia 判断文件是否需要展示
func (o *Option) IsSupportMedia(checkPath string) bool {
	for _, ex := range o.Cfg.GetSupportMediaType() {
		suffix := strings.ToLower(path.Ext(checkPath)) // strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportFile 判断压缩包内的文件是否需要展示
func (o *Option) IsSupportFile(checkPath string) bool {
	for _, ex := range o.Cfg.GetSupportFileType() {
		suffix := strings.ToLower(path.Ext(checkPath)) // strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportArchiver 是否是支持的压缩文件
func (o *Option) IsSupportArchiver(checkPath string) bool {
	for _, ex := range o.Cfg.GetSupportFileType() {
		suffix := path.Ext(checkPath)
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSkipDir  检查路径是否应该跳过（排除文件，文件夹列表）。
func (o *Option) IsSkipDir(path string) bool {
	for _, substr := range o.Cfg.GetExcludePath() {
		if strings.HasSuffix(path, substr) {
			return true
		}
	}
	return false
}
