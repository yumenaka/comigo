package file

import (
	"context"
	"errors"
	"os"

	"github.com/klauspost/compress/zip"
	"github.com/yumenaka/archives"
	"github.com/yumenaka/comigo/tools/encoding"
	"github.com/yumenaka/comigo/tools/logger"
)

// ScanNonUTF8Zip 扫描文件，初始化书籍用
func ScanNonUTF8Zip(filePath string, textEncoding string) (reader *zip.Reader, err error) {
	// 打开文件，只读模式
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o400) // Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		logger.Infof("%s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Infof("file.Close() Error:%s", err)
		}
	}(file)
	// 是否是压缩包
	format, _, err := archives.Identify(context.Background(), filePath, file)
	if err != nil {
		return nil, err
	}
	// 如果是zip
	if ex, ok := format.(archives.Zip); ok {
		if textEncoding != "" {
			ex.TextEncoding = encoding.ByName(textEncoding)
		}
		ctx := context.Background()
		reader, err := ex.CheckNonUTF8Zip(ctx, file, func(ctx context.Context, f archives.FileInfo) error {
			// logger.Infof(f.title())
			return nil
		})
		if err != nil {
			return nil, err
		}
		return reader, err
	}
	return nil, errors.New("扫描文件错误")
}
