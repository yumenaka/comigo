package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/your-project/model"
)

func scanDirGetBook(dir string) (*model.Book, error) {
	// 获取文件夹信息
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	// 使用文件夹的修改时间替代 time.Now()
	book := model.NewBook(
		filepath.Base(dir),
		dirInfo.ModTime(),  // 替换原来的 time.Now()
	)

	return book, nil
} 