package common

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	archiver "github.com/yumenaka/archiver/v4"

	"github.com/yumenaka/comi/arch"
	"github.com/yumenaka/comi/book"
	"github.com/yumenaka/comi/locale"
)

func AddBooksToStore(bookList []*book.Book, path string) {
	err := book.AddBooks(bookList, path, Config.MinImageNum)
	if err != nil {
		fmt.Println(locale.GetString("AddBook_error"), path)
	}
	//然后生成对应的虚拟书籍组
	if err := book.Stores.GenerateBookGroup(); err != nil {
		fmt.Println(err)
	}
}

// ScanAndGetBookList 扫描一个路径，并返回书籍列表
func ScanAndGetBookList(storePath string, skipPathList []string) (bookList []*book.Book, err error) {
	storePathAbs, err := filepath.Abs(storePath)
	if err != nil {
		storePathAbs = storePath
		fmt.Println(err)
	}
	err = filepath.Walk(storePathAbs, func(walkPath string, fileInfo os.FileInfo, err error) error {
		//是否需要跳过
		skip := false
		for _, p := range skipPathList {
			AbsW, err := filepath.Abs(walkPath) //取得绝对路径
			if err != nil {
				//因为权限问题，无法的情况下，用相对路径
				fmt.Println(err, AbsW)
			}
			if walkPath == p || AbsW == p {
				skip = true
			}
		}
		if skip {
			fmt.Println("Found in Database,Skip scan:" + walkPath)
			return nil
		}
		//路径深度
		depth := strings.Count(walkPath, "/") - strings.Count(storePathAbs, "/")
		if runtime.GOOS == "windows" {
			depth = strings.Count(walkPath, "\\") - strings.Count(storePathAbs, "\\")
		}
		if depth > Config.MaxDepth {
			fmt.Printf("超过最大搜索深度 %d，base：%s scan: %s:\n", Config.MaxDepth, storePathAbs, walkPath)
			return filepath.SkipDir //当WalkFunc的返回值是filepath.SkipDir时，Walk将会跳过这个目录，照常执行下一个文件。
		}
		if Config.IsSkipDir(walkPath) {
			fmt.Println("Skip Scan:" + walkPath)
			return filepath.SkipDir
		}
		if fileInfo == nil {
			return err
		}

		//如果不是文件夹
		if !fileInfo.IsDir() {
			if !Config.IsSupportArchiver(walkPath) {
				return nil
			}
			//得到书籍文件数据
			getBook, err := scanFileGetBook(walkPath, storePathAbs, depth)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			bookList = append(bookList, getBook)
		}

		//如果是文件夹
		if fileInfo.IsDir() {
			//得到书籍文件数据
			getBook, err := scanDirGetBook(walkPath, storePathAbs, depth)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			bookList = append(bookList, getBook)
		}
		return nil
	})
	//所有可用书籍，包括压缩包与文件夹
	return bookList, err
}

func scanDirGetBook(dirPath string, storePath string, depth int) (*book.Book, error) {
	//初始化，生成UUID
	newBook := book.New(dirPath, time.Now(), 0, storePath, depth)
	newBook.Type = book.TypeDir
	// 目录中的文件和子目录
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		// 跳过子目录, 只搜寻目录中的文件
		if file.IsDir() {
			continue
		}
		// 输出绝对路径
		strAbsPath, errPath := filepath.Abs(dirPath + "/" + file.Name())
		if errPath != nil {
			fmt.Println(errPath)
		}
		//fmt.Println(strAbsPath)
		if Config.IsSupportMedia(file.Name()) {
			TempURL := "api/getfile?id=" + newBook.BookID + "&filename=" + url.QueryEscape(file.Name())
			newBook.Pages = append(newBook.Pages, book.SinglePageInfo{RealImageFilePATH: strAbsPath, FileSize: file.Size(), ModeTime: file.ModTime(), NameInArchive: file.Name(), Url: TempURL})
		}
	}

	//不管页数，直接返回：在添加到书库时判断页数
	return newBook, err
}

// 扫描一个路径，并返回对应书籍
func scanFileGetBook(filePath string, storePath string, depth int) (*book.Book, error) {
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
	newBook := book.New(filePath, FileInfo.ModTime(), FileInfo.Size(), storePath, depth)
	//根据文件类型，走不同的初始化流程
	switch newBook.Type {
	//为解决archiver/v4的BUG “zip文件无法读取2级目录” 单独处理zip文件
	case book.TypeZip, book.TypeCbz, book.TypeEpub:
		//使用Archiver的虚拟文件系统，无法处理非UTF-8编码
		fsys, zipErr := zip.OpenReader(filePath)
		if zipErr != nil {
			fmt.Println(zipErr)
			return nil, errors.New("not a valid zip file:" + filePath)
		}
		err = walkUTF8ZipFs(fsys, "", ".", newBook)
		//如果扫描ZIP文件的时候遇到了 fs.PathError ，则扫描到NonUTF-8 ZIP文件，需要特殊处理
		if _, ok := err.(*fs.PathError); ok {
			if Config.Debug {
				fmt.Println("NonUTF-8 ZIP:" + filePath)
				fmt.Println("  Error:" + err.Error())
			}
			//忽略 fs.PathError 并换个方式扫描
			err = scanNonUTF8ZipFile(filePath, newBook)
		}
	//TODO:服务器解压速度太慢，网页用PDF.js解析？
	case book.TypePDF:
		newBook.AllPageNum = 1
		newBook.InitComplete = true
		newBook.Cover = book.SinglePageInfo{RealImageFilePATH: "", FileSize: FileInfo.Size(), ModeTime: FileInfo.ModTime(), NameInArchive: "", Url: "/images/pdf.png"}
		//pageCount, err := arch.CountPagesOfPDF(newBook.FilePath)
		//if err != nil {
		//	return nil, errors.New("PDF Error：" + newBook.FilePath)
		//}
		//newBook.AllPageNum = pageCount
		//for i := 0; i < pageCount; i++ {
		//	imageUrl := "api/get_pdf_image?id=" + newBook.BookID + "&filename=" + strconv.Itoa(i+1) + ".jpg"
		//	newBook.Pages = append(newBook.Pages, book.SinglePageInfo{RealImageFilePATH: "", FileSize: int64(i + 1), ModeTime: FileInfo.ModTime(), NameInArchive: "", Url: imageUrl})
		//}
		//newBook.Cover = newBook.Pages[0]
		//newBook.Cover.Url = "/images/pdf.png"
	//TODO：简单的网页播放器
	case book.TypeVideo:
		newBook.AllPageNum = 1
		newBook.InitComplete = true
		newBook.Cover = book.SinglePageInfo{NameInArchive: "video.png", Url: "/images/video.png"}
	//其他类型的压缩文件或文件夹
	default:
		fsys, err := archiver.FileSystem(filePath)
		if err != nil {
			return nil, err
		}
		err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if Config.IsSkipDir(path) {
				fmt.Println("Skip Scan:" + path)
				return fs.SkipDir
			}
			f, errInfo := d.Info()
			if errInfo != nil {
				fmt.Println(errInfo)
				return fs.SkipDir
			}
			if !Config.IsSupportMedia(path) {
				logrus.Debugf(locale.GetString("unsupported_file_type") + path)
			} else {
				u, ok := f.(archiver.File) //f.Name不包含路径信息.需要转换一下
				if !ok {
					//如果是文件夹+图片
					newBook.Type = book.TypeDir
					////用Archiver的虚拟文件系统提供图片文件，理论上现在不应该用到
					//newBook.Pages = append(newBook.Pages, SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: "", Url: "/cache/" + newBook.BookID + "/" + url.QueryEscape(path)})
					//实验：用getfile接口提供文件服务
					TempURL := "api/getfile?id=" + newBook.BookID + "&filename=" + url.QueryEscape(path)
					newBook.Pages = append(newBook.Pages, book.SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: "", Url: TempURL})
					//fmt.Println(locale.GetString("unsupported_extract")+" %s", f)
				} else {
					//替换特殊字符的时候，额外将“+替换成"%2b"，因为gin会将+解析为空格。
					TempURL := "api/getfile?id=" + newBook.BookID + "&filename=" + url.QueryEscape(u.NameInArchive)
					//不替换特殊字符
					//TempURL := "api/getfile?id=" + newBook.BookID + "&filename=" + u.NameInArchive
					newBook.Pages = append(newBook.Pages, book.SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: u.NameInArchive, Url: TempURL})
				}
			}
			return nil
		})
	}
	//不管页数，直接返回：在添加到书库时判断页数
	return newBook, err
}

func scanNonUTF8ZipFile(filePath string, b *book.Book) error {
	b.NonUTF8Zip = true
	reader, err := arch.ScanNonUTF8Zip(filePath, Config.ZipFileTextEncoding)
	if err != nil {
		return err
	}
	for _, f := range reader.File {
		if Config.IsSupportMedia(f.Name) {
			//如果是压缩文件
			//替换特殊字符的时候，额外将“+替换成"%2b"，因为gin会将+解析为空格。
			TempURL := "api/getfile?id=" + b.BookID + "&filename=" + url.QueryEscape(f.Name)
			b.Pages = append(b.Pages, book.SinglePageInfo{RealImageFilePATH: "", FileSize: f.FileInfo().Size(), ModeTime: f.FileInfo().ModTime(), NameInArchive: f.Name, Url: TempURL})
		} else {
			logrus.Debugf(locale.GetString("unsupported_file_type") + f.Name)
		}
	}
	return err
}

//手动写的递归查找，功能与fs.WalkDir()相同。发现一个Archiver/V4的BUG：zip文件的虚拟文件系统，找不到正确的多级文件夹？
// https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter06/06.3.html
func walkUTF8ZipFs(fsys fs.FS, parent, base string, b *book.Book) error {
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
			err = walkUTF8ZipFs(fsys, joinPath, base, b)
		} else if !Config.IsSupportMedia(name) {
			logrus.Debugf(locale.GetString("unsupported_file_type") + name)
		} else {
			inArchiveName := path.Join(parent, f.Name())
			TempURL := "api/getfile?id=" + b.BookID + "&filename=" + url.QueryEscape(inArchiveName)
			//替换特殊字符的时候,不要用url.PathEscape()，PathEscape不会把“+“替换成"%2b"，会导致BUG，让gin会将+解析为空格。
			b.Pages = append(b.Pages, book.SinglePageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: inArchiveName, Url: TempURL})
		}
	}
	return err
}
