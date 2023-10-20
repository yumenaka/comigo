package arch

//
//import (
//	"bytes"
//	"errors"
//	"fmt"
//	"image/jpeg"
//	"log"
//	"os"
//	"strconv"
//	"time"
//
//	"github.com/disintegration/imaging"
//	"github.com/pdfcpu/pdfcpu/pkg/api"
//	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
//)
//
//// CountPagesOfPDF 确定PDF的页数
//func CountPagesOfPDF(pdfFileName string) (int, error) {
//	// use default configuration for pdfcpu ("nil")
//	err := api.ValidateFile(pdfFileName, nil)
//	if err != nil {
//		return -1, errors.New("CountPagesOfPDF: invalid PDF: %v")
//	}
//	return api.PageCountFile(pdfFileName)
//}
//
//// GetImageFromPDF 从PDF里面取jpeg文件，pageNum从1开始
//func GetImageFromPDF(pdfFileName string, pageNum int) ([]byte, error) {
//	start := time.Now()
//	//打开PDF文件
//	file, err := os.Open(pdfFileName)
//	if err != nil {
//		logger.Info(err)
//	}
//	defer file.Close()
//
//	pdfSetting := pdfcpu.NewDefaultConfiguration()
//	pdfSetting.DecodeAllStreams = true
//
//	////api.ExtractImagesRaw(） 只能导出图片，无法输出文档
//	//imageList, err := api.ExtractImagesRaw(file, []string{strconv.Itoa(pageNum)}, pdfSetting)
//	//if err != nil {
//	//	logger.Info(err)
//	//}
//	//imageOut, err := imaging.Decode(imageList[0])
//	//if err != nil {
//	//	logger.Info(err)
//	//	return nil, errors.New("imaging.Decode() Error")
//	//}
//	//buffer := &bytes.Buffer{}
//	//err = jpeg.Encode(buffer, imageOut, &jpeg.Options{Quality: 75})
//	//if err != nil {
//	//	return nil, errors.New("jpeg.Encode( Error")
//	//}
//
//	//api.ExtractImagesRaw(） 另一种调用方式，效果相同，依然无法渲染文字为图片。
//	buffer := &bytes.Buffer{}
//	err = api.ExtractImages(file, []string{strconv.Itoa(pageNum)}, digestImage(buffer), pdfSetting)
//	if err != nil {
//		logger.Info(err)
//	}
//	//参考：https://github.com/pdfcpu/pdfcpu/issues/45
//	// https://github.com/pdfcpu/pdfcpu/issues/351
//	//将PDF转换为图像，相当于重写一个pdf查看器。
//	//虽然可以引用mupdf。但这会导致导入C代码，破坏跨平台兼容性。
//	//或许可以调用imagemagick，曲线救国？
//
//	logger.Info(time.Now().Sub(start))
//	return buffer.Bytes(), nil
//}
//
////自定义钩子函数参数的方法：输入自定义参数、输出符合要求的函数
//func digestImage(buff *bytes.Buffer) func(pdfcpu.Image, bool, int) error {
//	return func(img pdfcpu.Image, singleImgPerPage bool, maxPageDigits int) error {
//		imageOut, err := imaging.Decode(img)
//		if err != nil {
//			logger.Info(err)
//			return errors.New("imaging.Decode() Error")
//		}
//		err = jpeg.Encode(buff, imageOut, &jpeg.Options{Quality: 75})
//		if err != nil {
//			return errors.New("digestImage jpeg.Encode( Error")
//		}
//		return nil
//	}
//}
//
//func ExportImageFromPDF(pdfFile string, pageNum int) {
//	start := time.Now()
//	//打开PDF文件
//	file, err := os.Open(pdfFile)
//	if err != nil {
//		logger.Info(err)
//	}
//	defer file.Close()
//
//	imageList, err := api.ExtractImagesRaw(file, []string{strconv.Itoa(pageNum)}, pdfcpu.NewDefaultConfiguration())
//	if err != nil {
//		logger.Info(err)
//	}
//	imageOut, err := imaging.Decode(imageList[0])
//	if err != nil {
//		logger.Info(err)
//		return
//	}
//	//保存文件
//	err = imaging.Save(imageOut, "test/"+strconv.Itoa(pageNum)+".jpg")
//	if err != nil {
//		logger.Info(err)
//	}
//	logger.Info(time.Now().Sub(start))
//}
//
////PDF页面分辨率
//type dim struct {
//	width  float64
//	height float64
//}
//
////取得PDF页面分辨率
//func GetPageDimensions(fileName string) []dim {
//	pageCount, _ := CountPagesOfPDF(fileName)
//	log.Printf("pagecount of %v was %v", fileName, pageCount)
//	var pageDimensions []dim
//	var currentPageDim dim
//	pdfDims, err := api.PageDimsFile(fileName)
//	if err != nil {
//		log.Printf("Error %v", err)
//	}
//	for i := 0; i < pageCount; i++ {
//		currentPageDim.width = pdfDims[i].Width
//		currentPageDim.height = pdfDims[i].Height
//		pageDimensions = append(pageDimensions, currentPageDim)
//	}
//	return pageDimensions
//}
//
//func ExportAllImageFromPDF(pdfFile string) {
//	//打开PDF文件
//	file, err := os.Open(pdfFile)
//	if err != nil {
//		logger.Info(err)
//	}
//	defer file.Close()
//	pageCount, err := CountPagesOfPDF("01.pdf")
//	if err != nil {
//		logger.Info(err)
//	}
//	list := []string{}
//	for i := 0; i < pageCount; i++ {
//		list = append(list, strconv.Itoa(i+1))
//	}
//	imageList, err := api.ExtractImagesRaw(file, list, pdfcpu.NewDefaultConfiguration())
//	if err != nil {
//		logger.Info(err)
//	}
//	for k, i := range imageList {
//		imageOut, err := imaging.Decode(i)
//		if err != nil {
//			logger.Info(err)
//			return
//		}
//		//保存文件
//		err = imaging.Save(imageOut, "test/"+strconv.Itoa(k+1)+".jpg")
//		if err != nil {
//			logger.Info(err)
//		}
//	}
//}
