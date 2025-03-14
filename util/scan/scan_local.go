package scan

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/klauspost/compress/zip"
	"github.com/yumenaka/archives"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util"
	fileutil "github.com/yumenaka/comigo/util/file"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
)

// 全局变量：标记是否正在扫描，避免并发扫描
var (
	scanning  bool       // 标记是否正在扫描，避免并发扫描
	scanMutex sync.Mutex // 保护 scanning 标志的锁
)

// 扫描目录的核心函数：递归遍历目录，忽略指定名称的文件夹，收集图片文件信息
func ScanDirectoryNew(currentPath string, depth int, option Option) (model.DirNode, []string, []model.MediaFileInfo, error) {
	node := model.DirNode{
		Name: filepath.Base(currentPath), // filepath.Base():返回路径的最后一个元素
		Path: currentPath,
	}
	var foundDirs []string
	var foundFiles []model.MediaFileInfo

	// 如果超过最大深度限制，直接返回空节点
	if option.Cfg.GetMaxScanDepth() >= 0 && depth > option.Cfg.GetMaxScanDepth() {
		return node, foundDirs, foundFiles, nil
	}

	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return node, foundDirs, foundFiles, err
	}

	// 当前目录计入 foundDirs（用于记录树状结构）
	foundDirs = append(foundDirs, currentPath)

	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(currentPath, name)
		if entry.IsDir() {
			// 检查是否在忽略列表
			if option.IsSkipDir(name) {
				continue
			}
			// 递归扫描子目录
			subNode, subDirs, subFiles, subErr := ScanDirectoryNew(fullPath, depth+1, option)
			if subErr != nil {
				// 忽略单个子目录出错，继续扫描其他目录
				fmt.Println("扫描子目录出错:", subErr)
				continue
			}
			node.SubDirs = append(node.SubDirs, subNode)
			// 合并子目录扫描结果
			foundDirs = append(foundDirs, subDirs...)
			foundFiles = append(foundFiles, subFiles...)
		} else {
			// 文件：检查扩展名是否为支持的图片格式
			ext := strings.ToLower(filepath.Ext(name))
			// 非支持媒体或压缩包格式，跳过
			if (!option.IsSupportMedia(ext)) && (!option.IsSupportFile(ext)) {
				continue
			}
			// 获取文件信息
			info, err := entry.Info()
			if err != nil {
				fmt.Println("获取文件信息失败:", err)
				continue
			}
			size := info.Size()
			modTime := info.ModTime()
			// 说明: Go 的标准库 os.FileInfo 不直接提供创建时间。如果需要准确创建时间，可以使用第三方库获取。
			mediaFileInfo := model.MediaFileInfo{
				Name:    name,
				Path:    fullPath,
				Size:    size,
				ModTime: modTime,
			}
			node.Files = append(node.Files, mediaFileInfo)
			foundFiles = append(foundFiles, mediaFileInfo)
		}
	}
	return node, foundDirs, foundFiles, nil
}

// ScanDirectory 扫描本地路径，取得路径里的书籍
func ScanDirectory(storePath string, option Option) ([]*model.Book, error) {
	if !util.PathExists(storePath) {
		return nil, errors.New(locale.GetString("path_not_exist"))
	}
	storePathAbs, err := filepath.Abs(storePath)
	if err != nil {
		logger.Infof("Failed to get absolute path: %s", err)
		storePathAbs = storePath
	}
	logger.Infof(locale.GetString("scan_start_hint")+" %s", storePathAbs)

	// 已存在书籍的集合，跳过已有书籍，提高查找效率
	existingBooks := make(map[string]struct{})
	for _, book := range model.GetAllBookList() {
		existingBooks[book.FilePath] = struct{}{}
	}

	var newBookList []*model.Book
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
			logger.Infof(locale.GetString("found_in_bookstore")+" %s", walkPath)
			return nil
		}

		// 计算路径深度
		relPath, err := filepath.Rel(storePathAbs, walkPath)
		if err != nil {
			logger.Infof("Failed to get relative path: %s", err)
			return err
		}
		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > option.Cfg.GetMaxScanDepth() {
			logger.Infof(locale.GetString("exceeds_maximum_depth")+" %d, base: %s, scan: %s", option.Cfg.GetMaxScanDepth(), storePathAbs, walkPath)
			return filepath.SkipDir
		}

		if option.IsSkipDir(walkPath) {
			logger.Infof(locale.GetString("skip_path")+" %s", walkPath)
			return filepath.SkipDir
		}

		if fileInfo == nil {
			logger.Infof("MediaFileInfo is nil for path: %s", walkPath)
			return nil
		}

		// 如果是文件
		if !fileInfo.IsDir() {
			if !option.IsSupportArchiver(walkPath) {
				return nil
			}
			book, err := scanFileGetBook(walkPath, storePathAbs, depth, option)
			if err != nil {
				logger.Infof("Failed to scan file: %s, error: %v", walkPath, err)
				return nil
			}
			newBookList = append(newBookList, book)
		}
		// 如果是目录
		if fileInfo.IsDir() {
			// 如果是目录
			book, err := scanDirGetBook(walkPath, storePathAbs, depth, option)
			if err != nil {
				logger.Infof("Failed to scan directory: %s, error: %v", walkPath, err)
				return nil
			}
			newBookList = append(newBookList, book)
		}
		return nil
	})

	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("found_in_path"), len(newBookList), storePathAbs)
	}
	return newBookList, err
}

// 扫描本地文件，并返回对应书籍
func scanFileGetBook(filePath string, storePath string, depth int, scanOption Option) (*model.Book, error) {
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

	newBook, err := model.NewBook(filePath, fileInfo.ModTime(), fileInfo.Size(), storePath, depth, model.GetBookTypeByFilename(filePath))
	if err != nil {
		return nil, err
	}

	switch newBook.Type {
	case model.TypeZip, model.TypeCbz, model.TypeEpub:
		err = handleZipAndEpubFiles(filePath, newBook, scanOption)
		if err != nil {
			return nil, err
		}
	case model.TypePDF:
		err = handlePdfFiles(filePath, newBook)
		if err != nil {
			return nil, err
		}
	case model.TypeVideo:
		handleMediaFiles(newBook, "video.png", "/images/video.png")
	case model.TypeAudio:
		handleMediaFiles(newBook, "audio.png", "/images/audio.png")
	case model.TypeUnknownFile:
		handleMediaFiles(newBook, "unknown.png", "/images/unknown.png")
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
func handleZipAndEpubFiles(filePath string, newBook *model.Book, option Option) error {
	fsys, err := zip.OpenReader(filePath)
	if err != nil {
		return errors.New(locale.GetString("not_a_valid_zip_file") + filePath)
	}
	defer fsys.Close()

	err = walkUTF8ZipFs(fsys, "", ".", newBook, option)
	if err != nil {
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			if option.Cfg.GetDebug() {
				logger.Infof("NonUTF-8 ZIP: %s, Error: %s", filePath, err.Error())
			}
			err = scanNonUTF8ZipFile(filePath, newBook, option)
		} else {
			return err
		}
	}

	if newBook.Type == model.TypeEpub {
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
func handlePdfFiles(filePath string, newBook *model.Book) error {
	pageCount, err := fileutil.CountPagesOfPDF(filePath)
	if err != nil {
		return err
	}
	if pageCount < 1 {
		return errors.New(locale.GetString("no_pages_in_pdf") + filePath)
	}
	logger.Infof(locale.GetString("scan_pdf")+" %s: %d pages", filePath, pageCount)
	newBook.PageCount = pageCount
	newBook.InitComplete = true
	newBook.SetCover(model.MediaFileInfo{Url: "/images/pdf.png"})

	for i := 1; i <= pageCount; i++ {
		tempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + strconv.Itoa(i) + ".jpg"
		newBook.Pages.Images = append(newBook.Pages.Images, model.MediaFileInfo{
			Name: strconv.Itoa(i),
			Url:  tempURL,
		})
	}
	return nil
}

// 处理视频、音频等媒体文件
func handleMediaFiles(newBook *model.Book, imageName, imageUrl string) {
	newBook.PageCount = 1
	newBook.InitComplete = true
	newBook.SetCover(model.MediaFileInfo{Name: imageName, Url: imageUrl})
}

// 处理其他类型的压缩文件
func handleOtherArchiveFiles(filePath string, newBook *model.Book, option Option) error {
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

		if option.IsSkipDir(path) {
			logger.Infof("Skip Scan: %s", path)
			return fs.SkipDir
		}

		f, err := d.Info()
		if err != nil {
			logger.Infof("Failed to get file info in archive: %v", err)
			return fs.SkipDir
		}

		if option.IsSupportMedia(path) {
			archivedFile, ok := f.(archives.FileInfo)
			var tempURL string
			if ok {
				tempURL = "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(archivedFile.NameInArchive)
				newBook.Pages.Images = append(newBook.Pages.Images, model.MediaFileInfo{
					Name: archivedFile.NameInArchive,
					Url:  tempURL,
				})
			} else {
				tempURL = "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(path)
				newBook.Pages.Images = append(newBook.Pages.Images, model.MediaFileInfo{
					Url: tempURL,
				})
			}
		} else {
			if option.Cfg.GetDebug() {
				logger.Infof(locale.GetString("unsupported_file_type")+" %s", path)
			}
		}
		return nil
	})
	return err
}

// 扫描目录，并返回对应书籍
func scanDirGetBook(dirPath string, storePath string, depth int, option Option) (*model.Book, error) {
	// 获取文件夹信息
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}
	newBook, err := model.NewBook(dirPath, dirInfo.ModTime(), dirInfo.Size(), storePath, depth, model.TypeDir)
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
		if !option.IsSupportMedia(fileName) {
			continue
		}

		fileInfo, err := entry.Info()
		if err != nil {
			logger.Infof("Failed to get file info: %s, error: %v", fileName, err)
			continue
		}

		absPath := filepath.Join(dirPath, fileName)
		tempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(fileName)
		newBook.Pages.Images = append(newBook.Pages.Images, model.MediaFileInfo{
			Path:    absPath,
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime(),
			Name:    fileName,
			Url:     tempURL,
		})
	}

	newBook.SortPages("default")
	return newBook, nil
}
