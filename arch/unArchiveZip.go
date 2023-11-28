package arch

import (
	"context"
	"os"

	"github.com/yumenaka/archiver/v4"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/util"
)

// UnArchiveZip 一次性解压zip文件
func UnArchiveZip(filePath string, extractPath string, textEncoding string) error {
	extractPath = util.GetAbsPath(extractPath)
	//如果解压路径不存在，创建路径
	err := os.MkdirAll(extractPath, os.ModePerm)
	if err != nil {
		logger.Info(err)
	}
	//打开文件，只读模式
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		logger.Info(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Info(err)
		}
	}(file)
	//是否是压缩包
	format, _, err := archiver.Identify(filePath, file)
	if err != nil {
		return err
	}
	//如果是zip
	if ex, ok := format.(archiver.Zip); ok {
		ex.TextEncoding = textEncoding // “”  "shiftjis" "gbk"
		ctx := context.Background()
		//WithValue返回parent的一个副本，该副本保存了传入的key/value，而调用Context接口的Value(key)方法就可以得到val。注意在同一个context中设置key/value，若key相同，值会被覆盖。
		ctx = context.WithValue(ctx, "extractPath", extractPath)
		_, err := ex.LsAllFile(ctx, file, extractFileHandler)
		if err != nil {
			return err
		}
		logger.Info("zip文件解压完成：" + util.GetAbsPath(filePath) + " 解压到：" + util.GetAbsPath(extractPath))
	}
	return nil
}
