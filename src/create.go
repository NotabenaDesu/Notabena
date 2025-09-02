package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/rivo/tview"
)

func Create(file *os.File, db DB) {
	app := tview.NewApplication()
	textArea := tview.NewTextArea().SetWrap(true).SetPlaceholder("Write all your thoughts here! :D")
	textArea.SetTitle("New note").SetBorder(true)
	info := tview.NewTextView().SetText("Press Ctrl+X to save")
	position := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)
	pages := tview.NewPages()
	updateInfo := func() {
		fromRow, fromColumn, toRow, toColumn := textArea.GetCursor()
		if fromRow == toRow && fromColumn == toColumn {
			position.SetText(fmt.Sprintf("Note [yellow]#%d[white], Row: [yellow]%d[white], Column: [yellow]%d ", len(db.GetNotes()), fromRow, fromColumn))
		} else {
			position.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
		}
	}

	textArea.SetMovedFunc(updateInfo)
	updateInfo()
	mainView := tview.NewGrid().SetRows(0, 1).AddItem(textArea, 0, 0, 1, 2, 0, 0, true).AddItem(info, 1, 0, 1, 1, 0, 0, false).AddItem(position, 1, 1, 1, 1, 0, 0, false)
	saved := tview.NewTextView().SetDynamicColors(true).SetText(`[green]Ready to go!
[blue]Please give your note a name.`)

	savedPopup := tview.NewGrid()
	savedPopup.SetBorder(true).SetTitle("Success")
	savedPopup.AddItem(saved, 0, 0, 1, 2, 0, 0, false)
	titleInput := tview.NewTextArea().SetWrap(false).SetPlaceholder("Title here")
	titleInput.SetTitle("Title your note").SetBorder(true)
	savedPopup.AddItem(titleInput, 1, 0, 5, 3, 0, 0, true)
	savedPopup.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			app.Stop()
			notes := []*Note{}
			sqlscan.Select(context.Background(), db.Db, &notes, "SELECT id FROM saved_notes;")
			note := Note{
				Id:      uint32(len(notes)),
				Name:    titleInput.GetText(),
				Content: textArea.GetText(),
				Created: time.Now().Format("2006-01-02 15:04")}
			note.Save(file.Name())
			return nil
		}
		return event
	})

	pages.AddAndSwitchToPage("main", mainView, true).AddPage("saved", tview.NewGrid().SetColumns(0, 64, 0).SetRows(0, 22, 0).AddItem(savedPopup, 1, 1, 1, 1, 0, 0, true), true, false)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlX {
			pages.ShowPage("saved")
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		log.Fatalf("Error while starting Notabena: %s", err)
	}
}
