package data_api

import (
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

// GetQrcode 下载服务器配置
func GetQrcode(c echo.Context) error {
	// 通过参数传递自定义文本，生成二维码
	qrcodeStr := c.QueryParam("qrcode_str")
	base64Encode := getBoolQueryParam(c, "base64", false)
	if qrcodeStr != "" {
		png, err := qrcode.Encode(qrcodeStr, qrcode.Medium, 256)
		if err != nil {
			logger.Infof("%s", err)
			return err
		}
		if base64Encode {
			// 返回Base64编码的二维码
			base64Str := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
			return c.String(http.StatusOK, base64Str)
		}
		return c.Blob(http.StatusOK, "image/png", png)
	}
	// 根据配置文件中的URL，生成二维码
	qrcodeStr = config.GetQrcodeURL()
	png, err := qrcode.Encode(qrcodeStr, qrcode.Medium, 256)
	if err != nil {
		logger.Infof("%s", err)
		return err
	}
	if base64Encode {
		// 返回Base64编码的二维码
		base64Str := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
		return c.String(http.StatusOK, base64Str)
	}
	return c.Blob(http.StatusOK, "image/png", png)
}
