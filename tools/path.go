package tools

import (
	"path/filepath"
)

// NormalizeAbsPath 标准化文件路径为绝对路径
// 用于路径比较和冲突检测
func NormalizeAbsPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	abs = filepath.Clean(abs)
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
