package entity

import (
	"sync"
	"time"
)

type BookGroup struct {
	BookInfo
	ChildBook sync.Map //key：BookID,value: *BookInfo
}

// NewBookGroup   初始化BookGroup，设置文件路径、书名、BookID等等
func NewBookGroup(filePath string, modified time.Time, fileSize int64, storePath string, depth int, bookType SupportFileType) (*BookGroup, error) {
	//初始化书籍
	var group = BookGroup{
		BookInfo: BookInfo{
			Modified:      modified,
			FileSize:      fileSize,
			InitComplete:  false,
			Depth:         depth,
			BookStorePath: storePath,
			Type:          bookType},
	}
	//设置属性：
	group.setTitle(filePath).setFilePath(filePath).setAuthor().setParentFolder(filePath).initBookID()
	return &group, nil
}
