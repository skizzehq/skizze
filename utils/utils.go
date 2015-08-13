package utils

import "os"

/*
PanicOnError is a helper function to panic on Error
*/
func PanicOnError(e error) {
	if e != nil {
		logger.Error.Println(e)
		panic(e)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
