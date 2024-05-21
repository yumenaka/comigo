package scan

import (
	"context"
	"errors"
	"github.com/yumenaka/comi/util"
	fileutil "github.com/yumenaka/comi/util/file"
	"github.com/yumenaka/comi/util/locale"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/klauspost/compress/zip"
	"github.com/yumenaka/archiver/v4"
	"github.com/yumenaka/comi/entity"
	"github.com/yumenaka/comi/logger"
)

// Local 扫描路径，取得路径里的书籍
func Local(storePath string, scanOption Option) (newBookList []*entity.Book, err error) {
	// 路径不存在的Error，不过目前并不会打印出来
	if !util.PathExists(storePath) {
		return nil, errors.New(locale.GetString("PATH_NOT_EXIST"))
	}
	storePathAbs, err := filepath.Abs(storePath)
	if err != nil {
		storePathAbs = storePath
		logger.Infof("%s", err)
	}
	logger.Infof(locale.GetString("SCAN_START_HINT")+"%s", storePathAbs)
	err = filepath.Walk(storePathAbs, func(walkPath string, fileInfo os.FileInfo, err error) error {
		if !scanOption.ReScanFile {
			for _, p := range entity.GetArchiveBooks() {
				AbsW, err := filepath.Abs(walkPath) // 取得绝对路径
				if err != nil {
					// 无法取得的情况下，用相对路径
					AbsW = walkPath
					logger.Info(err, AbsW)
				}
				if walkPath == p.FilePath || AbsW == p.FilePath {
					//跳过已经在数据库里面的文件
					logger.Infof(locale.GetString("FoundInDatabase")+"%s", walkPath)
					return nil
				}
			}
		}
		// 路径深度
		depth := strings.Count(walkPath, "/") - strings.Count(storePathAbs, "/")
		if runtime.GOOS == "windows" {
			depth = strings.Count(walkPath, "\\") - strings.Count(storePathAbs, "\\")
		}
		if depth > scanOption.MaxScanDepth {
			logger.Infof(locale.GetString("ExceedsMaximumDepth")+" %d，base：%s scan: %s:", scanOption.MaxScanDepth, storePathAbs, walkPath)
			return filepath.SkipDir // 当WalkFunc的返回值是filepath.SkipDir时，Walk将会跳过这个目录，照常执行下一个文件。
		}
		if scanOption.IsSkipDir(walkPath) {
			logger.Infof(locale.GetString("SkipPath")+"%s", walkPath)
			return filepath.SkipDir
		}
		if fileInfo == nil {
			return err
		}
		// 如果不是文件夹
		if !fileInfo.IsDir() {
			if !scanOption.IsSupportArchiver(walkPath) {
				return nil
			}
			// 得到书籍文件数据
			getBook, err := scanFileGetBook(walkPath, storePathAbs, depth, scanOption)
			if err != nil {
				logger.Infof("%s", err)
				return nil
			}
			newBookList = append(newBookList, getBook)
		}
		// 如果是文件夹
		if fileInfo.IsDir() {
			// 得到书籍文件数据
			getBook, err := scanDirGetBook(walkPath, storePathAbs, depth, scanOption)
			if err != nil {
				logger.Infof("%s", err)
				return nil
			}
			newBookList = append(newBookList, getBook)
		}
		return nil
	})
	// 所有可用书籍，包括压缩包与文件夹
	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("FOUND_IN_PATH"), len(newBookList), storePathAbs)
	}
	return newBookList, err
}

// 扫描本地路径，并返回对应书籍
func scanFileGetBook(filePath string, storePath string, depth int, scanOption Option) (*entity.Book, error) {
	// 打开文件
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o400) // Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Infof("%s", err.Error())
		}
	}(file)
	FileInfo, err := file.Stat()
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	// 初始化一本书，设置文件路径等等
	newBook, err := entity.NewBook(filePath, FileInfo.ModTime(), FileInfo.Size(), storePath, depth, entity.GetBookTypeByFilename(filePath))
	if err != nil {
		return nil, err
	}
	// 根据文件类型，走不同的初始化流程
	switch newBook.Type {
	// 为解决archiver/v4的BUG “zip文件无法读取2级目录” 单独处理zip文件
	case entity.TypeZip, entity.TypeCbz, entity.TypeEpub:
		// 使用Archiver的虚拟文件系统，无法处理非UTF-8编码
		fsys, zipErr := zip.OpenReader(filePath)
		if zipErr != nil {
			// logger.Infof(zipErr)
			return nil, errors.New(locale.GetString("NOT_A_VALID_ZIP_FILE") + filePath)
		}
		err = walkUTF8ZipFs(fsys, "", ".", newBook, scanOption)
		// 如果扫描ZIP文件的时候遇到了 fs.PathError ，则扫描到NonUTF-8 ZIP文件，需要特殊处理
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			if scanOption.Debug {
				logger.Infof("NonUTF-8 ZIP:%s  Error:%s", filePath, err.Error())
			}
			// 忽略 fs.PathError 并换个方式扫描
			err = scanNonUTF8ZipFile(filePath, newBook, scanOption)
		}
		// epub文件，需要根据 META-INF/container.xml 里面定义的rootfile （.opf文件）来重新排序
		if newBook.Type == entity.TypeEpub {
			imageList, err := fileutil.GetImageListFromEpubFile(newBook.FilePath)
			if err != nil {
				logger.Infof("%s", err)
			} else {
				newBook.SortPagesByImageList(imageList)
			}
			// 根据metadata，改写书籍信息
			metaData, err := fileutil.GetEpubMetadata(newBook.FilePath)
			if err != nil {
				logger.Infof("%s", err)
			} else {
				newBook.Author = metaData.Creator
				newBook.Press = metaData.Publisher
			}
		}
	// TODO:服务器解压速度太慢，网页用PDF.js解析？
	case entity.TypePDF:
		pageCount, pdfErr := fileutil.CountPagesOfPDF(filePath)
		if pdfErr != nil {
			return nil, pdfErr
		}
		if pageCount < 1 {
			return nil, errors.New(locale.GetString("NO_PAGES_IN_PDF") + filePath)
		}
		logger.Infof(locale.GetString("SCAN_PDF")+"%s: %d", filePath, pageCount)
		newBook.PageCount = pageCount
		newBook.InitComplete = true
		newBook.SetClover(entity.ImageInfo{RealImageFilePATH: "", FileSize: FileInfo.Size(), ModeTime: FileInfo.ModTime(), NameInArchive: "", Url: "/images/pdf.png"})
		for i := 1; i <= pageCount; i++ {
			TempURL := "api/get_file?id=" + newBook.BookID + "&filename=" + strconv.Itoa(i) + ".jpg"
			newBook.Pages.Images = append(newBook.Pages.Images, entity.ImageInfo{RealImageFilePATH: "", FileSize: FileInfo.Size(), ModeTime: FileInfo.ModTime(), NameInArchive: strconv.Itoa(i), Url: TempURL})
		}
	// TODO：简单的网页播放器
	case entity.TypeVideo:
		newBook.PageCount = 1
		newBook.InitComplete = true
		newBook.SetClover(entity.ImageInfo{NameInArchive: "video.png", Url: "images/video.png"})
	case entity.TypeAudio:
		newBook.PageCount = 1
		newBook.InitComplete = true
		newBook.SetClover(entity.ImageInfo{NameInArchive: "audio.png", Url: "images/audio.png"})
	case entity.TypeUnknownFile:
		newBook.PageCount = 1
		newBook.InitComplete = true
		newBook.SetClover(entity.ImageInfo{NameInArchive: "unknown.png", Url: "images/unknown.png"})
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
			if scanOption.IsSkipDir(path) {
				logger.Infof("Skip Scan:", path)
				return fs.SkipDir
			}
			f, errInfo := d.Info()
			if errInfo != nil {
				logger.Info(errInfo)
				return fs.SkipDir
			}
			if !scanOption.IsSupportMedia(path) {
				if scanOption.Debug {
					logger.Infof(locale.GetString("unsupported_file_type")+"%s", path)
				}
			} else {
				u, ok := f.(archiver.File) // f.Name不包含路径信息.需要转换一下
				if !ok {
					// 如果是文件夹+图片
					newBook.Type = entity.TypeDir
					////用Archiver的虚拟文件系统提供图片文件，理论上现在不应该用到
					//newBook.Pages = append(newBook.Pages, ImageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: "", Url: "/cache/" + newBook.BookID + "/" + url.QueryEscape(path)})
					//实验：用get_file接口提供文件服务
					TempURL := "api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(path)
					newBook.Pages.Images = append(newBook.Pages.Images, entity.ImageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: "", Url: TempURL})
					// logger.Infof(locale.GetString("unsupported_extract")+" %s", f)
				} else {
					// 替换特殊字符的时候，额外将“+替换成"%2b"，因为gin会将+解析为空格。
					TempURL := "api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(u.NameInArchive)
					// 不替换特殊字符
					// TempURL := "api/get_file?id=" + newBook.BookID + "&filename=" + u.NameInArchive
					newBook.Pages.Images = append(newBook.Pages.Images, entity.ImageInfo{RealImageFilePATH: "", FileSize: f.Size(), ModeTime: f.ModTime(), NameInArchive: u.NameInArchive, Url: TempURL})
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

func scanDirGetBook(dirPath string, storePath string, depth int, scanOption Option) (*entity.Book, error) {
	// 初始化，生成UUID
	newBook, err := entity.NewBook(dirPath, time.Now(), 0, storePath, depth, entity.TypeDir)
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
			logger.Info(errPath)
		}
		if scanOption.IsSupportMedia(file.Name()) {
			TempURL := "api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(file.Name())
			newBook.Pages.Images = append(newBook.Pages.Images, entity.ImageInfo{RealImageFilePATH: strAbsPath, FileSize: file.Size(), ModeTime: file.ModTime(), NameInArchive: file.Name(), Url: TempURL})
		}
	}
	newBook.SortPages("default")
	// 在添加到书库时判断页数
	return newBook, err
}
