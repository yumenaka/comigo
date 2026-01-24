package tools

import (
	"path/filepath"
	"runtime"
	"strings"
)

// NormalizeAbsPath 标准化文件路径为绝对路径
// 将路径转换为绝对路径，清理路径格式，Windows 下统一转换为小写
// 用于路径比较和冲突检测
func NormalizeAbsPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	abs = filepath.Clean(abs)
	// Windows 文件系统不区分大小写，统一转换为小写进行比较
	if runtime.GOOS == "windows" {
		abs = strings.ToLower(abs)
	}
	return abs, nil
}

// NormalizeAbsPathNoError 标准化文件路径为绝对路径（忽略错误，返回原路径）
// 适用于已知路径有效的场景，或者错误可以忽略的场景
func NormalizeAbsPathNoError(path string) string {
	normalized, err := NormalizeAbsPath(path)
	if err != nil {
		return path
	}
	return normalized
}
