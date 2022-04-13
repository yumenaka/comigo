package arch

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

// CountPagesOfPDF 确定PDF的页数
func CountPagesOfPDF(pdfFileName string) (int, error) {
	// use default configuration for pdfcpu ("nil")
	err := api.ValidateFile(pdfFileName, nil)
	if err != nil {
		return -1, errors.New("CountPagesOfPDF: invalid PDF: %v")
	}
	return api.PageCountFile(pdfFileName)
}

// GetImageFromPDF 从PDF里面取jpeg文件，pageNum从1开始
func GetImageFromPDF(pdfFileName string, pageNum int) ([]byte, error) {
	start := time.Now()
	//打开PDF文件
	file, err := os.Open(pdfFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	imageList, err := api.ExtractImagesRaw(file, []string{strconv.Itoa(pageNum)}, pdfcpu.NewDefaultConfiguration())
	if err != nil {
		fmt.Println(err)
	}
	imageOut, err := imaging.Decode(imageList[0])
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("imaging.Decode() Error")
	}
	buf2 := &bytes.Buffer{}
	err = jpeg.Encode(buf2, imageOut, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, errors.New("jpeg.Encode( Error")
	}
	fmt.Println(time.Now().Sub(start))
	return buf2.Bytes(), nil
}

func ExportImageFromPDF(pdfFile string, pageNum int) {
	start := time.Now()
	//打开PDF文件
	file, err := os.Open(pdfFile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	imageList, err := api.ExtractImagesRaw(file, []string{strconv.Itoa(pageNum)}, pdfcpu.NewDefaultConfiguration())
	if err != nil {
		fmt.Println(err)
	}
	imageOut, err := imaging.Decode(imageList[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	//保存文件
	err = imaging.Save(imageOut, "test/"+strconv.Itoa(pageNum)+".jpg")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(time.Now().Sub(start))
}

//PDF页面分辨率
type dim struct {
	width  float64
	height float64
}

//取得PDF页面分辨率
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
	imageList, err := api.ExtractImagesRaw(file, list, pdfcpu.NewDefaultConfiguration())
	if err != nil {
		fmt.Println(err)
	}
	for k, i := range imageList {
		imageOut, err := imaging.Decode(i)
		if err != nil {
			fmt.Println(err)
			return
		}
		//保存文件
		err = imaging.Save(imageOut, "test/"+strconv.Itoa(k+1)+".jpg")
		if err != nil {
			fmt.Println(err)
		}
	}
}
