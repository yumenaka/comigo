package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/types"
	"github.com/yumenaka/comi/util"
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

func GetServerInfoPublic(c *gin.Context) {
	serverName := "Comigo " + config.Version
	//本机首选出站IP
	OutIP := util.GetOutboundIP().String()
	host := ""
	if config.Config.Host == "DefaultHost" {
		host = OutIP
	} else {
		host = config.Config.Host
	}
	var serverStatus = ServerStatus{
		ServerName:        serverName,
		ServerHost:        host,
		ServerPort:        config.Config.Port,
		SupportUploadFile: config.Config.EnableUpload,
		NumberOfBooks:     types.GetBooksNumber(),
	}
	c.PureJSON(http.StatusOK, serverStatus)
}

func GetServerInfo(c *gin.Context) {
	serverName := "Comigo " + config.Version
	//本机首选出站IP
	host := ""
	if config.Config.Host == "DefaultHost" {
		host = util.GetOutboundIP().String()
	} else {
		host = config.Config.Host
	}
	var serverStatus = ServerStatus{
		ServerName:            serverName,
		ServerHost:            host,
		ServerPort:            config.Config.Port,
		SupportUploadFile:     config.Config.EnableUpload,
		NumberOfBooks:         types.GetBooksNumber(),
		NumberOfOnLineUser:    1,
		NumberOfOnLineDevices: 1,
		ClientIP:              c.ClientIP(),
		OSInfo:                util.GetSystemStatus(),
	}
	c.PureJSON(http.StatusOK, serverStatus)
}
