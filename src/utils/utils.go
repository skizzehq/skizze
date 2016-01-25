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

// Stringp returns a pointer to a string, a convenience for dealing with protbuf generated code
func Stringp(s string) *string {
	return &s
}

// Int32p as above
func Int32p(i int32) *int32 {
	return &i
}

// Float32p as above
func Float32p(i float32) *float32 {
	return &i
}

// Int64p as above
func Int64p(i int64) *int64 {
	return &i
}

// Boolp as above
func Boolp(b bool) *bool {
	return &b
}
