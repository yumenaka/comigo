package file

import (
	"context"
	"errors"
	"os"

	"github.com/yumenaka/archives"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/encoding"
	"github.com/yumenaka/comigo/tools/logger"
)

// UnArchiveZip 一次性解压 ZIP 文件
func UnArchiveZip(filePath string, extractPath string, textEncoding string) error {
	extractPath = tools.GetAbsPath(extractPath)
	// 如果解压路径不存在，创建路径
	err := os.MkdirAll(extractPath, os.ModePerm)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_create_extract_path"), err)
		return err
	}

	// 打开文件，只读模式
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

	// 如果是 ZIP 文件
	if zipFormat, ok := format.(archives.Zip); ok {
		if textEncoding != "" {
			zipFormat.TextEncoding = encoding.ByName(textEncoding)
		}
		ctx := context.WithValue(context.Background(), "extractPath", extractPath)

		err := zipFormat.Extract(ctx, file, extractFileHandler)
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_extract_zip_file"), err)
			return err
		}
		logger.Infof(locale.GetString("log_zip_file_extracted"), tools.GetAbsPath(filePath), extractPath)
	} else {
		logger.Infof(locale.GetString("err_file_not_zip_archive"), filePath)
		return errors.New(locale.GetString("err_file_not_zip_archive"))
	}

	return nil
}
