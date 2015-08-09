package utils

/*
PanicOnError is a helper function to panic on Error
*/
func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}
