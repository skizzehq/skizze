package main

import (
	"counts/server"
	"counts/utils"
	"flag"
)

var logger = utils.GetLogger()

func main() {
	var port = flag.String("p", "7596", "specifies the port for Counts to run on")
	flag.Parse()

	logger.Info.Println("Starting counts...")
	config := utils.GetConfig()
	logger.Info.Println("Using data dir: ", config.GetDataDir())
	server := server.New()
	server.Run(*port)
}
