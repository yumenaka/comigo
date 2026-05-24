package data_api

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
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools/logger"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//go:embed Roboto-Medium.ttf
var fontBytes []byte

type generatedImageRequest struct {
	width    int
	height   int
	text     string
	fontSize float64
	bgColor  string
}

func GetGeneratedImage(c echo.Context) error {
	req, err := parseGeneratedImageRequest(c)
	if err != nil {
		return writeValidationError(c, err)
	}
	logger.Infof(locale.GetString("log_get_generated_image_params"),
		strconv.Itoa(req.height), strconv.Itoa(req.width), req.text, strconv.FormatFloat(req.fontSize, 'f', -1, 64))

	// 解析字体
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return apiresp.Error(c, http.StatusInternalServerError, "parse_font_failed", "Could not parse font", err.Error())
	}

	// 创建图像
	rgba := image.NewRGBA(image.Rect(0, 0, req.width, req.height))

	// 解析背景颜色
	bgColor := color.RGBA{
		R: 166,
		G: 166,
		B: 166,
		A: 255,
	} // 默认背景为灰色
	if req.bgColor != "" {
		bgColorParsed, err := ParseHexColor(req.bgColor)
		if err != nil {
			return apiresp.BadRequest(c, "invalid_bg_color", "Invalid bg_color", map[string]string{"bg_color": req.bgColor})
		}
		bgColor = bgColorParsed
	}
	// 填充背景颜色
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 设置字体
	face := truetype.NewFace(f, &truetype.Options{
		Size: req.fontSize,
		DPI:  72,
	})

	// 准备绘制文本
	d := &font.Drawer{
		Face: face,
	}

	// 计算文本的宽度和高度
	textWidth := d.MeasureString(req.text).Round()
	metrics := face.Metrics()
	ascent := metrics.Ascent.Round()
	descent := metrics.Descent.Round()

	// 计算文本绘制的起始坐标，使其在图像中居中
	x := (req.width - textWidth) / 2
	y := (req.height + ascent - descent) / 2

	d = &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot: fixed.Point26_6{
			X: fixed.Int26_6(x << 6),
			Y: fixed.Int26_6(y << 6),
		},
	}
	d.DrawString(req.text)

	// 设置响应头并返回图片
	c.Response().Header().Set("Content-Type", "image/jpeg")
	if err := jpeg.Encode(c.Response().Writer, rgba, nil); err != nil {
		return apiresp.Error(c, http.StatusInternalServerError, "encode_image_failed", "Failed to encode image", err.Error())
	}

	return nil
}

// parseGeneratedImageRequest 在分配图片前限制尺寸、面积、字体和文本长度，避免生成接口被超大参数拖垮。
func parseGeneratedImageRequest(c echo.Context) (generatedImageRequest, error) {
	width, err := parseRequiredBoundedInt(c, "width", 1, imageQueryMaxDimension)
	if err != nil {
		return generatedImageRequest{}, err
	}
	height, err := parseRequiredBoundedInt(c, "height", 1, imageQueryMaxDimension)
	if err != nil {
		return generatedImageRequest{}, err
	}
	fontSize, err := parseRequiredBoundedFloat(c, "font_size", generatedImageMinFontSize, generatedImageMaxFontSize)
	if err != nil {
		return generatedImageRequest{}, err
	}
	text := c.QueryParam("text")
	if err := validateGeneratedImageCost(width, height, text); err != nil {
		return generatedImageRequest{}, err
	}
	return generatedImageRequest{
		width:    width,
		height:   height,
		text:     text,
		fontSize: fontSize,
		bgColor:  c.QueryParam("bg_color"),
	}, nil
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
