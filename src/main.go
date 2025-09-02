package main

import (
	"errors"
	"log"
	"os"

	"github.com/kirsle/configdir"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
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
				log.Fatalf("Can't create file: %s", err)
			}
		} else {
			log.Fatalf("Can't open file: %s", err)
		}
	}
	defer file.Close()
	db := InitDB(file.Name())
	Create(file, db)
}
