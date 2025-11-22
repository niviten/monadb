package monadb

import (
	"fmt"

	"github.com/niviten/monadb/internal/util"
)

type TableInfo struct {
	Name        string
	ColumnInfos []*ColumnInfo
}

type DataType uint8

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

var InvalidDataType = fmt.Errorf("invalid data type")

func valToBytes(dt DataType, v any, size int) ([]byte, error) {
	switch dt {
	case INT32:
		val, ok := v.(int32)
		if !ok {
			return nil, InvalidDataType
		}
		return util.UInt32ToBytes(uint32(val)), nil
	case UINT32:
		val, ok := v.(uint32)
		if !ok {
			return nil, InvalidDataType
		}
		return util.UInt32ToBytes(val), nil
	case INT64:
		val, ok := v.(int64)
		if !ok {
			return nil, InvalidDataType
		}
		return util.UInt64ToBytes(uint64(val)), nil
	case UINT64:
		val, ok := v.(uint64)
		if !ok {
			return nil, InvalidDataType
		}
		return util.UInt64ToBytes(val), nil
	case VARCHAR:
		val, ok := v.(string)
		if !ok {
			return nil, InvalidDataType
		}
		return util.StringToBytes(val, size)
	}
	return nil, InvalidDataType
}
