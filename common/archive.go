package common

import (
	"archive/tar"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/klauspost/compress/zip"
	"github.com/mholt/archiver/v3"
	"github.com/nwaples/rardecode"
	"github.com/sirupsen/logrus"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/tools"
)

var (
	compressionLevel       int
	overwriteExisting      bool
	mkdirAll               bool
	selectiveCompression   bool
	implicitTopLevelFolder bool
	continueOnError        bool
	//filenameEncoding       string
)

func init() {
	mkdirAll = true
	overwriteExisting = false
	continueOnError = true
}
func ScanArchive(scanPath string) (*Book, error) {
	b := Book{AllPageNum: 0, FilePath: scanPath, IsFolder: false, FileSize: 0, ExtractComplete: false}
	// 获取支持的格式
	iface, err := getFormat(scanPath)
	if err != nil {
		return &b, err
	}
	_, ok := iface.(archiver.Extractor)
	if !ok {
		logrus.Debugf(locale.GetString("unsupported_extract")+"%s", iface)
		return &b, err
	} else {
		fmt.Println(locale.GetString("scan_ing"), scanPath)
	}
	err = archiver.Walk(scanPath, func(f archiver.File) error {
		inArchiveName := f.Name()
		if !isSupportMedia(inArchiveName) {
			if inArchiveName != scanPath {
				logrus.Debugf(locale.GetString("unsupported_file_type") + inArchiveName)
			}
		} else {
			b.AllPageNum++
		}
		return nil
	})
	return &b, err
}

////一次解压所有文件，还在测试中，无法正常工作
//func ExtractArchiveOnce(b *Book) (err error) {
//	// 获取支持的格式
//	iface, err := getFormat(b.FilePath)
//	if err != nil {
//		return err
//	}
//	u, ok := iface.(archiver.Unarchiver)
//	if !ok {
//		fmt.Println(locale.GetString("unsupported_extract")+" %s", iface)
//	}
//	if b.FileID == "" {
//		b.FileID = uuid.NewV4().String()
//	}
//	extraFolder := path.Join(TempDir, b.FileID)
//	fmt.Println(extraFolder)
//	err = u.Unarchive(b.FilePath, extraFolder)
//	if err != nil {
//		return err
//	}
//	b, err = ScanDirGetBook(extraFolder)
//	if err != nil {
//		return err
//	}
//	PictureDir = extraFolder
//	ReadingBook.ExtractComplete = true
//	ReadingBook.ExtractNum = ReadingBook.AllPageNum
//	return err
//}

//func md5file(fName string) string {
//	f, e := os.Open(fName)
//	if e != nil {
//		log.Fatal(e)
//	}
//	h := md5.New()
//	_, e = io.Copy(h, f)
//	if e != nil {
//		log.Fatal(e)
//	}
//	return hex.EncodeToString(h.Sum(nil))
//}

func md5string(s string) string {
	r := md5.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

func ExtractArchive(b *Book) (err error) {
	// 获取支持的格式
	iface, err := getFormat(b.FilePath)
	if err != nil {
		return err
	}
	e, ok := iface.(archiver.Extractor)
	if !ok {
		fmt.Println(locale.GetString("unsupported_extract")+"%s", iface)
		return err
	}
	if b.FileID == "" {
		//b.FileID = uuid.NewV4().String()
		b.SetFileID()
	}
	extractFolder := path.Join(TempDir, b.FileID)
	extractNum := 0
	Percent := 0
	tempPercent := 0
	fmt.Println(locale.GetString("start_extract"), b.FilePath)
	//var wg sync.WaitGroup
	err = archiver.Walk(b.FilePath, func(f archiver.File) error {
		//解压用
		inArchiveName := f.Name()
		modeTime := f.ModTime()
		fileSize := f.Size()
		////zip编码用
		//inArchiveNameZip := f.Name()
		switch h := f.Header.(type) {
		case zip.FileHeader: //Now zip not "archive/zip"
			logrus.Debugf("%s\t%d\t%d\t%s\t%s\n",
				f.Mode(),
				h.Method,
				f.Size(),
				f.ModTime(),
				h.Name,
			)
			b.FileType = ".zip"
			inArchiveName = h.Name
		case *tar.Header:
			logrus.Debugf("%s\t%s\t%s\t%d\t%s\t%s\n",
				f.Mode(),
				h.Uname,
				h.Gname,
				f.Size(),
				f.ModTime(),
				h.Name,
			)
			b.FileType = ".tar"
			inArchiveName = h.Name
		case *rardecode.FileHeader:
			logrus.Debugf("%s\t%d\t%d\t%s\t%s\n",
				f.Mode(),
				int(h.HostOS),
				f.Size(),
				f.ModTime(),
				h.Name,
			)
			b.FileType = ".rar"
			inArchiveName = h.Name
		default:
			fmt.Printf("%s\t%d\t%s\t?/%s\n",
				f.Mode(),
				f.Size(),
				f.ModTime(),
				f.Name(), // we don't know full path from this
			)
		}
		if !isSupportMedia(inArchiveName) {
			logrus.Debugf(locale.GetString("unsupported_file_type") + inArchiveName)
			return nil
		}
		//解压后的文件
		filePath := extractFolder + "/" + inArchiveName
		temp := SinglePageInfo{ModeTime: modeTime, FileSize: fileSize, LocalPath: filePath, Name: inArchiveName, Url: "cache/" + inArchiveName}

		//fix bugfix extract single file from zip，not use
		//if path.Ext(b.FilePath) == ".zip" {
		//	filePath = extractFolder + "/" + inArchiveName + "/" + inArchiveName
		//	temp = SinglePageInfo{ModeTime: modeTime, FileSize: fileSize, LocalPath: filePath, Name: inArchiveName, Url: "cache/" + inArchiveName}
		//}
		b.PageInfo = append(b.PageInfo, temp)
		//转义，避免特殊路径造成文件不能读取
		b.PageInfo[len(b.PageInfo)-1].Url = url.PathEscape(b.PageInfo[len(b.PageInfo)-1].Url)
		if tools.ChickFileExists(filePath) {
			logrus.Debugf(locale.GetString("file_exit") + filePath)
		} else {
			//解压文件
			if path.Ext(b.FilePath) == ".zip" {
				err := e.Extract(b.FilePath, inArchiveName, TempDir+"/"+b.FileID) //解压到临时文件夹
				if err != nil {
					logrus.Debugf(err.Error())
				}
			} else {
				err := e.Extract(b.FilePath, inArchiveName, TempDir+"/"+b.FileID) //解压到临时文件夹
				if err != nil {
					logrus.Debugf(err.Error())
				}
			}

		}
		//输出解压比例
		extractNum++
		if b.AllPageNum != 0 {
			Percent = int((float32(extractNum) / float32(b.AllPageNum)) * 100)
			if tempPercent != Percent {
				if (Percent % 10) == 0 { //换个行
					fmt.Println(strconv.Itoa(Percent) + "% ")
				} else {
					if Percent < 10 {
						fmt.Print("0" + strconv.Itoa(Percent) + "% ")
					} else {
						fmt.Print(strconv.Itoa(Percent) + "% ")
					}
				}
			}
			tempPercent = Percent
		}
		//因为有最大打开文件限制，暂不并发解压
		return err
	})
	//wg.Wait()
	fmt.Println(locale.GetString("completed_extract"), b.FilePath)
	return err
}

func getFormat(subcommand string) (interface{}, error) {
	// 通过文件扩展名获取格式
	f, err := archiver.ByExtension(subcommand)
	if err != nil {
		return nil, err
	}
	// 准备一个Tar，以备不时之需
	tar := &archiver.Tar{
		OverwriteExisting:      overwriteExisting,
		MkdirAll:               mkdirAll,
		ImplicitTopLevelFolder: implicitTopLevelFolder,
		ContinueOnError:        continueOnError,
	}
	// fully configure the new value
	switch v := f.(type) {
	case *archiver.Rar:
		v.OverwriteExisting = overwriteExisting
		v.MkdirAll = mkdirAll
		v.ImplicitTopLevelFolder = implicitTopLevelFolder
		v.ContinueOnError = continueOnError
		v.Password = os.Getenv("ARCHIVE_PASSWORD")
	case *archiver.Tar:
		v = tar
	case *archiver.TarBrotli:
		v.Tar = tar
		v.Quality = compressionLevel
	case *archiver.TarBz2:
		v.Tar = tar
		v.CompressionLevel = compressionLevel
	case *archiver.TarGz:
		v.Tar = tar
		v.CompressionLevel = compressionLevel
	case *archiver.TarLz4:
		v.Tar = tar
		v.CompressionLevel = compressionLevel
	case *archiver.TarSz:
		v.Tar = tar
	case *archiver.TarXz:
		v.Tar = tar
	case *archiver.TarZstd:
		v.Tar = tar
	case *archiver.Zip:
		v.CompressionLevel = compressionLevel
		v.OverwriteExisting = overwriteExisting
		v.MkdirAll = mkdirAll
		v.SelectiveCompression = selectiveCompression
		v.ImplicitTopLevelFolder = implicitTopLevelFolder
		v.ContinueOnError = continueOnError
	case *archiver.Gz:
		v.CompressionLevel = compressionLevel
	case *archiver.Brotli:
		v.Quality = compressionLevel
	case *archiver.Bz2:
		v.CompressionLevel = compressionLevel
	case *archiver.Lz4:
		v.CompressionLevel = compressionLevel
	case *archiver.Snappy:
		// nothing to customize
	case *archiver.Xz:
		// nothing to customize
	case *archiver.Zstd:
		// nothing to customize
	default:
		return nil, fmt.Errorf(locale.GetString("format_customization_error")+" %s", f)
	}
	return f, nil
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

//func GetBookPath(scanPath string) (bookPath string, err error) {
//	f, err := os.Stat(scanPath)
//	if err, ok := err.(*os.PathError); ok {
//		fmt.Println("File at path", err.Path, "failed to stat")
//		return bookPath, err
//	}
//	if f.IsDir() == true { //如果是文件夹
//		err := ScanBookPath(scanPath)
//		if err != nil {
//			return bookPath, err
//		}
//		if len(BookList) > 0 {
//			bookPath = BookList[0].FilePath
//		}
//	} else {
//		bookPath = scanPath
//	}
//	return bookPath, err
//}

func ScanBookPath(pathname string) (err error) {
	var fileList, dirList []string
	var bookList []Book
	err = filepath.Walk(pathname, func(path string, fileInfo os.FileInfo, err error) error {
		depth := strings.Count(path, "/") - strings.Count(pathname, "/")
		if runtime.GOOS == "windows" {
			depth = strings.Count(path, "\\") - strings.Count(pathname, "\\")
		}
		if depth > Config.MaxDepth {
			return filepath.SkipDir
		}
		if fileInfo == nil {
			return err
		}
		if fileInfo.IsDir() {
			dirList = append(dirList, path)
			return nil
		}
		if !isSupportArchiver(path) {
			return nil
		}
		fileList = append(fileList, path)
		return nil
	})
	for _, f := range fileList {
		book, err := ScanArchive(f)
		if err != nil {
			fmt.Println(err)
		}
		book.SetArchiveBookName(book.FilePath)
		if book.AllPageNum >= Config.MinImageNum {
			if book.FileID == "" {
				book.SetFileID()
			}
			bookList = append(bookList, *book)
		}
	}
	for _, f := range dirList {
		book, err := ScanDirGetBook(f)
		if err != nil {
			fmt.Println(err)
		}
		book.SetImageFolderBookName(book.FilePath)
		if book.AllPageNum >= Config.MinImageNum {
			if book.FileID == "" {
				book.SetFileID()
			}
			bookList = append(bookList, *book)
		}
	}
	BookList = bookList
	return err
}

func ScanDirGetBook(folder string) (*Book, error) {
	var book = Book{IsFolder: true, AllPageNum: 0, ExtractComplete: true}
	archiveNum := 0
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return &book, err
	}
	for _, file := range files {
		if file.IsDir() {
			//递归处理
			//ScanDirGetBook(folder + "/" + file.Name())
		} else {
			// 输出绝对路径
			strAbsPath, errPath := filepath.Abs(folder + "/" + file.Name())
			if errPath != nil {
				fmt.Println(errPath)
			}
			//fmt.Println(strAbsPath)
			if isSupportMedia(file.Name()) {
				book.AllPageNum += 1
				book.PageInfo = append(book.PageInfo, SinglePageInfo{LocalPath: strAbsPath, Url: "/cache/" + file.Name()})
			}
			if isSupportArchiver(file.Name()) {
				archiveNum += 1
			}
		}
	}
	book.FilePath = folder
	return &book, err
}
