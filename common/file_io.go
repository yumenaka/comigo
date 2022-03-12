package common

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/yumenaka/comi/locale"
)

func ScanArchiveOrFolder(FilePath string) (*Book, error) {
	//打开文件
	var file, err = os.OpenFile(FilePath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	FileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
	}
	//设置文件路径等等
	book := Book{AllPageNum: 0, FilePath: FilePath, Modified: FileInfo.ModTime(), IsDir: FileInfo.IsDir(), FileSize: FileInfo.Size(), ExtractComplete: false}
	//设置书籍UUID，根据路径算出
	book.InitBook(book.FilePath)
	//为了解决archiver/v4的BUG “zip文件无法读取2级目录” 单独处理zip文件
	ext := path.Ext(FilePath)
	if ext == ".zip" || ext == ".epub" {
		//建立一个zipfs，无法处理非UTF-8编码
		fsys, zip_err := zip.OpenReader(FilePath)
		if zip_err != nil {
			fmt.Println(zip_err)
		}
		////其他类型的压缩文件或文件夹
		//fsys, err := archiver.FileSystem(FilePath)
		//if err != nil {
		//	return nil, err
		//}
		err = walkZipFs(fsys, "", ".", &book)
		if _, ok := err.(*fs.PathError); ok {
			book.NonUTF8 = true
			fmt.Println("出现编码错误")
			err = walkZipFs(fsys, "", ".", &book)

			format, err := archiver.Identify(FilePath, file)
			if err != nil {
				fmt.Println(err)
			}
			if ex, ok := format.(archiver.Zip); ok {
				ex.TextEncoding = "shiftjis"
				ex.LsAllFile(context.Background(), file, func(_ context.Context, f archiver.File) error {
					// do something with the file ...
					fmt.Println(f)
					f.Open()
					return nil
				})
			}
		}
		if Config.SortImage != "none" {
			book.SortPages()
		}
		return &book, err
	}
	//其他类型的压缩文件或文件夹
	fsys, err := archiver.FileSystem(FilePath)
	if err != nil {
		return nil, err
	}
	////等效函数：
	//walkArchiveFs(fsys, "", ".", &book)
	// https://bitfieldconsulting.com/golang/filesystems
	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		var exclude = false
		for _, substr := range ExcludeFileOrFolders {
			if strings.Contains(path, substr) && path != "" {
				exclude = true
			}
		}
		if exclude {
			return fs.SkipDir
		}
		f, errInfo := d.Info()
		if errInfo != nil {
			fmt.Println(errInfo)
		}
		if isSupportMedia(path) {
			book.AllPageNum++
			u, ok := f.(archiver.File) //f.Name不包含路径信息.需要转换一下
			if !ok {
				//如果是文件夹中的图片
				book.PageInfo = append(book.PageInfo, SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: "", Url: "/cache/" + book.BookID + "/" + url.PathEscape(path)})
				//fmt.Println(locale.GetString("unsupported_extract")+" %s", f)
			} else {
				//如果是压缩文件
				TempURL := "/cache/" + book.BookID + "/" + url.PathEscape(u.NameInArchive)
				book.PageInfo = append(book.PageInfo, SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: u.NameInArchive, Url: TempURL})
			}
		} else {
			logrus.Debugf(locale.GetString("unsupported_file_type") + path)
		}
		return nil

	})
	if Config.SortImage != "none" {
		book.SortPages()
	}
	return &book, err
}

//手动写的递归查找，功能与fs.WalkDir()相同。发现zip文件的虚拟文件系统，似乎找不到正确的文件夹？
// https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter06/06.3.html
func walkZipFs(fsys fs.FS, parent, base string, book *Book) error {
	//fmt.Println("parent:" + parent + " base:" + base)
	dirName := path.Join(parent, base)
	dirEntries, err := fs.ReadDir(fsys, dirName)
	for _, dirEntry := range dirEntries {
		name := dirEntry.Name()
		f, errInfo := dirEntry.Info()
		if errInfo != nil {
			continue
		}
		if dirEntry.IsDir() == true {
			switch name {
			case ".comigo":
				return fs.SkipDir
			case "flutter_ui":
				return fs.SkipDir
			case "node_modules":
				return fs.SkipDir
			default:
			}
			join_path := path.Join(parent, name)
			err = walkZipFs(fsys, join_path, base, book)
		} else if !isSupportMedia(name) {
			logrus.Debugf(locale.GetString("unsupported_file_type") + name)
		} else {
			book.AllPageNum++
			inArchiveName := path.Join(parent, f.Name())
			book.PageInfo = append(book.PageInfo, SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: inArchiveName, Url: "/cache/" + book.BookID + "/" + url.PathEscape(inArchiveName)})
		}
	}
	return err
}

func isSupportMedia(checkPath string) bool {
	for _, ex := range SupportMediaType {
		suffix := path.Ext(checkPath)
		if ex == suffix {
			return true
		}
	}
	return false
}

func isSupportArchiver(checkPath string) bool {
	for _, ex := range SupportFileType {
		suffix := path.Ext(checkPath)
		if ex == suffix {
			return true
		}
	}
	return false
}

func ScanPath(path string) (err error) {
	//var fileList, dirList []string
	var pathList []string
	var bookList []Book
	err = filepath.Walk(path, func(path string, fileInfo os.FileInfo, err error) error {
		//路径深度
		depth := strings.Count(path, "/") - strings.Count(path, "/")
		if depth > Config.MaxDepth {
			return filepath.SkipDir
		}
		if fileInfo == nil {
			return err
		}
		if !isSupportArchiver(path) && !fileInfo.IsDir() {
			return nil
		}
		pathList = append(pathList, path) //文件与路径列表
		return nil
	})
	//分析所有的文件
	for _, f := range pathList {
		//得到书籍的文件数据
		book, err := ScanArchiveOrFolder(f)
		if err != nil {
			fmt.Println(err)
		}
		if book.AllPageNum >= Config.MinImageNum || book.NonUTF8 {
			bookList = append(bookList, *book)
		}
	}
	//最后的所有可用书籍，包括压缩包与文件夹
	BookList = bookList
	return err
}

func ScanDir_InitBook(dirPath string) (*Book, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var book = Book{FilePath: dirPath, IsDir: true, AllPageNum: 0, ExtractComplete: true}
	//初始化，生成UUID
	book.InitBook(book.FilePath)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			// 输出绝对路径
			strAbsPath, errPath := filepath.Abs(dirPath + "/" + file.Name())
			if errPath != nil {
				fmt.Println(errPath)
			}
			//fmt.Println(strAbsPath)
			if isSupportMedia(file.Name()) {
				book.AllPageNum += 1
				book.PageInfo = append(book.PageInfo, SinglePageInfo{RealImageFilePATH: strAbsPath, FileSize: file.Size(), ModeTime: file.ModTime(), NameInArchive: file.Name(), Url: "/cache/" + book.BookID + "/" + url.PathEscape(file.Name())})
			}
		}
	}
	return &book, err
}
