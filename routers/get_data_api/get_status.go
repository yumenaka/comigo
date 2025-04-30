package get_data_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util"
)

// ServerStatus 服务器当前状况
type ServerStatus struct {
	ServerName            string            // 服务器描述
	ServerHost            string            //
	ServerPort            int               //
	NumberOfBooks         int               // 当前拥有的书籍总数
	NumberOfOnLineUser    int               // TODO：在线用户数
	NumberOfOnLineDevices int               // TODO：在线设备数
	SupportUploadFile     bool              //
	ClientIP              string            // 客户端IP
	OSInfo                util.SystemStatus // 系统信息
}

func GetServerInfoHandler(c echo.Context) error {
	serverStatus := util.GetServerInfo(config.GetHost(), config.GetVersion(), config.GetPort(), config.GetEnableUpload(), model.GetBooksNumber())
	return c.JSON(http.StatusOK, serverStatus)
}

func GetAllServerInfoHandler(c echo.Context) error {
	serverStatus := util.GetAllServerInfo(config.GetHost(), config.GetVersion(), config.GetPort(), config.GetEnableUpload(), model.GetBooksNumber(), c.RealIP())
	return c.JSON(http.StatusOK, serverStatus)
}
