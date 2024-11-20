package file

import (
	"context"
	"errors"
	"os"

	"github.com/yumenaka/archiver/v4"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
)

// UnArchiveRar 一次性解压 RAR 文件
func UnArchiveRar(filePath string, extractPath string) error {
	extractPath = util.GetAbsPath(extractPath)
	// 如果解压路径不存在，创建路径
	err := os.MkdirAll(extractPath, os.ModePerm)
	if err != nil {
		logger.Infof("Failed to create extract path: %v", err)
		return err
	}

	// 打开文件，读模式
	file, err := os.Open(filePath)
	if err != nil {
		logger.Infof("Failed to open file: %v", err)
		return err
	}
	defer file.Close()

	// 确认文件格式
	format, _, err := archiver.Identify(filePath, file)
	if err != nil {
		logger.Infof("Failed to identify file format: %v", err)
		return err
	}

	// 如果是 RAR 文件
	if rarFormat, ok := format.(archiver.Rar); ok {
		ctx := context.WithValue(context.Background(), "extractPath", extractPath)

		err := rarFormat.LsAllFile(ctx, file, extractFileHandler)
		if err != nil {
			logger.Infof("Failed to extract RAR file: %v", err)
			return err
		}
		logger.Infof("RAR 文件解压完成：%s 解压到：%s", util.GetAbsPath(filePath), extractPath)
	} else {
		logger.Infof("File is not a RAR archive: %s", filePath)
		return errors.New("file is not a RAR archive")
	}

	return nil
}
