package main

import (
	"flag"
	"os"
	"strconv"

	_ "net/http/pprof"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/server"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()

func main() {
	var port uint
	flag.UintVar(&port, "p", 3596, "specifies the port for Counts to run on")
	flag.Parse()

	//TODO: Add arguments for dataDir and infoDir

	os.Setenv("COUNTS_PORT", strconv.Itoa(int(port)))

	logger.Info.Println("Starting counts...")
	conf := config.GetConfig()
	logger.Info.Println("Using data dir: ", conf.GetDataDir())
	server, err := server.New()
	if err != nil {
		panic(err)
	}
	server.Run()
}
