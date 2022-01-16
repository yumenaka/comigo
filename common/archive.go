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
	"os"
	"path"
	"path/filepath"
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
	book := Book{AllPageNum: 0, FilePath: scanPath, Modified: FileInfo.ModTime(), IsDir: FileInfo.IsDir(), FileSize: FileInfo.Size(), ExtractComplete: false}
	//设置书籍UUID，根据路径算出
	book.InitBook(book.FilePath)
	// 获取支持的格式
	iface, err := getFormat(scanPath)
	if err != nil {
		return &book, err
	}
	//判断是否可解压（判断是否可解析： w, ok := iface.(archiver.Walker)）
	_, ok := iface.(archiver.Extractor)
	if !ok {
		logrus.Debugf(locale.GetString("unsupported_extract")+"%s", iface)
		return &book, err
	} else {
		fmt.Println(locale.GetString("scan_ing"), scanPath)
	}

	//统计页数,分析有几张图片
	err = archiver.Walk(scanPath, func(f archiver.File) error {
		inArchiveName := "" //f.Name不包含路径信息.
		////zip编码用
		decodeFileName := ""
		//解压后的路径
		extractFolder := path.Join(CacheFilePath, book.GetBookID())
		//f.Name()不包括路径，inArchiveName需要从f.Header当中获取
		switch h := f.Header.(type) {
		case zip.FileHeader: //Now zip not "archive/zip"
			book.FileType = ".zip"
			inArchiveName = h.Name
			if Config.ZipFilenameEncoding != "" {
				decodeTool := archiver.Zip{FilenameEncoding: Config.ZipFilenameEncoding}
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
			book.FileType = ".tar"
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
			book.FileType = ".rar"
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
		//解压后的文件路径
		imageFilePath := extractFolder + "/" + inArchiveName
		//zip编码的额外处理
		if Config.ZipFilenameEncoding != "" {
			imageFilePath = extractFolder + "/" + decodeFileName
			inArchiveName = decodeFileName
		}
		if !isSupportMedia(inArchiveName) {
			logrus.Debugf(locale.GetString("unsupported_file_type") + inArchiveName)
		} else {
			book.AllPageNum++
			book.PageInfo = append(book.PageInfo, SinglePageInfo{RealImageFilePATH: imageFilePath, FileSize: f.Size(), ModeTime: f.ModTime(), Url: "/cache/" + book.BookID + "/" + inArchiveName})
		}
		return nil
	})
	return &book, err
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
	extraFolder := path.Join(CacheFilePath, b.GetBookID())
	fmt.Println(extraFolder)
	err = u.Unarchive(b.FilePath, extraFolder)
	if err != nil {
		return err
	}
	//解压完成提示
	fmt.Println(locale.GetString("completed_ls"), b.FilePath)
	//ExtractPath = extraFolder
	ReadingBook.ExtractComplete = true
	ReadingBook.ExtractNum = ReadingBook.AllPageNum
	return err
}

func md5string(s string) string {
	r := md5.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

// 解压单个文件
func DecompressionSingleFlie(b *Book, pageNum int) (err error) {
	// 获取支持的格式
	iface, err := getFormat(b.FilePath)
	if err != nil {
		return err
	}
	e, ok := iface.(archiver.Extractor)
	if !ok {
		fmt.Println(locale.GetString("unsupported_extract")+"%s", iface) //这个文件好像没法ls啊
		return err
	}
	extractFolder := path.Join(CacheFilePath, b.GetBookID())
	//fmt.Println(locale.GetString("start_ls"), b.FilePath)
	return e.Extract(b.FilePath, b.PageInfo[pageNum-1].InArchiveName, extractFolder)
}

// GetSingleFile 不解压，只是提取单个文件
func GetSingleFile(b *Book, pageNum int) (err error) {
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

func ScanBookPath(pathname string) (err error) {
	var fileList, dirList []string
	var bookList []Book
	err = filepath.Walk(pathname, func(path string, fileInfo os.FileInfo, err error) error {
		//路径深度
		depth := strings.Count(path, "/") - strings.Count(pathname, "/")
		//windows替换 “/” 为 “\”，去掉
		//if runtime.GOOS == "windows" {
		//	depth = strings.Count(path, "\\") - strings.Count(pathname, "\\")
		//}
		if depth > Config.MaxDepth {
			return filepath.SkipDir
		}
		if fileInfo == nil {
			return err
		}
		if fileInfo.IsDir() {
			dirList = append(dirList, path) //文件夹列表
			return nil
		}
		if !isSupportArchiver(path) {
			return nil
		}
		fileList = append(fileList, path) //文件列表
		return nil
	})
	//分析所有的文件
	for _, f := range fileList {
		//得到书籍的文件数据
		book, err := ScanArchive_InitBook(f)
		if err != nil {
			fmt.Println(err)
		}
		if book.AllPageNum >= Config.MinImageNum {
			bookList = append(bookList, *book)
		}
	}
	//分析所有的文件夹
	for _, f := range dirList {
		//得到书籍的文件数据
		book, err := ScanDir_InitBook(f)
		if err != nil {
			fmt.Println(err)
		}
		if book.AllPageNum >= Config.MinImageNum {
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
				book.PageInfo = append(book.PageInfo, SinglePageInfo{RealImageFilePATH: strAbsPath, FileSize: file.Size(), ModeTime: file.ModTime(), Url: "/cache/" + book.BookID + "/" + file.Name()})
			}
		}
	}
	return &book, err
}
