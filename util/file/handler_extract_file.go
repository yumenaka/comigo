package file

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/yumenaka/archiver/v4"
	"github.com/yumenaka/comigo/util/logger"
)

// extractFileHandler 解压文件的处理函数
func extractFileHandler(ctx context.Context, f archiver.File) error {
	// 从上下文中获取解压路径
	extractPath, ok := ctx.Value("extractPath").(string)
	if !ok {
		return errors.New("extractPath not found in context")
	}

	// 打开压缩文件中的当前文件
	fileReader, err := f.Open()
	if err != nil {
		logger.Infof("Failed to open file in archive: %v", err)
		return err
	}
	defer fileReader.Close()

	// 目标文件路径
	targetPath := filepath.Join(extractPath, f.NameInArchive)

	// 如果是目录，创建目录并返回
	if f.IsDir() {
		err := os.MkdirAll(targetPath, os.ModePerm)
		if err != nil {
			logger.Infof("Failed to create directory: %v", err)
			return err
		}
		return nil
	}

	// 确保目标文件所在的目录存在
	err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
	if err != nil {
		logger.Infof("Failed to create parent directory: %v", err)
		return err
	}

	// 创建目标文件
	destFile, err := os.Create(targetPath)
	if err != nil {
		logger.Infof("Failed to create file: %v", err)
		return err
	}
	defer destFile.Close()

	// 将文件内容从压缩包复制到目标文件
	_, err = io.Copy(destFile, fileReader)
	if err != nil {
		logger.Infof("Failed to copy file content: %v", err)
		return err
	}

	return nil
}
