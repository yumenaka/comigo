package model

import (
	"path"
	"strings"
)

type SupportFileType string

// 书籍类型
const (
	TypeDir         SupportFileType = "dir"
	TypeZip         SupportFileType = ".zip"
	TypeRar         SupportFileType = ".rar"
	TypeBooksGroup  SupportFileType = "book_group"
	TypeCbz         SupportFileType = ".cbz"
	TypeCbr         SupportFileType = ".cbr"
	TypeTar         SupportFileType = ".tar"
	TypeEpub        SupportFileType = ".epub"
	TypePDF         SupportFileType = ".pdf"
	TypeVideo       SupportFileType = "video"
	TypeAudio       SupportFileType = "audio"
	TypeUnknownFile SupportFileType = "unknown"
)

// GetBookTypeByFilename 初始化Book时，取得BookType
func GetBookTypeByFilename(filename string) SupportFileType {
	//获取文件后缀
	switch strings.ToLower(path.Ext(filename)) {
	case ".zip":
		return TypeZip
	case ".rar":
		return TypeRar
	case ".cbz":
		return TypeCbz
	case ".cbr":
		return TypeCbr
	case ".epub":
		return TypeEpub
	case ".tar":
		return TypeTar
	case ".pdf":
		return TypePDF
	case ".mp4", ".m4v", ".flv", ".avi", ".webm":
		return TypeVideo
	case ".mp3", ".wav", ".wma", ".ogg":
		return TypeAudio
	default:
		return TypeUnknownFile
	}
}
