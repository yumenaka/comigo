//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest

package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/yumenaka/comi/config"
)

// tview sample
// https://github.com/letientai299/7guis/tree/master/tui
// https://github.com/devhulk/golang-gui/blob/main/main.go

// termui sample
// https://github.com/xxxserxxx/gotop

var (
	app *tview.Application // The tview application.
	//pages       *tview.Pages       // The application pages.
	//finderFocus tview.Primitive    // The primitive in the Finder that last had focus.
)

// @title Comigo API Service API 文档
// @version 1.0 版本
// @description Comigo API Service API 文档
// @Path /api/  基础路径
// @query.collection.format multi
func main() {
	runView()
	//cmd.Execute()
}

func runView() {
	//app := tview.NewApplication()
	//flex := tview.NewFlex().
	//	AddItem(tview.NewBox().SetBorder(true).SetTitle("Left (1/2 x width of Top)"), 0, 1, false).
	//	AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
	//		AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
	//		AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
	//		AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
	//	AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)
	//if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
	//	panic(err)
	//}

	title := config.Version
	box := tview.NewBox().
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle("Comigo " + title)
	app := tview.NewApplication().SetRoot(box, true)

	if err := app.Run(); err != nil {
		panic(err)
		//logger.Infof("Error running application: %s\n", err)
	}
}
