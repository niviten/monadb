package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var conn *sql.DB

func Init() error {
	dsn := "file:db/stuffs.db?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000"
	dbconn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	dbconn.SetMaxOpenConns(1)
	dbconn.SetMaxIdleConns(1)
	if err := dbconn.Ping(); err != nil {
		return err
	}
	conn = dbconn
	return nil
}

func GetConn() *sql.DB {
	return conn
}
