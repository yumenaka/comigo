package file

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// sample code: https://github.com/pdfcpu/pdfcpu/blob/master/pkg/api/extract.go

// CountPagesOfPDF 确定PDF的页数
func CountPagesOfPDF(pdfFileName string) (count int, err error) {
	// 底层解析 panic 也必须作为失败返回，不能伪装成零页成功。
	defer func() {
		if r := recover(); r != nil {
			logger.Infof(locale.GetString("log_countpages_pdf_invalid_error"), pdfFileName, r)
			count = -1
			err = fmt.Errorf(locale.GetString("err_countpages_pdf_invalid"), pdfFileName, fmt.Sprint(r))
		}
	}()
	// use default configuration for pdfcpu ("nil")
	// pdfcpu 的文件接口需要上下文；本地同步读取使用后台上下文。
	ctx := context.Background()
	err = api.ValidateFile(ctx, pdfFileName, nil, nil)
	if err != nil {
		return -1, fmt.Errorf(locale.GetString("err_countpages_pdf_invalid"), pdfFileName, err.Error())
	}
	return api.PageCountFile(ctx, pdfFileName)
}

// GetImageFromPDF 从PDF里面取jpeg文件，pageNum从1开始
func GetImageFromPDF(pdfFileName string, pageNum int, Debug bool) ([]byte, error) {
	start := time.Now()
	// 打开PDF文件
	file, err := os.Open(pdfFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// use default configuration for pdfcpu
	pdfSetting := model.NewDefaultConfiguration()
	pdfSetting.DecodeAllStreams = true
	buffer := &bytes.Buffer{}
	err = api.ExtractImages(context.Background(), file, []string{strconv.Itoa(pageNum)}, digestImage(buffer), pdfSetting)
	if err != nil {
		return nil, err
	}
	// 参考：https://github.com/pdfcpu/pdfcpu/issues/45
	// https://github.com/pdfcpu/pdfcpu/issues/351
	// 将PDF转换为图像，相当于重写一个pdf查看器。
	// 虽然可以引用mupdf。但这会导致导入C代码，破坏跨平台兼容性。
	// 或许可以cli调用imagemagick，曲线救国？
	if Debug {
		logger.Infof(locale.GetString("log_getimagefrompdf_time"), time.Now().Sub(start))
	}

	return buffer.Bytes(), nil
}

// 自定义钩子函数参数的方法：输入自定义参数、输出符合要求的函数
func digestImage(buff *bytes.Buffer) func(model.Image, bool, int) error {
	return func(img model.Image, singleImgPerPage bool, maxPageDigits int) error {
		imageOut, err := imaging.Decode(img)
		if err != nil {
			return errors.New(locale.GetString("err_imaging_decode_error"))
		}
		err = jpeg.Encode(buff, imageOut, &jpeg.Options{Quality: 75})
		if err != nil {
			return errors.New(locale.GetString("err_jpeg_encode_error"))
		}
		return nil
	}
}
