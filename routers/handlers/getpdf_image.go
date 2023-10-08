package handlers

//
//import (
//	"fmt"
//	"net/http"
//	"strconv"
//	"strings"
//
//	"github.com/gin-gonic/gin"
//
//	"github.com/yumenaka/comi/arch"
//	"github.com/yumenaka/comi/book"
//	"github.com/yumenaka/comi/common"
//)
//
//// GetPdfImageHandler 示例 URL： 127.0.0.1:1234/get_pdf_image?id=2b17a13&filename=1.jpg
//// 相关参数：
//// id：书籍ID，必须项目       							&id=2B17a
//// filename:图片文件名  							&filename=01.jpg
//////可选参数： no-cache
//func GetPdfImageHandler(c *gin.Context) {
//	id := c.DefaultQuery("id", "")
//	needFile := c.DefaultQuery("filename", "")
//	//没有指定这两项，直接返回
//	if id == "" && needFile == "" {
//		return
//	}
//	noCache := c.DefaultQuery("no-cache", "false")
//	//如果启用了本地缓存
//	if common.Config.UseCache && noCache == "false" {
//		//获取所有的参数键值对
//		query := c.Request.URL.Query()
//		//如果有缓存，直接读取本地获取缓存文件并返回
//		cacheData, ct, errGet := getFileFromCache(id, needFile, query, c.DefaultQuery("thumbnail_mode", "false") == "true")
//		if errGet == nil && cacheData != nil {
//			c.Data(http.StatusOK, ct, cacheData)
//			return
//		}
//	}
//	bookByID, err := book.GetBookByID(id, false)
//	if err != nil {
//		fmt.Println(err)
//	}
//	//如果不是PDF文件
//	if bookByID.Type != book.TypePDF {
//		return
//	}
//	bookPath := bookByID.GetFilePath()
//	//fmt.Println(bookPath)
//
//	//将1.jpg转换为1
//	pageNum, err := strconv.Atoi(strings.TrimSuffix(needFile, ".jpg"))
//	if err != nil {
//		return
//	}
//
//	var imgData []byte
//	//如果是PDF文件
//	if bookByID.Type == book.TypePDF {
//		imgData, err = arch.GetImageFromPDF(bookPath, pageNum)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//	//如果启用了本地缓存
//	if common.Config.UseCache && noCache == "false" {
//		//获取所有的参数键值对
//		query := c.Request.URL.Query()
//		//缓存文件到本地，避免重复解压
//		errSave := saveFileToCache(id, needFile, imgData, query, "image/jpeg", c.DefaultQuery("thumbnail_mode", "false") == "true")
//		if errSave != nil {
//			fmt.Println(errSave)
//		}
//	}
//	c.Data(http.StatusOK, "image/jpeg", imgData)
//}
