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
