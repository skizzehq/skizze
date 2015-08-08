package main

import (
	"counts/counters"
	"counts/server"
	"counts/utils"
	"flag"
)

var logger = utils.GetLogger()

func main() {
	var port = flag.String("p", "7596", "specifies the port for Counts to run on")
	flag.Parse()

	logger.Info.Println("Starting counts...")
	manager := counters.GetManager()
	server := server.New(manager)
	server.Run(*port)
}
