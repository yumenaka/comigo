package main

import (
	"fmt"
	"os"

	"github.com/yumenaka/comigo/cmd/tui"
)

// main 是默认的无托盘入口，适合 CLI、Linux 服务和 Docker 场景。
// 带系统托盘的桌面入口保留在 cmd/comigo/main.go。
func main() {
	if err := tui.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
