package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/logger"
)

// GetQrcode 下载服务器配置
func GetQrcode(c echo.Context) error {
	// 通过参数传递自定义文本，生成二维码
	qrcodeStr := c.QueryParam("qrcode_str")
	if qrcodeStr != "" {
		png, err := qrcode.Encode(qrcodeStr, qrcode.Medium, 256)
		if err != nil {
			logger.Infof("%s", err)
			return err
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
	return c.Blob(http.StatusOK, "image/png", png)
}
