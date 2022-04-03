package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/arch"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/tools"
	"io/ioutil"
	"net/http"
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
	if id != "" && needFile != "" {
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
			//自动切白边
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
			if gray == "true" {
				imgData = tools.ImageGray(imgData)
				contentType = tools.GetContentTypeByFileName(".jpg")
			}
			//获取对应图片的blurhash字符串并返回，不是图片
			blurhash, blurErr := strconv.Atoi(c.DefaultQuery("blurhash", "0"))
			if blurErr != nil {
				fmt.Println(blurErr)
			}
			//虽然blurhash components 理论上最大可以设到9，但反应速度太慢，毫无实用性、建议最大为2
			if blurhash >= 1 && blurhash <= 2 && blurErr == nil {
				////使用异步处理，避免程序阻塞
				//// goroutine 中只能使用只读的上下文 c.Copy()
				//cCp := c.Copy()
				//go func() {
				//	hash := tools.GetImageDataBlurHash(imgData, blurhash)
				//	contentType = tools.GetContentTypeByFileName(".txt")
				//	imgData = []byte(hash)
				//	// 读上下文，写会出错╮(￣▽￣")╭
				//	//cCp.Data(http.StatusOK, contentType, imgData)
				//	log.Println("GetBlurHash Done! in path " + cCp.Request.URL.BasePath)
				//	fmt.Println(hash)
				//	return
				//}()
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
			c.Data(http.StatusOK, contentType, imgData)
		}
	}
}
