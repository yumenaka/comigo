//go:build wails

package tailscale_plugin

import (
	"context"
	"net/netip"
	"runtime"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type TailsClientInfo struct {
	LoginUserName     string    `json:"login_user_name"`     // 登录用户身份
	NodeComputedName  string    `json:"node_computed_name"`  // 节点的机器名
	RequestRemoteAddr string    `json:"request_remote_addr"` // 远程访问地址
	AccessTime        time.Time `json:"access_time"`         // 客户端访问开始、时间
}

// TailscaleStatus 保存 Tailscale 服务的状态信息。
type TailscaleStatus struct {
	AuthURL          string
	BackendState     string
	Clients          []TailsClientInfo
	OS               string
	Online           bool
	FQDN             string
	TailscaleIPs     []netip.Addr
	Version          string
	FunnelCapability string
	mu               sync.Mutex
}

type TailscaleConfig struct {
	Hostname   string
	Port       uint16
	FunnelMode bool
	ConfigDir  string
	AuthKey    string
}

// RunTailscale 在 Wails 构建中禁用 Tailscale，避免桌面壳引入 tsnet 运行时依赖。
// ponytail: Wails 需要真实 Tailscale 时，再把 tsnet 依赖适配到桌面打包链路。
func RunTailscale(e *echo.Echo, c TailscaleConfig) error {
	return nil
}

// GetTailscaleIP 返回当前 Tailscale IP；Wails 空实现固定为空。
func (t *TailscaleStatus) GetTailscaleIP() string {
	return ""
}

// StopTailscale 停止Tailscale服务器。
func StopTailscale() error {
	return nil
}

// GetTailscaleStatus 返回 Wails 构建下的离线状态。
func GetTailscaleStatus(ctx context.Context) (*TailscaleStatus, error) {
	return &TailscaleStatus{
		BackendState:     "Stopped",
		OS:               runtime.GOOS,
		FunnelCapability: "unknown",
	}, nil
}
