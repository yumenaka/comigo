package types

import (
	"bytes"
	"github.com/bbrks/go-blurhash"
	"github.com/disintegration/imaging"
	"github.com/yumenaka/comi/fileutil"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/logger"
	"image"
	"log"
	"time"
)

// ImageInfo 单张漫画图片
type ImageInfo struct {
	PageNum           int       `json:"-"`        //这个字段不解析
	NameInArchive     string    `json:"filename"` //用于解压的压缩文件内文件路径，或图片名，为了适应特殊字符，经过一次转义
	Url               string    `json:"url"`      //远程用户读取图片的URL，为了适应特殊字符，经过一次转义
	Blurhash          string    `json:"-"`        //`json:"blurhash"` //blurhash占位符。需要扫描图片生成（util.GetImageDataBlurHash）
	Height            int       `json:"-"`        //暂时用不着 这个字段不解析`json:"height"`   //blurhash用，图片高
	Width             int       `json:"-"`        //暂时用不着 这个字段不解析`json:"width"`    //blurhash用，图片宽
	ModeTime          time.Time `json:"-"`        //这个字段不解析
	FileSize          int64     `json:"-"`        //这个字段不解析
	RealImageFilePATH string    `json:"-"`        //这个字段不解析  书籍为文件夹的时候，实际图片的路径
	ImgType           string    `json:"-"`        //这个字段不解析
}

// analyzeImage 获取某页漫画的分辨率与blurhash
func (i *ImageInfo) analyzeImage(bookPath string) (err error) {
	var img image.Image
	//img, err = imaging.Open(i.RealImageFilePATH)

	imgData, err := fileutil.GetSingleFile(bookPath, i.NameInArchive, "gbk")
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
		//很耗费服务器资源，以后再研究。
		str, err := blurhash.Encode(1, 1, img)
		if err != nil {
			// Handle errors
			log.Printf(locale.GetString("check_image_error")+" %v\n", err)
		}
		i.Blurhash = str
	}
	return err
}
