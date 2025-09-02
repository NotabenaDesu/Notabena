package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/georgysavva/scany/v2/sqlscan"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Note struct {
	Id      uint32
	Name    string
	Content string
	Created string
}

func InitDb(file string) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("Error while opening database file: %s", err)
	}
	defer db.Close()
	db.Exec(`CREATE TABLE IF NOT EXISTS saved_notes (
		id INTEGER PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		content TEXT NOT NULL,
		created TEXT NOT NULL
	);`)
}

func GetNotes(file string) []*Note {
	// TODO: finish rewriting this from 0.2
	// reference: https://github.com/The-Notabena-Organization/notabena-public-archive/blob/dev/src/api.rs
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("Error while opening database file: %s", err)
	}
	defer db.Close()
	notes := []*Note{}
	sqlscan.Select(context.Background(), db, &notes, "SELECT * FROM saved_notes;")
	return notes
}

func (note Note) Save(file string) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("Error while opening database file: %s", err)
	}
	defer db.Close()
	db.Exec("INSERT OR REPLACE INTO saved_notes (id, name, content, created) VALUES (?1, ?2, ?3, ?4);", note.Id, note.Name, note.Content, note.Created)
}

func (note Note) Delete(file string) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("Error while opening database file: %s", err)
	}
	defer db.Close()
	db.Exec("DELETE FROM saved_notes WHERE id = ?1;", note.Id)
}
