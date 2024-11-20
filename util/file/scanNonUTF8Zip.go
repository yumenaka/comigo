package file

import (
	"context"
	"errors"
	"os"

	"github.com/klauspost/compress/zip"
	"github.com/yumenaka/archiver/v4"
	"github.com/yumenaka/comigo/util/logger"
)

// ScanNonUTF8Zip 扫描文件，初始化书籍用
func ScanNonUTF8Zip(filePath string, textEncoding string) (reader *zip.Reader, err error) {
	//打开文件，只读模式
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		logger.Infof("%s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Infof("file.Close() Error:%s", err)
		}
	}(file)
	//是否是压缩包
	format, _, err := archiver.Identify(filePath, file)
	if err != nil {
		return nil, err
	}
	//如果是zip
	if ex, ok := format.(archiver.Zip); ok {
		ex.TextEncoding = textEncoding // “”  "shiftjis" "gbk"
		ctx := context.Background()
		////WithValue返回parent的一个副本，该副本保存了传入的key/value，而调用Context接口的Value(key)方法就可以得到val。注意在同一个context中设置key/value，若key相同，值会被覆盖。
		//ctx = context.WithValue(ctx, "extractPath", extractPath)
		reader, err := ex.LsAllFile(ctx, file, func(ctx context.Context, f archiver.File) error {
			//logger.Infof(f.title())
			return nil
		})
		if err != nil {
			return nil, err
		}
		return reader, err
	}
	return nil, errors.New("扫描文件错误")
}
