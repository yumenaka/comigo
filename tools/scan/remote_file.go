package scan

import (
	"context"
	"fmt"
	"io"
	fslib "io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/yumenaka/archives"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/vfs"
)

// handleRemoteZipFile 处理远程 ZIP/EPUB 文件
// 直接从 WebDAV 流式读取，不下载到本地
func handleRemoteZipFile(vfsInstance vfs.FileSystem, filePath string, newBook *model.Book) error {
	// 打开支持 Seek 的 Reader
	reader, err := vfsInstance.OpenReaderAtSeeker(filePath)
	if err != nil {
		return fmt.Errorf("无法打开远程文件: %w", err)
	}
	defer func() {
		if closer, ok := reader.(io.Closer); ok {
			_ = closer.Close()
		}
	}()

	// 获取文件名用于格式识别
	fileName := filepath.Base(filePath)

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用 archives.FileSystem 从流中读取
	archiveFS, err := archives.FileSystem(ctx, fileName, reader)
	if err != nil {
		return fmt.Errorf("无法创建压缩包文件系统: %w", err)
	}

	// 遍历压缩包内容（使用标准库的 fs.WalkDir）
	pageNum := 1
	err = fslib.WalkDir(archiveFS, ".", func(p string, d fslib.DirEntry, walkErr error) error {
		if walkErr != nil {
			logger.Infof(locale.GetString("log_failed_to_access_path_in_archive"), p, walkErr)
			return walkErr
		}

		if IsSkipDir(p) {
			logger.Infof(locale.GetString("log_skip_scan_path"), p)
			return fslib.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		if IsSupportMedia(p) {
			f, err := d.Info()
			if err != nil {
				return nil
			}

			archivedFile, ok := f.(archives.FileInfo)
			var tempURL string
			if ok {
				tempURL = "/api/get-file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(archivedFile.NameInArchive)
				newBook.PageInfos = append(newBook.PageInfos, model.PageInfo{
					Name:    archivedFile.NameInArchive,
					Url:     tempURL,
					PageNum: pageNum,
				})
			} else {
				tempURL = "/api/get-file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(p)
				newBook.PageInfos = append(newBook.PageInfos, model.PageInfo{
					Name:    p,
					Url:     tempURL,
					PageNum: pageNum,
				})
			}
			pageNum++
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("遍历压缩包失败: %w", err)
	}

	// EPUB 特殊处理
	if newBook.Type == model.TypeEpub {
		// EPUB 元数据需要读取文件，暂时跳过（或使用流式读取）
		// 这里可以后续优化
		logger.Infof("EPUB 元数据提取暂不支持远程流式读取")
	}

	newBook.SortPages("default")
	return nil
}

// handleRemotePdfFile 处理远程 PDF 文件
// PDF 需要随机访问，所以需要下载到本地缓存
// 注意：PDF 处理库（pdfcpu）需要文件路径进行随机访问，无法流式读取
func handleRemotePdfFile(vfsInstance vfs.FileSystem, filePath string, newBook *model.Book) error {
	localPath, err := downloadToCache(vfsInstance, filePath, newBook.BookID)
	if err != nil {
		return fmt.Errorf("下载远程 PDF 失败: %w", err)
	}

	// 使用本地文件处理逻辑
	err = handlePdfFiles(localPath, newBook)
	if err != nil {
		return err
	}

	return nil
}

// handleRemoteOtherArchiveFile 处理远程其他压缩包文件
// 直接从 WebDAV 流式读取，不下载到本地
func handleRemoteOtherArchiveFile(vfsInstance vfs.FileSystem, filePath string, newBook *model.Book) error {
	// 打开支持 Seek 的 Reader
	reader, err := vfsInstance.OpenReaderAtSeeker(filePath)
	if err != nil {
		return fmt.Errorf("无法打开远程文件: %w", err)
	}
	defer func() {
		if closer, ok := reader.(io.Closer); ok {
			_ = closer.Close()
		}
	}()

	// 获取文件名用于格式识别
	fileName := filepath.Base(filePath)

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用 archives.FileSystem 从流中读取
	archiveFS, err := archives.FileSystem(ctx, fileName, reader)
	if err != nil {
		return fmt.Errorf("无法创建压缩包文件系统: %w", err)
	}

	// 遍历压缩包内容（使用标准库的 fs.WalkDir）
	pageNum := 1
	err = fslib.WalkDir(archiveFS, ".", func(p string, d fslib.DirEntry, walkErr error) error {
		if walkErr != nil {
			logger.Infof(locale.GetString("log_failed_to_access_path_in_archive"), p, walkErr)
			return walkErr
		}

		if IsSkipDir(p) {
			logger.Infof(locale.GetString("log_skip_scan_path"), p)
			return fslib.SkipDir
		}

		f, err := d.Info()
		if err != nil {
			logger.Infof(locale.GetString("log_failed_to_get_file_info_in_archive"), err)
			return fslib.SkipDir
		}

		if IsSupportMedia(p) {
			archivedFile, ok := f.(archives.FileInfo)
			var tempURL string
			if ok {
				tempURL = "/api/get-file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(archivedFile.NameInArchive)
				newBook.PageInfos = append(newBook.PageInfos, model.PageInfo{
					Name:    archivedFile.NameInArchive,
					Url:     tempURL,
					PageNum: pageNum,
				})
			} else {
				tempURL = "/api/get-file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(p)
				newBook.PageInfos = append(newBook.PageInfos, model.PageInfo{
					Name:    p,
					Url:     tempURL,
					PageNum: pageNum,
				})
			}
			pageNum++
		} else {
			if cfg.GetDebug() {
				logger.Infof(locale.GetString("unsupported_file_type")+" %s", p)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("遍历压缩包失败: %w", err)
	}

	newBook.SortPages("default")
	return nil
}

// downloadToCache 下载远程文件到本地缓存目录（仅用于 PDF 等需要随机访问的文件）
// 返回本地缓存文件路径
func downloadToCache(fs vfs.FileSystem, remotePath string, bookID string) (string, error) {
	// 获取缓存目录
	cacheDir := config.GetCfg().CacheDir
	if cacheDir == "" {
		cacheDir = os.TempDir()
	}

	// 创建书籍专用缓存目录
	bookCacheDir := filepath.Join(cacheDir, "remote_books", bookID)
	if err := os.MkdirAll(bookCacheDir, 0o755); err != nil {
		return "", fmt.Errorf("创建缓存目录失败: %w", err)
	}

	// 获取文件名
	fileName := path.Base(remotePath)
	localPath := filepath.Join(bookCacheDir, fileName)

	// 检查缓存是否已存在
	if info, err := os.Stat(localPath); err == nil {
		// 获取远程文件信息
		remoteInfo, remoteErr := fs.Stat(remotePath)
		if remoteErr == nil {
			// 如果本地缓存文件的修改时间不早于远程文件，直接使用缓存
			if !info.ModTime().Before(remoteInfo.ModTime()) {
				logger.Infof("使用缓存文件: %s", localPath)
				return localPath, nil
			}
		}
	}

	// 下载文件
	logger.Infof("下载远程文件到缓存: %s -> %s", remotePath, localPath)
	data, err := fs.ReadFile(remotePath)
	if err != nil {
		return "", fmt.Errorf("读取远程文件失败: %w", err)
	}

	// 保存到本地
	if err := os.WriteFile(localPath, data, 0o644); err != nil {
		return "", fmt.Errorf("保存缓存文件失败: %w", err)
	}

	return localPath, nil
}

// handleRemoteDirectoryImages 处理远程目录类型书籍的图片
// 这种情况不需要下载整个压缩包，而是按需读取图片
func handleRemoteDirectoryImages(vfsInstance vfs.FileSystem, dirPath string, newBook *model.Book) error {
	entries, err := vfsInstance.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("读取远程目录失败: %w", err)
	}

	pageNum := 1
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !IsSupportMedia(fileName) {
			continue
		}

		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := vfsInstance.JoinPath(dirPath, fileName)
		tempURL := "/api/get-file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(fileName)

		newBook.PageInfos = append(newBook.PageInfos, model.PageInfo{
			Path:    fullPath,
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime(),
			Name:    fileName,
			Url:     tempURL,
			PageNum: pageNum,
		})
		pageNum++
	}

	return nil
}

// GetRemoteImageData 获取远程图片数据（用于目录类型书籍）
func GetRemoteImageData(book *model.Book, imageName string) ([]byte, error) {
	if !book.IsRemote {
		// 本地书籍，直接读取
		imagePath := filepath.Join(book.BookPath, imageName)
		return os.ReadFile(imagePath)
	}

	// 远程书籍
	fs, err := vfs.GetOrCreate(book.RemoteURL, vfs.Options{
		CacheEnabled: true,
		Timeout:      30,
	})
	if err != nil {
		return nil, fmt.Errorf("无法连接远程书库: %w", err)
	}

	// 读取远程图片
	imagePath := fs.JoinPath(book.BookPath, imageName)
	return fs.ReadFile(imagePath)
}
