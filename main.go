//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest

package main

import (
	"github.com/yumenaka/comi/cmd"
)

// @title Comigo API Service API 文档
// @version 1.0 版本
// @description Comigo API Service API 文档
// @BasePath /api/  基础路径
// @query.collection.format multi
func main() {
	cmd.Execute()
}
