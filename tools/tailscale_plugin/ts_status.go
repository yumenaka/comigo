package tailscale_plugin

import (
	"net/netip"
	"sync"
	"time"
)

type TailsClientInfo struct {
	LoginUserName     string    `json:"login_user_name"`     // 登录用户身份
	NodeComputedName  string    `json:"node_computed_name"`  // 节点的机器名
	RequestRemoteAddr string    `json:"request_remote_addr"` // 远程访问地址
	AccessTime        time.Time `json:"access_time"`         // 客户端访问开始、时间
}

// TailscaleStatus 保存 Tailscale 服务的状态信息
type TailscaleStatus struct {
	AuthURL          string            // 客户端授权用URL。如果节点已授权，则为空字符串。另外手动移除节点，不会重新生成新的AuthURL，直到下次重新运行tailscaled。
	BackendState     string            // BackendState: "NoState", "NeedsLogin", "NeedsMachineAuth", "Stopped", "Starting", "Running".	// 常见状态为 "Running"  "NeedsLogin"
	Clients          []TailsClientInfo // 曾与此节点连接的Tailscale客户端信息列表。
	OS               string            // "macOS" "windows" "linux" “android”，可能为空""（funnel-ingress-node）
	Online           bool              // 在线状态
	FQDN             string            // FQDN（Fully Qualified Domain Name 完全限定域名） = 子域名 + Domain → www.example.co.jp
	TailscaleIPs     []netip.Addr      // 分配给此节点的 Tailscale IP(s)，第一个是 IPv4地址，第二个是 IPv6 地址
	Version          string            // 当前Tailscale版本
	FunnelCapability string            // 是否支持 Funnel 功能 ,字符串：“true” “false” “unknown”
	mu               sync.Mutex        // Mutex lock
}

//
//func (t *TailscaleStatus) GetTailscaleIP() string {
//	t.mu.Lock()
//	defer t.mu.Unlock()
//	if len(t.TailscaleIPs) == 1 {
//		return t.TailscaleIPs[0].String()
//	}
//	if len(t.TailscaleIPs) >= 2 {
//		// 优先返回 IPv4 地址
//		for _, ip := range t.TailscaleIPs {
//			if ip.Is4() {
//				return ip.String()
//			}
//		}
//	}
//	return ""
//}
//
//// AddClientInfo 添加客户端信息到状态中，需要避免重复添加
//func (t *TailscaleStatus) AddClientInfo(tsClient *TailsClientInfo) {
//	t.mu.Lock()
//	defer t.mu.Unlock()
//	// 已存在相同的用户和节点名称，跳过
//	for _, client := range t.Clients {
//		if client.LoginUserName == tsClient.LoginUserName && client.NodeComputedName == tsClient.NodeComputedName {
//			return
//		}
//	}
//	t.Clients = append(t.Clients, *tsClient)
//}
//
//func (t *TailscaleStatus) CheckClientInfoExists(loginUserName, nodeComputedName string, requestRemoteAddr string) bool {
//	t.mu.Lock()
//	defer t.mu.Unlock()
//	// 已存在相同的用户和节点名称+访问地址，返回 true
//	for _, client := range t.Clients {
//		if client.LoginUserName == loginUserName && client.NodeComputedName == nodeComputedName && client.RequestRemoteAddr == requestRemoteAddr {
//			return true
//		}
//	}
//	return false
//}
//
//// GetTailscaleStatus 获取Tailscale服务状态 https://github.com/tailscale/golink/blob/b54cbbbb609ce8425193e7171a35af023cb5066d/golink.go#L787
//func GetTailscaleStatus(ctx context.Context) (*TailscaleStatus, error) {
//	if localClient == nil {
//		// Tailscale 未启用或尚未初始化时，返回离线状态而不是错误，避免上层页面 500
//		nowStatus.AuthURL = ""
//		nowStatus.BackendState = "Stopped"
//		nowStatus.TailscaleIPs = nil
//		nowStatus.Version = ""
//		nowStatus.OS = runtime.GOOS
//		nowStatus.FQDN = ""
//		nowStatus.Online = false
//		nowStatus.FunnelCapability = "unknown"
//		return nowStatus, nil
//	}
//	// *ipnstate.Status
//	st, err := localClient.Status(ctx)
//	if err != nil {
//		return nil, err
//	}
//	// 刷新状态
//	nowStatus.AuthURL = st.AuthURL
//	nowStatus.BackendState = st.BackendState
//	nowStatus.TailscaleIPs = st.TailscaleIPs
//	nowStatus.Version = st.Version
//	if st.Self != nil {
//		fullyQualifiedDomainName := st.Self.DNSName
//		if strings.HasSuffix(fullyQualifiedDomainName, ".") {
//			fullyQualifiedDomainName = fullyQualifiedDomainName[:len(fullyQualifiedDomainName)-1]
//		}
//		nowStatus.OS = st.Self.OS
//		nowStatus.FQDN = fullyQualifiedDomainName
//		nowStatus.Online = st.Self.Online
//		nowStatus.FunnelCapability = "false"
//		// 检查是否支持 Funnel 功能
//		if st.Self.CapMap != nil && st.Self.CapMap.Contains(tailcfg.NodeAttrFunnel) {
//			// 检查 st.Self.CapMap[tailcfg.NodeAttrFunnel] 是否存在
//			nowStatus.FunnelCapability = "true"
//		}
//	}
//	if st.Self == nil {
//		nowStatus.OS = runtime.GOOS
//		nowStatus.FQDN = ""
//		nowStatus.Online = false
//	}
//	return nowStatus, nil
//}
