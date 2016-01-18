package utils

import (
	"os"

	"github.com/njpatel/loggo"
)

var logger = loggo.GetLogger("util")

// PanicOnError is a helper function to panic on Error
func PanicOnError(err error) {
	if err != nil {
		logger.Errorf("%v", err)
		panic(err)
	}
}

// Exists returns if path exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// CloseFile ...
func CloseFile(file *os.File) {
	err := file.Close()
	PanicOnError(err)
}

// GetFileSize ...
func GetFileSize(file *os.File) (int64, error) {
	stat, err := file.Stat()
	if err != nil {
		return -1, err
	}
	return stat.Size(), nil
}
