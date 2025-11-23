package monadb

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/niviten/monadb/internal/fileops"
	"github.com/niviten/monadb/internal/sqlite"
)

type mcolumn struct {
	name string
	size int
	dtid uint8
}

func Insert(tableName string, data map[string]any) (uint64, error) {
	conn := sqlite.GetConn()
	query := `
		SELECT c.name, c.size, c.data_type_id
		FROM mcolumn c
			LEFT JOIN mtable t
				ON c.table_id = t.id
		WHERE t.name = ?`
	rows, err := conn.Query(query, tableName)
	if err != nil {
		log.Printf("Error while querying columns for table name %s: %s\n", tableName, err.Error())
		return 0, fmt.Errorf("internal server error")
	}
	defer rows.Close()
	columnList := make([]*mcolumn, 0)
	for rows.Next() {
		mc := &mcolumn{}
		err = rows.Scan(&mc.name, &mc.size, &mc.dtid)
		if err != nil {
			log.Printf("Error while scanning table column for table name %s: %s\n", tableName, err.Error())
			return 0, fmt.Errorf("internal server error")
		}
		if mc.size == -1 {
			mc.size = sizeOfDataType(DataType(mc.dtid))
		}
		columnList = append(columnList, mc)
	}
	blist := make([][]byte, 0)
	totalSize := 0
	for _, col := range columnList {
		val := data[col.name]
		dataBytes, err := valToBytes(DataType(col.dtid), val, col.size)
		if err == ErrInvalidDataType {
			log.Printf("InvalidDataType error for %s :: datatype: %d :: val: %v\n", col.name, col.dtid, val)
			return 0, err
		}
		if err != nil {
			log.Printf("Error while converting val to bytes for column: %s, err: %s\n", col.name, err.Error())
			return 0, fmt.Errorf("internal server error")
		}
		blist = append(blist, dataBytes)
		totalSize = totalSize + len(dataBytes)
	}
	rowBytes := mergeBytes(totalSize, blist)

	tablePath := filepath.Join(dataDirPath, "default", tableName, "table_data")
	if err = fileops.Append(tablePath, rowBytes); err != nil {
		log.Printf("Error while inserting rowbytes to file path %s, error: %s\n", tablePath, err.Error())
		return 0, fmt.Errorf("internal server error")
	}

	return 0, nil
}

func mergeBytes(n int, blist [][]byte) []byte {
	out := make([]byte, n)
	off := 0
	for _, b := range blist {
		off = off + copy(out[off:], b)
	}
	return out
}
