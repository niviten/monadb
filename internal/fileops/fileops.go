package fileops

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func Create(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file already exists: %s", path)
		}
		return err
	}
	defer file.Close()
	return nil
}

// TODO: write FileReader
// instead of opening and closing file for every row, we can open -> read all rows SELECTED -> close
func Read(path string, offset int64, n int) ([]byte, error) {
	if offset < 0 {
		return nil, errors.New("offset cannot be negative")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := make([]byte, n)
	readBytes, err := file.ReadAt(buf, offset)
	if err != nil {
		if errors.Is(err, io.EOF) && readBytes < n {
			return nil, fmt.Errorf("requested %d bytes, only %d available", n, readBytes)
		}
		return nil, err
	}

	return buf, err
}

func Write(path string, offset int64, data []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return err
	}

	fileSize := info.Size()
	writeLen := int64(len(data))

	if offset < 0 {
		return errors.New("offset cannot be negative")
	}
	if offset+writeLen > fileSize {
		return fmt.Errorf("write range (%d to %d) exceeds file size (%d)", offset, offset+writeLen, fileSize)
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := f.WriteAt(data, offset)
	if err != nil {
		return err
	}
	if int64(n) != writeLen {
		return fmt.Errorf("partial write: expected %d bytes, wrote %d", writeLen, n)
	}

	return nil
}

func Append(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return err
	}
	defer f.Close()

	n, err := f.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("partial write: expected %d bytes, wrote %d", len(data), n)
	}

	return nil
}

func Truncate(path string) error {
	return os.Truncate(path, 0)
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// not sure about this
// func DeleteDirectoryRecursive(path string) error {
// 	return os.RemoveAll(path)
// }
