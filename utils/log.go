package utils

import (
	"log"
	"os"
)

/*
Logger is responsible for populating the logs
*/
type Logger struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

var logger *Logger

/*
GetLog returns a singleton Logger
*/
func GetLog() *Logger {

	if logger == nil {
		info := log.New(os.Stdout,
			"INFO: ",
			log.Ldate|log.Ltime)

		warning := log.New(os.Stdout,
			"WARNING: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		err := log.New(os.Stderr,
			"ERROR: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		logger = &Logger{info, warning, err}
	}
	return logger

}
