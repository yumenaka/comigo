//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest

package main

import (
	"github.com/yumenaka/comigo/cmd"
)

func main() {
	cmd.Execute()
}
