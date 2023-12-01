package util

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"strconv"

	"github.com/bbrks/go-blurhash"
	"github.com/disintegration/imaging"
	"github.com/mandykoh/autocrop"
	"github.com/yumenaka/comi/logger"
)

// GetImageDataBlurHash  获取图片的BlurHash
func GetImageDataBlurHash(loadedImage []byte, components int) string {
	// Generate the BlurHash for a given image
	buf := bytes.NewBuffer(loadedImage)
	imageData, err := imaging.Decode(buf)
	if err != nil {
		logger.Info(err)
		return "error blurhash!"
	}
	str, err := blurhash.Encode(components, components, imageData)
	if err != nil {
		// Handle errors
		logger.Info(err)
		return "error blurhash!"
	}
	logger.Infof("Hash: %s\n", str)
	return str
}

// GetImageDataBlurHashImage 获取图片的BlurHash图
func GetImageDataBlurHashImage(loadedImage []byte, components int) []byte {
	// Generate the BlurHash for a given image
	buf := bytes.NewBuffer(loadedImage)
	imageData, err := imaging.Decode(buf)
	if err != nil {
		logger.Info(err)
	}
	str, err := blurhash.Encode(components, components, imageData)
	if err != nil {
		logger.Info(err)
	}
	// Generate an imageData for a given BlurHash
	// Punch specifies the contrasts and defaults to 1
	img, err := blurhash.Decode(str, imageData.Bounds().Dx(), imageData.Bounds().Dy(), 1)
	if err != nil {
		logger.Info(err)
	}
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, img, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageResizeByWidth 根据一个固定宽度缩放图片
func ImageResizeByWidth(loadedImage []byte, width int) []byte {
	buf := bytes.NewBuffer(loadedImage)
	decode, err := imaging.Decode(buf)
	if err != nil {
		logger.Info(err)
		return loadedImage
	}
	sourceWidth := decode.Bounds().Dx()
	scalingRatio := float64(width) / float64(sourceWidth)
	height := int(float64(decode.Bounds().Dy()) * scalingRatio)
	//生成缩略图
	decode = imaging.Resize(decode, width, height, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, decode, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageResizeByMaxWidth  设定一个图片宽度上限，大于这个宽度就缩放
func ImageResizeByMaxWidth(loadedImage []byte, maxWidth int) ([]byte, error) {
	buf := bytes.NewBuffer(loadedImage)
	decode, err := imaging.Decode(buf)
	if err != nil {
		logger.Info(err)
		return nil, errors.New("imaging.Decode() Error")
	}
	sourceWidth := decode.Bounds().Dx()
	if maxWidth > sourceWidth {
		return nil, errors.New("ImageResizeByMaxWidth Error maxWidth(" + strconv.Itoa(maxWidth) + ")> sourceWidth(" + strconv.Itoa(sourceWidth) + ")")
	}
	scalingRatio := float64(maxWidth) / float64(sourceWidth)
	height := int(float64(decode.Bounds().Dy()) * scalingRatio)
	//生成缩略图
	decode = imaging.Resize(decode, maxWidth, height, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, decode, imaging.JPEG)
	if err != nil {
		return nil, errors.New("imaging.Encode() Error")
	}
	return buf2.Bytes(), nil
}

// ImageResizeByMaxHeight  设定一个图片高度上限，大于这个高度就缩放
func ImageResizeByMaxHeight(loadedImage []byte, maxHeight int) ([]byte, error) {
	buf := bytes.NewBuffer(loadedImage)
	decode, err := imaging.Decode(buf)
	if err != nil {
		logger.Info(err)
		return nil, errors.New("imaging.Decode() Error")
	}
	sourceHeight := decode.Bounds().Dy()
	if maxHeight > sourceHeight {
		return nil, errors.New("ImageResizeByMaxHeight Error maxWidth(" + strconv.Itoa(maxHeight) + ")> sourceWidth(" + strconv.Itoa(sourceHeight) + ")")
	}
	scalingRatio := float64(maxHeight) / float64(sourceHeight)
	width := int(float64(decode.Bounds().Dx()) * scalingRatio)
	decode = imaging.Resize(decode, width, maxHeight, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, decode, imaging.JPEG)
	if err != nil {
		return nil, errors.New("imaging.Encode() Error")
	}
	return buf2.Bytes(), nil
}

// ImageResizeByHeight 根据一个固定 Height 缩放图片
func ImageResizeByHeight(loadedImage []byte, height int) []byte {
	buf := bytes.NewBuffer(loadedImage)
	decode, err := imaging.Decode(buf)
	if err != nil {
		logger.Info(err)
		return loadedImage
	}
	sourceHeight := decode.Bounds().Dy()
	scalingRatio := float64(height) / float64(sourceHeight)
	width := int(float64(decode.Bounds().Dx()) * scalingRatio)
	decode = imaging.Resize(decode, width, height, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, decode, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageResize 重设图片分辨率
func ImageResize(loadedImage []byte, width int, height int) []byte {
	//loadedImage, _ := ioutil.ReadFile("d:/1.jpg")
	buf := bytes.NewBuffer(loadedImage)
	decode, err := imaging.Decode(buf)
	if err != nil {
		logger.Info(err)
		return loadedImage
	}
	//生成缩略图，尺寸width*height
	decode = imaging.Resize(decode, width, height, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, decode, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageThumbnail 根据设定的图片大小,剪裁图片
func ImageThumbnail(loadedImage []byte, width int, height int) []byte {
	buf := bytes.NewBuffer(loadedImage)
	imageData, err := imaging.Decode(buf)
	if err != nil {
		logger.Info(err)
		return loadedImage
	}
	//生成缩略图，尺寸width*height
	imageData = imaging.Thumbnail(imageData, width, height, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, imageData, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageAutoCrop  自动裁白边
func ImageAutoCrop(loadedImage []byte, energyThreshold float32) []byte {
	////读取本地文件，本地文件尺寸300*400
	//loadedImage, _ := ioutil.ReadFile("d:/1.jpg")
	buf := bytes.NewBuffer(loadedImage)
	img, err := imaging.Decode(buf)
	if err != nil {
		logger.Info(err)
		return loadedImage
	}
	//使用 BoundsForThreshold 查找图像的自动裁剪边界
	//croppedBounds := autocrop.BoundsForThreshold(image, energyThreshold/100)

	nRGBAImg := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(nRGBAImg, nRGBAImg.Bounds(), img, img.Bounds().Min, draw.Src)
	result := autocrop.ToThreshold(nRGBAImg, energyThreshold/100)
	//如果不需要边界，可以使用ToThreshold函数方便地获得裁剪图像
	//croppedImg := autocrop.ToThreshold(image, energyThreshold)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, result, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageGray 转换为黑白图片
func ImageGray(loadedImage []byte) []byte {
	////读取本地文件，本地文件尺寸300*400
	//loadedImage, _ := ioutil.ReadFile("d:/1.jpg")
	buf := bytes.NewBuffer(loadedImage)
	img, err := imaging.Decode(buf)
	if err != nil {
		logger.Info(err)
		return loadedImage
	}
	result := imaging.Grayscale(img)
	//如果不需要边界，可以使用ToThreshold函数方便地获得裁剪图像
	//croppedImg := autocrop.ToThreshold(image, energyThreshold)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, result, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}
