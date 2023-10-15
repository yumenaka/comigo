package arch

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/yumenaka/archiver/v4"
)

// GetSingleFile  获取单个文件
// TODO:大文件需要针对性优化，可能需要保持打开状态、或通过持久化的虚拟文件系统获取
// TODO:可选择文件缓存功能，一旦解压，下次直接读缓存文件
func GetSingleFile(filePath string, NameInArchive string, textEncoding string) ([]byte, error) {
	//必须传值
	if NameInArchive == "" {
		return nil, errors.New("NameInArchive is empty")
	}
	//打开文件，只读模式
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		fmt.Println(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("file.Close() Error:", err)
		}
	}(file)
	//是否是压缩包
	format, sourceArchive, err := archiver.Identify(filePath, file)
	if err != nil {
		return nil, err
	}
	var data []byte
	//如果是zip文件,文件编码为UTF-8时textEncoding为空,其他特殊编码的zip文件根据设定指定（默认GBK）
	if ex, ok := format.(archiver.Zip); ok {
		//特殊编码
		ex.TextEncoding = textEncoding // “”  "shiftjis" "gbk"
		ctx := context.Background()
		//这里是file，而不是sourceArchive，否则会出错。
		err := ex.Extract(ctx, file, []string{NameInArchive}, func(ctx context.Context, f archiver.File) error {
			// 取得特定压缩文件
			file, err := f.Open()
			if err != nil {
				fmt.Println(err)
			}
			//defer file.Close()
			content, err := io.ReadAll(file)
			if err != nil {
				fmt.Println(err)
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
		//会引发500错误，原因似乎是10秒后的timeout？
		////archiver.FileSystem可以配合ctx了，加个默认超时时间
		//const shortDuration = 10 * 1000 * time.Millisecond //超时时间，10秒
		//ctx, cancel := context.WithTimeout(context.Background(), shortDuration)
		//defer cancel()

		//如果从来没保存过这个文件系统
		temp, errFS := archiver.FileSystem(context.Background(), filePath)
		if errFS == nil {
			//将文件系统加入到sync.Map
			mapBookFS.Store(filePath, temp) //因为被gin并发调用，需要考虑并发读写问题
			fsys = temp
		} else {
			fmt.Println(errFS)
		}
	}

	//通过虚拟文件系统打开特定文件
	fileInRarFS, errFSOpen := fsys.Open(NameInArchive)
	if errFSOpen != nil {
		fmt.Println(errFSOpen)
	}
	//defer fileInRarFS.Close()
	if errFSOpen == nil {
		content, err := io.ReadAll(fileInRarFS)
		if err != nil {
			fmt.Println(err)
		}
		data = content
		return data, nil
	}

	//TODO:Rar密码
	//如果是 Rar 文件
	if ex, ok := format.(archiver.Rar); ok {
		//如果虚拟FS方案无效，继续用Extract方案
		ctx := context.Background()
		err := ex.Extract(ctx, sourceArchive, []string{NameInArchive}, func(ctx context.Context, f archiver.File) error {
			// 取得特定压缩文件
			fileInRar, err := f.Open()
			if err != nil {
				fmt.Println(err)
			}
			defer func(fileInRar io.ReadCloser) {
				err := fileInRar.Close()
				if err != nil {
					fmt.Println("fileInRar.Close() Error:", err)
				}
			}(fileInRar)
			content, err := io.ReadAll(fileInRar)
			if err != nil {
				fmt.Println(err)
			}
			data = content
			return err
		})
		return data, err
	}

	//其他格式的压缩包，正常情况下不应该用到
	if ex, ok := format.(archiver.Extractor); ok {
		ctx := context.Background()
		err := ex.Extract(ctx, sourceArchive, []string{NameInArchive}, func(ctx context.Context, f archiver.File) error {
			// 取得特定压缩文件
			file, err := f.Open()
			if err != nil {
				fmt.Println(err)
			}
			defer func(file io.ReadCloser) {
				err := file.Close()
				if err != nil {
					fmt.Println("file.Close() Error:", err)
				}
			}(file)
			content, err := io.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}
			data = content
			return err
		})
		return data, err
	}
	return nil, errors.New("2,not Found " + NameInArchive + " in " + filePath)
}
