package monadb

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/niviten/monadb/internal/fileops"
	"github.com/niviten/monadb/internal/sqlite"
	"github.com/niviten/monadb/internal/util"
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

func (db *Database) CreateTable(tableInfo *TableInfo) error {
	if err := util.ValidateTableName(tableInfo.Name); err != nil {
		log.Println("Error: invalid table name ", err)
		return fmt.Errorf("invalid table name")
	}

	conn := sqlite.GetConn()

	row := conn.QueryRow("SELECT id FROM mtable WHERE name = ?", tableInfo.Name)

	var _id uint64
	err := row.Scan(&_id)
	if err == nil {
		return fmt.Errorf("table already exists")
	}
	if err != sql.ErrNoRows {
		log.Println("Error while query table details for table name "+tableInfo.Name+" ", err)
		return fmt.Errorf("internal server error")
	}

	tx, err := conn.Begin()
	if err != nil {
		log.Println("Error while sqlite transaction begin ", err)
		return fmt.Errorf("internal server error")
	}

	res, err := tx.Exec("INSERT INTO mtable (name) VALUES (?)", tableInfo.Name)
	if err != nil {
		log.Println("Error while inserting table "+tableInfo.Name+" into db", err)
		return fmt.Errorf("internal server error")
	}
	tableId, err := res.LastInsertId()
	if err != nil {
		log.Println("Error while inserting table "+tableInfo.Name+" into db", err)
		return fmt.Errorf("internal server error")
	}

	columnInsertQuery := "INSERT INTO mcolumn (name, size, data_type_id, table_id) VALUES (?, ?, ?, ?)"
	for _, columnInfo := range tableInfo.ColumnInfos {
		size := columnInfo.Size
		if size <= 0 {
			size = sizeOfDataType(columnInfo.DataType)
		}
		if size <= 0 {
			rollbackTx(tx)
			return fmt.Errorf("invalid size for column name = %s, size = %d, data type = %d", columnInfo.Name, columnInfo.Size, columnInfo.DataType)
		}
		_, err = tx.Exec(columnInsertQuery, columnInfo.Name, columnInfo.Size, columnInfo.DataType, tableId)
		if err != nil {
			rollbackTx(tx)
			log.Printf("Error while inserting column %s: %s\n", columnInfo.Name, err.Error())
			return fmt.Errorf("internal server error")
		}
	}

	tablePath := filepath.Join(db.dataDir, dbName, tableInfo.Name)
	if fileops.DirExists(tablePath) {
		rollbackTx(tx)
		return fmt.Errorf("table already exists")
	}

	if err := fileops.CreateDir(tablePath); err != nil {
		rollbackTx(tx)
		log.Printf("Error: creating directory %s: %s\n", tablePath, err.Error())
		return fmt.Errorf("internal server error")
	}

	tableFilePath := filepath.Join(db.dataDir, dbName, tableInfo.Name, "table_data")
	if err = fileops.Create(tableFilePath); err != nil {
		rollbackTx(tx)
		log.Printf("Error while creating file %s: %s\n", tableFilePath, err.Error())
		return fmt.Errorf("internal server error")
	}

	tx.Commit()

	return nil
}

func rollbackTx(tx *sql.Tx) {
	if tx == nil {
		return
	}
	if err := tx.Rollback(); err != nil {
		log.Println("Error while rollback ", err)
	}
}
