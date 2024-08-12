package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/entity"
	fileutil "github.com/yumenaka/comi/util/file"
	"github.com/yumenaka/comi/util/logger"
)

// GetFile 示例 URL： 127.0.0.1:1234/get_file?id=2b17a13&filename=1.jpg
// 缩放文件，会转化为jpeg：http://127.0.0.1:1234/api/get_file?resize_width=300&resize_height=400&id=597e06&filename=01.jpeg
// 相关参数：
// id：书籍的ID，必须项目       							&id=2B17a
// filename:获取的文件名，必须项目   							&filename=01.jpg
// //可选参数：
// resize_width:指定宽度，缩放图片  							&resize_width=300
// resize_height:指定高度，缩放图片 							&resize_height=300
// thumbnail_mode:缩略图模式，同时指定宽高的时候要不要剪切图片		&thumbnail_mode=true
// resize_max_width:指定宽度上限，图片宽度大于这个上限时缩小图片  	&resize_max_width=740
// resize_max_height:指定高度上限，图片高度大于这个上限时缩小图片  	&resize_max_height=300
// auto_crop:自动切白边，数字是切白边的阈值，范围是0~100 			&auto_crop=10
// gray:黑白化												&gray=true
// blurhash:获取对应图片的blurhash，不是原始图片 				&blurhash=3
// blurhash_image:获取对应图片的blurhash图片，不是原始图片  	    &blurhash_image=3
// blurhash_image:获取对应图片的blurhash图片，不是原始图片  	    &base64=true
func GetFile(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	needFile := c.DefaultQuery("filename", "")
	//没有指定这两项，直接返回
	if id == "" && needFile == "" {
		return
	}
	noCache := c.DefaultQuery("no-cache", "false")
	//如果启用了本地缓存
	if config.Config.UseCache && noCache == "false" {
		//获取所有的参数键值对
		query := c.Request.URL.Query()
		//如果有缓存，直接读取本地获取缓存文件并返回
		cacheData, ct, errGet := fileutil.GetFileFromCache(
			id,
			needFile,
			fileutil.GetQueryString(query),
			c.DefaultQuery("thumbnail_mode", "false") == "true",
			config.Config.CachePath,
			config.Config.Debug,
		)
		if errGet == nil && cacheData != nil {
			c.Data(http.StatusOK, ct, cacheData)
			return
		}
	}

	bookByID, err := entity.GetBookByID(id, "")
	if err != nil {
		logger.Infof("%s", err)
	}

	//读取图片Resize用的resizeWidth
	resizeWidth, errX := strconv.Atoi(c.DefaultQuery("resize_width", "0"))
	if errX != nil {
		resizeWidth = 0
	}
	//读取图片Resize用的resizeHeight
	resizeHeight, errY := strconv.Atoi(c.DefaultQuery("resize_height", "0"))
	if errY != nil {
		resizeHeight = 0
	}
	//自动切白边参数
	autoCrop, errCrop := strconv.Atoi(c.DefaultQuery("auto_crop", "-1"))
	if errCrop != nil {
		autoCrop = -1
	}
	//图片Resize, 按照 maxWidth 限制大小
	resizeMaxWidth, errMX := strconv.Atoi(c.DefaultQuery("resize_max_width", "0"))
	if errMX != nil {
		resizeMaxWidth = 0
	}
	//图片Resize, 按照 MaxHeight 限制大小
	resizeMaxHeight, errMY := strconv.Atoi(c.DefaultQuery("resize_max_height", "0"))
	if errMY != nil {
		resizeMaxHeight = 0
	}

	blurhash, blurErr := strconv.Atoi(c.DefaultQuery("blurhash", "0"))
	if blurErr != nil {
		blurhash = 0
	}

	blurhashImage, blurImageErr := strconv.Atoi(c.DefaultQuery("blurhash_image", "0"))
	if blurImageErr != nil {
		blurhashImage = 0
	}

	option := fileutil.GetPictureDataOption{
		PictureName:      needFile,
		BookIsPDF:        bookByID.Type == entity.TypePDF,
		BookIsDir:        bookByID.Type == entity.TypeDir,
		BookIsNonUTF8Zip: bookByID.NonUTF8Zip,
		BookFilePath:     bookByID.FilePath,
		Debug:            config.Config.Debug,
		UseCache:         config.Config.UseCache,
		ResizeWidth:      resizeWidth,
		ResizeHeight:     resizeHeight,
		ResizeMaxWidth:   resizeMaxWidth,
		ResizeMaxHeight:  resizeMaxHeight,
		ThumbnailMode:    c.DefaultQuery("thumbnail_mode", "false") == "true",
		AutoCrop:         autoCrop,
		Gray:             c.DefaultQuery("gray", "false") == "true",
		BlurHash:         blurhash,
		BlurHashImage:    blurhashImage,
	}
	imgData, contentType, err := fileutil.GetPictureData(option)
	if err != nil {
		c.String(http.StatusBadRequest, "GetPictureData error:%s", err)
	}

	//如果启用了本地缓存
	if config.Config.UseCache && noCache == "false" && bookByID.Type != entity.TypeDir {
		//获取所有的参数键值对
		query := c.Request.URL.Query()
		//缓存文件到本地，避免重复解压
		errSave := fileutil.SaveFileToCache(
			id,
			needFile,
			imgData,
			fileutil.GetQueryString(query),
			contentType,
			c.DefaultQuery("thumbnail_mode", "false") == "true",
			config.Config.CachePath,
			config.Config.Debug,
		)
		if errSave != nil {
			logger.Info(errSave)
		}
	}
	// 是否需要base64编码
	// https://freshman.tech/snippets/go/image-to-base64/
	base64Encode := c.DefaultQuery("base64Encode", "false") == "true"
	if base64Encode {
		var base64Encoding string
		// Determine the content type of the image file
		mimeType := http.DetectContentType(imgData)
		// Prepend the appropriate URI scheme header depending
		// on the MIME type
		switch mimeType {
		case "image/jpeg":
			base64Encoding += "data:image/jpeg;base64Encode,"
		case "image/png":
			base64Encoding += "data:image/png;base64Encode,"
		case "image/gif":
			base64Encoding += "data:image/gif;base64Encode,"
		case "image/webp":
			base64Encoding += "data:image/webp;base64Encode,"
		case "image/svg+xml":
			base64Encoding += "data:image/svg+xml;base64Encode,"
		case "image/bmp":
			base64Encoding += "data:image/bmp;base64Encode,"
		case "image/avif":
			base64Encoding += "data:image/avif;base64Encode,"
		}
		base64Encoding += base64.StdEncoding.EncodeToString(imgData)
		// 返回包含 base64 编码数据的 img 标签
		c.String(http.StatusOK, fmt.Sprintf(`<img class="m-2 max-w-full lg:max-w-[800px] rounded shadow-lg" src="%s" alt="Base64 Image"/>`, base64Encoding))
	}

	c.Data(http.StatusOK, contentType, imgData)
}
