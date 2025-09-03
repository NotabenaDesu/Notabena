package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func View(file *os.File, db DB, id uint32) {
	app := tview.NewApplication()
	textArea := tview.NewTextArea().SetWrap(true)
	textArea.SetDisabled(false)
	note := db.GetNote(id)
	textArea.SetText(note.Content, false)
	textArea.SetTitle("Viewing " + note.Name).SetBorder(true)
	info := tview.NewTextView().SetDynamicColors(true).SetText("Press Ctrl+X to exit [gray]Edits don't save, btw[white]")
	position := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)
	pages := tview.NewPages()
	updateInfo := func() {
		fromRow, fromColumn, toRow, toColumn := textArea.GetCursor()
		if fromRow == toRow && fromColumn == toColumn {
			position.SetText(fmt.Sprintf("Note [yellow]#%d[white], Created [yellow]%s[white], Row: [yellow]%d[white], Column: [yellow]%d ", len(db.GetNotes()), note.Created, fromRow, fromColumn))
		} else {
			position.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
		}
	}

	textArea.SetMovedFunc(updateInfo)
	updateInfo()
	mainView := tview.NewGrid().SetRows(0, 1).AddItem(textArea, 0, 0, 1, 2, 0, 0, true).AddItem(info, 1, 0, 1, 1, 0, 0, false).AddItem(position, 1, 1, 1, 1, 0, 0, false)

	pages.AddAndSwitchToPage("main", mainView, true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlX {
			app.Stop()
			List(file, db)
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		log.Fatalf("Error while starting Notabena: %s", err)
	}
}
