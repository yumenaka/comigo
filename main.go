//+build windows
//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"github.com/yumenaka/comi/cmd"
)

func main() {
	cmd.Execute()
}
