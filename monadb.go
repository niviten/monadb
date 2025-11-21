package monadb

import (
	"log"

	"github.com/niviten/monadb/internal/sqlite"
)

type Database struct {
	dataDir string
}

var db *Database = nil

func Start() {
	if db != nil {
		return
	}
	if err := sqlite.Init(); err != nil {
		log.Fatalln("Error connecting to sqlite db: ", err)
	}
	db = &Database{dataDir: dataDirPath}
	log.Println("monadb welcomes you")
}

func GetDB() *Database {
	return db
}
