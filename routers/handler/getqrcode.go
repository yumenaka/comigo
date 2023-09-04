package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"

	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/tools"
)

// GetQrcodeHandler 下载服务器配置
func GetQrcodeHandler(c *gin.Context) {
	//通过参数过去自定义文本的二维码，更通用
	qrcodeStr := c.DefaultQuery("qrcode_str", "")
	if qrcodeStr != "" {
		png, err := qrcode.Encode(qrcodeStr, qrcode.Medium, 256)
		if err != nil {
			fmt.Println(err)
		}
		c.Data(http.StatusOK, "image/png", png)
		return
	}

	//cmd打印链接二维码
	enableTls := common.Config.CertFile != "" && common.Config.KeyFile != ""
	protocol := "http://"
	if enableTls {
		protocol = "https://"
	}
	//取得本机的首选出站IP
	OutIP := tools.GetOutboundIP().String()
	if common.Config.Host == "DefaultHost" {
		var png []byte
		readURL := protocol + OutIP + ":" + strconv.Itoa(common.Config.Port)
		png, err := qrcode.Encode(readURL, qrcode.Medium, 256)
		if err != nil {
			fmt.Println(err)
		}
		c.Data(http.StatusOK, "image/png", png)
	} else {
		var png []byte
		readURL := protocol + common.Config.Host + ":" + strconv.Itoa(common.Config.Port)
		png, err := qrcode.Encode(readURL, qrcode.Medium, 256)
		if err != nil {
			fmt.Println(err)
		}
		c.Data(http.StatusOK, "image/png", png)
	}

}
