package common

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/mholt/archiver/v4"
	"github.com/sirupsen/logrus"
	"github.com/yumenaka/comi/arch"
	"github.com/yumenaka/comi/locale"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// ScanAndGetBookList 扫描一个路径，并返回书籍列表
func ScanAndGetBookList(path string) (bookList []*Book, err error) {
	//var fileList, dirList []string
	var pathList []string
	err = filepath.Walk(path, func(path string, fileInfo os.FileInfo, err error) error {
		//路径深度
		depth := strings.Count(path, "/") - strings.Count(path, "/")
		if depth > Config.MaxDepth {
			fmt.Println("超过最大搜索深度，path:" + path)
			return filepath.SkipDir //当WalkFunc的返回值是filepath.SkipDir时，Walk将会跳过这个目录，照常执行下一个文件。
		}
		if CheckPathSkip(path) {
			fmt.Println("Skip Scan:" + path)
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
		//得到书籍文件数据
		book, err := scanAndGetBook(f)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//多级路径的图片文件夹，避免重复
		if !book.IsDir {
			bookList = append(bookList, book)
		}
		//多级路径的图片文件夹，避免重复添加
		if book.IsDir {
			added := false
			for _, b := range bookList {
				if strings.HasPrefix(book.GetFilePath(), b.GetFilePath()) {
					added = true
				}
			}
			if !added {
				bookList = append(bookList, book)
			}
		}
	}
	//所有可用书籍，包括压缩包与文件夹
	return bookList, err
}

// 扫描一个路径，并返回对应书籍
func scanAndGetBook(filePath string) (*Book, error) {
	//打开文件
	var file, err = os.OpenFile(filePath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	FileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
	}
	//初始化一本书，设置文件路径等等
	book := InitBook(0, filePath, FileInfo.ModTime(), FileInfo.IsDir(), FileInfo.Size(), false)
	ext := path.Ext(filePath)
	//为解决archiver/v4的BUG “zip文件无法读取2级目录” 单独处理zip文件
	if ext == ".zip" || ext == ".cbz" || ext == ".epub" {
		//使用Archiver的虚拟文件系统，无法处理非UTF-8编码
		fsys, zipErr := zip.OpenReader(filePath)
		if zipErr != nil {
			fmt.Println(zipErr)
		}
		err = walkUTF8ZipFs(fsys, "", ".", book)
		//如果扫描ZIP文件的时候遇到了 fs.PathError ，则扫描到NonUTF-8 ZIP文件，需要特殊处理
		if _, ok := err.(*fs.PathError); ok {
			//忽略 fs.PathError 并换个方式扫描
			fmt.Println("扫描到NonUTF-8 ZIP文件:" + filePath + "  Error:" + err.Error())
			err = scanNonUTF8ZipFile(filePath, book)
		}
	} else {
		//其他类型的压缩文件或文件夹
		fsys, err := archiver.FileSystem(filePath)
		if err != nil {
			return nil, err
		}
		err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if CheckPathSkip(path) {
				fmt.Println("Skip Scan:" + path)
				return fs.SkipDir
			}
			f, errInfo := d.Info()
			if errInfo != nil {
				fmt.Println(errInfo)
				return fs.SkipDir
			}
			if !isSupportMedia(path) {
				logrus.Debugf(locale.GetString("unsupported_file_type") + path)
			} else {
				u, ok := f.(archiver.File) //f.Name不包含路径信息.需要转换一下
				if !ok {
					//如果是文件夹中的图片
					////用Archiver的虚拟文件系统提供图片文件
					//book.Pages = append(book.Pages, SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: "", Url: "/cache/" + book.BookID + "/" + url.PathEscape(path)})

					//实验：用getfile接口提供文件服务
					TempURL := "api/getfile?id=" + book.BookID + "&filename=" + url.PathEscape(path)
					book.Pages = append(book.Pages, SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: "", Url: TempURL})
					//fmt.Println(locale.GetString("unsupported_extract")+" %s", f)
				} else {
					////用Archiver的虚拟文件系统提供图片文件
					//TempURL := "/cache/" + book.BookID + "/" + url.PathEscape(u.NameInArchive)
					//book.Pages = append(book.Pages, SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: u.NameInArchive, Url: TempURL})

					//实验：用getfile接口提供提供图片文件
					TempURL := "api/getfile?id=" + book.BookID + "&filename=" + url.PathEscape(u.NameInArchive)
					book.Pages = append(book.Pages, SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: u.NameInArchive, Url: TempURL})
				}
			}
			return nil
		})
	}
	//根据文件数决定是否返回这本书
	totalPageHint := "filePath: " + filePath + " Total number of pages in the book:" + strconv.Itoa(book.GetAllPageNum())
	if book.GetAllPageNum() >= Config.MinImageNum {
		fmt.Println(totalPageHint)
		return book, err
	} else {
		return nil, errors.New(totalPageHint)
	}
}

func scanNonUTF8ZipFile(filePath string, book *Book) error {
	book.NonUTF8Zip = true
	reader, err := arch.ScanNonUTF8Zip(filePath, Config.ZipFileTextEncoding)
	if err != nil {
		return err
	}
	for _, f := range reader.File {
		if isSupportMedia(f.Name) {
			//如果是压缩文件
			TempURL := "api/getfile?id=" + book.BookID + "&filename=" + url.PathEscape(f.Name)
			book.Pages = append(book.Pages, SinglePageInfo{RealImageFilePATH: "", FileSize: f.FileInfo().Size(), ModeTime: f.FileInfo().ModTime(), NameInArchive: f.Name, Url: TempURL})
		} else {
			logrus.Debugf(locale.GetString("unsupported_file_type") + f.Name)
		}
	}
	return err
}

//手动写的递归查找，功能与fs.WalkDir()相同。发现一个Archiver/V4的BUG：zip文件的虚拟文件系统，找不到正确的多级文件夹？
// https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter06/06.3.html
func walkUTF8ZipFs(fsys fs.FS, parent, base string, book *Book) error {
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
			joinPath := path.Join(parent, name)
			err = walkUTF8ZipFs(fsys, joinPath, base, book)
		} else if !isSupportMedia(name) {
			logrus.Debugf(locale.GetString("unsupported_file_type") + name)
		} else {
			inArchiveName := path.Join(parent, f.Name())
			//book.Pages = append(book.Pages, SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: inArchiveName, Url: "/cache/" + book.BookID + "/" + url.PathEscape(inArchiveName)})
			TempURL := "api/getfile?id=" + book.BookID + "&filename=" + url.PathEscape(inArchiveName)
			book.Pages = append(book.Pages, SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: inArchiveName, Url: TempURL})
		}
	}
	return err
}

func isSupportMedia(checkPath string) bool {
	for _, ex := range Config.SupportMediaType {
		//strings.ToLower():某些文件会用大写文件名
		suffix := strings.ToLower(path.Ext(checkPath))
		if ex == suffix {
			return true
		}
	}
	return false
}

func isSupportArchiver(checkPath string) bool {
	for _, ex := range Config.SupportFileType {
		suffix := path.Ext(checkPath)
		if ex == suffix {
			return true
		}
	}
	return false
}

//func ScandirInitbook(dirPath string) (*Book, error) {
//	files, err := ioutil.ReadDir(dirPath)
//	if err != nil {
//		return nil, err
//	}
//	//初始化，生成UUID
//	book := InitBook(0, dirPath, time.Now(), true, 0, true)
//	for _, file := range files {
//		if file.IsDir() {
//			continue
//		} else {
//			// 输出绝对路径
//			strAbsPath, errPath := filepath.Abs(dirPath + "/" + file.Name())
//			if errPath != nil {
//				fmt.Println(errPath)
//			}
//			//fmt.Println(strAbsPath)
//			if isSupportMedia(file.Name()) {
//				book.Pages = append(book.Pages, SinglePageInfo{RealImageFilePATH: strAbsPath, FileSize: file.Size(), ModeTime: file.ModTime(), NameInArchive: file.Name(), Url: "/cache/" + book.BookID + "/" + url.PathEscape(file.Name())})
//			}
//		}
//	}
//	//根据文件数决定是否返回这本书
//	totalPageHint := "filePath: " + dirPath + " Total number of pages in the book:" + strconv.Itoa(book.GetAllPageNum())
//	if book.GetAllPageNum() >= Config.MinImageNum {
//		fmt.Println(totalPageHint)
//		return book, err
//	} else {
//		return nil, errors.New(totalPageHint)
//	}
//}
