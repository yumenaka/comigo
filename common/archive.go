package common

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	archiver "github.com/mholt/archiver/v3"
	"github.com/nwaples/rardecode"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
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
	b := Book{PageNum: 0, FilePath: scanPath, IsFolder: false, FileSize: 0, ExtractComplete: false}
	// 获取支持的格式
	iface, err := getFormat(scanPath)
	if err != nil {
		return &b, err
	}
	_, ok := iface.(archiver.Extractor)
	if !ok {
		fmt.Println("不支持解压：%s", iface)
		return &b, err
	}
	err = archiver.Walk(scanPath, func(f archiver.File) error {
		inArchiveName := f.Name()

		if !checkPicExt(inArchiveName) {
			if inArchiveName != scanPath {
				fmt.Println("不支持：" + inArchiveName)
			}
		} else {
			b.PageNum++
		}
		return nil
	})
	fmt.Println("正在扫描文件", scanPath)
	return &b, err
}

//一次解压所有文件，还在测试中，无法正常工作
func ExtractArchiveOnce(b *Book) (err error) {
	// 获取支持的格式
	iface, err := getFormat(b.FilePath)
	if err != nil {
		return err
	}
	u, ok := iface.(archiver.Unarchiver)
	if !ok {
		fmt.Println("不支持解压缩： %s", iface)
	}
	//b.FileType = ".zip"
	if b.UUID == "" {
		b.UUID = uuid.NewV4().String()
	}
	extraFolder := path.Join(TempDir, b.UUID)
	fmt.Println(extraFolder)
	err = u.Unarchive(b.FilePath, extraFolder)
	if err != nil {
		return err
	}
	b, err = ScanDirGetBook(extraFolder)
	if err != nil {
		return err
	}
	PictureDir = extraFolder
	ReadingBook.ExtractComplete = true
	ReadingBook.ExtractNum = ReadingBook.PageNum
	return err
}

func ExtractArchive(b *Book) (err error) {
	// 获取支持的格式
	iface, err := getFormat(b.FilePath)
	if err != nil {
		return err
	}
	e, ok := iface.(archiver.Extractor)
	if !ok {
		fmt.Println("不支持解压缩：%s", iface)
		return err
	}
	if b.UUID == "" {
		b.UUID = uuid.NewV4().String()
	}
	extraFolder := path.Join(TempDir, b.UUID)
	extractNum := 0
	//var wg sync.WaitGroup
	err = archiver.Walk(b.FilePath, func(f archiver.File) error {
		//解压用
		inArchiveName := f.Name()
		////zip编码用
		//inArchiveNameZip := f.Name()
		switch h := f.Header.(type) {
		case zip.FileHeader:
			fmt.Printf("%s\t%d\t%d\t%s\t%s\n",
				f.Mode(),
				h.Method,
				f.Size(),
				f.ModTime(),
				h.Name,
			)
			b.FileType = ".zip"
			inArchiveName = h.Name
			////手动指定zip编码
			//if Config.ZipFilenameEncoding != "" {
			//	inArchiveNameZip = DecodeFileName(h.Name)
			//}
		case *tar.Header:
			fmt.Printf("%s\t%s\t%s\t%d\t%s\t%s\n",
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
			fmt.Printf("%s\t%d\t%d\t%s\t%s\n",
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
		if !checkPicExt(inArchiveName) {
			fmt.Println("不支持的格式：" + inArchiveName)
			return nil
		}
		extractNum++
		//解压后的文件
		filePath := extraFolder + "/" + inArchiveName
		temp := ImageInfo{LocalPath: filePath, InArchiveName: inArchiveName, UrlPath: "cache/" + b.UUID + "/" + inArchiveName}
		if b.FileType == ".zip" {
			filePath = extraFolder + "/" + inArchiveName + "/" + inArchiveName
			temp = ImageInfo{LocalPath: filePath, UrlPath: "cache/" + b.UUID + "/" + inArchiveName + "/" + inArchiveName}
		}
		if ChickFileExists(filePath) {
			fmt.Println("文件已存在，跳过解压步骤：" + filePath)
			return err
		}
		b.PageInfo = append(b.PageInfo, temp)
		//转义，避免特殊路径造成文件不能读取
		b.PageInfo[len(b.PageInfo)-1].UrlPath = url.PathEscape(b.PageInfo[len(b.PageInfo)-1].UrlPath)
		//解压文件
		err := e.Extract(b.FilePath, inArchiveName, TempDir+"/"+b.UUID) //解压到临时文件夹
		if err != nil {
			fmt.Println(err)
		}
		//因为有最大打开文件限制，暂不并发解压
		return err
	})
	//wg.Wait()
	return err
}

func getFormat(subcommand string) (interface{}, error) {
	// 通过文件扩展名获取格式
	f, err := archiver.ByExtension(subcommand)
	if err != nil {
		return nil, err
	}
	// 准备一个Tar，以备不时之需
	mytar := &archiver.Tar{
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
		v = mytar
	case *archiver.TarBrotli:
		v.Tar = mytar
		v.Quality = compressionLevel
	case *archiver.TarBz2:
		v.Tar = mytar
		v.CompressionLevel = compressionLevel
	case *archiver.TarGz:
		v.Tar = mytar
		v.CompressionLevel = compressionLevel
	case *archiver.TarLz4:
		v.Tar = mytar
		v.CompressionLevel = compressionLevel
	case *archiver.TarSz:
		v.Tar = mytar
	case *archiver.TarXz:
		v.Tar = mytar
	case *archiver.TarZstd:
		v.Tar = mytar
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
		return nil, fmt.Errorf("format does not support customization: %s", f)
	}
	return f, nil
}

func checkPicExt(checkPath string) bool {
	for _, ex := range SupportPicType {
		filesuffix := path.Ext(checkPath)
		if ex == filesuffix {
			return true
		}
	}
	return false
}

func checkArchiverExt(checkPath string) bool {
	for _, ex := range SupportFileType {
		filesuffix := path.Ext(checkPath)
		if ex == filesuffix {
			return true
		}
	}
	return false
}
func GetBookPath(scanPath string) (bookPath string, err error) {
	f, err := os.Stat(scanPath)
	if err, ok := err.(*os.PathError); ok {
		fmt.Println("File at path", err.Path, "failed to stat")
		return bookPath, err
	}
	if f.IsDir() == true { //如果是文件夹
		err := ScanBookPath(scanPath)
		if err != nil {
			return bookPath, err
		}
		if len(BookList) > 0 {
			bookPath = BookList[0].FilePath
		}
	} else {
		bookPath = scanPath
	}
	return bookPath, err
}

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
		if !checkArchiverExt(path) {
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
		if book.PageNum > Config.MinImageNum {
			if book.UUID == "" {
				book.UUID = uuid.NewV4().String()
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
		if book.PageNum >= Config.MinImageNum {
			if book.UUID == "" {
				book.UUID = uuid.NewV4().String()
			}
			bookList = append(bookList, *book)
		}
	}
	BookList = bookList
	return err
}

func ScanDirGetBook(folder string) (*Book, error) {
	var book = Book{IsFolder: true, PageNum: 0, ExtractComplete: true}
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
			fmt.Println(strAbsPath)
			if checkPicExt(file.Name()) {
				book.PageNum += 1
				book.PageInfo = append(book.PageInfo, ImageInfo{LocalPath: strAbsPath, UrlPath: "/cache/" + file.Name()})
			}
			if checkArchiverExt(file.Name()) {
				archiveNum += 1
			}
		}
	}
	book.FilePath = folder
	return &book, err
}
