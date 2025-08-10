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
	ServerHost            string            // 服务器地址
	ServerPort            int               // 服务器端口（程序运行的端口）
	NumberOfBooks         int               // 当前拥有的书籍总数
	NumberOfOnLineUser    int               // 在线用户数（未实现）
	NumberOfOnLineDevices int               // 在线设备数（未实现）
	SupportUploadFile     bool              // 是否支持上传文件
	ClientIP              string            // 客户端IP
	OSInfo                util.SystemStatus // 系统信息
}

func GetServerInfoHandler(c echo.Context) error {
	serverStatus := util.GetServerInfo(config.GetHost(), config.GetVersion(), config.GetPort(), config.GetEnableUpload(), model.MainStoreGroup.GetBooksNumber())
	return c.JSON(http.StatusOK, serverStatus)
}

func GetAllServerInfoHandler(c echo.Context) error {
	serverStatus := util.GetAllServerInfo(config.GetHost(), config.GetVersion(), config.GetPort(), config.GetEnableUpload(), model.MainStoreGroup.GetBooksNumber(), c.RealIP())
	return c.JSON(http.StatusOK, serverStatus)
}
