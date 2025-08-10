// This program demonstrates how to use tsnet as a library.
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

	"tailscale.com/tsnet"
)

var (
	addr = flag.String("addr", ":443", "address to listen on")
)

func main() {
	flag.Parse()
	srv := new(tsnet.Server)
	defer srv.Close()
	// 设置 Tailscale 服务器的主机名。此名称将用于 Tailscale 网络中的节点标识。
	// 如果未设置，Tailscale 将使用机器的主机名。默认将会是tsnet。
	srv.Hostname = "comigo" // Set your desired hostname here
	// 监听器 ln 是一个 net.Listener 对象，它将处理来自 Tailscale网络的 TCP 连接。
	ln, err := srv.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	// LocalClient 返回一个与 s 通信的 LocalClient 对象。
	// 如果服务器尚未启动，它将启动服务器。如果服务器已成功启动，
	// 它不会返回错误。
	lc, err := srv.LocalClient()
	if err != nil {
		log.Fatal(err)
	}
	// 自动设置监听器 ln 的 TLS 证书。
	if *addr == ":443" {
		ln = tls.NewListener(ln, &tls.Config{
			GetCertificate: lc.GetCertificate,
		})
	}
	// 使用 http.Serve 来启动 HTTP 服务器，处理来自 Tailscale 网络的请求。
	// 这里使用 http.HandlerFunc 来处理 HTTP 请求。
	// 当有请求到达时，它将调用 WhoIs 方法来获取请求者的信息。
	// 如果成功，将返回一个 HTML 页面，显示请求者的用户名、节点名称和远程地址。如果获取信息失败，将返回 500 错误。
	log.Fatal(http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		who, err := lc.WhoIs(r.Context(), r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintf(w, "<html><body><h1>Hello, world!</h1>\n")
		fmt.Fprintf(w, "<p>You are <b>%s</b> from <b>%s</b> (%s)</p>",
			html.EscapeString(who.UserProfile.LoginName),
			html.EscapeString(firstLabel(who.Node.ComputedName)),
			r.RemoteAddr)
	})))
}

func firstLabel(s string) string {
	s, _, _ = strings.Cut(s, ".")
	return s
}
