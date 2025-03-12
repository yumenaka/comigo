package embed

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io/fs"
	"mime"
	"path/filepath"
)

//go:embed all:static
var Static embed.FS
var StaticFS fs.FS

//go:embed all:images
var Images embed.FS
var ImagesFS fs.FS

// GetFileStr 从Static获取字符串形式的脚本
func GetFileStr(filePath string) string {
	// 使用ReadFile从嵌入文件系统中读取文件内容
	data, err := Static.ReadFile(filePath)
	if err != nil {
		return "Not Found Script:" + filePath
	}
	// 将内容转换为字符串并返回
	return string(data)
}

// GetImageSrc 从Static获取Base64编码的图片的src属性
func GetImageSrc(filePath string) string {
	// 使用ReadFile从嵌入文件系统中读取文件内容
	data, err := Static.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return "Not Found Image:" + filePath
	}

	// 获取文件扩展名，并猜测MIME类型
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Base64 编码图片数据
	base64Data := base64.StdEncoding.EncodeToString(data)

	// 生成图片的 src 属性
	src := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

	return src
}

// GetData 从Static获取字节切片形式的数据
func GetData(filePath string) []byte {
	// 使用ReadFile从嵌入文件系统中读取文件内容
	data, err := Static.ReadFile(filePath)
	if err != nil {
		// 如果有错误发生，返回空的字节切片，并输出错误信息
		fmt.Println("Error:", err)
		return []byte{}
	}
	// 返回文件内容作为字节切片
	return data
}
