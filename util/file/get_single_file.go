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
	"github.com/yumenaka/comigo/util/encoding"
	"github.com/yumenaka/comigo/util/logger"
)

// 使用sync.Map代替map，保证并发情况下的读写安全
var mapBookFS sync.Map

// GetSingleFile 获取单个文件
func GetSingleFile(filePath, nameInArchive, textEncoding string) ([]byte, error) {
	if nameInArchive == "" {
		return nil, errors.New("nameInArchive is empty")
	}

	// 创建一个30秒超时的Context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		logger.Infof("无法打开文件 %s: %v", filePath, err)
		return nil, err
	}
	defer file.Close()

	// 识别压缩格式
	format, sourceArchiveReader, err := archives.Identify(ctx, filePath, file)
	if err != nil {
		// 检查是否是超时错误
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Infof("操作超时：识别压缩格式花费了超过30秒")
		} else {
			logger.Infof("识别压缩格式失败: %v", err)
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
				logger.Infof("操作超时：创建文件系统花费了超过30秒")
			} else {
				logger.Infof("创建文件系统失败: %v", err)
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
				logger.Infof("操作超时：读取文件内容花费了超过30秒")
			} else {
				logger.Infof("读取文件内容失败: %v", err)
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

	return nil, errors.New("不支持的压缩格式或在压缩包中未找到文件")
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
				logger.Infof("操作超时：打开压缩包内文件花费了超过30秒")
			} else {
				logger.Infof("打开压缩包内文件失败: %v", err)
			}
			return err
		}
		defer readCloser.Close()
		data, err = io.ReadAll(readCloser)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Infof("操作超时：读取文件内容花费了超过30秒")
			} else {
				logger.Infof("读取文件内容失败: %v", err)
			}
		}
		return err
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Infof("操作超时：提取文件花费了超过30秒")
		} else {
			logger.Infof("提取文件失败: %v", err)
		}
		return nil, err
	}
	if data != nil {
		return data, nil
	}
	return nil, errors.New("在压缩包中未找到文件")
}
