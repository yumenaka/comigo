package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"

	"github.com/mholt/archives"
	"github.com/yumenaka/comigo/util/encoding"
	"github.com/yumenaka/comigo/util/logger"
)

// 使用sync.Map代替map，保证并发情况下的读写安全
var mapBookFS sync.Map

// GetSingleFile  获取单个文件
// TODO:大文件需要针对性优化，可能需要保持打开状态、或通过持久化的虚拟文件系统获取
// TODO:测试文件缓存功能，一旦解压，下次直接读缓存文件
func GetSingleFile(filePath string, NameInArchive string, textEncoding string) ([]byte, error) {
	//必须传值
	if NameInArchive == "" {
		return nil, errors.New("NameInArchive is empty")
	}
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
	//判断是否是压缩包
	//context 用来取消，这里用不到，所以传了一个空的Context
	format, sourceArchive, err := archives.Identify(context.Background(), filePath, file)
	if err != nil {
		return nil, err
	}
	var data []byte
	//如果是zip文件,文件编码为UTF-8时textEncoding为空,其他特殊编码的zip文件根据设定指定（默认GBK）
	if ex, ok := format.(archives.Zip); ok {
		ex.TextEncoding = encoding.GetEncodingByName(textEncoding) // “”  "shiftjis" "gbk"
		ctx := context.Background()
		// 在这里对文件进行处理；比如如果只需要特定的文件或目录，则只需返回所需的 f.NameInArchive 值即可。 https://github.com/mholt/archives#extract-archive
		err := ex.Extract(ctx, file, func(ctx context.Context, f archives.FileInfo) error {
			if f.NameInArchive != NameInArchive {
				fmt.Println("Skip " + f.NameInArchive + "|" + NameInArchive + " in " + filePath)
				return nil
			}
			//if f.NameInArchive == NameInArchive {
			//	fmt.Println("Get " + NameInArchive + " in " + filePath)
			//}
			readCloser, err := f.Open()
			if err != nil {
				logger.Infof("%s", err)
			}
			defer readCloser.Close()
			content, err := io.ReadAll(readCloser)
			if err != nil {
				logger.Infof("%s", err)
			}
			data = content
			return err
		})
		return data, err
	}

	//通过一个持久化的虚拟文件系统读取文件（加快rar文件的解压速度），key是文件路径
	var fsys fs.FS
	fsysAny, fsOK := mapBookFS.Load(filePath)
	if fsOK {
		fsys = fsysAny.(fs.FS)
	} else {
		//如果从来没保存过这个文件系统
		temp, errFS := archives.FileSystem(context.Background(), filePath, nil)
		if errFS == nil {
			//将文件系统加入到sync.Map
			mapBookFS.Store(filePath, temp) //因为被gin并发调用，需要考虑并发读写问题
			fsys = temp
		} else {
			logger.Infof("%s", errFS)
		}
	}

	//通过虚拟文件系统打开特定文件
	fileInRarFS, errFSOpen := fsys.Open(NameInArchive)
	if errFSOpen != nil {
		logger.Infof("%s", errFSOpen)
	}
	//defer fileInRarFS.Close()
	if errFSOpen == nil {
		content, err := io.ReadAll(fileInRarFS)
		if err != nil {
			logger.Infof("%s", err)
		}
		data = content
		return data, nil
	}

	//TODO:Rar密码
	//如果是 Rar 文件
	if ex, ok := format.(archives.Rar); ok {
		//如果虚拟FS方案无效，继续用Extract方案
		ctx := context.Background()
		err := ex.Extract(ctx, sourceArchive, func(ctx context.Context, f archives.FileInfo) error {
			// 取得特定压缩文件
			fileInRar, err := f.Open()
			if err != nil {
				logger.Infof("%s", err)
			}
			defer func(fileInRar io.ReadCloser) {
				err := fileInRar.Close()
				if err != nil {
					logger.Infof("fileInRar.Close() Error:%s", err)
				}
			}(fileInRar)
			content, err := io.ReadAll(fileInRar)
			if err != nil {
				logger.Infof("%s", err)
			}
			data = content
			return err
		})
		return data, err
	}

	//其他格式的压缩包，正常情况下不应该用到
	if ex, ok := format.(archives.Extractor); ok {
		ctx := context.Background()
		err := ex.Extract(ctx, sourceArchive, func(ctx context.Context, f archives.FileInfo) error {
			// 取得特定压缩文件
			file, err := f.Open()
			if err != nil {
				logger.Infof("%s", err)
			}
			defer func(file io.ReadCloser) {
				err := file.Close()
				if err != nil {
					logger.Infof("file.Close() Error:%s", err)
				}
			}(file)
			content, err := io.ReadAll(file)
			if err != nil {
				logger.Infof("%s", err)
			}
			data = content
			return err
		})
		return data, err
	}
	return nil, errors.New("2,not Found " + NameInArchive + " in " + filePath)
}
