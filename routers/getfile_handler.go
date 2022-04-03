package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/arch"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/tools"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// 示例 URL： 127.0.0.1:1234/getfile?id=2b17a13&filename=1.jpg
// 缩放文件，会转化为jpeg：http://127.0.0.1:1234/api/getfile?resize_width=300&resize_height=400&id=597e06&filename=01.jpeg
// 相关参数：
// id：书籍的ID，必须项目       							&id=2B17a
// filename:获取的文件名，必须项目   							&filename=01.jpg
////可选参数：
// resize_width:指定宽度，缩放图片  							&resize_width=300
// resize_height:指定高度，缩放图片 							&resize_height=300
// thumbnail_mode:缩略图模式，同时指定宽高的时候要不要剪切图片		&thumbnail_mode=true
// resize_max_width:指定宽度上限，图片宽度大于这个上限时缩小图片  	&resize_max_width=740
// resize_max_height:指定高度上限，图片高度大于这个上限时缩小图片  	&resize_max_height=300
// auto_crop:自动切白边，数字是切白边的阈值，范围是0~100 			&auto_crop=10
// gray:黑白化												&gray=true
// blurhash:获取对应图片的blurhash，而不是原始图片 				&blurhash=3
// blurhash_image:获取对应图片的blurhash图片，而不是原始图片  	&blurhash_image=3
//TODO：生成临时文件加速下一次访问，并在退出后清理
func getFileHandler(c *gin.Context) {
	//time.Sleep(5 * time.Second)
	id := c.DefaultQuery("id", "")
	needFile := c.DefaultQuery("filename", "")
	//没有指定这两项，直接返回
	if id == "" && needFile == "" {
		return
	}
	//获取所有的参数键值对
	query := c.Request.URL.Query()
	//fmt.Println(query)
	//如果有缓存，直接读取本地获取缓存文件并返回
	cacheData, ct, errGet := getFileFromCache(id, needFile, query, c.DefaultQuery("thumbnail_mode", "false") == "true")
	if errGet == nil {
		c.Data(http.StatusOK, ct, cacheData)
		return
	}
	book, err := common.GetBookByID(id, false)
	if err != nil {
		fmt.Println(err)
	}
	bookPath := book.GetFilePath()
	//fmt.Println(bookPath)
	var imgData []byte
	//如果是特殊编码的ZIP文件
	if book.NonUTF8Zip && book.BookType != common.BookTypeDir {
		imgData, err = arch.GetSingleFile(bookPath, needFile, "gbk")
		if err != nil {
			fmt.Println(err)
		}
	}
	//如果是一般压缩文件
	if !book.NonUTF8Zip && book.BookType != common.BookTypeDir {
		imgData, err = arch.GetSingleFile(bookPath, needFile, "")
		if err != nil {
			fmt.Println(err)
		}
	}
	//如果是本地文件夹
	if book.BookType == common.BookTypeDir {
		//直接读取磁盘文件
		imgData, err = ioutil.ReadFile(filepath.Join(bookPath, needFile))
		if err != nil {
			fmt.Println(err)
		}
	}
	//默认的媒体类型，默认值根据文件后缀设定。
	contentType := tools.GetContentTypeByFileName(needFile)
	canConvert := false
	for _, ext := range []string{".jpg", ".jpeg", ".gif", ".png", ".bmp"} {
		if strings.HasSuffix(strings.ToLower(needFile), ext) {
			canConvert = true
		}
	}
	//不支持类型的图片直接返回原始数据
	if !canConvert {
		c.Data(http.StatusOK, contentType, imgData)
		return
	}
	//处理图像文件
	if imgData != nil {
		//读取图片Resize用的resizeWidth
		resizeWidth, errX := strconv.Atoi(c.DefaultQuery("resize_width", "0"))
		if errX != nil {
			fmt.Println(errX)
		}
		//读取图片Resize用的resizeHeight
		resizeHeight, errY := strconv.Atoi(c.DefaultQuery("resize_height", "0"))
		if errY != nil {
			fmt.Println(errY)
		}
		//图片Resize, 按照固定的width height缩放
		if errX == nil && errY == nil && resizeWidth > 0 && resizeHeight > 0 {
			//是否要用缩略图模式
			thumbnailMode := c.DefaultQuery("thumbnail_mode", "false")
			if thumbnailMode == "true" {
				imgData = tools.ImageThumbnail(imgData, resizeWidth, resizeHeight)
			} else {
				imgData = tools.ImageResize(imgData, resizeWidth, resizeHeight)
			}
			contentType = tools.GetContentTypeByFileName(".jpg")
		}
		//图片Resize, 按照 width 等比例缩放
		if errX == nil && errY != nil && resizeWidth > 0 {
			imgData = tools.ImageResizeByWidth(imgData, resizeWidth)
			contentType = tools.GetContentTypeByFileName(".jpg")
		}
		//图片Resize, 按照 height 等比例缩放
		if errY == nil && errX != nil && resizeHeight > 0 {
			imgData = tools.ImageResizeByHeight(imgData, resizeHeight)
			contentType = tools.GetContentTypeByFileName(".jpg")
		}
		//图片Resize, 按照 maxWidth 限制大小
		resizeMaxWidth, errMX := strconv.Atoi(c.DefaultQuery("resize_max_width", "0"))
		if errMX != nil {
			fmt.Println(errMX)
		}
		if errMX == nil && resizeMaxWidth > 0 {
			tempData, limitErr := tools.ImageResizeByMaxWidth(imgData, resizeMaxWidth)
			if limitErr != nil {
				fmt.Println(limitErr)
			} else {
				imgData = tempData
				contentType = tools.GetContentTypeByFileName(".jpg")
			}
		}
		//图片Resize, 按照 MaxHeight 限制大小
		resizeMaxHeight, errMY := strconv.Atoi(c.DefaultQuery("resize_max_height", "0"))
		if errMY != nil {
			fmt.Println(errMY)
		}
		if errMY == nil && resizeMaxHeight > 0 {
			tempData, limitErr := tools.ImageResizeByMaxHeight(imgData, resizeMaxHeight)
			if limitErr != nil {
				fmt.Println(limitErr)
			} else {
				imgData = tempData
				contentType = tools.GetContentTypeByFileName(".jpg")
			}
		}
		//自动切白边参数
		autoCrop, errCrop := strconv.Atoi(c.DefaultQuery("auto_crop", "-1"))
		if errCrop != nil {
			fmt.Println(errCrop)
		}

		if errCrop == nil && autoCrop > 0 && autoCrop <= 100 {
			imgData = tools.ImageAutoCrop(imgData, float32(autoCrop))
			contentType = tools.GetContentTypeByFileName(".jpg")
		}
		//转换为黑白图片
		gray := c.DefaultQuery("gray", "false")
		//获取对应图片的blurhash字符串并返回，不是图片
		blurhash, blurErr := strconv.Atoi(c.DefaultQuery("blurhash", "0"))
		if blurErr != nil {
			fmt.Println(blurErr)
		}
		if gray == "true" {
			imgData = tools.ImageGray(imgData)
			contentType = tools.GetContentTypeByFileName(".jpg")
		}

		//虽然blurhash components 理论上最大可以设到9，但反应速度太慢，毫无实用性、建议最大为2
		if blurhash >= 1 && blurhash <= 2 && blurErr == nil {
			hash := tools.GetImageDataBlurHash(imgData, blurhash)
			contentType = tools.GetContentTypeByFileName(".txt")
			imgData = []byte(hash)
		}
		//返回图片的blurhash图
		blurhashImage, blurImageErr := strconv.Atoi(c.DefaultQuery("blurhash_image", "0"))
		if blurImageErr != nil {
			fmt.Println(blurImageErr)
		}
		//虽然blurhash components 理论上最大可以设到9，但反应速度太慢，毫无实用性、建议为1（最大为2）
		if blurhashImage >= 1 && blurhashImage <= 2 && blurErr == nil {
			imgData = tools.GetImageDataBlurHashImage(imgData, blurhashImage)
			contentType = tools.GetContentTypeByFileName(".jpg")
		}

		//缓存文件到本地，避免重复解压
		errSave := saveFileToCache(id, needFile, imgData, query, contentType, c.DefaultQuery("thumbnail_mode", "false") == "true")
		if errSave != nil {
			fmt.Println(errSave)
		}
		c.Data(http.StatusOK, contentType, imgData)
	}

}

type cacheKey struct {
	bookID      string
	queryString string
}

//需要一个Map保存ContentType
var mapContentType map[cacheKey]string

func init() {
	mapContentType = make(map[cacheKey]string)
}

//读取过一次的文件，就保存到硬盘上加快读取
func saveFileToCache(id string, filename string, data []byte, query url.Values, contentType string, isCover bool) error {
	err := os.MkdirAll(filepath.Join(common.Config.CacheFilePath, id), os.ModePerm)
	if err != nil {
		println(locale.GetString("saveFileToCache_error"))
	}
	//特殊字符转义，避免保存不了
	filename = url.PathEscape(filename)
	//如果是封面，另存为cover.png、cover.jpeg
	if isCover {
		filename = "cover" + path.Ext(filename)
	}
	err = ioutil.WriteFile(filepath.Join(common.Config.CacheFilePath, id, filename), data, 0644)
	if err != nil {
		fmt.Println(err)
	}
	qS := ""
	for k, values := range query {
		temp := ""
		for _, v := range values {
			temp = temp + v
		}
		qS = qS + k + temp
	}
	key := cacheKey{bookID: id, queryString: qS}
	//将ContentType存入Map
	mapContentType[key] = contentType
	return err
}

//读取缓存，加快第二次访问的速度
func getFileFromCache(id string, filename string, query url.Values, isCover bool) ([]byte, string, error) {
	//文件名经过转义，避免保存不了，所以这里也必须转义才能取到本地文件
	filename = url.PathEscape(filename)
	//如果是封面，另存为cover.png、cover.jpeg
	if isCover {
		filename = "cover" + path.Ext(filename)
	}
	loadedImage, err := ioutil.ReadFile(filepath.Join(common.Config.CacheFilePath, id, filename))
	if err != nil {
		fmt.Println("getFileFromCache Error:" + filename)
	}
	qS := ""
	for k, values := range query {
		temp := ""
		for _, v := range values {
			temp = temp + v
		}
		qS = qS + k + temp
	}
	key := cacheKey{bookID: id, queryString: qS}
	contentType := mapContentType[key]
	return loadedImage, contentType, err
}
