package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/yumenaka/comi/storage"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/klauspost/compress/zip"
	"github.com/yumenaka/comi/tools"

	"github.com/sirupsen/logrus"
	"github.com/yumenaka/archiver/v4"

	"github.com/yumenaka/comi/arch"
	"github.com/yumenaka/comi/book"
	"github.com/yumenaka/comi/locale"
)

// ScanStorePath 3、扫描配置文件指定的的书籍库
func ScanStorePath() error {
	if len(Config.StoresPath) > 0 {
		for _, p := range Config.StoresPath {
			addList, err := ScanAndGetBookList(p, DatabaseBookList)
			if err != nil {
				fmt.Println(locale.GetString("scan_error"), p, err)
				return err
			} else {
				AddBooksToStore(addList, p)
			}
		}
	}
	return nil
}

// SaveResultsToDatabase 4，保存扫描结果到数据库，并清理不存在的书籍
func SaveResultsToDatabase() error {
	AllBook := book.GetAllBookList()
	//设置清理数据库的时候，是否清理没扫描到的书籍信息
	if Config.ClearDatabaseWhenExit {
		for _, checkBook := range DatabaseBookList {
			needClear := true //这条数据是否需要清理
			for _, b := range AllBook {
				if b.BookID == checkBook.BookID {
					needClear = false //如果扫到了这本书,就不清理相关数据
				}
			}
			if needClear {
				storage.ClearBookData(checkBook, Config.Debug)
			}
		}
	}
	saveErr := storage.SaveBookListToDatabase(AllBook)
	if saveErr != nil {
		fmt.Println(saveErr)
		return saveErr
	}
	return nil
}

// AddBooksToStore 添加一组书到书库
func AddBooksToStore(bookList []*book.Book, path string) {
	err := book.AddBooks(bookList, path, Config.MinImageNum)
	if err != nil {
		fmt.Println(locale.GetString("AddBook_error"), path)
	}
	// 然后生成对应的虚拟书籍组
	if err := book.Stores.GenerateBookGroup(); err != nil {
		fmt.Println(err)
	}
}

// ScanAndGetBookList 扫描路径，取得路径里的书籍
func ScanAndGetBookList(storePath string, DatabaseBookList []*book.Book) (newBookList []*book.Book, err error) {
	// 路径不存在的Error，不过目前并不会打印出来
	if !tools.PathExists(storePath) {
		return nil, errors.New(locale.GetString("PATH_NOT_EXIST"))
	}
	storePathAbs, err := filepath.Abs(storePath)
	if err != nil {
		storePathAbs = storePath
		fmt.Println(err)
	}
	fmt.Println(locale.GetString("SCAN_START_HINT") + storePathAbs)
	err = filepath.Walk(storePathAbs, func(walkPath string, fileInfo os.FileInfo, err error) error {
		// 是否需要跳过
		skip := false
		for _, p := range DatabaseBookList {
			AbsW, err := filepath.Abs(walkPath) // 取得绝对路径
			if err != nil {
				// 无法取得的情况下，用相对路径
				AbsW = walkPath
				fmt.Println(err, AbsW)
			}
			if walkPath == p.FilePath || AbsW == p.FilePath {
				skip = true
				newBookList = append(newBookList, p)
			}
		}
		if skip {
			fmt.Println(locale.GetString("FoundInDatabase") + walkPath)
			return nil
		}
		// 路径深度
		depth := strings.Count(walkPath, "/") - strings.Count(storePathAbs, "/")
		if runtime.GOOS == "windows" {
			depth = strings.Count(walkPath, "\\") - strings.Count(storePathAbs, "\\")
		}
		if depth > Config.MaxScanDepth {
			fmt.Printf(locale.GetString("ExceedsMaximumDepth")+" %d，base：%s scan: %s:\n", Config.MaxScanDepth, storePathAbs, walkPath)
			return filepath.SkipDir // 当WalkFunc的返回值是filepath.SkipDir时，Walk将会跳过这个目录，照常执行下一个文件。
		}
		if Config.IsSkipDir(walkPath) {
			fmt.Println(locale.GetString("SkipPath") + walkPath)
			return filepath.SkipDir
		}
		if fileInfo == nil {
			return err
		}
		// 如果不是文件夹
		if !fileInfo.IsDir() {
			if !Config.IsSupportArchiver(walkPath) {
				return nil
			}
			// 得到书籍文件数据
			getBook, err := scanFileGetBook(walkPath, storePathAbs, depth)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			newBookList = append(newBookList, getBook)
		}
		// 如果是文件夹
		if fileInfo.IsDir() {
			// 得到书籍文件数据
			getBook, err := scanDirGetBook(walkPath, storePathAbs, depth)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			newBookList = append(newBookList, getBook)
		}
		return nil
	})
	// 所有可用书籍，包括压缩包与文件夹
	if len(newBookList) > 0 {
		fmt.Printf(locale.GetString("FOUND_IN_PATH"), len(newBookList), storePathAbs)
	}
	return newBookList, err
}

func scanDirGetBook(dirPath string, storePath string, depth int) (*book.Book, error) {
	// 初始化，生成UUID
	newBook, err := book.New(dirPath, time.Now(), 0, storePath, depth, book.TypeDir)
	if err != nil {
		return nil, err
	}
	//// 获取目录中的文件和子目录的详细信息
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	for _, file := range infos {
		// 跳过子目录, 只搜寻目录中的文件
		if file.IsDir() {
			continue
		}
		// 输出绝对路径
		strAbsPath, errPath := filepath.Abs(dirPath + "/" + file.Name())
		if errPath != nil {
			fmt.Println(errPath)
		}
		if Config.IsSupportMedia(file.Name()) {
			TempURL := "api/getfile?id=" + newBook.BookID + "&filename=" + url.QueryEscape(file.Name())
			newBook.Pages.Images = append(newBook.Pages.Images, book.ImageInfo{RealImageFilePATH: strAbsPath, FileSize: file.Size(), ModeTime: file.ModTime(), NameInArchive: file.Name(), Url: TempURL})
		}
	}
	newBook.SortPages("default")
	// 在添加到书库时判断页数
	return newBook, err
}

// 扫描一个路径，并返回对应书籍
func scanFileGetBook(filePath string, storePath string, depth int) (*book.Book, error) {
	// 打开文件
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o400) // Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		fmt.Println(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(file)
	FileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
	}
	// 初始化一本书，设置文件路径等等
	newBook, err := book.New(filePath, FileInfo.ModTime(), FileInfo.Size(), storePath, depth, book.GetBookTypeByFilename(filePath))
	if err != nil {
		return nil, err
	}
	// 根据文件类型，走不同的初始化流程
	switch newBook.Type {
	// 为解决archiver/v4的BUG “zip文件无法读取2级目录” 单独处理zip文件
	case book.TypeZip, book.TypeCbz, book.TypeEpub:
		// 使用Archiver的虚拟文件系统，无法处理非UTF-8编码
		fsys, zipErr := zip.OpenReader(filePath)
		if zipErr != nil {
			// fmt.Println(zipErr)
			return nil, errors.New(locale.GetString("NOT_A_VALID_ZIP_FILE") + filePath)
		}
		err = walkUTF8ZipFs(fsys, "", ".", newBook)
		// 如果扫描ZIP文件的时候遇到了 fs.PathError ，则扫描到NonUTF-8 ZIP文件，需要特殊处理
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			if Config.Debug {
				fmt.Println("NonUTF-8 ZIP:" + filePath + "  Error:" + err.Error())
			}
			// 忽略 fs.PathError 并换个方式扫描
			err = scanNonUTF8ZipFile(filePath, newBook)
		}
		// epub文件，需要根据 META-INF/container.xml 里面定义的rootfile （.opf文件）来重新排序
		if newBook.Type == book.TypeEpub {
			imageList, err := arch.GetImageListFromEpubFile(newBook.FilePath)
			if err != nil {
				fmt.Println(err)
			} else {
				newBook.SortPagesByImageList(imageList)
			}
			// 根据metadata，改写书籍信息
			metaData, err := arch.GetEpubMetadata(newBook.FilePath)
			if err != nil {
				fmt.Println(err)
			} else {
				newBook.Author = metaData.Creator
				newBook.Press = metaData.Publisher
			}
		}
	// TODO:服务器解压速度太慢，网页用PDF.js解析？
	case book.TypePDF:
		newBook.AllPageNum = 1
		newBook.InitComplete = true
		newBook.Cover = book.ImageInfo{RealImageFilePATH: "", FileSize: FileInfo.Size(), ModeTime: FileInfo.ModTime(), NameInArchive: "", Url: "/images/pdf.png"}
	// TODO：简单的网页播放器
	case book.TypeVideo:
		newBook.AllPageNum = 1
		newBook.InitComplete = true
		newBook.Cover = book.ImageInfo{NameInArchive: "video.png", Url: "/images/video.png"}
	case book.TypeAudio:
		newBook.AllPageNum = 1
		newBook.InitComplete = true
		newBook.Cover = book.ImageInfo{NameInArchive: "audio.png", Url: "/images/audio.png"}
	case book.TypeUnknownFile:
		newBook.AllPageNum = 1
		newBook.InitComplete = true
		newBook.Cover = book.ImageInfo{NameInArchive: "unknown.png", Url: "/images/unknown.png"}
	// 其他类型的压缩文件或文件夹
	default:
		// archiver.FileSystem可以配合ctx了，加个默认超时时间
		const shortDuration = 10 * 1000 * time.Millisecond // 超时时间，10秒
		ctx, cancel := context.WithTimeout(context.Background(), shortDuration)
		defer cancel()
		fsys, err := archiver.FileSystem(ctx, filePath)
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
				u, ok := f.(archiver.File) // f.Name不包含路径信息.需要转换一下
				if !ok {
					// 如果是文件夹+图片
					newBook.Type = book.TypeDir
					////用Archiver的虚拟文件系统提供图片文件，理论上现在不应该用到
					//newBook.Pages = append(newBook.Pages, ImageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: "", Url: "/cache/" + newBook.BookID + "/" + url.QueryEscape(path)})
					//实验：用getfile接口提供文件服务
					TempURL := "api/getfile?id=" + newBook.BookID + "&filename=" + url.QueryEscape(path)
					newBook.Pages.Images = append(newBook.Pages.Images, book.ImageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: "", Url: TempURL})
					// fmt.Println(locale.GetString("unsupported_extract")+" %s", f)
				} else {
					// 替换特殊字符的时候，额外将“+替换成"%2b"，因为gin会将+解析为空格。
					TempURL := "api/getfile?id=" + newBook.BookID + "&filename=" + url.QueryEscape(u.NameInArchive)
					// 不替换特殊字符
					// TempURL := "api/getfile?id=" + newBook.BookID + "&filename=" + u.NameInArchive
					newBook.Pages.Images = append(newBook.Pages.Images, book.ImageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: u.NameInArchive, Url: TempURL})
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	// 不管页数，直接返回：在添加到书库时判断页数
	newBook.SortPages("default")
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
			// 如果是压缩文件
			// 替换特殊字符的时候，额外将“+替换成"%2b"，因为gin会将+解析为空格。
			TempURL := "api/getfile?id=" + b.BookID + "&filename=" + url.QueryEscape(f.Name)
			b.Pages.Images = append(b.Pages.Images, book.ImageInfo{RealImageFilePATH: "", FileSize: f.FileInfo().Size(), ModeTime: f.FileInfo().ModTime(), NameInArchive: f.Name, Url: TempURL})
		} else {
			logrus.Debugf(locale.GetString("unsupported_file_type") + f.Name)
		}
	}
	b.SortPages("default")
	return err
}

// 手动写的递归查找，功能与fs.WalkDir()相同。发现一个Archiver/V4的BUG：zip文件的虚拟文件系统，找不到正确的多级文件夹？
// https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter06/06.3.html
func walkUTF8ZipFs(fsys fs.FS, parent, base string, b *book.Book) error {
	// 一般zip文件的处理流程
	// fmt.Println("parent:" + parent + " base:" + base)
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
			// 替换特殊字符的时候,不要用url.PathEscape()，PathEscape不会把“+“替换成"%2b"，会导致BUG，让gin会将+解析为空格。
			b.Pages.Images = append(b.Pages.Images, book.ImageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: inArchiveName, Url: TempURL})
		}
	}
	b.SortPages("default")
	return err
}
