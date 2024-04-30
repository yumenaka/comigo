package types

import (
	"sync"
	"time"
)

type BookGroup struct {
	BookInfo
	ChildBook sync.Map //key：BookID,value: *BookInfo
}

// NewBook  初始化Book，设置文件路径、书名、BookID等等
func NewBookGroup(filePath string, modified time.Time, fileSize int64, storePath string, depth int, bookType SupportFileType) (*BookGroup, error) {
	//初始化书籍
	var group = BookGroup{
		BookInfo: BookInfo{
			Author:        "",
			Modified:      modified,
			FileSize:      fileSize,
			InitComplete:  false,
			Depth:         depth,
			BookStorePath: storePath,
			Type:          bookType},
	}
	//设置属性：
	//FilePath，转换为绝对路径
	group.setFilePath(filePath)
	group.setTitle(filePath)
	group.setAuthor()
	//设置属性：父文件夹
	group.setParentFolder(filePath)
	group.setBookID()
	return &group, nil
}
