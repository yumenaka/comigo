package arch

import (
	"context"
	"fmt"
	"github.com/yumenaka/archiver/v4"
	"github.com/yumenaka/comi/tools"
	"io"
	"os"
	"path/filepath"
)

// 解压文件的函数
func extractFileHandler(ctx context.Context, f archiver.File) error {
	extractPath := ""
	if e, ok := ctx.Value("extractPath").(string); ok {
		extractPath = e
	}
	// 取得压缩文件
	file, err := f.Open()
	if err != nil {
		fmt.Println(err)
	}
	defer func(file io.ReadCloser) {
		err := file.Close()
		if err != nil {
			fmt.Println("file.Close() Error:", err)
		}
	}(file)
	//如果是文件夹，直接创建文件夹
	if f.IsDir() {
		err = os.MkdirAll(filepath.Join(extractPath, f.NameInArchive), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		return err
	}
	//如果是一般文件，将文件写入磁盘
	writeFilePath := filepath.Join(extractPath, f.NameInArchive)
	//写文件前，如果对应文件夹不存在，就创建对应文件夹
	checkDir := filepath.Dir(writeFilePath)
	if !tools.IsExist(checkDir) {
		err = os.MkdirAll(checkDir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		return err
	}
	//具体内容
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	//写入文件
	err = os.WriteFile(writeFilePath, content, 0644)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
