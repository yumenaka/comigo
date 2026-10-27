//go:build wails && !js

package wails_systray

import (
	"embed"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/routers/data_api"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/service"
	"github.com/yumenaka/comigo/tools/traffic"
)

//go:embed panel.html
var panelFiles embed.FS

// PanelService 只绑定本进程的面板操作，不提供任意路径或任意 URL 执行入口。
type PanelService struct{ tray *Tray }

// PanelState 仅公开面板需要的非敏感状态，绝不发送完整配置或密码。
type PanelState struct {
	Version        string                `json:"version"`
	URLs           []string              `json:"urls"`
	ReadingURL     string                `json:"readingURL"`
	QR             string                `json:"qr"`
	Books          int                   `json:"books"`
	Users          int                   `json:"users"`
	Devices        int                   `json:"devices"`
	Connections    int                   `json:"connections"`
	Traffic        traffic.Snapshot      `json:"traffic"`
	ExternalAccess bool                  `json:"externalAccess"`
	ReadOnly       bool                  `json:"readOnly"`
	ListenAddress  string                `json:"listenAddress"`
	Config         config.ConfigFileInfo `json:"config"`
	Stores         []string              `json:"stores"`
	Language       string                `json:"language"`
	Labels         map[string]string     `json:"labels"`
}

// ServePanel 仅由 Wails 资产处理器调用；公网 Echo 不注册面板和绑定接口。
func ServePanel(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.FS(panelFiles)).ServeHTTP(w, r)
}

// AttachPanel 将独立面板附着托盘，左键打开面板，右键保留原生菜单。
func AttachPanel(t *Tray, service *PanelService) {
	service.tray = t
	if t.icon == nil {
		return
	}
	panel := t.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name: "tray-panel", Title: "Comigo", URL: "/panel.html", Width: 420, Height: 590,
		Hidden: true, Frameless: true, AlwaysOnTop: true, DisableResize: true,
		Windows:     application.WindowsWindow{HiddenOnTaskbar: true},
		KeyBindings: map[string]func(application.Window){"Escape": func(w application.Window) { w.Hide() }},
	})
	// 使用原生窗口失焦事件，切换到阅读器或其他应用时收起面板。
	panel.OnWindowEvent(events.Common.WindowLostFocus, func(*application.WindowEvent) { panel.Hide() })
	// macOS 切换应用可能保留 key window，应用失活时也收起面板。
	t.app.Event.OnApplicationEvent(events.Mac.ApplicationDidResignActive, func(*application.ApplicationEvent) { panel.Hide() })
	panel.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) { panel.Hide(); e.Cancel() })
	t.icon.AttachWindow(panel).WindowOffset(6)
}

// readingURLs 为二维码提供本机已知地址；不接受任意网络目标。
func readingURLs() []string {
	result := []string{config.GetQrcodeURL()}
	if config.GetCfg().DisableLAN {
		return result
	}
	ips, _ := tools.GetIPList()
	for _, ip := range ips {
		parsed := net.ParseIP(ip)
		if parsed == nil || !parsed.IsGlobalUnicast() {
			continue
		}
		u, err := url.Parse(result[0])
		if err != nil {
			continue
		}
		port := u.Port()
		u.Host = net.JoinHostPort(ip, port)
		if port == "" {
			u.Host = strings.TrimSuffix(u.Host, ":")
		}
		if !slices.Contains(result, u.String()) {
			result = append(result, u.String())
		}
	}
	return result
}

// Snapshot 从内存取状态；二维码在本地编码，不请求网络或泄露阅读地址。
func (p *PanelService) Snapshot(selectedURL string) (PanelState, error) {
	cfg := config.GetCfg()
	state := PanelState{Version: config.GetVersion(), URLs: readingURLs(), Books: model.GetAllBooksNumber(), Traffic: traffic.Default.Snapshot(), ExternalAccess: !cfg.DisableLAN, ReadOnly: cfg.ReadOnlyMode, ListenAddress: config.GetListenHost(), Stores: slices.Clone(cfg.StoreUrls), Language: cfg.Language, Labels: map[string]string{}}
	for i, path := range state.Stores {
		state.Stores[i] = tools.LibraryLabel(path)
	}
	if !slices.Contains(state.URLs, selectedURL) {
		selectedURL = state.URLs[0]
	}
	state.ReadingURL = selectedURL
	png, err := qrcode.Encode(selectedURL, qrcode.Medium, 256)
	if err != nil {
		return state, err
	}
	state.QR = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	state.Connections, state.Users, state.Devices = data_api.ConnectionCounts()
	var status config.Status
	if err := status.SetConfigStatus(); err != nil {
		return state, err
	}
	state.Config = status.Current
	for _, key := range strings.Fields("nav_home nav_status nav_service nav_config reading qr open copy refresh server_status running version books online_users connections received sent total traffic_note external_access external_note listen_address config_file config_memory config_missing manage_config github website language libraries show_reader rescan config_directory error done auto close") {
		state.Labels[key] = locale.GetString("tray_panel_" + key)
	}
	return state, nil
}

// SetSharing 只改变局域网共享，其持久化与回滚沿用服务端事务。
func (p *PanelService) SetSharing(enabled bool) error { return routers.SetLANSharing(enabled) }

// SetLanguage 只允许界面支持的三种语言和系统自动选择。
func (p *PanelService) SetLanguage(language string) error {
	if !slices.Contains([]string{"auto", "zh-CN", "en-US", "ja-JP"}, language) {
		return errors.New("unsupported language")
	}
	return p.tray.setLanguage(language)
}

// Action 限定可执行操作；目录仅来自当前配置，不接收前端传入的磁盘路径。
func (p *PanelService) Action(action string, index int) error {
	switch action {
	case "reader":
		p.tray.Show()
	case "browser":
		return p.tray.app.Browser.OpenURL(config.GetLocalBrowserURL())
	case "copy":
		urls := readingURLs()
		if index < 0 || index >= len(urls) {
			return errors.New("invalid reading address")
		}
		if !p.tray.app.Clipboard.SetText(urls[index]) {
			return errors.New(locale.GetString("tray_panel_error"))
		}
	case "settings":
		p.tray.Show()
		p.tray.window.SetURL(config.PrefixPath("/settings"))
	case "rescan":
		if config.GetCfg().ReadOnlyMode {
			return errors.New(locale.GetString("err_config_locked"))
		}
		service.StartReScan()
	case "config":
		path, err := config.GetConfigDir()
		if err != nil {
			return err
		}
		return tools.OpenURL(path)
	case "store":
		stores := config.GetCfg().StoreUrls
		if index < 0 || index >= len(stores) {
			return errors.New("invalid library")
		}
		return tools.OpenURL(stores[index])
	case "github":
		return p.tray.app.Browser.OpenURL("https://github.com/yumenaka/comigo")
	case "website":
		return p.tray.app.Browser.OpenURL("https://comigo.xyz/")
	default:
		return errors.New("unknown panel action")
	}
	return nil
}
