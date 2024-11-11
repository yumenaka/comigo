package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/util"
)

// ServerStatus 服务器当前状况
type ServerStatus struct {
	ServerName            string            //服务器描述
	ServerHost            string            //
	ServerPort            int               //
	NumberOfBooks         int               //当前拥有的书籍总数
	NumberOfOnLineUser    int               //TODO：在线用户数
	NumberOfOnLineDevices int               //TODO：在线设备数
	SupportUploadFile     bool              //
	ClientIP              string            //客户端IP
	OSInfo                util.SystemStatus //系统信息
}

func GetServerInfoHandler(c *gin.Context) {
	serverStatus := util.GetServerInfo(config.Config.Host, config.Version, config.Config.Port, config.Config.EnableUpload, entity.GetBooksNumber())
	c.PureJSON(http.StatusOK, serverStatus)
}

func GetAllServerInfoHandler(c *gin.Context) {
	serverStatus := util.GetAllServerInfo(config.Config.Host, config.Version, config.Config.Port, config.Config.EnableUpload, entity.GetBooksNumber(), c.ClientIP())
	c.PureJSON(http.StatusOK, serverStatus)
}
