package file

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/yumenaka/archives"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/encoding"
	"github.com/yumenaka/comigo/tools/logger"
)

// 使用sync.Map代替map，保证并发情况下的读写安全
var mapBookFS sync.Map

// GetSingleFile 获取单个文件
func GetSingleFile(filePath, nameInArchive, textEncoding string) ([]byte, error) {
	if nameInArchive == "" {
		return nil, errors.New(locale.GetString("err_name_in_archive_empty"))
	}
	// 创建一个30秒超时的Context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_open_file_get_single"), filePath, err)
		return nil, err
	}
	defer file.Close()

	// 识别压缩格式
	format, sourceArchiveReader, err := archives.Identify(ctx, filePath, file)
	if err != nil {
		// 检查是否是超时错误
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Info(locale.GetString("log_timeout_identify_archive_format"))
		} else {
			logger.Infof(locale.GetString("log_failed_to_identify_archive_format"), err)
		}
		return nil, err
	}

	// 处理 ZIP 文件
	if zipFormat, ok := format.(archives.Zip); ok {
		if textEncoding != "" {
			zipFormat.TextEncoding = encoding.ByName(textEncoding)
		}
		return extractFileFromArchive(ctx, zipFormat, file, nameInArchive)
	}

	// 从缓存中获取虚拟文件系统
	var fileSystem fs.FS
	if fsInterface, ok := mapBookFS.Load(filePath); ok {
		fileSystem = fsInterface.(fs.FS)
	} else {
		fileSystem, err = archives.FileSystem(ctx, filePath, nil)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Info(locale.GetString("log_timeout_create_filesystem"))
			} else {
				logger.Infof(locale.GetString("log_failed_to_create_filesystem"), err)
			}
			return nil, err
		}
		mapBookFS.Store(filePath, fileSystem)
	}

	// 从虚拟文件系统中读取文件
	fileInArchive, err := fileSystem.Open(nameInArchive)
	if err == nil {
		defer fileInArchive.Close()
		data, err := io.ReadAll(fileInArchive)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Info(locale.GetString("log_timeout_read_file_content"))
			} else {
				logger.Infof(locale.GetString("log_failed_to_read_file_content"), err)
			}
			return nil, err
		}
		return data, nil
	}

	// 处理 RAR 文件
	if rarFormat, ok := format.(archives.Rar); ok {
		return extractFileFromArchive(ctx, rarFormat, sourceArchiveReader, nameInArchive)
	}

	// 处理其他格式的压缩文件
	if extractor, ok := format.(archives.Extractor); ok {
		return extractFileFromArchive(ctx, extractor, sourceArchiveReader, nameInArchive)
	}

	return nil, errors.New(locale.GetString("err_unsupported_archive_format"))
}

func extractFileFromArchive(ctx context.Context, extractor archives.Extractor, sourceArchive io.Reader, nameInArchive string) ([]byte, error) {
	var data []byte
	err := extractor.Extract(ctx, sourceArchive, func(ctx context.Context, f archives.FileInfo) error {
		if f.NameInArchive != nameInArchive {
			return nil
		}
		readCloser, err := f.Open()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Info(locale.GetString("log_timeout_open_file_in_archive"))
			} else {
				logger.Infof(locale.GetString("log_failed_to_open_file_in_archive"), err)
			}
			return err
		}
		defer readCloser.Close()
		data, err = io.ReadAll(readCloser)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Info(locale.GetString("log_timeout_read_file_content"))
			} else {
				logger.Infof(locale.GetString("log_failed_to_read_file_content"), err)
			}
		}
		return err
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Info(locale.GetString("log_timeout_extract_file"))
		} else {
			logger.Infof(locale.GetString("log_failed_to_extract_file"), err)
		}
		return nil, err
	}
	if data != nil {
		return data, nil
	}
	return nil, errors.New(locale.GetString("err_file_not_found_in_archive"))
}
