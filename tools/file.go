package tools

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/yumenaka/comigo/tools/logger"
)

// 规范化路径，确保路径格式一致
func NormalizePath(path string) string {
	// 替换反斜杠为正斜杠
	path = strings.ReplaceAll(path, "\\", "/")
	// 移除末尾的斜杠
	path = strings.TrimSuffix(path, "/")
	// 确保路径不为空
	if path == "" {
		path = "."
	}
	return path
}

// 验证路径是否有效
func IsValidPath(path string) bool {
	// 路径不能为空
	if path == "" {
		return false
	}
	// 路径不能包含特殊字符
	if strings.ContainsAny(path, "?*<>|\"") {
		return false
	}
	return true
}

// RemoveExtension 从文件名中去除扩展名
func RemoveExtension(filename string) string {
	extension := filepath.Ext(filename)
	return strings.TrimSuffix(filename, extension)
}

// IsExist 判断文件夹或文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		logger.Infof("%s", err)
		return false
	}
	return true
}

// 删除文件
func DeleteFileIfExist(filePath string) error {
	// 使用os.Stat检查文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		// 文件存在，尝试删除
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		return errors.New("File does not exist:" + filePath)
	} else {
		return err
	}
	return nil
}

// ChickIsDir 判断所给路径是否为文件夹
func ChickIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// GetMainName 取得无后缀的文件名
func GetMainName(filename string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(filename)
	main := strings.TrimSuffix(base, ext)
	return main
}

// GetAbsPath 获取绝对路径
func GetAbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return abs
}

// PathExists reports whether the named file or directory exists.
func PathExists(path string) bool {
	if path == "" {
		return false
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsDir reports whether the named directory exists.
func IsDir(path string) bool {
	if path == "" {
		return false
	}

	if fi, err := os.Stat(path); err == nil {
		return fi.IsDir()
	}
	return false
}

// FileExists reports whether the named file or directory exists.
func FileExists(path string) bool {
	return IsFile(path)
}

// IsFile reports whether the named file or directory exists.
func IsFile(path string) bool {
	if path == "" {
		return false
	}

	if fi, err := os.Stat(path); err == nil {
		return !fi.IsDir()
	}
	return false
}

// IsAbsPath is abs path.
func IsAbsPath(aPath string) bool {
	return path.IsAbs(aPath)
}

// GetContentTypeByFileName https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
func GetContentTypeByFileName(fileName string) (contentType string) {
	ext := strings.ToLower(path.Ext(fileName))
	switch {
	case ext == ".png":
		contentType = "image/png"
	case ext == ".jpg" || ext == ".jpeg":
		contentType = "image/jpeg"
	case ext == ".webp":
		contentType = "image/webp"
	case ext == ".gif":
		contentType = "image/gif"
	case ext == ".bmp":
		contentType = "image/bmp"
	case ext == ".heif":
		contentType = "image/heif"
	case ext == ".ico":
		contentType = "image/image/vnd.microsoft.icon"
	case ext == ".zip":
		contentType = "application/zip"
	case ext == ".rar":
		contentType = "application/x-rar-compressed"
	case ext == ".pdf":
		contentType = "application/pdf"
	case ext == ".txt":
		contentType = "text/plain"
	case ext == ".tar":
		contentType = "application/x-tar"
	case ext == ".epub":
		contentType = "application/epub+zip"
	default:
		contentType = "application/octet-stream"
	}
	return contentType
}
