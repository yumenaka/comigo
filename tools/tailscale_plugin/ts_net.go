package tailscale_plugin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/tools/logger"
	"tailscale.com/client/local"
	"tailscale.com/tsnet"
)

// 参考文档： https://tailscale.com/kb/1244/tsnet
// https://github.com/tailscale/tailscale/tree/v1.86.5/tsnet/example

// 比较复杂的例子：
// https://github.com/tailscale/golink/blob/main/golink.go

// Create AUTH_KEY：
// https://login.tailscale.com/admin/settings/keys
// TS_AUTHKEY="tskey-auth-<key>" go run ./cmd/comi

// Tailscale 数据文件存储在UserConfigDir （https://pkg.go.dev/os#UserConfigDir）中。只要可以保存，就只需要在首次运行时提供身份验证密钥。
// Linux： $HOME/.config.
// Darwin  $HOME/Library/Application Support.
// Windows： %AppData% （like： C:\Users\[ユーザー名]\AppData）

var (
	tsServer    *tsnet.Server
	netListener net.Listener
	localClient *local.Client
	httpServer  *http.Server
	nowStatus   *TailscaleStatus
)

func init() {
	nowStatus = &TailscaleStatus{
		Online: false,
	}
}

// ctxKeyIsTailscale 用于标记来自 Tailscale 监听器的请求
const ctxKeyIsTailscale = "is_tailscale_request"

// withTailscaleContext 为经过 Tailscale 监听器的请求注入上下文标记
// 这样中间件可以识别这些请求，让中间件仅在检测到该标记时执行 WhoIs 查询（普通请求触发 Tailscale WhoIs 查询会导致2秒左右的延迟）
func withTailscaleContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxKeyIsTailscale, true)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RunTailscale 启动Tailscale网络服务器，统一使用echo处理请求
func RunTailscale(e *echo.Echo, c InitConfig) error {
	// 初始化Tailscale服务器
	if err := InitTailscale(c); err != nil {
		return fmt.Errorf("failed to initialize Tailscale server: %w", err)
	}

	// 添加Tailscale身份验证信息中间件
	e.Use(TailscaleAuthMiddleware())

	// 创建HTTP服务器，使用echo.Echo作为处理器
	httpServer = &http.Server{
		Addr:    c.Hostname + ":" + fmt.Sprint(c.Port),
		Handler: withTailscaleContext(e), // echo.Echo (多包一层)
	}
	// 使用Tailscale网络监听器启动服务器
	logger.Infof("Starting Tailscale HTTP server on %s:%d", c.Hostname, c.Port)
	go func() {
		if err := httpServer.Serve(netListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				logger.Errorf("Tailscale HTTP server error: %v", err)
			}
		}
	}()
	return nil
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
	return nil
}

// firstLabel 提取主机名的第一个标签
func firstLabel(s string) string {
	s, _, _ = strings.Cut(s, ".")
	return s
}
