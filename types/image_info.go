package types

import "time"

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
