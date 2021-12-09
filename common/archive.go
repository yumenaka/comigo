package common

import (
	"archive/tar"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/klauspost/compress/zip"
	"github.com/nwaples/rardecode"
	"github.com/sirupsen/logrus"
	"github.com/yumenaka/comi/locale"
)

var (
	compressionLevel       int
	overwriteExisting      bool
	mkdirAll               bool
	selectiveCompression   bool
	implicitTopLevelFolder bool
	continueOnError        bool
)

func init() {
	mkdirAll = true
	overwriteExisting = false
	continueOnError = true
}
func ScanArchive_InitBook(scanPath string) (*Book, error) {
	//打开文件
	var file, err = os.OpenFile(scanPath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	FileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
	}
	//设置文件路径等等
	b := Book{AllPageNum: 0, FilePath: scanPath, Modified: FileInfo.ModTime(), IsDir: FileInfo.IsDir(), FileSize: FileInfo.Size(), ExtractComplete: false}
	//设置书籍ID，根据路径算出
	b.InitBook(b.FilePath)
	// 获取支持的格式
	iface, err := getFormat(scanPath)
	if err != nil {
		return &b, err
	}
	//判断是否可解压
	_, ok := iface.(archiver.Extractor)
	if !ok {
		logrus.Debugf(locale.GetString("unsupported_extract")+"%s", iface)
		return &b, err
	} else {
		fmt.Println(locale.GetString("scan_ing"), scanPath)
	}
	//可解压就统计页数
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

// UnArchive 一次解压所有文件
func UnArchive(b *Book) (err error) {
	// 获取支持的格式
	iface, err := getFormat(b.FilePath)
	if err != nil {
		return err
	}
	u, ok := iface.(archiver.Unarchiver)
	if !ok {
		fmt.Println(locale.GetString("unsupported_extract")+" %s", iface)
	}
	extraFolder := path.Join(ComigoCacheFilePath, b.GetBookID())
	fmt.Println(extraFolder)
	err = u.Unarchive(b.FilePath, extraFolder)
	if err != nil {
		return err
	}
	//解压完成提示
	fmt.Println(locale.GetString("completed_ls"), b.FilePath)
	WebImagePath = extraFolder
	ReadingBook.ExtractComplete = true
	ReadingBook.ExtractNum = ReadingBook.AllPageNum
	return err
}

func md5string(s string) string {
	r := md5.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

// LsArchive 分析压缩包文件，不解压
func LsArchive(b *Book) (err error) {
	// 获取支持的格式
	iface, err := getFormat(b.FilePath)
	if err != nil {
		return err
	}
	w, ok := iface.(archiver.Walker)
	if !ok {
		fmt.Println(locale.GetString("unsupported_extract")+"%s", iface) //这个文件好像没法ls啊
		return err
	}
	extractFolder := path.Join(ComigoCacheFilePath, b.GetBookID())
	fmt.Println(locale.GetString("start_ls"), b.FilePath)
	//// Console progress bar
	//bar := pb.StartNew(b.AllPageNum)
	err = w.Walk(b.FilePath, func(f archiver.File) error {
		//解压用
		inArchiveName := f.Name()
		modeTime := f.ModTime()
		fileSize := f.Size()
		////zip编码用
		decodeFileName := ""
		decodeTool := archiver.Zip{FilenameEncoding: Config.ZipFilenameEncoding}
		switch h := f.Header.(type) {
		case zip.FileHeader: //Now zip not "archive/zip"
			b.FileType = ".zip"
			inArchiveName = h.Name
			if Config.ZipFilenameEncoding != "" {
				decodeFileName = decodeTool.DecodeFileName(h)
			}
			logrus.Debugf("%s\t%d\t%d\t%s\t%s\n",
				f.Mode(),
				h.Method,
				f.Size(),
				f.ModTime(),
				h.Name,
			)
		case *tar.Header:
			b.FileType = ".tar"
			inArchiveName = h.Name
			logrus.Debugf("%s\t%s\t%s\t%d\t%s\t%s\n",
				f.Mode(),
				h.Uname,
				h.Gname,
				f.Size(),
				f.ModTime(),
				h.Name,
			)
		case *rardecode.FileHeader:
			b.FileType = ".rar"
			inArchiveName = h.Name
			logrus.Debugf("%s\t%d\t%d\t%s\t%s\n",
				f.Mode(),
				int(h.HostOS),
				f.Size(),
				f.ModTime(),
				h.Name,
			)
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
		//解压后的文件路径
		imageFilePath := extractFolder + "/" + inArchiveName
		temp := SinglePageInfo{ModeTime: modeTime, FileSize: fileSize, ImageFilePATH: imageFilePath, ImageFileName: inArchiveName, Url: "cache/" + inArchiveName}
		//zip编码处理
		if Config.ZipFilenameEncoding != "" {
			imageFilePath = extractFolder + "/" + decodeFileName
			temp.ImageFilePATH = imageFilePath
			temp.ImageFileName = decodeFileName
			temp.Url = "cache/" + decodeFileName
		}
		b.PageInfo = append(b.PageInfo, temp)
		//转义，避免特殊路径造成文件不能读取
		b.PageInfo[len(b.PageInfo)-1].Url = url.PathEscape(b.PageInfo[len(b.PageInfo)-1].Url)
		////进度条计数
		//bar.Increment()
		return err
	})
	//// 进度条跑完
	//bar.Finish()
	fmt.Println(locale.GetString("completed_ls"), b.FilePath)
	return err
}

// GetSingleFile 不解压，只是提取单个文件的
func GetSingleFile(b *Book) (err error) {
	// 参考
	// https://pkg.go.dev/github.com/mholt/archiver/v3#example-Zip-StreamingRead
	iface, err := getFormat(b.FilePath)
	if err != nil {
		return err
	}
	req := new(http.Request)
	contentLen, err := strconv.Atoi(req.Header.Get("Content-Length"))
	if err != nil {
		log.Fatal(err)
	}
	reader, ok := iface.(archiver.Reader)
	if !ok {
		fmt.Println("unsupported_getSingleFile", iface)
		return err
	}
	//Zip 格式需要知道流的长度，但其他格式通常不需要它，因此在使用它们时可以将其保留为 0
	err = reader.Open(req.Body, int64(contentLen))
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	return err
}

func getFormat(subcommand string) (interface{}, error) {
	// 通过文件扩展名获取格式
	f, err := archiver.ByExtension(subcommand)
	if err != nil {
		return nil, err
	}
	// 准备一个Tar，下面要用到
	tarball := &archiver.Tar{
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
		v = tarball
	case *archiver.TarBrotli:
		v.Tar = tarball
		v.Quality = compressionLevel
	case *archiver.TarBz2:
		v.Tar = tarball
		v.CompressionLevel = compressionLevel
	case *archiver.TarGz:
		v.Tar = tarball
		v.CompressionLevel = compressionLevel
	case *archiver.TarLz4:
		v.Tar = tarball
		v.CompressionLevel = compressionLevel
	case *archiver.TarSz:
		v.Tar = tarball
	case *archiver.TarXz:
		v.Tar = tarball
	case *archiver.TarZstd:
		v.Tar = tarball
	case *archiver.Zip:
		v.CompressionLevel = compressionLevel
		v.OverwriteExisting = overwriteExisting
		v.MkdirAll = mkdirAll
		v.SelectiveCompression = selectiveCompression
		v.ImplicitTopLevelFolder = implicitTopLevelFolder
		v.ContinueOnError = continueOnError
		v.FilenameEncoding = Config.ZipFilenameEncoding
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
		book, err := ScanArchive_InitBook(f)
		if err != nil {
			fmt.Println(err)
		}

		if book.AllPageNum >= Config.MinImageNum {
			bookList = append(bookList, *book)
		}
	}
	for _, f := range dirList {
		book, err := ScanDirGetBook(f)
		if err != nil {
			fmt.Println(err)
		}
		if book.AllPageNum >= Config.MinImageNum {
			bookList = append(bookList, *book)
		}
	}
	BookList = bookList
	return err
}

func ScanDirGetBook(filePath string) (*Book, error) {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		return nil, err
	}
	var book = Book{FilePath: filePath, IsDir: true, AllPageNum: 0, ExtractComplete: true}
	book.InitBook(book.FilePath)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			// 输出绝对路径
			strAbsPath, errPath := filepath.Abs(filePath + "/" + file.Name())
			if errPath != nil {
				fmt.Println(errPath)
			}
			//fmt.Println(strAbsPath)
			if isSupportMedia(file.Name()) {
				book.AllPageNum += 1
				book.PageInfo = append(book.PageInfo, SinglePageInfo{ImageFilePATH: strAbsPath, Url: "/cache/" + file.Name()})
			}
		}
	}
	return &book, err
}
