package util

import (
	"encoding/binary"
	"fmt"
)

func UInt8ToBytes(val uint8) []byte {
	return []byte{val}
}

func BytesToUInt8(b []byte) (uint8, error) {
	if err := checkSize(b, 1); err != nil {
		return 0, err
	}
	return uint8(b[0]), nil
}

func UInt16ToBytes(val uint16) []byte {
	var out [2]byte
	binary.BigEndian.PutUint16(out[:], val)
	return out[:]
}

func BytesToUInt16(b []byte) (uint16, error) {
	if err := checkSize(b, 2); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b), nil
}

func UInt32ToBytes(val uint32) []byte {
	var out [4]byte
	binary.BigEndian.PutUint32(out[:], val)
	return out[:]
}

func BytesToUInt32(b []byte) (uint32, error) {
	if err := checkSize(b, 4); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}

func UInt64ToBytes(val uint64) []byte {
	var out [8]byte
	binary.BigEndian.PutUint64(out[:], val)
	return out[:]
}

func BytesToUInt64(b []byte) (uint64, error) {
	if err := checkSize(b, 8); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}

func StringToBytes(val string, n int) ([]byte, error) {
	if n < 0 {
		return nil, fmt.Errorf("negative length: %d", n)
	}
	if len(val) > n {
		val = val[:n]
	}
	b := make([]byte, n)
	copy(b, val)
	return b, nil
}

func BytesToString(b []byte, n int) (string, error) {
	if err := checkSize(b, n); err != nil {
		return "", err
	}
	return string(b[:n]), nil
}

func checkSize(b []byte, want int) error {
	if len(b) < want {
		return fmt.Errorf("short buffer: need %d bytes, got %d", want, len(b))
	}
	return nil
}
