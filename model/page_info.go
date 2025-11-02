package model

import (
	"time"
)

// PageInfo 单个媒体文件的信息
type PageInfo struct {
	Name       string    `json:"name"`     // 用于解压的压缩文件内文件路径，或图片名，为了适应特殊字符，经过一次转义
	Path       string    `json:"path"`     // 文件路径
	Size       int64     `json:"size"`     // 文件大小
	ModTime    time.Time `json:"mod_time"` // 修改时间
	Url        string    `json:"url"`      // 远程用户读取图片的URL，为了适应特殊字符，经过转义
	PageNum    int       `json:"page_num"` // 图片在原始文件中的页码位置，这个字段不解析。用来按照原始顺序排序
	Blurhash   string    `json:"-"`        // `json:"blurhash"` //blurhash占位符。扫描图片生成（tools.GetImageDataBlurHash）
	Height     int       `json:"-"`        // 暂时用不着 这个字段不解析`json:"height"`   //blurhash用，图片高
	Width      int       `json:"-"`        // 暂时用不着 这个字段不解析`json:"width"`    //blurhash用，图片宽
	ImgType    string    `json:"-"`        // 这个字段不解析
	InsertHtml string    `json:"-"`        // 这个字段不解析
}
