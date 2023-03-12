//go:build windows
// +build windows

//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"github.com/yumenaka/comi/cmd"
)

// @title IvanApi Swagger 标题
// @version 1.0 版本
// @description IvanApi Service 描述
// @BasePath /api/v1  基础路径
// @query.collection.format multi
func main() {
	cmd.Execute()
}
