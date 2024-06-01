package file

import (
	"github.com/yumenaka/comi/util"
	"github.com/yumenaka/comi/util/locale"
	"github.com/yumenaka/comi/util/logger"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type GetPictureDataOption struct {
	PictureName      string
	BookIsDir        bool
	BookIsPDF        bool
	BookIsNonUTF8Zip bool
	BookFilePath     string
	Debug            bool
	UseCache         bool
	ResizeWidth      int
	ResizeHeight     int
	ResizeMaxWidth   int
	ResizeMaxHeight  int
	ThumbnailMode    bool
	AutoCrop         int
	Gray             bool
	BlurHash         int
	BlurHashImage    int
}

func GetPictureData(option GetPictureDataOption) (imgData []byte, contentType string, err error) {
	pictureName := option.PictureName
	bookFilePath := option.BookFilePath
	//如果是特殊编码的ZIP文件
	if option.BookIsNonUTF8Zip {
		imgData, err = GetSingleFile(bookFilePath, pictureName, "gbk")
		if err != nil {
			return nil, "", err
		}
	}
	//如果是一般压缩文件，如zip、rar。epub
	if !option.BookIsNonUTF8Zip && !option.BookIsDir && !option.BookIsPDF {
		imgData, err = GetSingleFile(bookFilePath, pictureName, "")
		if err != nil {
			return nil, "", err
		}
	}
	//图片媒体类型，默认根据文件后缀设定。
	contentType = util.GetContentTypeByFileName(pictureName)
	//如果是PDF
	if option.BookIsPDF {
		//获取PDF的第几页
		page, err := strconv.Atoi(util.RemoveExtension(pictureName))
		if err != nil {
			return nil, "", err
		}
		imgData, err = GetImageFromPDF(bookFilePath, page, option.Debug)
		if err != nil {
			return nil, "", err
		}
		if imgData == nil {
			logger.Infof("GetImageFromPDF: imgData is nil")
			imgData, err = util.GenerateImage("Page " + util.RemoveExtension(pictureName) + ": " + locale.GetString("UnableToExtractImagesFromPDF"))
			if err != nil {
				return nil, "", err
			}
		}
		contentType = util.GetContentTypeByFileName(".jpg")
	}
	//如果是本地文件夹
	if option.BookIsDir {
		//直接读取磁盘文件
		imgData, err = os.ReadFile(filepath.Join(bookFilePath, pictureName))
		if err != nil {
			return nil, "", err
		}
	}
	canConvert := false
	for _, ext := range []string{".jpg", ".jpeg", ".gif", ".png", ".bmp"} {
		if strings.HasSuffix(strings.ToLower(pictureName), ext) {
			canConvert = true
		}
	}
	//不支持类型的图片直接返回原始数据
	if !canConvert {
		return imgData, contentType, nil
	}
	//处理图像文件
	//图片Resize, 按照固定的width height缩放
	if option.ResizeWidth > 0 && option.ResizeHeight > 0 {
		if option.ThumbnailMode {
			imgData = util.ImageThumbnail(imgData, option.ResizeWidth, option.ResizeHeight)
		} else {
			imgData = util.ImageResize(imgData, option.ResizeWidth, option.ResizeHeight)
		}
		contentType = util.GetContentTypeByFileName(".jpg")
	}
	//图片Resize, 按照 width 等比例缩放
	if option.ResizeHeight == 0 && option.ResizeWidth > 0 {
		imgData = util.ImageResizeByWidth(imgData, option.ResizeWidth)
		contentType = util.GetContentTypeByFileName(".jpg")
	}
	//图片Resize, 按照 height 等比例缩放
	if option.ResizeHeight > 0 && option.ResizeWidth == 0 {
		imgData = util.ImageResizeByHeight(imgData, option.ResizeHeight)
		contentType = util.GetContentTypeByFileName(".jpg")
	}
	//图片Resize, 按照 maxWidth 限制大小
	if option.ResizeMaxWidth > 0 {
		tempData, limitErr := util.ImageResizeByMaxWidth(imgData, option.ResizeMaxWidth)
		if limitErr != nil {
			logger.Info(limitErr)
		} else {
			imgData = tempData
		}
		contentType = util.GetContentTypeByFileName(".jpg")
	}
	//图片Resize, 按照 MaxHeight 限制大小
	if option.ResizeMaxHeight > 0 {
		tempData, limitErr := util.ImageResizeByMaxHeight(imgData, option.ResizeMaxHeight)
		if limitErr != nil {
			logger.Info(limitErr)
		} else {
			imgData = tempData
		}
		contentType = util.GetContentTypeByFileName(".jpg")
	}
	//自动切白边
	if option.AutoCrop > 0 && option.AutoCrop <= 100 {
		imgData = util.ImageAutoCrop(imgData, float32(option.AutoCrop))
		contentType = util.GetContentTypeByFileName(".jpg")
	}
	//转换为黑白图片
	if option.Gray {
		imgData = util.ImageGray(imgData)
		contentType = util.GetContentTypeByFileName(".jpg")
	}
	////获取对应图片的blurhash字符串(!)
	if option.BlurHash >= 1 && option.BlurHash <= 2 {
		hash := util.GetImageDataBlurHash(imgData, option.BlurHash)
		contentType = util.GetContentTypeByFileName(".txt")
		imgData = []byte(hash)
	}
	//返回blurhash图片 虽然blurhash components 理论上最大可以设到9，但速度太慢，限定为1或2
	if option.BlurHashImage >= 1 && option.BlurHashImage <= 2 {
		imgData = util.GetImageDataBlurHashImage(imgData, option.BlurHashImage)
		contentType = util.GetContentTypeByFileName(".jpg")
	}
	return imgData, contentType, nil
}
