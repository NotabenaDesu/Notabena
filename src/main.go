package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/kirsle/configdir"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/rivo/tview"
)

func main() {
	configPath := configdir.LocalConfig("Notabena")
	err := configdir.MakePath(configPath)
	if err != nil {
		log.Fatalf("No config folder found: %s", err)
	}
	path := configPath + "/notes.db"
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			file, err = os.Create(path)
			if err != nil {
				log.Fatalf("Can't create path: %s", err)
			}
		} else {
			log.Fatalln(err)
		}
	}
	defer file.Close()
	InitDb(file.Name())

	app := tview.NewApplication()
	textArea := tview.NewTextArea().SetWrap(true).SetPlaceholder("Write all your thoughts here! :D")
	textArea.SetTitle("New note").SetBorder(true)
	info := tview.NewTextView().SetText("Press Ctrl+X to save")
	position := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)
	pages := tview.NewPages()
	updateInfo := func() {
		fromRow, fromColumn, toRow, toColumn := textArea.GetCursor()
		if fromRow == toRow && fromColumn == toColumn {
			position.SetText(fmt.Sprintf("Note [yellow]#%d[white], Row: [yellow]%d[white], Column: [yellow]%d ", len(GetNotes(path)), fromRow, fromColumn))
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
			db, err := sql.Open("sqlite3", file.Name())
			if err != nil {
				log.Fatalf("Error while opening database file: %s", err)
			}
			defer db.Close()
			notes := []*Note{}
			sqlscan.Select(context.Background(), db, &notes, "SELECT id FROM saved_notes;")
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
