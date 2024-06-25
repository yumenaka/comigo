package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/util/logger"
	"net/http"
)

// GetQrcode 下载服务器配置
func GetQrcode(c *gin.Context) {
	//通过参数传递自定义文本，生成二维码
	qrcodeStr := c.DefaultQuery("qrcode_str", "")
	if qrcodeStr != "" {
		png, err := qrcode.Encode(qrcodeStr, qrcode.Medium, 256)
		if err != nil {
			logger.Infof("%s", err)
		}
		c.Data(http.StatusOK, "image/png", png)
		return
	}
	//根据配置文件中的URL，生成二维码
	qrcodeStr = config.GetQrcodeURL()
	png, err := qrcode.Encode(qrcodeStr, qrcode.Medium, 256)
	if err != nil {
		logger.Infof("%s", err)
	}
	c.Data(http.StatusOK, "image/png", png)
}
