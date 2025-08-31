package main

import (
	"database/sql"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
)

type Note struct {
	Id      uint32
	Name    string
	Content string
	Created string
}

func InitDb(path string) {
	db, _ := sql.Open("sqlite3", path+"/notes.db")
	db.Exec(`CREATE TABLE IF NOT EXISTS saved_notes (
            id INTEGER PRIMARY KEY NOT NULL,
            name TEXT NOT NULL,
            content TEXT NOT NULL,
            created TEXT NOT NULL
    );`)
}

func GetNotes(path string) {
	// TODO: finish rewriting this from 0.2
	// reference: https://github.com/The-Notabena-Organization/notabena-public-archive/blob/dev/src/api.rs
	db, _ := sql.Open("sqlite3", path+"/notes.db")
	stmt, _ := db.Prepare("SELECT id, name, content, created FROM saved_notes;")
	noteIter, _ := stmt.Query()
	fmt.Println(noteIter)
}

func SaveNote(note Note, path string) {
	db, _ := sql.Open("sqlite3", path+"/notes.db")
	db.Exec("INSERT OR REPLACE INTO saved_notes (id, name, content, created) VALUES (?1, ?2, ?3, ?4);", note.Id, note.Name, note.Content, note.Created)
}

func DeleteNotes(notes []Note, path string) {
	db, _ := sql.Open("sqlite3", path+"/notes.db")
	for _, note := range notes {
		db.Exec("DELETE FROM saved_notes WHERE id = ?1;", note.Id)
	}
}
