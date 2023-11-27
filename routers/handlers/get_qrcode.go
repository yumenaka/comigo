package handlers

import (
	"github.com/yumenaka/comi/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"

	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/util"
)

// GetQrcodeHandler 下载服务器配置
func GetQrcodeHandler(c *gin.Context) {
	//通过参数过去自定义文本的二维码，更通用
	qrcodeStr := c.DefaultQuery("qrcode_str", "")
	if qrcodeStr != "" {
		png, err := qrcode.Encode(qrcodeStr, qrcode.Medium, 256)
		if err != nil {
			logger.Info(err)
		}
		c.Data(http.StatusOK, "image/png", png)
		return
	}

	//cmd打印链接二维码
	enableTls := config.Config.CertFile != "" && config.Config.KeyFile != ""
	protocol := "http://"
	if enableTls {
		protocol = "https://"
	}
	//取得本机的首选出站IP
	OutIP := util.GetOutboundIP().String()
	if config.Config.Host == "DefaultHost" {
		var png []byte
		readURL := protocol + OutIP + ":" + strconv.Itoa(config.Config.Port)
		png, err := qrcode.Encode(readURL, qrcode.Medium, 256)
		if err != nil {
			logger.Info(err)
		}
		c.Data(http.StatusOK, "image/png", png)
	} else {
		var png []byte
		readURL := protocol + config.Config.Host + ":" + strconv.Itoa(config.Config.Port)
		png, err := qrcode.Encode(readURL, qrcode.Medium, 256)
		if err != nil {
			logger.Info(err)
		}
		c.Data(http.StatusOK, "image/png", png)
	}

}
