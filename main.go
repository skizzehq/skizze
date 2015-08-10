package main

import (
	"counts/server"
	"counts/utils"
	"flag"
	"os"
	"strconv"
)

var logger = utils.GetLogger()

func main() {
	var port uint
	flag.UintVar(&port, "p", 7596, "specifies the port for Counts to run on")
	flag.Parse()

	//TODO: Add arguments for dataDir and infoDir

	os.Setenv("COUNTS_PORT", strconv.Itoa(int(port)))

	logger.Info.Println("Starting counts...")
	config := utils.GetConfig()
	logger.Info.Println("Using data dir: ", config.GetDataDir())
	server := server.New()
	server.Run()
}
