package monadb

type TableInfo struct {
	Name        string
	ColumnInfos []*ColumnInfo
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

func sizeOfDataType(dt DataType) int {
	switch dt {
	case INT32:
		return 4
	case UINT32:
		return 4
	case INT64:
		return 8
	case UINT64:
		return 8
	case VARCHAR:
		return -1
	}
	return -1
}
