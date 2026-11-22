//go:build !wails

package tailscale_plugin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
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
	lifecycleMu   sync.RWMutex
	cancelStartup context.CancelFunc
	tsServer      *tsnet.Server
	netListener   net.Listener
	localClient   *local.Client
	tsHttpServer  *http.Server
	nowStatus     = &TailscaleStatus{}
)

// tailscaleClientKey 让请求始终使用接收该连接的实例查询身份。
type tailscaleClientKey struct{}

// RunTailscale 启动Tailscale网络服务器，统一使用echo处理请求
func RunTailscale(e *echo.Echo, c TailscaleConfig) error {
	lifecycleMu.Lock()
	// CLI reload 和网页设置共用此入口；先释放旧实例，再按新配置创建监听。
	if err := stopTailscale(); err != nil {
		lifecycleMu.Unlock()
		return err
	}
	if c.FunnelMode && c.Port != 443 && c.Port != 8443 && c.Port != 10000 {
		lifecycleMu.Unlock()
		return errors.New(locale.GetString("err_funnel_mode_ports_only"))
	}
	instance := &tsnet.Server{Hostname: c.Hostname, Port: c.Port, Dir: c.ConfigDir, AuthKey: c.AuthKey, UserLogf: tailscaleUserLogf}
	tsServer = instance
	client, err := instance.LocalClient() // 启动完成后才允许 Stop 关闭 tsnet。
	if err != nil {
		err = errors.Join(err, stopTailscale())
		lifecycleMu.Unlock()
		return err
	}
	localClient = client
	// 等待 Funnel 登录时允许状态查询、停止和重启；只发布仍属于本次启动的监听。
	ctx, cancel := context.WithCancel(context.Background())
	cancelStartup = cancel
	lifecycleMu.Unlock()
	ln, err := listenTailscale(ctx, instance, client, c)
	cancel()
	lifecycleMu.Lock()
	defer lifecycleMu.Unlock()
	if tsServer != instance {
		if ln != nil {
			err = errors.Join(err, ln.Close())
		}
		return errors.Join(context.Canceled, err)
	}
	cancelStartup = nil
	if err != nil {
		return errors.Join(fmt.Errorf("failed to initialize Tailscale listener: %w", err), stopTailscale())
	}
	netListener = ln
	tsHttpServer = &http.Server{
		Handler: e,
		BaseContext: func(net.Listener) context.Context {
			return context.WithValue(context.Background(), tailscaleClientKey{}, client)
		},
	}
	// 使用Tailscale网络监听器启动服务器
	logger.Infof(locale.GetString("log_starting_tailscale_http_server"), c.Hostname, c.Port)
	// 捕获本次实例，避免重载改写全局变量后，旧协程误用新监听器。
	server, listener := tsHttpServer, netListener
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				logger.Errorf(locale.GetString("err_tailscale_http_server_error"), err)
			}
		}
	}()
	return nil
}

// StopTailscale 停止Tailscale服务器
func StopTailscale() error {
	lifecycleMu.Lock()
	defer lifecycleMu.Unlock()
	return stopTailscale()
}

// stopTailscale 在持有生命周期锁时释放全部资源，单项关闭失败不跳过其余清理。
func stopTailscale() error {
	if cancelStartup != nil {
		cancelStartup()
		cancelStartup = nil
	}
	if tsServer == nil && tsHttpServer == nil && netListener == nil && localClient == nil {
		return nil
	}
	var errs []error
	if tsHttpServer != nil {
		errs = append(errs, tsHttpServer.Close())
		tsHttpServer = nil
	}
	if netListener != nil {
		if err := netListener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			errs = append(errs, err)
		}
		netListener = nil
	}
	if tsServer != nil {
		errs = append(errs, tsServer.Close())
		tsServer = nil
	}
	localClient = nil
	err := errors.Join(errs...)
	if err == nil {
		logger.Info(locale.GetString("log_tailscale_server_stopped_successfully"))
	}
	return err
}

// firstLabel 提取主机名的第一个标签
func firstLabel(s string) string {
	s, _, _ = strings.Cut(s, ".")
	return s
}
