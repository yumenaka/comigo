package model

import (
	"bytes"
	"image"
	"log"
	"time"

	"github.com/bbrks/go-blurhash"
	"github.com/disintegration/imaging"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/util/file"
	"github.com/yumenaka/comigo/util/logger"
)

// MediaFileInfo 单个媒体文件的信息
type MediaFileInfo struct {
	Name       string    `json:"name"`     // 用于解压的压缩文件内文件路径，或图片名，为了适应特殊字符，经过一次转义
	Path       string    `json:"path"`     // 文件路径
	Size       int64     `json:"size"`     // 文件大小
	ModTime    time.Time `json:"mod_time"` // 修改时间
	Url        string    `json:"url"`      // 远程用户读取图片的URL，为了适应特殊字符，经过转义
	PageNum    int       `json:"-"`        // 这个字段不解析
	Blurhash   string    `json:"-"`        // `json:"blurhash"` //blurhash占位符。扫描图片生成（util.GetImageDataBlurHash）
	Height     int       `json:"-"`        // 暂时用不着 这个字段不解析`json:"height"`   //blurhash用，图片高
	Width      int       `json:"-"`        // 暂时用不着 这个字段不解析`json:"width"`    //blurhash用，图片宽
	ImgType    string    `json:"-"`        // 这个字段不解析
	InsertHtml string    `json:"-"`        // 这个字段不解析
}

// analyzeImage 获取某页漫画的分辨率与blurhash
func (i *MediaFileInfo) analyzeImage(bookPath string) (err error) {
	var img image.Image
	// img, err = imaging.Open(i.RealImageFilePATH)

	imgData, err := file.GetSingleFile(bookPath, i.Name, "gbk")
	if err != nil {
		logger.Infof("%s", err)
	}
	buf := bytes.NewBuffer(imgData)
	img, err = imaging.Decode(buf)
	if err != nil {
		log.Printf(locale.GetString("check_image_error")+" %v\n", err)
	} else {
		i.Width = img.Bounds().Dx()
		i.Height = img.Bounds().Dy()
		// 很耗费服务器资源，以后再研究。
		str, err := blurhash.Encode(1, 1, img)
		if err != nil {
			// Handle errors
			log.Printf(locale.GetString("check_image_error")+" %v\n", err)
		}
		i.Blurhash = str
	}
	return err
}
