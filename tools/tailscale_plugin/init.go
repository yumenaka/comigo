package tailscale_plugin

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
	"tailscale.com/tsnet"
)

type TailscaleConfig struct {
	Hostname   string
	Port       uint16
	FunnelMode bool
	ConfigDir  string
	AuthKey    string
}

func InitTailscale(c TailscaleConfig) error {
	// 设置 Tailscale 服务器的主机名。此名称将用于 Tailscale 网络中的节点标识。
	// 如果未设置，Tailscale 将使用机器的主机名。默认将会是二进制文件名。
	tsServer = new(tsnet.Server)
	// 用于Tailscale网络中的标识节点，不影响监听地址
	tsServer.Hostname = c.Hostname
	tsServer.Port = c.Port
	// 如果提供了配置目录路径，则将其分配给 Tailscale 服务器的 Dir 字段。该目录用于存储 Tailscale 的状态和配置文件。
	if c.ConfigDir != "" {
		tsServer.Dir = c.ConfigDir
	}
	// 如果提供了身份验证密钥，则将其分配给 Tailscale 服务器的 AuthKey 字段。
	if c.AuthKey != "" {
		tsServer.AuthKey = c.AuthKey
	}

	// 监听器 ln 是一个 net.Listener 对象，它将处理来自 Tailscale网络的 TCP 连接。
	// Tailscale的Listen方法要求host部分必须是空的或者是IP字面量，不能使用主机名
	var err error
	listenAddr := ":" + fmt.Sprint(c.Port)

	if c.FunnelMode {
		// Funnel only supports TCP on ports 443, 8443, and 10000
		if c.Port != 443 && c.Port != 8443 && c.Port != 10000 {
			logger.Errorf(locale.GetString("err_funnel_mode_ports_only"))
			return errors.New(locale.GetString("err_funnel_mode_ports_only"))
		}
		netListener, err = tsServer.ListenFunnel("tcp", listenAddr)
		if err != nil {
			logger.Errorf(locale.GetString("err_failed_to_create_tailscale_funnel_listener"), listenAddr, err)
			return err
		}
	}
	if !c.FunnelMode {
		netListener, err = tsServer.Listen("tcp", listenAddr)
		if err != nil {
			logger.Errorf(locale.GetString("err_failed_to_create_tailscale_listener"), listenAddr, err)
			return err
		}
	}

	// LocalClient 返回一个与 s 通信的 LocalClient 对象。
	// 如果服务器尚未启动，它将启动服务器。如果服务器已成功启动，
	// 它不会返回错误。
	localClient, err = tsServer.LocalClient()
	if err != nil {
		logger.Errorf(locale.GetString("err_failed_to_create_tailscale_local_client"), err)
		if netListener != nil {
			_ = netListener.Close()
		}
		return err
	}
	// // 自动设置监听器 ln 的 TLS 证书
	if c.Port == 443 && !c.FunnelMode {
		netListener = tls.NewListener(netListener, &tls.Config{
			GetCertificate: localClient.GetCertificate,
		})
	}
	logger.Infof(locale.GetString("log_tailscale_server_initialized"), c.Hostname, c.Port)
	return nil
}
