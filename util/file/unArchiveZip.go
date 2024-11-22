package file

import (
	"context"
	"errors"
	"os"

	"github.com/mholt/archives"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/encoding"
	"github.com/yumenaka/comigo/util/logger"
)

// UnArchiveZip 一次性解压 ZIP 文件
func UnArchiveZip(filePath string, extractPath string, textEncoding string) error {
	extractPath = util.GetAbsPath(extractPath)
	// 如果解压路径不存在，创建路径
	err := os.MkdirAll(extractPath, os.ModePerm)
	if err != nil {
		logger.Infof("Failed to create extract path: %v", err)
		return err
	}

	// 打开文件，只读模式
	file, err := os.Open(filePath)
	if err != nil {
		logger.Infof("Failed to open file: %v", err)
		return err
	}
	defer file.Close()

	// 确认文件格式
	format, _, err := archives.Identify(context.Background(), filePath, file)
	if err != nil {
		logger.Infof("Failed to identify file format: %v", err)
		return err
	}

	// 如果是 ZIP 文件 TODO:测试文件解压
	if zipFormat, ok := format.(archives.Zip); ok {
		if textEncoding != "" {
			zipFormat.TextEncoding = encoding.ByName(textEncoding)
		}
		ctx := context.WithValue(context.Background(), "extractPath", extractPath)

		err := zipFormat.Extract(ctx, file, extractFileHandler)
		if err != nil {
			logger.Infof("Failed to extract zip file: %v", err)
			return err
		}
		logger.Infof("ZIP 文件解压完成：%s 解压到：%s", util.GetAbsPath(filePath), extractPath)
	} else {
		logger.Infof("File is not a ZIP archive: %s", filePath)
		return errors.New("file is not a ZIP archive")
	}

	return nil
}
