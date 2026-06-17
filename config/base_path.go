package config

import (
	pathpkg "path"
	"strings"
)

// NormalizeBasePath 将用户配置的反向代理基础路径规范为 "" 或 "/some/path"。
func NormalizeBasePath(basePath string) string {
	basePath = strings.TrimSpace(basePath)
	if basePath == "" || basePath == "/" {
		return ""
	}
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	basePath = pathpkg.Clean(basePath)
	if basePath == "." || basePath == "/" {
		return ""
	}
	return strings.TrimRight(basePath, "/")
}

// GetBasePath 返回当前运行时使用的基础路径；空字符串表示根路径。
func GetBasePath() string {
	return cfg.GetBasePath()
}

// PrefixPath 为站内绝对路径加上基础路径，外部 URL 和 data/blob URL 会原样返回。
func PrefixPath(urlPath string) string {
	basePath := GetBasePath()
	if urlPath == "" {
		if basePath == "" {
			return "/"
		}
		return basePath + "/"
	}
	if strings.HasPrefix(urlPath, "http://") ||
		strings.HasPrefix(urlPath, "https://") ||
		strings.HasPrefix(urlPath, "data:") ||
		strings.HasPrefix(urlPath, "blob:") ||
		strings.HasPrefix(urlPath, "#") {
		return urlPath
	}
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	if basePath == "" {
		return urlPath
	}
	// PrefixPath 可能被通用组件二次调用；已带基础路径的站内 URL 原样返回，避免生成 /base/base/...
	if urlPath == basePath || strings.HasPrefix(urlPath, basePath+"/") {
		return urlPath
	}
	if urlPath == "/" {
		return basePath + "/"
	}
	return basePath + urlPath
}

// StripBasePath 去掉请求路径中的基础路径，便于已有页面类型判断复用。
func StripBasePath(urlPath string) string {
	basePath := GetBasePath()
	if basePath == "" || urlPath == "" {
		return urlPath
	}
	if urlPath == basePath {
		return "/"
	}
	if strings.HasPrefix(urlPath, basePath+"/") {
		return strings.TrimPrefix(urlPath, basePath)
	}
	return urlPath
}
