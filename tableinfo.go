package monadb

type TableInfo struct {
	Name    string
	Columns []*ColumnInfo
}

type DataType int

const (
	INT32   = 1
	UINT32  = 2
	INT64   = 3
	UINT64  = 4
	VARCHAR = 5
)

type ColumnInfo struct {
	Name     string
	DataType DataType
	Size     int
}
