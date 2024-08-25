package file

import (
	"context"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
	"io"
	"os"
	"path/filepath"

	"github.com/yumenaka/archiver/v4"
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
		logger.Infof("%s", err)
	}
	defer func(file io.ReadCloser) {
		err := file.Close()
		if err != nil {
			logger.Infof("file.Close() Error:%s", err)
		}
	}(file)
	//如果是文件夹，直接创建文件夹
	if f.IsDir() {
		err = os.MkdirAll(filepath.Join(extractPath, f.NameInArchive), os.ModePerm)
		if err != nil {
			logger.Infof("%s", err)
		}
		return err
	}
	//如果是一般文件，将文件写入磁盘
	writeFilePath := filepath.Join(extractPath, f.NameInArchive)
	//写文件前，如果对应文件夹不存在，就创建对应文件夹
	checkDir := filepath.Dir(writeFilePath)
	if !util.IsExist(checkDir) {
		err = os.MkdirAll(checkDir, os.ModePerm)
		if err != nil {
			logger.Infof("%s", err)
		}
		return err
	}
	//具体内容
	content, err := io.ReadAll(file)
	if err != nil {
		logger.Infof("%s", err)
	}
	//写入文件
	err = os.WriteFile(writeFilePath, content, 0644)
	if err != nil {
		logger.Infof("%s", err)
	}
	return err
}
