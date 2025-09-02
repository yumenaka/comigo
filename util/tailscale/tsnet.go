package tailscale

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/util/logger"
	"tailscale.com/client/local"
	"tailscale.com/tsnet"
)

var (
	tsServer    *tsnet.Server
	netListener net.Listener
	localClient *local.Client
	httpServer  *http.Server
)

func InitTailscale(hostname string, port uint16) error {
	// 设置 Tailscale 服务器的主机名。此名称将用于 Tailscale 网络中的节点标识。
	// 如果未设置，Tailscale 将使用机器的主机名。默认将会是二进制文件名。
	tsServer = new(tsnet.Server)
	// 用于Tailscale网络中的标识节点，不影响监听地址
	tsServer.Hostname = hostname

	if port == 0 {
		port = 443
	}
	tsServer.Port = port

	// 监听器 ln 是一个 net.Listener 对象，它将处理来自 Tailscale网络的 TCP 连接。
	// Tailscale的Listen方法要求host部分必须是空的或者是IP字面量，不能使用主机名
	var err error
	listenAddr := ":" + fmt.Sprint(port)
	netListener, err = tsServer.Listen("tcp", listenAddr)
	if err != nil {
		logger.Errorf("Failed to create Tailscale listener on %s: %v", listenAddr, err)
		return err
	}

	// LocalClient 返回一个与 s 通信的 LocalClient 对象。
	// 如果服务器尚未启动，它将启动服务器。如果服务器已成功启动，
	// 它不会返回错误。
	localClient, err = tsServer.LocalClient()
	if err != nil {
		logger.Errorf("Failed to create Tailscale local client: %v", err)
		netListener.Close()
		return err
	}

	// 自动设置监听器 ln 的 TLS 证书
	if port == 443 {
		netListener = tls.NewListener(netListener, &tls.Config{
			GetCertificate: localClient.GetCertificate,
		})
	}

	logger.Infof("Tailscale server initialized successfully on %s:%d", hostname, port)
	return nil
}

// RunTailscale 启动Tailscale网络服务器，统一使用echo处理请求
func RunTailscale(e *echo.Echo, hostname string, port int) error {
	// 初始化Tailscale服务器
	if err := InitTailscale(hostname, uint16(port)); err != nil {
		return fmt.Errorf("failed to initialize Tailscale server: %w", err)
	}

	// 添加Tailscale身份验证中间件
	e.Use(TailscaleAuthMiddleware())

	// 添加Tailscale信息路由
	e.GET("/tailscale/whois", TailscaleWhoIsHandler)

	// 添加根路径路由，提供基本的欢迎页面
	e.GET("/tailscale", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `<html><body>
			<h1>欢迎使用 Tailscale 网络服务</h1>
			<p>服务正在运行中...</p>
			<p><a href="/tailscale/whois">查看身份信息</a></p>
		</body></html>`)
	})

	// 创建HTTP服务器，使用echo.Echo作为处理器
	httpServer = &http.Server{
		Addr:    hostname + ":" + fmt.Sprint(port),
		Handler: e, // echo.Echo 实现了 http.Handler 接口
	}
	// 使用Tailscale网络监听器启动服务器
	logger.Infof("Starting Tailscale HTTP server on %s:%d", hostname, port)
	go func() {
		if err := httpServer.Serve(netListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				logger.Errorf("Tailscale HTTP server error: %v", err)
			}
		}
	}()
	return nil
}

// TailscaleAuthMiddleware 创建Tailscale身份验证中间件
func TailscaleAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 获取请求者的Tailscale身份信息
			who, err := localClient.WhoIs(c.Request().Context(), c.Request().RemoteAddr)
			// 如果可以获取身份信息
			if err == nil {
				// 将访问者的Tailscale信息存储到上下文
				c.Set("tailscale_user", who.UserProfile.LoginName)
				c.Set("tailscale_node", firstLabel(who.Node.ComputedName))
				c.Set("tailscale_remote_addr", c.Request().RemoteAddr)
			}
			// // 如果可以获取不到访问者的Tailscale信息
			// if err != nil  {
			// 	logger.Infof("Failed to get Tailscale identity for %s: %v", c.Request().RemoteAddr, err)
			// }
			// 继续处理请求
			return next(c)
		}
	}
}

// TailscaleWhoIsHandler 处理Tailscale身份信息查询请求
func TailscaleWhoIsHandler(c echo.Context) error {
	// 从上下文获取Tailscale身份信息
	user, ok := c.Get("tailscale_user").(string)
	if !ok {
		return c.String(http.StatusInternalServerError, "无法获取Tailscale用户信息")
	}

	node, ok := c.Get("tailscale_node").(string)
	if !ok {
		return c.String(http.StatusInternalServerError, "无法获取Tailscale节点信息")
	}

	remoteAddr, ok := c.Get("tailscale_remote_addr").(string)
	if !ok {
		return c.String(http.StatusInternalServerError, "无法获取远程地址信息")
	}

	// 返回HTML格式的身份信息
	html := fmt.Sprintf(`<html><body><h1>Hello, Tailscale!</h1>
<p>您是以 <b>%s</b> 身份从 <b>%s</b> (%s) 访问的</p>
<p>欢迎使用Tailscale网络服务！</p>
</body></html>`, user, node, remoteAddr)

	return c.HTML(http.StatusOK, html)
}

// StopTailscale 停止Tailscale服务器
func StopTailscale() error {
	if httpServer != nil {
		logger.Infof("Stopping Tailscale HTTP server...")
		if err := httpServer.Close(); err != nil {
			logger.Errorf("Error closing HTTP server: %v", err)
			return err
		}
		httpServer = nil
	}

	if netListener != nil {
		if err := netListener.Close(); err != nil {
			logger.Errorf("Error closing network listener: %v", err)
			return err
		}
		netListener = nil
	}

	if tsServer != nil {
		tsServer.Close()
		tsServer = nil
	}

	logger.Infof("Tailscale server stopped successfully")
	return nil
}

// GetTailscaleStatus 获取Tailscale服务器状态
func GetTailscaleStatus() map[string]interface{} {
	status := map[string]interface{}{
		"server_running":   httpServer != nil,
		"listener_active":  netListener != nil,
		"ts_server_active": tsServer != nil,
		"local_client_ok":  localClient != nil,
	}

	if tsServer != nil {
		status["hostname"] = tsServer.Hostname
		status["port"] = tsServer.Port
	}

	return status
}

// firstLabel 提取主机名的第一个标签
func firstLabel(s string) string {
	s, _, _ = strings.Cut(s, ".")
	return s
}
