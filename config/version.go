package config

import "golang.org/x/mod/semver"

var version = "v1.4.0"

// minSupportedVersion 是 metadata 与远程 Comigo 的最小支持版本，仅在数据格式不兼容时提高。
var minSupportedVersion = "v1.4.0"

// GetVersion 返回当前程序版本。
func GetVersion() string {
	return version
}

// GetMinSupportedVersion 返回索引与远程 Comigo 的最小支持版本，与程序版本分开维护。
func GetMinSupportedVersion() string {
	return minSupportedVersion
}

// SupportsBookVersion 判断书籍索引是否可以直接复用；旧算法的索引需要重扫。
func SupportsBookVersion(version string) bool {
	return semver.IsValid(version) && semver.Compare(version, minSupportedVersion) >= 0
}
