package get_data_api

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/tools/logger"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//go:embed Roboto-Medium.ttf
var fontBytes []byte

func GetGeneratedImage(c echo.Context) error {
	// 获取请求参数
	heightStr := c.QueryParam("height")
	widthStr := c.QueryParam("width")
	text := c.QueryParam("text")
	fontSizeStr := c.QueryParam("font_size")
	logger.Infof("GetGeneratedImage: height=%s, width=%s, text=%s, font_size=%s",
		heightStr, widthStr, text, fontSizeStr)
	bgColorStr := c.QueryParam("bg_color") // 背景颜色参数

	// 将高度和宽度转换为整数
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid height")
	}
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid width")
	}

	// 将字体大小转换为浮点数
	fontSize, err := strconv.ParseFloat(fontSizeStr, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid font_size")
	}

	// 解析字体
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not parse font")
	}

	// 创建图像
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	// 解析背景颜色
	bgColor := color.RGBA{
		R: 166,
		G: 166,
		B: 166,
		A: 255,
	} // 默认背景为灰色
	if bgColorStr != "" {
		bgColorParsed, err := ParseHexColor(bgColorStr)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid bg_color")
		}
		bgColor = bgColorParsed
	}
	// 填充背景颜色
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 设置字体
	face := truetype.NewFace(f, &truetype.Options{
		Size: fontSize,
		DPI:  72,
	})

	// 准备绘制文本
	d := &font.Drawer{
		Face: face,
	}

	// 计算文本的宽度和高度
	textWidth := d.MeasureString(text).Round()
	metrics := face.Metrics()
	ascent := metrics.Ascent.Round()
	descent := metrics.Descent.Round()

	// 计算文本绘制的起始坐标，使其在图像中居中
	x := (width - textWidth) / 2
	y := (height + ascent - descent) / 2

	d = &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot: fixed.Point26_6{
			X: fixed.Int26_6(x << 6),
			Y: fixed.Int26_6(y << 6),
		},
	}
	d.DrawString(text)

	// 设置响应头并返回图片
	c.Response().Header().Set("Content-Type", "image/jpeg")
	if err := jpeg.Encode(c.Response().Writer, rgba, nil); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to encode image")
	}

	return nil
}

// ParseHexColor 用于解析如 "#FFFFFF" 或 "FFFFFF" 格式的十六进制颜色字符串
func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	if strings.HasPrefix(s, "#") {
		s = s[1:]
	}
	switch len(s) {
	case 6:
		_, err = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 3:
		var r, g, b uint8
		_, err = fmt.Sscanf(s, "%1x%1x%1x", &r, &g, &b)
		c.R = r * 17
		c.G = g * 17
		c.B = b * 17
	default:
		err = fmt.Errorf("invalid length")
	}
	return
}
