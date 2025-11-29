package file

import (
	"context"
	"errors"
	"os"

	"github.com/yumenaka/archives"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// UnArchiveRar 一次性解压 RAR 文件
func UnArchiveRar(filePath string, extractPath string) error {
	extractPath = tools.GetAbsPath(extractPath)
	// 如果解压路径不存在，创建路径
	err := os.MkdirAll(extractPath, os.ModePerm)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_create_extract_path"), err)
		return err
	}

	// 打开文件，读模式
	file, err := os.Open(filePath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_open_file_unarchive"), err)
		return err
	}
	defer file.Close()

	// 确认文件格式
	format, _, err := archives.Identify(context.Background(), filePath, file)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_identify_file_format"), err)
		return err
	}

	// 如果是 RAR 文件
	if rarFormat, ok := format.(archives.Rar); ok {
		ctx := context.WithValue(context.Background(), "extractPath", extractPath)

		err := rarFormat.Extract(ctx, file, extractFileHandler)
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_extract_rar_file"), err)
			return err
		}
		logger.Infof(locale.GetString("log_rar_file_extracted"), tools.GetAbsPath(filePath), extractPath)
	} else {
		logger.Infof(locale.GetString("err_file_not_rar_archive"), filePath)
		return errors.New(locale.GetString("err_file_not_rar_archive"))
	}

	return nil
}
