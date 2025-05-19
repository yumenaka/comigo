package main

import (
	"github.com/yumenaka/comigo/routers"
)

// env GOOS=js GOARCH=wasm CGO_ENABLED=0 go build -o your.wasm cmd/wasm/main.go

// package command-line-arguments
// imports github.com/yumenaka/comigo/routers
// imports github.com/yumenaka/comigo/routers/config_api
// imports github.com/yumenaka/comigo/util/scan
// imports github.com/yumenaka/comigo/internal/database
// imports modernc.org/sqlite
// imports modernc.org/libc
// imports modernc.org/libc/errno: build constraints exclude all Go files in /Users/bai/soft/gopath/pkg/mod/modernc.org/libc@v1.65.6/errno
func main() {
	// 初始化命令行flag，环境变量与配置文件
	// cmd.Execute()
	// 启动网页服务器（不阻塞）
	routers.StartWebServer()
	// 扫描书库（命令行指定）
	// cmd.ScanStore(cmd.Args)
	// 在命令行显示QRCode
	// cmd.ShowQRCode()
	// 退出时清理临时文件的处理函数
	// cmd.SetShutdownHandler()
}
