package arch

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/mholt/archiver/v4"
	"io/ioutil"
	"os"
	"path/filepath"
)

//var (
//	compressionLevel       int
//	overwriteExisting      bool
//	mkdirAll               bool
//	selectiveCompression   bool
//	implicitTopLevelFolder bool
//	continueOnError        bool
//)
//
//func init() {
//	mkdirAll = true
//	overwriteExisting = false
//	continueOnError = true
//}

// ScanNonUTF8Zip 扫描文件，初始化书籍用
func ScanNonUTF8Zip(filePath string, textEncoding string) (reader *zip.Reader, err error) {
	//打开文件，只读模式
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	//是否是压缩包
	format, err := archiver.Identify(filePath, file)
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
			//fmt.Println(f.Name())
			return nil
		})
		if err != nil {
			return nil, err
		}
		return reader, err
	}
	return nil, errors.New("扫描文件错误")
}

// UnArchiveZip 一次性解压zip文件
func UnArchiveZip(filePath string, extractPath string, textEncoding string) error {
	extractPath = getAbsPath(extractPath)
	//如果解压路径不存在，创建路径
	err := os.MkdirAll(extractPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	//打开文件，只读模式
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	//是否是压缩包
	format, err := archiver.Identify(filePath, file)
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
		fmt.Println("zip文件解压完成：" + getAbsPath(filePath) + " 解压到：" + getAbsPath(extractPath))
	}
	return nil
}

// UnArchiveRar  一次性解压rar文件
func UnArchiveRar(filePath string, extractPath string) error {
	extractPath = getAbsPath(extractPath)
	//如果解压路径不存在，创建路径
	err := os.MkdirAll(extractPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	//打开文件，只读模式
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	//是否是压缩包
	format, err := archiver.Identify(filePath, file)
	if err != nil {
		return err
	}
	//如果是rar
	if ex, ok := format.(archiver.Rar); ok {
		ctx := context.Background()
		//WithValue返回parent的一个副本，该副本保存了传入的key/value，而调用Context接口的Value(key)方法就可以得到val。注意在同一个context中设置key/value，若key相同，值会被覆盖。
		ctx = context.WithValue(ctx, "extractPath", extractPath)
		err := ex.LsAllFile(ctx, file, extractFileHandler)
		if err != nil {
			return err
		}
		fmt.Println("rar文件解压完成：" + getAbsPath(filePath) + " 解压到：" + getAbsPath(extractPath))
	}
	return nil
}

//解压文件的函数
func extractFileHandler(ctx context.Context, f archiver.File) error {
	extractPath := ""
	if e, ok := ctx.Value("extractPath").(string); ok {
		extractPath = e
	}
	// 取得压缩文件
	file, err := f.Open()
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	//如果是文件夹，直接创建文件夹
	if f.IsDir() {
		err = os.MkdirAll(filepath.Join(extractPath, f.NameInArchive), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		return err
	}
	//如果是一般文件，将文件写入磁盘
	writeFilePath := filepath.Join(extractPath, f.NameInArchive)
	//写文件前，如果对应文件夹不存在，就创建对应文件夹
	checkDir := filepath.Dir(writeFilePath)
	if !isExist(checkDir) {
		err = os.MkdirAll(checkDir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		return err
	}
	//具体内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	//写入文件
	err = ioutil.WriteFile(writeFilePath, content, 0644)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// GetSingleFile  获取单个文件
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
	defer file.Close()
	//是否是压缩包
	format, err := archiver.Identify(filePath, file)
	if err != nil {
		return nil, err
	}

	var data []byte
	//如果是特殊编码的zip文件
	if textEncoding != "" {
		if ex, ok := format.(archiver.Zip); ok {
			ex.TextEncoding = textEncoding // “”  "shiftjis" "gbk"
			ctx := context.Background()
			err := ex.Extract(ctx, file, []string{NameInArchive}, func(ctx context.Context, f archiver.File) error {
				// 取得特定压缩文件
				file, err := f.Open()
				if err != nil {
					fmt.Println(err)
				}
				defer file.Close()
				content, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println(err)
				}
				data = content
				return err
			})
			return data, err
		}
	} else {
		//如果不是特殊编码的zip文件
		if ex, ok := format.(archiver.Extractor); ok {

			ctx := context.Background()
			err := ex.Extract(ctx, file, []string{NameInArchive}, func(ctx context.Context, f archiver.File) error {
				// 取得特定压缩文件
				file, err := f.Open()
				if err != nil {
					fmt.Println(err)
				}
				defer file.Close()
				content, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println(err)
				}
				data = content
				return err
			})
			return data, err
		}
	}
	return nil, errors.New("2,not Found " + NameInArchive + " in " + filePath)
}

//判断文件夹或文件是否存在
func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

//获取绝对路径
func getAbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return abs
}
