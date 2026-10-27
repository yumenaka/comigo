package assets

import (
	"strings"
	"testing"

	"github.com/yumenaka/comigo/config"
)

// TestWailsWebSocketURL 验证桌面地址包含实际端口、基础路径和 TLS，普通 Web 不注入本机地址。
func TestWailsWebSocketURL(t *testing.T) {
	previous := *config.GetCfg()
	t.Cleanup(func() { config.ReplaceConfig(previous) })
	for _, tt := range []struct {
		name, basePath, host, want string
		tls, autoTLS               bool
	}{
		{name: "自定义端口", want: "ws://127.0.0.1:23451/api/ws"},
		{name: "反代基础路径", basePath: "/reader", want: "ws://127.0.0.1:23451/reader/api/ws"},
		{name: "TLS", tls: true, basePath: "/reader", want: "wss://127.0.0.1:23451/reader/api/ws"},
		{name: "自动证书域名", autoTLS: true, host: "reader.example.com", want: "wss://reader.example.com:443/api/ws"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg := previous
			cfg.Port, cfg.BasePath, cfg.Host = 23451, tt.basePath, tt.host
			cfg.EnableTLS, cfg.AutoTLSCertificate = tt.tls, tt.autoTLS
			config.ReplaceConfig(cfg)
			if !strings.Contains(GetBasePathScript(true), `window.ComiGoWailsWebSocketURL = "`+tt.want+`";`) {
				t.Fatalf("桌面 WebSocket 地址应为 %s", tt.want)
			}
			if !strings.Contains(GetBasePathScript(false), `window.ComiGoWailsWebSocketURL = "";`) {
				t.Fatal("普通 Web 不应使用桌面 WebSocket 地址")
			}
		})
	}
}
