package main

import (
	"counts/counters"
	"counts/server"
	"counts/utils"
)

var logger = utils.GetLogger()

func main() {
	logger.Info.Println("Starting counts...")
	manager := counters.GetManager()
	server := server.New(manager)
	server.Run()
}
