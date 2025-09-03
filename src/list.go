package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func List(file *os.File, db DB) {
	app := tview.NewApplication()

	mainView := tview.NewTreeNode("Welcome to Notabena!").SetColor(tcell.ColorMediumPurple)
	noteTree := tview.NewTreeView().
		SetRoot(mainView).
		SetCurrentNode(mainView)
	for _, v := range db.GetNotes() {
		stringId := strconv.FormatUint(uint64(v.Id), 10)
		node := tview.NewTreeNode(v.Name + " [grey]#" + stringId + "[white]")
		node.
			SetReference(v.Id).
			SetSelectable(true)
		node.SetExpanded(false)
		mainView.AddChild(node)
		node.AddChild(
			tview.NewTreeNode("Edit").SetReference("EDT+" + stringId).SetSelectable(true).SetColor(tcell.ColorLightCyan),
		)
		node.AddChild(
			tview.NewTreeNode("View").SetReference("VWR+" + stringId).SetSelectable(true),
		)
		node.AddChild(
			tview.NewTreeNode("Delete").SetReference("DEL+" + stringId).SetSelectable(true).SetColor(tcell.ColorRed),
		)
	}

	mainView.AddChild(
		tview.NewTreeNode("Create note!").SetReference("NEW").SetSelectable(true).SetColor(tcell.ColorBlue),
	)
	mainView.AddChild(
		tview.NewTreeNode("(Use Enter + arrow keys to navigate)").SetReference(nil).SetSelectable(false).SetColor(tcell.ColorDarkCyan),
	)

	noteTree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		if reference == "NEW" {
			app.Stop()
			Create(file, db, 0)
		}
		str, ok := reference.(string)
		if ok {
			if strings.HasPrefix(str, "DEL") {
				// this is a bit of a workaround
				part := strings.Split(str, "+")[1]
				num, err := strconv.ParseUint(part, 10, 32)
				if err != nil {
					panic(err)
				}
				db.DeleteNote(uint32(num))
				app.Stop()
				List(file, db)
			}
			if strings.HasPrefix(str, "EDT") {
				// this is a bit of a workaround
				part := strings.Split(str, "+")[1]
				num, err := strconv.ParseUint(part, 10, 32)
				if err != nil {
					panic(err)
				}
				app.Stop()
				Create(file, db, uint32(num))
			}
			if strings.HasPrefix(str, "VWR") {
				// this is a bit of a workaround
				part := strings.Split(str, "+")[1]
				num, err := strconv.ParseUint(part, 10, 32)
				if err != nil {
					panic(err)
				}
				app.Stop()
				View(file, db, uint32(num))
			}
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	if err := app.SetRoot(noteTree, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		log.Fatalf("Error while starting Notabena: %s", err)
	}
}
