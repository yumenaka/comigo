package config

var version = "v1.3.8"

// minSupportedVersion 是可直接加载的 metadata 最小版本，仅在数据格式不兼容时提高。
var minSupportedVersion = "v1.3.0"

// GetVersion 返回当前程序版本。
func GetVersion() string {
	return version
}

// GetMinSupportedVersion 返回旧 JSON 数据的最小支持版本，与程序版本分开维护。
func GetMinSupportedVersion() string {
	return minSupportedVersion
}
