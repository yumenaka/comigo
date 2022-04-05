package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/tools"
	"net/http"
	"strconv"
)

// ServerStatus 服务器当前状况
type ServerStatus struct {
	ServerName            string //服务器描述
	ServerHostAndPort     string
	NumberOfBooks         int                //当前拥有的书籍总数
	NumberOfOnLineUser    int                //TODO：在线用户数
	NumberOfOnLineDevices int                //TODO：在线设备数
	OSInfo                tools.SystemStatus //系统信息
}

func serverStatusHandler(c *gin.Context) {
	serverName := "Comigo " + common.Version
	//取得本机的首选出站IP
	OutIP := tools.GetOutboundIP().String()
	hostAndPort := ""
	if common.Config.Host == "" {
		hostAndPort = OutIP + ":" + strconv.Itoa(common.Config.Port)
	} else {
		hostAndPort = common.Config.Host + ":" + strconv.Itoa(common.Config.Port)
	}
	var serverStatus = ServerStatus{
		ServerName:            serverName,
		ServerHostAndPort:     hostAndPort,
		NumberOfBooks:         common.GetBooksNumber(),
		NumberOfOnLineUser:    1,
		NumberOfOnLineDevices: 1,
		OSInfo:                tools.GetSystemStatus(),
	}
	c.PureJSON(http.StatusOK, serverStatus)
}
