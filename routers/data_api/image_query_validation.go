package data_api

import (
	"fmt"
	"math"
	"strconv"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/routers/apiresp"
)

const (
	imageQueryMaxDimension      = 4096
	generatedImageMaxArea       = imageQueryMaxDimension * imageQueryMaxDimension
	generatedImageMaxTextRunes  = 200
	generatedImageMaxFontSize   = 512
	generatedImageMinFontSize   = 1
	imageQueryMaxAutoCrop       = 100
	imageQueryMaxBlurComponents = 2
)

type requestValidationError struct {
	code    string
	message string
	details any
}

func (err requestValidationError) Error() string {
	return err.message
}

// writeValidationError 把 parse 阶段的验证错误统一输出成 API JSON，避免 helper 里提前写响应后继续执行。
func writeValidationError(c echo.Context, err error) error {
	if validationErr, ok := err.(requestValidationError); ok {
		return apiresp.BadRequest(c, validationErr.code, validationErr.message, validationErr.details)
	}
	return err
}

// parseOptionalBoundedInt 解析可选整数参数；参数一旦出现就必须合法，避免非法值静默退回默认值。
func parseOptionalBoundedInt(c echo.Context, key string, defaultValue, minValue, maxValue int) (int, error) {
	valueStr := c.QueryParam(key)
	if valueStr == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue, invalidImageQueryParam(key, valueStr, minValue, maxValue)
	}
	if value < minValue || value > maxValue {
		return defaultValue, invalidImageQueryParam(key, valueStr, minValue, maxValue)
	}
	return value, nil
}

// parseRequiredBoundedInt 解析必填整数参数，并统一限制图片生成尺寸。
func parseRequiredBoundedInt(c echo.Context, key string, minValue, maxValue int) (int, error) {
	valueStr := c.QueryParam(key)
	if valueStr == "" {
		return 0, requestValidationError{
			code:    "missing_param",
			message: key + " is required",
			details: map[string]string{"param": key},
		}
	}
	return parseOptionalBoundedInt(c, key, 0, minValue, maxValue)
}

// parseRequiredBoundedFloat 解析必填浮点参数，用于限制动态图片的字体大小。
func parseRequiredBoundedFloat(c echo.Context, key string, minValue, maxValue float64) (float64, error) {
	valueStr := c.QueryParam(key)
	if valueStr == "" {
		return 0, requestValidationError{
			code:    "missing_param",
			message: key + " is required",
			details: map[string]string{"param": key},
		}
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value < minValue || value > maxValue {
		return 0, invalidFloatImageQueryParam(key, valueStr, minValue, maxValue)
	}
	return value, nil
}

// validateGeneratedImageCost 限制动态图片生成的面积和文本长度，避免单个请求造成高内存/CPU 消耗。
func validateGeneratedImageCost(width, height int, text string) error {
	if width*height > generatedImageMaxArea {
		return requestValidationError{code: "image_too_large", message: "Generated image is too large", details: map[string]int{
			"width":    width,
			"height":   height,
			"max_area": generatedImageMaxArea,
		}}
	}
	if utf8.RuneCountInString(text) > generatedImageMaxTextRunes {
		return requestValidationError{code: "text_too_long", message: "Generated image text is too long", details: map[string]int{
			"max_runes": generatedImageMaxTextRunes,
		}}
	}
	return nil
}

func invalidImageQueryParam(key, value string, minValue, maxValue int) error {
	return requestValidationError{code: "invalid_image_param", message: "Invalid image query parameter", details: map[string]any{
		"param": key,
		"value": value,
		"min":   minValue,
		"max":   maxValue,
	}}
}

func invalidFloatImageQueryParam(key, value string, minValue, maxValue float64) error {
	return requestValidationError{code: "invalid_image_param", message: "Invalid image query parameter", details: map[string]any{
		"param": key,
		"value": value,
		"min":   fmt.Sprintf("%g", minValue),
		"max":   fmt.Sprintf("%g", maxValue),
	}}
}
