package arch

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/yumenaka/comi/logger"
)

//sample code: https://github.com/pdfcpu/pdfcpu/blob/master/pkg/api/extract.go

// CountPagesOfPDF 确定PDF的页数
func CountPagesOfPDF(pdfFileName string) (int, error) {
	//设置一个defer语句来捕获并处理潜在的panic。defer语句会确保在函数返回之前执行其中的代码，而recover函数用于捕获并恢复panic，防止panic向上传播并导致整个程序崩溃
	defer func() {
		if r := recover(); r != nil {
			logger.Infof("CountPagesOfPDF: invalid PDF: %v Error:%v", pdfFileName, r)
			// 这里可以根据需要进行错误处理，比如返回特定的错误值给调用者
		}
	}()
	// use default configuration for pdfcpu ("nil")
	err := api.ValidateFile(pdfFileName, nil)
	if err != nil {
		return -1, errors.New("CountPagesOfPDF: invalid PDF: " + pdfFileName + " " + err.Error())
	}
	return api.PageCountFile(pdfFileName)
}

// GetImageFromPDF 从PDF里面取jpeg文件，pageNum从1开始
func GetImageFromPDF(pdfFileName string, pageNum int, Debug bool) ([]byte, error) {
	start := time.Now()
	//打开PDF文件
	file, err := os.Open(pdfFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	pdfSetting := model.NewDefaultConfiguration()
	pdfSetting.DecodeAllStreams = true

	////api.ExtractImagesRaw(） 只能导出图片，无法输出文档
	//pageImagesMap, err := api.ExtractImagesRaw(file, []string{strconv.Itoa(pageNum)}, pdfSetting)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//images := make([]model.Image, 0)
	//for _, pageImages := range pageImagesMap {
	//	for _, img := range pageImages {
	//		images = append(images, img)
	//	}
	//}
	//var imgBytes []byte = nil
	//for i := range images {
	//	b, err := io.ReadAll(images[i])
	//	if err != nil {
	//		continue
	//	}
	//	imgBytes = b
	//	fmt.Println(time.Now().Sub(start))
	//	return b, nil
	//}
	//return imgBytes, nil

	//api.ExtractImagesRaw(） 另一种调用方式，效果相同，依然无法渲染文字为图片。
	buffer := &bytes.Buffer{}
	err = api.ExtractImages(file, []string{strconv.Itoa(pageNum)}, digestImage(buffer), pdfSetting)
	if err != nil {
		fmt.Println(err)
	}
	//参考：https://github.com/pdfcpu/pdfcpu/issues/45
	// https://github.com/pdfcpu/pdfcpu/issues/351
	//将PDF转换为图像，相当于重写一个pdf查看器。
	//虽然可以引用mupdf。但这会导致导入C代码，破坏跨平台兼容性。
	//或许可以cli调用imagemagick，曲线救国？
	if Debug {
		logger.Infof("GetImageFromPDF: %v", time.Now().Sub(start))
	}

	return buffer.Bytes(), nil
}

// 自定义钩子函数参数的方法：输入自定义参数、输出符合要求的函数
func digestImage(buff *bytes.Buffer) func(model.Image, bool, int) error {
	return func(img model.Image, singleImgPerPage bool, maxPageDigits int) error {
		imageOut, err := imaging.Decode(img)
		if err != nil {
			fmt.Println(err)
			return errors.New("imaging.Decode() Error")
		}
		err = jpeg.Encode(buff, imageOut, &jpeg.Options{Quality: 75})
		if err != nil {
			return errors.New("digestImage jpeg.Encode( Error")
		}
		return nil
	}
}

func ExportImageFromPDF(pdfFile string, pageNum int) {
	start := time.Now()
	//打开PDF文件
	file, err := os.Open(pdfFile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	pageImagesMap, err := api.ExtractImagesRaw(file, []string{strconv.Itoa(pageNum)}, model.NewDefaultConfiguration())
	if err != nil {
		fmt.Println(err)
	}
	images := make([]model.Image, 0)
	for _, pageImages := range pageImagesMap {
		for _, img := range pageImages {
			images = append(images, img)
		}
	}

	for i := range images {
		imgBytes, err := io.ReadAll(images[i])
		if err != nil {
			continue
		}
		// 写入文件，如果文件不存在则创建，文件权限设置为 0644
		err = os.WriteFile("example.jpeg", imgBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(time.Now().Sub(start))
}

// PDF页面分辨率
type dim struct {
	width  float64
	height float64
}

// 取得PDF页面分辨率
func GetPageDimensions(fileName string) []dim {
	pageCount, _ := CountPagesOfPDF(fileName)
	log.Printf("pagecount of %v was %v", fileName, pageCount)
	var pageDimensions []dim
	var currentPageDim dim
	pdfDims, err := api.PageDimsFile(fileName)
	if err != nil {
		log.Printf("Error %v", err)
	}
	for i := 0; i < pageCount; i++ {
		currentPageDim.width = pdfDims[i].Width
		currentPageDim.height = pdfDims[i].Height
		pageDimensions = append(pageDimensions, currentPageDim)
	}
	return pageDimensions
}

func ExportAllImageFromPDF(pdfFile string) {
	//打开PDF文件
	file, err := os.Open(pdfFile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	pageCount, err := CountPagesOfPDF("01.pdf")
	if err != nil {
		fmt.Println(err)
	}
	list := []string{}
	for i := 0; i < pageCount; i++ {
		list = append(list, strconv.Itoa(i+1))
	}

	pageImagesMap, err := api.ExtractImagesRaw(file, []string{strconv.Itoa(pageCount)}, model.NewDefaultConfiguration())
	if err != nil {
		fmt.Println(err)
	}
	images := make([]model.Image, 0)
	for _, pageImages := range pageImagesMap {
		for _, img := range pageImages {
			images = append(images, img)
		}
	}

	for i := range images {
		imgBytes, err := io.ReadAll(images[i])
		if err != nil {
			continue
		}
		// 写入文件，如果文件不存在则创建，文件权限设置为 0644
		err = os.WriteFile("test/"+strconv.Itoa(i+1)+".jpg", imgBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

}
