//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest

package main

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/yumenaka/comi/cmd"
)

// tview sample
// https://github.com/rivo/tview/wiki/Postgres#source-code
// https://github.com/devhulk/golang-gui/blob/main/main.go

var (
	app *tview.Application // The tview application.
	//pages       *tview.Pages       // The application pages.
	//finderFocus tview.Primitive    // The primitive in the Finder that last had focus.
)

// @title Comigo API Service API 文档
// @version 1.0 版本
// @description Comigo API Service API 文档
// @BasePath /api/  基础路径
// @query.collection.format multi
func main() {
	box := tview.NewBox().
		SetBorder(true).
		SetTitle("Comigo Reader")
	app = tview.NewApplication().SetRoot(box, true)
	go cmd.Execute()
	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
}
