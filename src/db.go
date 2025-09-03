package main

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"github.com/georgysavva/scany/v2/sqlscan"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type DB struct {
	File string
	Db   *sql.DB
}

type Note struct {
	Id      uint32
	Name    string
	Content string
	Created string
}

func InitDB(file string) DB {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("Error while opening database file: %s", err)
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS saved_notes (
		id INTEGER PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		content TEXT NOT NULL,
		created TEXT NOT NULL
	);`)
	return DB{File: file, Db: db}
}

func (db DB) GetNotes() []*Note {
	notes := []*Note{}
	sqlscan.Select(context.Background(), db.Db, &notes, "SELECT * FROM saved_notes;")
	return notes
}

func (db DB) GetNote(id uint32) Note {
	var notes = db.GetNotes()
	for _, v := range notes {
		if v.Id == id {
			return *v
		}
	}
	panic("No note found with ID " + strconv.FormatUint(uint64(id), 10))
}

func (db DB) DeleteNote(id uint32) {
	db.Db.Exec(`DELETE FROM saved_notes WHERE id = ?1;`, id)
}

func (note Note) Save(file string) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("Error while opening database file: %s", err)
	}
	db.Exec("INSERT OR REPLACE INTO saved_notes (id, name, content, created) VALUES (?1, ?2, ?3, ?4);", note.Id, note.Name, note.Content, note.Created)
}

func (note Note) Delete(file string) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("Error while opening database file: %s", err)
	}
	db.Exec("DELETE FROM saved_notes WHERE id = ?1;", note.Id)
}
