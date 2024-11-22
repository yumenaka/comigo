package scan

import (
	"context"
	"errors"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/klauspost/compress/zip"
	"github.com/mholt/archives"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/util"
	fileutil "github.com/yumenaka/comigo/util/file"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
)

// Local 扫描路径，取得路径里的书籍
func Local(storePath string, scanOption Option) ([]*entity.Book, error) {
	if !util.PathExists(storePath) {
		return nil, errors.New(locale.GetString("PATH_NOT_EXIST"))
	}

	storePathAbs, err := filepath.Abs(storePath)
	if err != nil {
		logger.Infof("Failed to get absolute path: %s", err)
		storePathAbs = storePath
	}
	logger.Infof(locale.GetString("SCAN_START_HINT")+" %s", storePathAbs)

	// 创建已存在书籍的集合，提高查找效率
	existingBooks := make(map[string]struct{})
	if !scanOption.ReScanFile {
		for _, book := range entity.GetArchiveBooks() {
			existingBooks[book.FilePath] = struct{}{}
		}
	}

	var newBookList []*entity.Book
	err = filepath.Walk(storePathAbs, func(walkPath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			logger.Infof("Failed to access path %s: %v", walkPath, err)
			return err
		}

		absWalkPath, err := filepath.Abs(walkPath)
		if err != nil {
			logger.Infof("Failed to get absolute path: %s", err)
			absWalkPath = walkPath
		}

		// 跳过已存在的书籍
		if _, exists := existingBooks[absWalkPath]; exists {
			logger.Infof(locale.GetString("FoundInDatabase")+" %s", walkPath)
			return nil
		}

		// 计算路径深度
		relPath, err := filepath.Rel(storePathAbs, walkPath)
		if err != nil {
			logger.Infof("Failed to get relative path: %s", err)
			return err
		}
		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > scanOption.MaxScanDepth {
			logger.Infof(locale.GetString("ExceedsMaximumDepth")+" %d, base: %s, scan: %s", scanOption.MaxScanDepth, storePathAbs, walkPath)
			return filepath.SkipDir
		}

		if scanOption.IsSkipDir(walkPath) {
			logger.Infof(locale.GetString("SkipPath")+" %s", walkPath)
			return filepath.SkipDir
		}

		if fileInfo == nil {
			logger.Infof("FileInfo is nil for path: %s", walkPath)
			return nil
		}

		// 如果是文件
		if !fileInfo.IsDir() {
			if !scanOption.IsSupportArchiver(walkPath) {
				return nil
			}
			book, err := scanFileGetBook(walkPath, storePathAbs, depth, scanOption)
			if err != nil {
				logger.Infof("Failed to scan file: %s, error: %v", walkPath, err)
				return nil
			}
			newBookList = append(newBookList, book)
		}
		// 如果是目录
		if fileInfo.IsDir() {
			// 如果是目录
			book, err := scanDirGetBook(walkPath, storePathAbs, depth, scanOption)
			if err != nil {
				logger.Infof("Failed to scan directory: %s, error: %v", walkPath, err)
				return nil
			}
			newBookList = append(newBookList, book)
		}
		return nil
	})

	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("FOUND_IN_PATH"), len(newBookList), storePathAbs)
	}
	return newBookList, err
}

// 扫描本地文件，并返回对应书籍
func scanFileGetBook(filePath string, storePath string, depth int, scanOption Option) (*entity.Book, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Infof("Failed to open file: %s, error: %v", filePath, err)
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logger.Infof("Failed to get file info: %s, error: %v", filePath, err)
		return nil, err
	}

	newBook, err := entity.NewBook(filePath, fileInfo.ModTime(), fileInfo.Size(), storePath, depth, entity.GetBookTypeByFilename(filePath))
	if err != nil {
		return nil, err
	}

	switch newBook.Type {
	case entity.TypeZip, entity.TypeCbz, entity.TypeEpub:
		err = handleZipAndEpubFiles(filePath, newBook, scanOption)
		if err != nil {
			return nil, err
		}
	case entity.TypePDF:
		err = handlePdfFiles(filePath, newBook)
		if err != nil {
			return nil, err
		}
	case entity.TypeVideo:
		handleMediaFiles(newBook, "video.png", "images/video.png")
	case entity.TypeAudio:
		handleMediaFiles(newBook, "audio.png", "images/audio.png")
	case entity.TypeUnknownFile:
		handleMediaFiles(newBook, "unknown.png", "images/unknown.png")
	default:
		err = handleOtherArchiveFiles(filePath, newBook, scanOption)
		if err != nil {
			return nil, err
		}
	}

	newBook.SortPages("default")
	return newBook, nil
}

// 处理 ZIP 和 EPUB 文件
func handleZipAndEpubFiles(filePath string, newBook *entity.Book, scanOption Option) error {
	fsys, err := zip.OpenReader(filePath)
	if err != nil {
		return errors.New(locale.GetString("NOT_A_VALID_ZIP_FILE") + filePath)
	}
	defer fsys.Close()

	err = walkUTF8ZipFs(fsys, "", ".", newBook, scanOption)
	if err != nil {
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			if scanOption.Debug {
				logger.Infof("NonUTF-8 ZIP: %s, Error: %s", filePath, err.Error())
			}
			err = scanNonUTF8ZipFile(filePath, newBook, scanOption)
		} else {
			return err
		}
	}

	if newBook.Type == entity.TypeEpub {
		imageList, err := fileutil.GetImageListFromEpubFile(newBook.FilePath)
		if err == nil {
			newBook.SortPagesByImageList(imageList)
		} else {
			logger.Infof("Failed to get image list from EPUB: %s, error: %v", newBook.FilePath, err)
		}

		metaData, err := fileutil.GetEpubMetadata(newBook.FilePath)
		if err == nil {
			newBook.Author = metaData.Creator
			newBook.Press = metaData.Publisher
		} else {
			logger.Infof("Failed to get metadata from EPUB: %s, error: %v", newBook.FilePath, err)
		}
	}
	return nil
}

// 处理 PDF 文件
func handlePdfFiles(filePath string, newBook *entity.Book) error {
	pageCount, err := fileutil.CountPagesOfPDF(filePath)
	if err != nil {
		return err
	}
	if pageCount < 1 {
		return errors.New(locale.GetString("NO_PAGES_IN_PDF") + filePath)
	}
	logger.Infof(locale.GetString("SCAN_PDF")+" %s: %d pages", filePath, pageCount)
	newBook.PageCount = pageCount
	newBook.InitComplete = true
	newBook.SetCover(entity.ImageInfo{Url: "/images/pdf.png"})

	for i := 1; i <= pageCount; i++ {
		tempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + strconv.Itoa(i) + ".jpg"
		newBook.Pages.Images = append(newBook.Pages.Images, entity.ImageInfo{
			NameInArchive: strconv.Itoa(i),
			Url:           tempURL,
		})
	}
	return nil
}

// 处理视频、音频等媒体文件
func handleMediaFiles(newBook *entity.Book, imageName, imageUrl string) {
	newBook.PageCount = 1
	newBook.InitComplete = true
	newBook.SetCover(entity.ImageInfo{NameInArchive: imageName, Url: imageUrl})
}

// 处理其他类型的压缩文件
func handleOtherArchiveFiles(filePath string, newBook *entity.Book, scanOption Option) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fsys, err := archives.FileSystem(ctx, filePath, nil)
	if err != nil {
		return err
	}

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Infof("Failed to access path %s in archive: %v", path, err)
			return err
		}

		if scanOption.IsSkipDir(path) {
			logger.Infof("Skip Scan: %s", path)
			return fs.SkipDir
		}

		f, err := d.Info()
		if err != nil {
			logger.Infof("Failed to get file info in archive: %v", err)
			return fs.SkipDir
		}

		if scanOption.IsSupportMedia(path) {
			archivedFile, ok := f.(archives.FileInfo)
			var tempURL string
			if ok {
				tempURL = "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(archivedFile.NameInArchive)
				newBook.Pages.Images = append(newBook.Pages.Images, entity.ImageInfo{
					NameInArchive: archivedFile.NameInArchive,
					Url:           tempURL,
				})
			} else {
				tempURL = "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(path)
				newBook.Pages.Images = append(newBook.Pages.Images, entity.ImageInfo{
					Url: tempURL,
				})
			}
		} else {
			if scanOption.Debug {
				logger.Infof(locale.GetString("unsupported_file_type")+" %s", path)
			}
		}
		return nil
	})
	return err
}

// 扫描目录，并返回对应书籍
func scanDirGetBook(dirPath string, storePath string, depth int, scanOption Option) (*entity.Book, error) {
	newBook, err := entity.NewBook(dirPath, time.Now(), 0, storePath, depth, entity.TypeDir)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		logger.Infof("Failed to read directory: %s, error: %v", dirPath, err)
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !scanOption.IsSupportMedia(fileName) {
			continue
		}

		fileInfo, err := entry.Info()
		if err != nil {
			logger.Infof("Failed to get file info: %s, error: %v", fileName, err)
			continue
		}

		absPath := filepath.Join(dirPath, fileName)
		tempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(fileName)
		newBook.Pages.Images = append(newBook.Pages.Images, entity.ImageInfo{
			RealImageFilePATH: absPath,
			FileSize:          fileInfo.Size(),
			ModeTime:          fileInfo.ModTime(),
			NameInArchive:     fileName,
			Url:               tempURL,
		})
	}

	newBook.SortPages("default")
	return newBook, nil
}
