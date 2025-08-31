package main

import (
	"github.com/kirsle/configdir"
	"github.com/rivo/tview"
)

func main() {
	path := configdir.LocalConfig("Notabena")
	err := configdir.MakePath(path)
	if err != nil {
		panic(err)
	}

	InitDb(path)
	box := tview.NewBox().SetBorder(true).SetTitle("hello world")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
