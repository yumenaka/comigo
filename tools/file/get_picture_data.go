package file

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/vfs"
)

type GetPictureDataOption struct {
	PictureName      string
	BookID           string // 书籍 ID，用于计算远程 PDF 的缓存路径（压缩包类型不需要缓存）
	BookIsDir        bool
	BookIsPDF        bool
	BookIsNonUTF8Zip bool
	BookPath         string
	Debug            bool
	UseCache         bool
	ResizeWidth      int
	ResizeHeight     int
	ResizeMaxWidth   int
	ResizeMaxHeight  int
	AutoCrop         int
	Gray             bool
	BlurHash         int
	BlurHashImage    int
	// 远程书籍支持
	IsRemote  bool   // 是否为远程书籍
	RemoteURL string // 远程存储的 URL
}

func GetPictureData(option GetPictureDataOption) (imgData []byte, contentType string, err error) {
	pictureName := option.PictureName
	bookPath := option.BookPath

	// 如果是远程压缩包，直接从 WebDAV 流式读取
	if option.IsRemote && !option.BookIsDir && !option.BookIsPDF && option.RemoteURL != "" {
		// 获取 VFS 实例
		fs, fsErr := vfs.GetOrCreate(option.RemoteURL, vfs.Options{
			CacheEnabled: false, // 不缓存，直接流式读取
			Timeout:      30,
		})
		if fsErr != nil {
			return nil, "", fmt.Errorf("无法连接远程书库: %w", fsErr)
		}

		// 打开支持 Seek 的 Reader
		reader, readerErr := fs.OpenReaderAtSeeker(bookPath)
		if readerErr != nil {
			return nil, "", fmt.Errorf("无法打开远程文件: %w", readerErr)
		}
		defer func() {
			if closer, ok := reader.(io.Closer); ok {
				_ = closer.Close()
			}
		}()

		// 获取文件名用于格式识别
		fileName := filepath.Base(bookPath)

		// 创建超时上下文
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// 从流中读取文件
		textEncoding := ""
		if option.BookIsNonUTF8Zip {
			textEncoding = "gbk"
		}
		imgData, err = GetSingleFileFromStream(ctx, fileName, reader, pictureName, textEncoding)
		if err != nil {
			return nil, "", err
		}
	} else {
		// 本地文件或目录类型
		actualBookPath := bookPath

		// 如果是特殊编码的ZIP文件
		if option.BookIsNonUTF8Zip {
			imgData, err = GetSingleFile(actualBookPath, pictureName, "gbk")
			if err != nil {
				return nil, "", err
			}
		}
		// 如果是一般压缩文件，如zip、rar。epub
		if !option.BookIsNonUTF8Zip && !option.BookIsDir && !option.BookIsPDF {
			imgData, err = GetSingleFile(actualBookPath, pictureName, "")
			if err != nil {
				return nil, "", err
			}
		}
	}
	// 图片媒体类型，默认根据文件后缀设定。
	contentType = tools.GetContentTypeByFileName(pictureName)
	// 如果是PDF
	if option.BookIsPDF {
		// 获取PDF的第几页
		page, err := strconv.Atoi(tools.RemoveExtension(pictureName))
		if err != nil {
			return nil, "", err
		}

		// 远程 PDF：需要从缓存读取（PDF 处理需要随机访问）
		actualPdfPath := bookPath
		if option.IsRemote && option.RemoteURL != "" && option.BookID != "" {
			// PDF 需要下载到缓存（因为 PDF 处理库需要随机访问）
			actualPdfPath = getRemotePdfCachePath(option.BookID, bookPath, option.RemoteURL)
		}

		imgData, err = GetImageFromPDF(actualPdfPath, page, option.Debug)
		if err != nil {
			return nil, "", err
		}
		if imgData == nil {
			logger.Info(locale.GetString("log_getimagefrompdf_imgdata_nil"))
			imgData, err = tools.GenerateImage("Page " + tools.RemoveExtension(pictureName) + ": " + locale.GetString("unable_to_extract_images_from_pdf"))
			if err != nil {
				return nil, "", err
			}
		}
		contentType = tools.GetContentTypeByFileName(".jpg")
	}
	// 如果是文件夹类型
	if option.BookIsDir {
		if option.IsRemote && option.RemoteURL != "" {
			// 远程目录：直接从 WebDAV 读取图片文件
			fs, fsErr := vfs.GetOrCreate(option.RemoteURL, vfs.Options{
				CacheEnabled: true, // 启用 VFS 层面的文件读取缓存（用于优化性能）
				Timeout:      30,
			})
			if fsErr != nil {
				return nil, "", fsErr
			}
			imagePath := fs.JoinPath(bookPath, pictureName)
			imgData, err = fs.ReadFile(imagePath)
			if err != nil {
				return nil, "", err
			}
		} else {
			// 本地目录：直接读取磁盘文件
			imgData, err = os.ReadFile(filepath.Join(bookPath, pictureName))
			if err != nil {
				return nil, "", err
			}
		}
	}
	canConvert := false
	for _, ext := range []string{".jpg", ".jpeg", ".gif", ".png", ".bmp"} {
		if strings.HasSuffix(strings.ToLower(pictureName), ext) {
			canConvert = true
		}
	}
	// 不支持类型的图片直接返回原始数据
	if !canConvert {
		return imgData, contentType, nil
	}
	// 处理图像文件
	// 图片Resize, 按照固定的width height缩放
	if option.ResizeWidth > 0 && option.ResizeHeight > 0 {
		imgData = tools.ImageResize(imgData, option.ResizeWidth, option.ResizeHeight)
		contentType = tools.GetContentTypeByFileName(".jpg")
	}
	// 图片Resize, 按照 width 等比例缩放
	if option.ResizeHeight == 0 && option.ResizeWidth > 0 {
		imgData = tools.ImageResizeByWidth(imgData, option.ResizeWidth)
		contentType = tools.GetContentTypeByFileName(".jpg")
	}
	// 图片Resize, 按照 height 等比例缩放
	if option.ResizeHeight > 0 && option.ResizeWidth == 0 {
		imgData = tools.ImageResizeByHeight(imgData, option.ResizeHeight)
		contentType = tools.GetContentTypeByFileName(".jpg")
	}
	// 图片Resize, 按照 maxWidth 限制大小
	if option.ResizeMaxWidth > 0 {
		tempData, limitErr := tools.ImageResizeByMaxWidth(imgData, option.ResizeMaxWidth)
		if limitErr != nil {
			logger.Info(limitErr)
		} else {
			imgData = tempData
		}
		contentType = tools.GetContentTypeByFileName(".jpg")
	}
	// 图片Resize, 按照 MaxHeight 限制大小
	if option.ResizeMaxHeight > 0 {
		tempData, limitErr := tools.ImageResizeByMaxHeight(imgData, option.ResizeMaxHeight)
		if limitErr != nil {
			logger.Info(limitErr)
		} else {
			imgData = tempData
		}
		contentType = tools.GetContentTypeByFileName(".jpg")
	}
	// 自动切白边
	if option.AutoCrop > 0 && option.AutoCrop <= 100 {
		imgData = tools.ImageAutoCrop(imgData, float32(option.AutoCrop))
		contentType = tools.GetContentTypeByFileName(".jpg")
	}
	// 转换为黑白图片
	if option.Gray {
		imgData = tools.ImageGray(imgData)
		contentType = tools.GetContentTypeByFileName(".jpg")
	}
	// //获取对应图片的blurhash字符串(!)
	if option.BlurHash >= 1 && option.BlurHash <= 2 {
		hash := tools.GetImageDataBlurHash(imgData, option.BlurHash)
		contentType = tools.GetContentTypeByFileName(".txt")
		imgData = []byte(hash)
	}
	// 返回blurhash图片 虽然blurhash components 理论上最大可以设到9，但速度太慢，限定为1或2
	if option.BlurHashImage >= 1 && option.BlurHashImage <= 2 {
		imgData = tools.GetImageDataBlurHashImage(imgData, option.BlurHashImage)
		contentType = tools.GetContentTypeByFileName(".jpg")
	}
	return imgData, contentType, nil
}

// getRemotePdfCachePath 获取远程 PDF 的本地缓存路径
// PDF 处理库（pdfcpu）需要文件路径进行随机访问，无法流式读取，所以需要下载到本地
// 注意：PDF 文件在扫描时应该已经下载到缓存，这里作为后备机制按需下载
func getRemotePdfCachePath(bookID string, remotePath string, remoteURL string) string {
	cacheDir := os.Getenv("COMIGO_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = os.TempDir()
	}
	fileName := filepath.Base(remotePath)
	cachePath := filepath.Join(cacheDir, "remote_books", bookID, fileName)

	// 如果缓存不存在，按需下载（PDF 在扫描时应该已经下载，这里作为后备）
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		logger.Infof("按需下载远程 PDF: %s", remotePath)
		fs, err := vfs.GetOrCreate(remoteURL, vfs.Options{
			CacheEnabled: false,
			Timeout:      30,
		})
		if err == nil {
			data, readErr := fs.ReadFile(remotePath)
			if readErr == nil {
				_ = os.MkdirAll(filepath.Dir(cachePath), 0o755)
				_ = os.WriteFile(cachePath, data, 0o644)
			}
		}
	}

	return cachePath
}
